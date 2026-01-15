package cli

import "flag"

type Config struct {
	Lang    string
	Mode    string
	Seed    string
	Count   int
	Explain bool
}

func ParseFlags() Config {
	var cfg Config

	// Language flag: determines which language-specific word files to load
	flag.StringVar(&cfg.Lang, "lang", "en", "language (en, de)")

	// Mode flag: controls the complexity and style of generated names
	flag.StringVar(&cfg.Mode, "mode", "startup", "mode (startup, enterprise, bullshit, minimal)")

	// Seed flag: when provided, ensures deterministic name generation
	flag.StringVar(&cfg.Seed, "seed", "", "deterministic seed")

	// Count flag: allows batch generation of multiple names
	flag.IntVar(&cfg.Count, "count", 1, "number of names")

	// Explain flag: enables verbose output showing how each name was generated
	flag.BoolVar(&cfg.Explain, "explain", false, "explain how the name was generated")

	flag.Parse()
	return cfg
}
