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
	var level []bool
	return dirTreeLevel(out, path, printFiles, level)

}

func dirTreeLevel(out io.Writer, path string, printFiles bool, level []bool) error {
	f, _ := os.Open(path)
	files, _ := f.Readdir(0)

	sort.Slice(files, func(l, r int) bool {
		return files[l].Name() < files[r].Name()
	})

	var dirs []os.FileInfo
	if (!printFiles) {
		for _, v := range files {
			if (v.IsDir()) {
				dirs = append(dirs,v)
			}
		}
	}

	level = append(level,false)

	var res *[]os.FileInfo
	if (printFiles) {
		res = &files
	} else {
		res = &dirs
	}
		

	for i, file := range *res {
		name := file.Name()
		var prefix string
		for i := 0; i < len(level) - 1; i++ {
			if level[i] {
				prefix = prefix + "│\t"
			} else {
				prefix = prefix + "\t"
			}
		}
		if i < len(*res)-1 {
			name = prefix + "├───" + name
		} else {
			name = prefix + "└───" + name
		}
		if (i == len(*res)-1) {
			level[len(level)-1] = false
		} else {
			level[len(level)-1] = true
		}
		if file.IsDir() {
			io.WriteString(out,name+"\n")
			dirTreeLevel(out, path+"/"+file.Name(), printFiles, level)
		} else if printFiles {
			size := file.Size()
			if (size != 0) {
				name = name + " (" + fmt.Sprintf("%d", size) + "b)"
			} else {
				name = name + " (empty)"
			}
			io.WriteString(out,name+"\n")
		}

	}
	return nil
}
