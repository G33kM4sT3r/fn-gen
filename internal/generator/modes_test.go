package generator

import "testing"

func TestPattern_Lengths(t *testing.T) {
	tests := []struct {
		mode Mode
		want int
	}{
		{Minimal, 2},
		{Startup, 3},
		{Enterprise, 4},
		{Bullshit, 5},
	}

	for _, tt := range tests {
		t.Run(string(tt.mode), func(t *testing.T) {
			p := Pattern(tt.mode)
			if len(p) != tt.want {
				t.Errorf("Pattern(%q) has %d elements, want %d", tt.mode, len(p), tt.want)
			}
		})
	}
}

func TestPattern_AllStartWithAdjectives(t *testing.T) {
	modes := []Mode{Minimal, Startup, Enterprise, Bullshit}

	for _, m := range modes {
		p := Pattern(m)
		if p[0] != "adjectives" {
			t.Errorf("Pattern(%q)[0] = %q, want %q", m, p[0], "adjectives")
		}
	}
}

func TestPattern_AllExceptMinimalEndWithSuffix(t *testing.T) {
	modes := []Mode{Startup, Enterprise, Bullshit}

	for _, m := range modes {
		p := Pattern(m)
		if p[len(p)-1] != "suffix" {
			t.Errorf("Pattern(%q) last element = %q, want %q", m, p[len(p)-1], "suffix")
		}
	}
}

func TestPattern_UnknownModeFallback(t *testing.T) {
	p := Pattern("unknown")

	if len(p) != 2 {
		t.Errorf("unknown mode: got %d elements, want 2 (minimal fallback)", len(p))
	}
}
