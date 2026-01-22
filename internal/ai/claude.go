package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"recipe_box/internal/recipe"
)

const (
	defaultAPIURL = "https://api.anthropic.com/v1/messages"
	apiVersion    = "2023-06-01"
	model         = "claude-sonnet-4-20250514"
)

// getAPIURL returns the API URL, allowing override via environment variable for testing
func getAPIURL() string {
	if url := os.Getenv("ANTHROPIC_API_URL"); url != "" {
		return url
	}
	return defaultAPIURL
}

var (
	ErrNoAPIKey       = errors.New("API key not configured - run: recipe_box config set api_key <your-key>")
	ErrGenerateFailed = errors.New("failed to generate recipe")
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}, nil
}

type GenerateOptions struct {
	Prompt      string
	Ingredients []string
	Cuisine     string
	Vegetarian  bool
	Quick       bool   // under 30 min
	Language    string // "en" or "de"
	Servings    int    // target servings (0 = default/4)
}

func (c *Client) GenerateRecipe(opts GenerateOptions) (*recipe.Recipe, error) {
	// Default to English if not specified
	if opts.Language == "" {
		opts.Language = "en"
	}

	prompt := buildPrompt(opts)

	reqBody := messagesRequest{
		Model:     model,
		MaxTokens: 2048,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", getAPIURL(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s (status %d)", ErrGenerateFailed, string(respBody), resp.StatusCode)
	}

	var messagesResp messagesResponse
	if err := json.Unmarshal(respBody, &messagesResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(messagesResp.Content) == 0 {
		return nil, fmt.Errorf("%w: empty response", ErrGenerateFailed)
	}

	return parseRecipeFromResponse(messagesResp.Content[0].Text, opts.Language)
}

func buildPrompt(opts GenerateOptions) string {
	var constraints []string

	if opts.Cuisine != "" {
		constraints = append(constraints, fmt.Sprintf("Cuisine: %s", opts.Cuisine))
	}
	if opts.Vegetarian {
		constraints = append(constraints, "Must be vegetarian (no meat or fish)")
	}
	if opts.Quick {
		constraints = append(constraints, "Total time (prep + cook) must be under 30 minutes")
	}
	if len(opts.Ingredients) > 0 {
		constraints = append(constraints, fmt.Sprintf("Must use these ingredients: %s", strings.Join(opts.Ingredients, ", ")))
	}
	if opts.Servings > 0 {
		constraints = append(constraints, fmt.Sprintf("Recipe must be for exactly %d servings", opts.Servings))
	}

	var userRequest string
	if opts.Prompt != "" {
		userRequest = opts.Prompt
	} else {
		userRequest = "a delicious recipe"
	}

	constraintText := ""
	if len(constraints) > 0 {
		constraintText = "\n\nConstraints:\n- " + strings.Join(constraints, "\n- ")
	}

	// Language instruction
	langInstruction := "Write all text (title, description, ingredients, steps) in English."
	if opts.Language == "de" {
		langInstruction = "Write all text (title, description, ingredients, steps) in German (Deutsch)."
	}

	return fmt.Sprintf(`Generate a recipe for: %s%s

%s

Return ONLY valid JSON in this exact format (no markdown, no explanation):
{
  "title": "Recipe Title",
  "description": "Brief description of the dish",
  "servings": 4,
  "prep_time": 15,
  "cook_time": 30,
  "ingredients": [
    {"name": "ingredient name", "quantity": 2.0, "unit": "cups", "category": "produce"},
    {"name": "another ingredient", "quantity": 1.5, "unit": "tbsp", "category": "pantry"}
  ],
  "steps": [
    "Step 1 description",
    "Step 2 description"
  ],
  "tags": ["tag1", "tag2"],
  "ascii_art": "small ASCII art of the dish (8-12 lines, max 40 chars wide)"
}

Categories for ingredients must be one of: produce, dairy, meat, pantry, frozen, other
Units should be metric where appropriate (g, ml, etc.) or common cooking units (cups, tbsp, tsp)
The ascii_art should be a small, simple ASCII art representation of the finished dish using basic characters. Use \n for newlines within the string.`, userRequest, constraintText, langInstruction)
}

func parseRecipeFromResponse(text string, language string) (*recipe.Recipe, error) {
	// Try to extract JSON from the response (Claude might wrap it)
	text = strings.TrimSpace(text)

	// If wrapped in markdown code block, extract it
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		var jsonLines []string
		inBlock := false
		for _, line := range lines {
			if strings.HasPrefix(line, "```") {
				inBlock = !inBlock
				continue
			}
			if inBlock {
				jsonLines = append(jsonLines, line)
			}
		}
		text = strings.Join(jsonLines, "\n")
	}

	var r recipe.Recipe
	if err := json.Unmarshal([]byte(text), &r); err != nil {
		return nil, fmt.Errorf("failed to parse recipe JSON: %w\nResponse was: %s", err, text)
	}

	r.Source = recipe.SourceAI
	if language == "de" {
		r.Language = recipe.LangDE
	} else {
		r.Language = recipe.LangEN
	}

	return &r, nil
}

// API request/response types
type messagesRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type messagesResponse struct {
	Content []contentBlock `json:"content"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
