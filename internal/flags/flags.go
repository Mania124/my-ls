package flags

import (
	"flag"
	"os"
	"strings"
)

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
// Supports both individual flags (-l -R) and combined flags (-lR)
func Parse() *Config {
	config := &Config{}

	// First, preprocess arguments to expand combined flags
	expandedArgs := expandCombinedFlags(os.Args[1:])

	// Temporarily replace os.Args to use our expanded arguments
	originalArgs := os.Args
	os.Args = append([]string{os.Args[0]}, expandedArgs...)

	flag.BoolVar(&config.Recursive, "R", false, "List directories recursively")
	flag.BoolVar(&config.ReverseSort, "r", false, "Reverse the order")
	flag.BoolVar(&config.SortByTime, "t", false, "Sort by modification time")
	flag.BoolVar(&config.ShowAll, "a", false, "Show all files (including hidden)")
	flag.BoolVar(&config.LongFormat, "l", false, "Use long listing format")
	flag.BoolVar(&config.Directories, "d", false, "List directories themselves, not their contents")
	flag.BoolVar(&config.HumanReadable, "h", false, "Display human-readable file sizes")

	flag.Parse()

	// Restore original os.Args
	os.Args = originalArgs

	return config
}

// expandCombinedFlags expands combined flags like -lR into -l -R
func expandCombinedFlags(args []string) []string {
	var expanded []string
	validFlags := map[rune]bool{
		'R': true, 'r': true, 't': true, 'a': true,
		'l': true, 'd': true, 'h': true,
	}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 2 && !strings.HasPrefix(arg, "--") {
			// This is a potential combined flag like -lR
			flagChars := arg[1:] // Remove the leading '-'
			allValid := true

			// Check if all characters are valid flags
			for _, char := range flagChars {
				if !validFlags[char] {
					allValid = false
					break
				}
			}

			if allValid {
				// Expand combined flags
				for _, char := range flagChars {
					expanded = append(expanded, "-"+string(char))
				}
			} else {
				// Not a valid combined flag, keep as is
				expanded = append(expanded, arg)
			}
		} else {
			// Not a flag or already a single flag, keep as is
			expanded = append(expanded, arg)
		}
	}

	return expanded
}
