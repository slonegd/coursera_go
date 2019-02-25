package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	// "io/ioutil"
	// "path/filepath"
	// "strings"
)

func dirTree (out io.Writer, path string, printFiles bool) error {
	f, _ := os.Open (path)
	files, _ := f.Readdir(0)

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			fmt.Println(name)
			dirTree(out, path + "/" + name, printFiles)
		} else {
			name = name + " (" + fmt.Sprintf("%d", file.Size()) + ")" 
			fmt.Println(name)
		}
		
	}
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
