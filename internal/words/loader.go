package words

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type WordSet struct {
	Adjectives []string `json:"adjectives"` // Descriptive words (Smart, Dynamic, Scalable, ...)
	Buzzwords  []string `json:"buzzwords"`  // Trendy tech terms (Cloud, AI-Assisted, Serverless, ...)
	Core       []string `json:"core"`       // Central concept words (Workflow, Data, Integration, ...)
	Suffix     []string `json:"suffix"`     // Ending words (Hub, Engine, Platform, ...)
}

// Load reads a word set from a JSON file based on language and mode.
// The file path is constructed as: internal/words/data/{lang}/{mode}.json
//
// Parameters:
//   - lang: Language code (e.g., "en", "de")
//   - mode: Generation mode (e.g., "startup", "enterprise")
//
// Returns an error if the file cannot be read or parsed.
//
// Example file paths:
//   - internal/words/data/en/startup.json
//   - internal/words/data/de/enterprise.json
func Load(lang, mode string) (WordSet, error) {
	// Construct the path to the word file
	path := filepath.Join("internal", "words", "data", lang, mode+".json")

	// Read the entire file content
	data, err := os.ReadFile(path)
	if err != nil {
		return WordSet{}, fmt.Errorf("cannot load words: %w", err)
	}

	// Parse JSON into WordSet struct
	var ws WordSet
	if err := json.Unmarshal(data, &ws); err != nil {
		return WordSet{}, err
	}

	return ws, nil
}

// Get retrieves the word list for a given category key.
// This provides a dynamic way to access word pools by name,
// which is used by the generator when iterating through patterns.
//
// Valid keys: "adjectives", "buzzwords", "core", "suffix"
// Returns nil for unknown keys.
func (w WordSet) Get(key string) []string {
	switch key {
	case "adjectives":
		return w.Adjectives
	case "buzzwords":
		return w.Buzzwords
	case "core":
		return w.Core
	case "suffix":
		return w.Suffix
	default:
		return nil // Unknown category
	}
}
