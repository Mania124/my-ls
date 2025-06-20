package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"

	"ls/internal/flags"
)

// FilterFiles filters files based on flags
func FilterFiles(files *[]os.DirEntry, config *flags.Config) {
	filtered := []os.DirEntry{}
	for _, file := range *files {
		// `-a` flag: show hidden files
		if !config.ShowAll && file.Name()[0] == '.' {
			continue
		}
		// `-d` flag: show directories only
		if config.Directories && !file.IsDir() {
			continue
		}
		filtered = append(filtered, file)
	}
	*files = filtered
}

// SortFiles sorts files based on flags
func SortFiles(files []os.DirEntry, dir string, config *flags.Config) {
	sort.Slice(files, func(i, j int) bool {
		infoI, err1 := os.Stat(filepath.Join(dir, files[i].Name()))
		infoJ, err2 := os.Stat(filepath.Join(dir, files[j].Name()))

		if err1 != nil || err2 != nil {
			return false
		}

		// Sort by modification time (-t)
		if config.SortByTime {
			if config.ReverseSort {
				return infoI.ModTime().Before(infoJ.ModTime())
			}
			return infoI.ModTime().After(infoJ.ModTime())
		}

		// Sort hidden files ('.' by second character)
		nameI, nameJ := files[i].Name(), files[j].Name()
		if nameI[0] == '.' && len(nameI) > 1 {
			nameI = nameI[1:]
		}
		if nameJ[0] == '.' && len(nameJ) > 1 {
			nameJ = nameJ[1:]
		}

		// Default alphabetical sorting
		if config.ReverseSort {
			return nameI > nameJ
		}
		return nameI < nameJ
	})
}

// ComputeTotalBlocks computes total blocks for `ls -l` output
func ComputeTotalBlocks(files []os.DirEntry, dir string) int {
	total := 0
	for _, file := range files {
		fullPath := filepath.Join(dir, file.Name())
		info, err := os.Lstat(fullPath)
		if err == nil {
			stat := info.Sys().(*syscall.Stat_t)
			total += int(stat.Blocks)
		}
	}
	return total / 2 // Adjust block size like `ls -l`
}

// HandleRecursive handles recursive (-R) listing
func HandleRecursive(files []os.DirEntry, dir string, config *flags.Config, displayFunc func([]os.DirEntry, string, *flags.Config)) {
	// ✅ If the base directory is ".", print ".:"
	if dir == "." {
		fmt.Println(".:")
		displayFunc(files, dir, config) // List files in the primary directory
	}

	// Loop through and process subdirectories
	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())

		// Only process directories recursively
		if file.IsDir() {
			fmt.Printf("\n%s:\n", filePath) // Print directory name before listing
			subFiles, err := os.ReadDir(filePath)
			if err == nil {
				FilterFiles(&subFiles, config)
				SortFiles(subFiles, filePath, config)
				displayFunc(subFiles, filePath, config)
				HandleRecursive(subFiles, filePath, config, displayFunc) // Recursive call
			}
		}
	}
}
