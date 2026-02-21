package words

import (
	"os"
	"path/filepath"
	"testing"
)

// Tests run from the package directory, but Load expects paths relative to project root.
// chdir to project root before each test that calls Load.
func chdirToRoot(t *testing.T) {
	t.Helper()
	// Walk up from internal/words to project root
	root := filepath.Join("..", "..")
	if err := os.Chdir(root); err != nil {
		t.Fatalf("cannot chdir to project root: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(filepath.Join("internal", "words"))
	})
}

func TestLoad_EnglishMinimal(t *testing.T) {
	chdirToRoot(t)

	ws, err := Load("en", "minimal")
	if err != nil {
		t.Fatalf("Load(en, minimal) error: %v", err)
	}

	if len(ws.Adjectives) == 0 {
		t.Error("adjectives should not be empty")
	}
	if len(ws.Core) == 0 {
		t.Error("core should not be empty")
	}
	// Minimal mode has empty buzzwords and suffix
	if len(ws.Buzzwords) != 0 {
		t.Errorf("minimal buzzwords should be empty, got %d", len(ws.Buzzwords))
	}
	if len(ws.Suffix) != 0 {
		t.Errorf("minimal suffix should be empty, got %d", len(ws.Suffix))
	}
}

func TestLoad_AllLanguagesAndModes(t *testing.T) {
	chdirToRoot(t)

	langs := []string{"en", "de"}
	modes := []string{"minimal", "startup", "enterprise", "bullshit"}

	for _, lang := range langs {
		for _, mode := range modes {
			t.Run(lang+"/"+mode, func(t *testing.T) {
				ws, err := Load(lang, mode)
				if err != nil {
					t.Fatalf("Load(%s, %s) error: %v", lang, mode, err)
				}
				if len(ws.Adjectives) == 0 {
					t.Error("adjectives should not be empty")
				}
				if len(ws.Core) == 0 {
					t.Error("core should not be empty")
				}
			})
		}
	}
}

func TestLoad_NonMinimalModesHaveSuffix(t *testing.T) {
	chdirToRoot(t)

	modes := []string{"startup", "enterprise", "bullshit"}

	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			ws, err := Load("en", mode)
			if err != nil {
				t.Fatalf("Load(en, %s) error: %v", mode, err)
			}
			if len(ws.Suffix) == 0 {
				t.Errorf("mode %s should have suffix words", mode)
			}
		})
	}
}

func TestLoad_EnterpriseAndBullshitHaveBuzzwords(t *testing.T) {
	chdirToRoot(t)

	for _, mode := range []string{"enterprise", "bullshit"} {
		t.Run(mode, func(t *testing.T) {
			ws, err := Load("en", mode)
			if err != nil {
				t.Fatalf("Load(en, %s) error: %v", mode, err)
			}
			if len(ws.Buzzwords) == 0 {
				t.Errorf("mode %s should have buzzwords", mode)
			}
		})
	}
}

func TestLoad_InvalidLanguage(t *testing.T) {
	chdirToRoot(t)

	_, err := Load("xx", "startup")
	if err == nil {
		t.Error("expected error for invalid language, got nil")
	}
}

func TestLoad_InvalidMode(t *testing.T) {
	chdirToRoot(t)

	_, err := Load("en", "nonexistent")
	if err == nil {
		t.Error("expected error for invalid mode, got nil")
	}
}

func TestGet_ValidKeys(t *testing.T) {
	ws := WordSet{
		Adjectives: []string{"Smart"},
		Buzzwords:  []string{"Cloud"},
		Core:       []string{"Engine"},
		Suffix:     []string{"Hub"},
	}

	tests := []struct {
		key  string
		want string
	}{
		{"adjectives", "Smart"},
		{"buzzwords", "Cloud"},
		{"core", "Engine"},
		{"suffix", "Hub"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := ws.Get(tt.key)
			if len(got) != 1 || got[0] != tt.want {
				t.Errorf("Get(%q) = %v, want [%q]", tt.key, got, tt.want)
			}
		})
	}
}

func TestGet_UnknownKey(t *testing.T) {
	ws := WordSet{Adjectives: []string{"Smart"}}

	got := ws.Get("unknown")
	if got != nil {
		t.Errorf("Get(unknown) = %v, want nil", got)
	}
}
