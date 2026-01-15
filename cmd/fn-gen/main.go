package main

import (
	"fmt"
	"os"

	"fn-gen/internal/cli"
	"fn-gen/internal/generator"
	"fn-gen/internal/words"
)

func main() {
	// Parse command-line flags to get configuration
	cfg := cli.ParseFlags()

	// Load the word set for the specified language and mode
	// Each language/mode combination has its own JSON file with word pools
	wordSet, err := words.Load(cfg.Lang, cfg.Mode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Initialize the generator with the loaded words and configuration
	gen := generator.New(wordSet, cfg)

	// Generate the requested number of names
	// Each iteration uses a unique index for seed differentiation
	for i := 0; i < cfg.Count; i++ {
		if cfg.Explain {
			// Explain mode: show detailed breakdown of how the name was generated
			result := gen.GenerateExplained(i)

			fmt.Println(result.Name)
			fmt.Println("— explanation —")
			fmt.Printf("seed: %s\n", result.Seed)
			fmt.Printf("pattern: %v\n", result.Pattern)

			// Print details for each word part showing the hash calculation
			for _, p := range result.Parts {
				fmt.Printf(
					"- %s: %q (hash=%d index=%d/%d)\n",
					p.Category,
					p.Word,
					p.Hash,
					p.Index,
					p.ListSize,
				)
			}
			fmt.Println()
		} else {
			// Standard mode: just print the generated name
			fmt.Println(gen.Generate(i))
		}
	}
}
