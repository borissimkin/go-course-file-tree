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

func renderTree(out io.Writer, path string, printFiles bool, level int, isLast bool) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	dirs, _ := file.ReadDir(0)

	sort.Slice(dirs, func(i int, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	for index, dir := range dirs {
		info, err := dir.Info()

		if err != nil {
			return err
		}

		if !printFiles && !info.IsDir() {
			continue
		}

		if level > 0 {
			for i := 0; i < level; i++ {
				if !isLast {
					fmt.Fprintf(out, "%s%c", separatorChar, tabChar)
				} else {
					if i < level-1 {
						fmt.Fprintf(out, "%s%c", separatorChar, tabChar)
					} else {
						fmt.Fprintf(out, "%c", tabChar)
					}
				}
			}
		}

		var nodeChar string

		if index == len(dirs)-1 {
			nodeChar = lastRoodNodeChar
		} else {
			nodeChar = rootNodeChar
		}

		fmt.Fprintf(out, "%s%s\n", nodeChar, getDirName(dir))

		if info.IsDir() {
			renderTree(out, fmt.Sprintf("%s%c%s", path, os.PathSeparator, dir.Name()), printFiles, level+1, index == len(dirs)-1)
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := renderTree(out, path, printFiles, 0, false)

	if err != nil {
		return err
	}
	// file, err := os.Open(path)

	// if err != nil {
	// 	return fmt.Errorf("path is not corrected")
	// }

	// dirs, _ := file.ReadDir(0)

	// for _, dir := range dirs {
	// 	renderTree(out, path, dir, 0)
	// }

	return nil
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
