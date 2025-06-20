package display

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"ls/internal/flags"
	"ls/internal/permissions"
)

// Files displays files with color & format
func Files(files []os.DirEntry, dir string, config *flags.Config) {
	for i, file := range files {
		fullPath := filepath.Join(dir, file.Name()) // Use full path
		info, err := os.Lstat(fullPath)             // ✅ Use Lstat to detect symlinks
		if err != nil {
			continue
		}

		if config.LongFormat {
			PrintLongFormat(info, fullPath, config) // ✅ Pass file info and full path
			if i != len(files)-1 {
				fmt.Println()
			}
		} else {
			printShortFormat(file, info, fullPath)
		}
	}
	fmt.Println()
}

// printShortFormat prints short format (default)
func printShortFormat(file os.DirEntry, info os.FileInfo, fullPath string) {
	color := getFileColor(info)
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(fullPath)
		if err == nil {
			targetInfo, err := os.Stat(fullPath)
			if err == nil {
				color = getFileColor(targetInfo)
			}
			fmt.Printf("\033[1;36m%s\033[0m -> %s%s\033[0m  ", file.Name(), color, filepath.Base(target)) // Cyan for symlinks, proper color for target
			return
		}
	}
	fmt.Printf("%s%s\033[0m  ", color, file.Name())
}

// PrintLongFormat prints long format (-l flag)
func PrintLongFormat(info os.FileInfo, fullPath string, config *flags.Config) {
	stat := info.Sys().(*syscall.Stat_t)
	nLinks := stat.Nlink
	owner, err := user.LookupId(strconv.Itoa(int(stat.Uid)))
	if err != nil {
		owner = &user.User{Username: strconv.Itoa(int(stat.Uid))} // Fallback to UID if lookup fails
	}
	group, err := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))
	if err != nil {
		group = &user.Group{Name: strconv.Itoa(int(stat.Gid))} // Fallback to GID if lookup fails
	}

	size := info.Size()
	if config.HumanReadable {
		sizeS := permissions.HumanReadableSize(size)
		sizes, _ := strconv.Atoi(sizeS)
		size = int64(sizes)
	}

	linkTarget := ""
	fileColor := getFileColor(info)
	permColor := "\033[0;37m" // Default color for permissions, user, group, and date
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(fullPath)
		if err == nil {
			targetInfo, err := os.Stat(fullPath)
			if err == nil {
				fileColor = getFileColor(targetInfo)
			}
			linkTarget = " -> " + fileColor + filepath.Base(target) + "\033[0m"
		}
	}

	fmt.Printf("%s%s %2d %s%s %s%s %6d %s%s %s%s%s",
		permColor, permissions.GetPermissions(info.Mode()),
		nLinks,
		permColor, owner.Username,
		permColor, group.Name,
		size,
		permColor, info.ModTime().Format("Jan 02 15:04"),
		fileColor, filepath.Base(fullPath),
		linkTarget,
	)
}

// getFileColor determines file color based on type and permissions
func getFileColor(info os.FileInfo) string {
	if info.IsDir() {
		return "\033[1;34m" // Blue for directories
	} else if info.Mode()&0o111 != 0 {
		return "\033[1;32m" // Green for executables
	} else if info.Mode()&os.ModeSymlink != 0 {
		return "\033[1;36m" // Cyan for symlinks
	}
	return ""
}
