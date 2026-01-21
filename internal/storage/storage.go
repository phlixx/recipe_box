package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const dirName = ".recipe_box"

type Storage struct {
	baseDir string
}

func New() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(homeDir, dirName)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	return &Storage{baseDir: baseDir}, nil
}

func (s *Storage) Save(filename string, data any) error {
	path := filepath.Join(s.baseDir, filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (s *Storage) Load(filename string, dest any) error {
	path := filepath.Join(s.baseDir, filename)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(dest)
}

func (s *Storage) Exists(filename string) bool {
	path := filepath.Join(s.baseDir, filename)
	_, err := os.Stat(path)
	return err == nil
}

func (s *Storage) Delete(filename string) error {
	path := filepath.Join(s.baseDir, filename)
	return os.Remove(path)
}

func (s *Storage) Path(filename string) string {
	return filepath.Join(s.baseDir, filename)
}
