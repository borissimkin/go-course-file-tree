package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

var (
	rootNodeChar     = "├───"
	lastRoodNodeChar = "└───"
	separatorChar    = "│"
	tabChar          = '\t'
)

func getDirName(dir os.DirEntry) string {
	if dir.IsDir() {
		return dir.Name()
	}

	info, err := dir.Info()

	if err != nil {
		return dir.Name()
	}

	size := info.Size()

	var sizeText string

	if size > 0 {
		sizeText = fmt.Sprintf("%vb", size)
	} else {
		sizeText = "empty"
	}

	return fmt.Sprintf("%s (%s)", dir.Name(), sizeText)
}

func filterFiles(dirs []os.DirEntry) []os.DirEntry {
	newDirs := make([]os.DirEntry, 0, len(dirs))

	for _, dir := range dirs {
		if dir.IsDir() {
			newDirs = append(newDirs, dir)
		}
	}

	return newDirs
}

func getNodeChar(isLast bool) string {
	if isLast {
		return lastRoodNodeChar
	}

	return rootNodeChar
}

func renderTree(out io.Writer, path string, printFiles bool, level int, isLast bool, prevPrefix string) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	dirs, _ := file.ReadDir(0)

	if !printFiles {
		dirs = filterFiles(dirs)
	}

	sort.Slice(dirs, func(i int, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	for index, dir := range dirs {
		info, err := dir.Info()

		if err != nil {
			return err
		}

		var prefix string = prevPrefix

		if level > 0 {
			if !isLast {
				prefix = fmt.Sprintf("%s%s%c", prefix, separatorChar, tabChar)
			} else {
				prefix = fmt.Sprintf("%s%c", prefix, tabChar)
			}
		}

		currentIsLast := index == len(dirs)-1

		nodeChar := getNodeChar(currentIsLast)

		fmt.Fprint(out, prefix)
		fmt.Fprintf(out, "%s%s\n", nodeChar, getDirName(dir))

		if info.IsDir() {
			renderTree(out, fmt.Sprintf("%s%c%s", path, os.PathSeparator, dir.Name()), printFiles, level+1, currentIsLast, prefix)
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return renderTree(out, path, printFiles, 0, false, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
