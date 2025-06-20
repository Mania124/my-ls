package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"

	"ls/internal/display"
	"ls/internal/flags"
	"ls/internal/utils"
)

func main() {
	// Parse command-line flags
	config := flags.Parse()

	// Get directory from args or use current dir
	var dirs []string
	if flag.NArg() > 0 {
		dirs = flag.Args()
	} else {
		dirs = []string{"."}
	}
	sort.Slice(dirs, func(i, j int) bool {
		if config.ReverseSort {
			return dirs[i][0] > dirs[j][0]
		}
		return dirs[i][0] < dirs[j][0]
	})
	for i, dir := range dirs {
		// Expand `~` to the home directory
		if dir == "~" || (len(dir) >= 2 && dir[:2] == "~/") {
			usr, _ := user.Current()
			dir = filepath.Join(usr.HomeDir, dir[1:])
		}
		info, e := os.Lstat(dir)
		if e != nil {
			fmt.Printf("my-ls: cannot access '%s': No such file or directory\n", dir)
			continue
		}
		if !info.IsDir() {
			dir = filepath.Base(dir)
			if config.LongFormat {
				// f,_:=os.ReadDir(dir)
				display.PrintLongFormat(info, dir, config)
				fmt.Println()
			} else {
				fmt.Print(dir)
				if i != len(dirs)-1 {
					fmt.Print(" ")
				} else {
					fmt.Println()
				}
			}
			continue
		}
		// Read files
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Fatalf("Error reading directory: %v", err)
		}
		// Compute and print total block count
		if config.LongFormat {
			total := utils.ComputeTotalBlocks(files, dir)
			fmt.Printf("total %d\n", total)
		}
		// Filter, sort, and display files
		utils.FilterFiles(&files, config)
		utils.SortFiles(files, dir, config)
		// Handle recursive (-R)
		if config.Recursive {
			utils.HandleRecursive(files, dir, config, display.Files)
		} else {
			display.Files(files, dir, config)
		}
	}
}
