package flags

import "flag"

// Config holds all command-line flag values
type Config struct {
	Recursive     bool
	ReverseSort   bool
	SortByTime    bool
	ShowAll       bool
	LongFormat    bool
	Directories   bool
	HumanReadable bool
}

// Parse parses command-line flags and returns a Config struct
func Parse() *Config {
	config := &Config{}
	
	flag.BoolVar(&config.Recursive, "R", false, "List directories recursively")
	flag.BoolVar(&config.ReverseSort, "r", false, "Reverse the order")
	flag.BoolVar(&config.SortByTime, "t", false, "Sort by modification time")
	flag.BoolVar(&config.ShowAll, "a", false, "Show all files (including hidden)")
	flag.BoolVar(&config.LongFormat, "l", false, "Use long listing format")
	flag.BoolVar(&config.Directories, "d", false, "List directories themselves, not their contents")
	flag.BoolVar(&config.HumanReadable, "h", false, "Display human-readable file sizes")
	
	flag.Parse()
	
	return config
}
