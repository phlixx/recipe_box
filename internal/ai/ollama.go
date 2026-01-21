package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"recipe_box/internal/recipe"
)

const (
	defaultOllamaURL   = "http://localhost:11434/api/chat"
	defaultOllamaModel = "llama3.2"
)

var (
	ErrOllamaNotRunning = errors.New("Ollama not running - start it with: ollama serve")
	ErrGenerateFailed   = errors.New("failed to generate recipe")
)

type Client struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewClient(model string) *Client {
	if model == "" {
		model = defaultOllamaModel
	}
	return &Client{
		baseURL:    defaultOllamaURL,
		model:      model,
		httpClient: &http.Client{},
	}
}

type GenerateOptions struct {
	Prompt      string
	Ingredients []string
	Cuisine     string
	Vegetarian  bool
	Quick       bool // under 30 min
}

func (c *Client) GenerateRecipe(opts GenerateOptions) (*recipe.Recipe, error) {
	prompt := buildPrompt(opts)

	reqBody := ollamaRequest{
		Model:  c.model,
		Stream: false,
		Messages: []ollamaMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOllamaNotRunning, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s (status %d)", ErrGenerateFailed, string(respBody), resp.StatusCode)
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return parseRecipeFromResponse(ollamaResp.Message.Content)
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

	return fmt.Sprintf(`Generate a recipe for: %s%s

Return ONLY valid JSON in this exact format (no markdown, no explanation, no extra text):
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
  "tags": ["tag1", "tag2"]
}

Categories for ingredients must be one of: produce, dairy, meat, pantry, frozen, other
Units should be metric where appropriate (g, ml, etc.) or common cooking units (cups, tbsp, tsp)

IMPORTANT: Return ONLY the JSON object, nothing else.`, userRequest, constraintText)
}

func parseRecipeFromResponse(text string) (*recipe.Recipe, error) {
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

	// Try to find JSON object if there's extra text
	startIdx := strings.Index(text, "{")
	endIdx := strings.LastIndex(text, "}")
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		text = text[startIdx : endIdx+1]
	}

	var r recipe.Recipe
	if err := json.Unmarshal([]byte(text), &r); err != nil {
		return nil, fmt.Errorf("failed to parse recipe JSON: %w\nResponse was: %s", err, text)
	}

	r.Source = recipe.SourceAI
	r.Language = recipe.LangEN

	return &r, nil
}

// Ollama API types
type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message ollamaMessage `json:"message"`
}
