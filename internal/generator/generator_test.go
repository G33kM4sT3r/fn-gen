package generator

import (
	"fmt"
	"testing"

	"fn-gen/internal/cli"
	"fn-gen/internal/words"
)

func testWordSet() words.WordSet {
	return words.WordSet{
		Adjectives: []string{"Smart", "Fast", "Bold"},
		Buzzwords:  []string{"Cloud", "AI", "Quantum"},
		Core:       []string{"Engine", "Pipeline", "Gateway"},
		Suffix:     []string{"Hub", "Pro", "Plus"},
	}
}

func largeWordSet() words.WordSet {
	adj := make([]string, 50)
	buz := make([]string, 50)
	cor := make([]string, 50)
	suf := make([]string, 50)
	for i := range 50 {
		adj[i] = fmt.Sprintf("Adj%d", i)
		buz[i] = fmt.Sprintf("Buzz%d", i)
		cor[i] = fmt.Sprintf("Core%d", i)
		suf[i] = fmt.Sprintf("Suf%d", i)
	}
	return words.WordSet{
		Adjectives: adj,
		Buzzwords:  buz,
		Core:       cor,
		Suffix:     suf,
	}
}

func testConfig(mode, seed string) cli.Config {
	return cli.Config{
		Lang: "en",
		Mode: mode,
		Seed: seed,
	}
}

func TestGenerate_Deterministic(t *testing.T) {
	ws := testWordSet()
	cfg := testConfig("startup", "test-seed")

	g1 := New(ws, cfg)
	g2 := New(ws, cfg)

	name1 := g1.Generate(0)
	name2 := g2.Generate(0)

	if name1 != name2 {
		t.Errorf("same seed produced different names: %q vs %q", name1, name2)
	}
}

func TestGenerate_DifferentSeeds(t *testing.T) {
	ws := largeWordSet()

	g1 := New(ws, testConfig("startup", "seed-a"))
	g2 := New(ws, testConfig("startup", "seed-b"))

	name1 := g1.Generate(0)
	name2 := g2.Generate(0)

	if name1 == name2 {
		t.Errorf("different seeds produced same name: %q", name1)
	}
}

func TestGenerate_DifferentIndicesWithAutoSeed(t *testing.T) {
	ws := largeWordSet()
	// Index only differentiates names when using auto-seed (empty seed).
	// With a custom seed, the index parameter is not incorporated.
	cfg := testConfig("startup", "")
	g := New(ws, cfg)

	name0 := g.Generate(0)
	name1 := g.Generate(1)

	if name0 == name1 {
		t.Errorf("different indices with auto-seed produced same name: %q", name0)
	}
}

func TestGenerate_CustomSeedIgnoresIndex(t *testing.T) {
	ws := testWordSet()
	cfg := testConfig("startup", "fixed-seed")
	g := New(ws, cfg)

	name0 := g.Generate(0)
	name1 := g.Generate(1)

	if name0 != name1 {
		t.Errorf("custom seed should ignore index, got %q and %q", name0, name1)
	}
}

func TestGenerate_WordCountPerMode(t *testing.T) {
	ws := testWordSet()

	tests := []struct {
		mode      string
		wantWords int
	}{
		{"minimal", 2},
		{"startup", 3},
		{"enterprise", 4},
		{"bullshit", 5},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			g := New(ws, testConfig(tt.mode, "count-test"))
			result := g.GenerateExplained(0)

			if len(result.Parts) != tt.wantWords {
				t.Errorf("mode %q: got %d words, want %d", tt.mode, len(result.Parts), tt.wantWords)
			}
		})
	}
}

func TestGenerateExplained_Metadata(t *testing.T) {
	ws := testWordSet()
	cfg := testConfig("startup", "meta-test")
	g := New(ws, cfg)

	result := g.GenerateExplained(0)

	if result.Seed != "meta-test" {
		t.Errorf("seed = %q, want %q", result.Seed, "meta-test")
	}

	wantPattern := []string{"adjectives", "core", "suffix"}
	if len(result.Pattern) != len(wantPattern) {
		t.Fatalf("pattern length = %d, want %d", len(result.Pattern), len(wantPattern))
	}
	for i, p := range result.Pattern {
		if p != wantPattern[i] {
			t.Errorf("pattern[%d] = %q, want %q", i, p, wantPattern[i])
		}
	}

	for _, part := range result.Parts {
		if part.Word == "" {
			t.Error("part has empty word")
		}
		if part.ListSize != 3 {
			t.Errorf("category %q: ListSize = %d, want 3", part.Category, part.ListSize)
		}
		if part.Index >= uint64(part.ListSize) {
			t.Errorf("category %q: Index %d >= ListSize %d", part.Category, part.Index, part.ListSize)
		}
	}
}

func TestGenerate_EmptyCategory(t *testing.T) {
	ws := words.WordSet{
		Adjectives: []string{"Smart"},
		Core:       []string{"Engine"},
		Suffix:     []string{"Hub"},
	}
	cfg := testConfig("enterprise", "empty-test")
	g := New(ws, cfg)

	result := g.GenerateExplained(0)

	// Enterprise pattern has 4 categories, but buzzwords is empty → 3 parts
	if len(result.Parts) != 3 {
		t.Errorf("got %d parts, want 3 (empty buzzwords should be skipped)", len(result.Parts))
	}
}

func TestGenerate_AutoSeed(t *testing.T) {
	ws := testWordSet()
	cfg := testConfig("minimal", "")
	g := New(ws, cfg)

	result := g.GenerateExplained(0)

	if result.Seed == "" {
		t.Error("auto seed should not be empty")
	}
	if result.Name == "" {
		t.Error("name should not be empty with auto seed")
	}
}

func TestGenerate_UnknownModeFallsBackToMinimal(t *testing.T) {
	ws := testWordSet()
	cfg := testConfig("nonexistent", "fallback-test")
	g := New(ws, cfg)

	result := g.GenerateExplained(0)

	if len(result.Parts) != 2 {
		t.Errorf("unknown mode: got %d parts, want 2 (minimal fallback)", len(result.Parts))
	}
}
