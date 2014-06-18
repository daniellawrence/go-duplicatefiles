package main
 
import (
   "path/filepath"
   "fmt"
   "os"
   "strings"
   "crypto/sha256"
   "io/ioutil"
   "encoding/hex"
)
 
func main(){
     location := "."
     GoWalk(location)
     chann := GoWalk(location)
     for msg := range chann {
     	 fmt.Println("# reviewed", msg)
     }
     return
}

func checkDuplicates(files []string) (result bool) {
     stored_hashs  := make(map[string][]string)
     for _, filename := range files {
	 hasher := sha256.New()
	 file_content, _ := ioutil.ReadFile(filename)
	 hasher.Write(file_content)
	 file_hash := hex.EncodeToString(hasher.Sum(nil))
	 if _, ok := stored_hashs[file_hash]; ok {
	    fmt.Printf("%s is a Duplicate of %s\n", filename, stored_hashs[file_hash])
	 }
	 stored_hashs[file_hash] = append(stored_hashs[file_hash], filename)
     }
     return true

}
 
func GoWalk(location string) (chan string) {
     chann := make(chan string)
     dict  := make(map[int64][]string)

     go func(){
     	filepath.Walk(location, func(path string, fileinfo os.FileInfo, _ error)(err error){
	    // skip over hidden files
	    if strings.HasPrefix(path, ".") {
	       return
	    }
	    // Skip over tmp files
	    if strings.HasSuffix(path, "~") {
	       return
	    }
	    file_size := fileinfo.Size()
	    dict[file_size] = append(dict[file_size], path)
	    chann <- path
	    return
	})
	for _, v := range dict {
	    if len(v) == 1 {
	       continue
	    }
	    // fmt.Println("possible duplicates:", v)
	    checkDuplicates(v)
	}
	defer close(chann)
     }()
	return chann
}