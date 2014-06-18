package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	args := os.Args[1:]
	location := "."
	if len(args) > 0 {
		location = args[0]
	}
	GoWalk(location)
	return
}

func checkDuplicates(files []string) (result bool) {
	stored_hashs := make(map[string][]string)
	wasted_space := 0

	for _, filename := range files {
		hasher := sha256.New()
		file_content, _ := ioutil.ReadFile(filename)
		hasher.Write(file_content)
		file_hash := hex.EncodeToString(hasher.Sum(nil))

		// If the file_hash is already in the stored_hashs
		// then we know we have a duplicate
		if _, ok := stored_hashs[file_hash]; ok {
			fmt.Printf("%s is a Duplicate of %s\n", filename, stored_hashs[file_hash])
			wasted_space = len(stored_hashs[file_hash]) * len(file_content)
		}
		stored_hashs[file_hash] = append(stored_hashs[file_hash], filename)
	}
	if wasted_space > 0 {
		fmt.Println(files[0], "Wasted space:", wasted_space)
	}
	return true

}

func GoWalk(location string) {
	dict := make(map[int64][]string)

	filepath.Walk(location, func(path string, fileinfo os.FileInfo, _ error) (err error) {
		// skip over hidden files
		if strings.Contains(path, "/.") {
			return
		}
		// Skip over tmp files
		if strings.HasSuffix(path, "~") {
			return
		}
		// Duplicate dirs... you mean emtpy dirs
		if fileinfo.IsDir() {
			return
		}
		file_size := fileinfo.Size()
		dict[file_size] = append(dict[file_size], path)
		return
	})

	for _, v := range dict {
		// only check if the file is a duplicate if they are the same size
		if len(v) == 1 {
			continue
		}
		checkDuplicates(v)
	}
}
