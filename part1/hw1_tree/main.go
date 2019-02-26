package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

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

// ---

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeLevel(out, path, printFiles, nil)
}

func dirTreeLevel(out io.Writer, path string, printFiles bool, level []bool) error {
	f, _ := os.Open(path)
	defer f.Close()
	files, _ := f.Readdir(0)

	sort.Slice(files, func(l, r int) bool {
		return files[l].Name() < files[r].Name()
	})

	var toOut *[]os.FileInfo
	if printFiles {
		toOut = &files
	} else {
		toOut = new([]os.FileInfo)
		*toOut = copyIf(files, func(v os.FileInfo) bool {
			return v.IsDir()
		})
	}

	var prefix string
	for _, v := range level {
		if v {
			prefix = prefix + "│\t"
		} else {
			prefix = prefix + "\t"
		}
	}

	level = append(level, false)

	for i, file := range *toOut {
		name := file.Name()

		if i < len(*toOut)-1 {
			name = prefix + "├───" + name
			level[len(level)-1] = true
		} else {
			name = prefix + "└───" + name
			level[len(level)-1] = false
		}

		if file.IsDir() {
			io.WriteString(out, name+"\n")
			dirTreeLevel(out, path+"/"+file.Name(), printFiles, level)
		} else if printFiles {
			size := file.Size()
			if size != 0 {
				name = name + " (" + fmt.Sprintf("%d", size) + "b)"
			} else {
				name = name + " (empty)"
			}
			io.WriteString(out, name+"\n")
		}

	}
	return nil
}

func copyIf(in []os.FileInfo, predicate func(v os.FileInfo) bool) (out []os.FileInfo) {
	for _, v := range in {
		if predicate(v) {
			out = append(out, v)
		}
	}
	return
}

