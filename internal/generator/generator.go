package generator

import (
	"fmt"
	"strings"
	"time"

	"fn-gen/internal/cli"
	"fn-gen/internal/words"
)

type Generator struct {
	words words.WordSet // Word pools for each category (adjectives, buzzwords, etc.)
	cfg   cli.Config    // User configuration from CLI flags
}

type ExplainedPart struct {
	Category string // Word category (e.g., "adjectives", "core", "suffix")
	Word     string // The selected word from the category
	Hash     uint64 // Raw hash value computed from the seed
	Index    uint64 // Array index after modulo operation (Hash % ListSize)
	ListSize int    // Total number of words available in this category
}

type ExplainedResult struct {
	Name    string          // The final generated feature name
	Seed    string          // The seed used for generation (auto or user-provided)
	Pattern []string        // The word category pattern used (e.g., ["adjectives", "core", "suffix"])
	Parts   []ExplainedPart // Detailed breakdown of each word selection
}

// New creates a new Generator instance with the given word set and configuration.
// The generator is ready to produce names immediately after creation.
func New(words words.WordSet, cfg cli.Config) *Generator {
	return &Generator{words: words, cfg: cfg}
}

// Generate produces a single feature name for the given index.
// This is a convenience wrapper around GenerateExplained that returns only the name.
//
// The index parameter differentiates names when generating multiple names
// in a single run (used with -count flag).
func (g *Generator) Generate(index int) string {
	return g.GenerateExplained(index).Name
}

// GenerateExplained produces a feature name with full generation details.
// It returns an ExplainedResult containing the name and metadata about
// how each word was selected.
//
// The generation process:
//  1. Determine the word pattern based on the configured mode
//  2. Construct the seed (use provided seed or generate automatic one)
//  3. For each word category in the pattern:
//     a. Compute a unique hash using the seed, position, and category
//     b. Use the hash to select a word from the category's word list
//  4. Join all selected words with spaces to form the final name
func (g *Generator) GenerateExplained(index int) ExplainedResult {
	// Get the word pattern for the current mode (e.g., ["adjectives", "core", "suffix"])
	pattern := Pattern(Mode(g.cfg.Mode))

	// Determine the seed to use for hash generation
	baseSeed := g.cfg.Seed
	if baseSeed == "" {
		// No user seed provided: generate automatic seed from config + date
		// Format: "{lang}-{mode}-{index}-{date}"
		// This makes names reproducible within the same day
		baseSeed = fmt.Sprintf(
			"%s-%s-%d-%s",
			g.cfg.Lang,
			g.cfg.Mode,
			index,
			time.Now().Format("2006-01-02"),
		)
	}

	var parts []ExplainedPart
	var nameParts []string

	// Iterate through each word category in the pattern
	for i, key := range pattern {
		// Get the word list for this category
		list := g.words.Get(key)
		if len(list) == 0 {
			continue // Skip empty categories
		}

		// Compute a unique hash for this word position
		// Format: "{seed}-{position}-{category}"
		// This ensures each position gets a different word even with the same seed
		hash := HashToUint64(fmt.Sprintf("%s-%d-%s", baseSeed, i, key))

		// Use modulo to convert the hash to a valid array index
		idx := hash % uint64(len(list))

		// Select the word at the computed index
		word := list[idx]

		nameParts = append(nameParts, word)

		// Store detailed information for explain mode
		parts = append(parts, ExplainedPart{
			Category: key,
			Word:     word,
			Hash:     hash,
			Index:    idx,
			ListSize: len(list),
		})
	}

	return ExplainedResult{
		Name:    strings.Join(nameParts, " "), // Combine words with spaces
		Seed:    baseSeed,
		Pattern: pattern,
		Parts:   parts,
	}
}
