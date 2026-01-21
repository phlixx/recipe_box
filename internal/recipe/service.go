package recipe

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"

	"recipe_box/internal/storage"
)

const recipesFile = "recipes.json"

var (
	ErrNotFound = errors.New("recipe not found")
)

// Storage defines the interface for recipe persistence
type Storage interface {
	Save(filename string, data any) error
	Load(filename string, dest any) error
	Exists(filename string) bool
	Delete(filename string) error
}

type Service struct {
	storage Storage
}

func NewService(s *storage.Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) Add(r *Recipe) error {
	if r.ID == "" {
		r.ID = generateID()
	}

	recipes, err := s.List()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	recipes = append(recipes, *r)
	return s.storage.Save(recipesFile, recipes)
}

func (s *Service) List() ([]Recipe, error) {
	var recipes []Recipe

	if !s.storage.Exists(recipesFile) {
		return recipes, nil
	}

	err := s.storage.Load(recipesFile, &recipes)
	return recipes, err
}

func (s *Service) Get(id string) (*Recipe, error) {
	recipes, err := s.List()
	if err != nil {
		return nil, err
	}

	for _, r := range recipes {
		if r.ID == id {
			return &r, nil
		}
	}

	return nil, ErrNotFound
}

func (s *Service) Delete(id string) error {
	recipes, err := s.List()
	if err != nil {
		return err
	}

	found := false
	filtered := make([]Recipe, 0, len(recipes))
	for _, r := range recipes {
		if r.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, r)
	}

	if !found {
		return ErrNotFound
	}

	return s.storage.Save(recipesFile, filtered)
}

func generateID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
