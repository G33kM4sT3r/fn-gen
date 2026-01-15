package generator

type Mode string

const (
	Minimal    Mode = "minimal"    // Concise, no-frills names
	Startup    Mode = "startup"    // Balanced startup-style names
	Enterprise Mode = "enterprise" // Corporate-sounding names with buzzwords
	Bullshit   Mode = "bullshit"   // Over-the-top buzzword-heavy names
)

// Pattern returns the ordered list of word categories for a given mode.
// Each category corresponds to a key in the WordSet struct and determines
// which word pool is used for that position in the generated name.
//
// Word categories:
//   - "adjectives": Descriptive words (Smart, Dynamic, Scalable, ...)
//   - "buzzwords":  Trendy tech terms (Cloud, AI-Assisted, Serverless, ...)
//   - "core":       Central concept words (Workflow, Data, Integration, ...)
//   - "suffix":     Ending words (Hub, Engine, Platform, ...)
//
// Example patterns:
//
//	Minimal:    ["adjectives", "core"]                              → "Scalable Core"
//	Startup:    ["adjectives", "core", "suffix"]                    → "Dynamic Workflow Hub"
//	Enterprise: ["adjectives", "buzzwords", "core", "suffix"]       → "Unified Cloud Integration Platform"
//	Bullshit:   ["adjectives", "buzzwords", "buzzwords", ...suffix] → "Synergized AI-Powered Blockchain Data Engine"
func Pattern(mode Mode) []string {
	switch mode {
	case Minimal:
		// Two words: simple and clean
		return []string{"adjectives", "core"}
	case Startup:
		// Three words: the standard startup name formula
		return []string{"adjectives", "core", "suffix"}
	case Enterprise:
		// Four words: adds a buzzword for that corporate feel
		return []string{"adjectives", "buzzwords", "core", "suffix"}
	case Bullshit:
		// Five words: double buzzwords for maximum buzzword density
		return []string{"adjectives", "buzzwords", "buzzwords", "core", "suffix"}
	default:
		// Fallback to minimal for unknown modes
		return []string{"adjectives", "core"}
	}
}
