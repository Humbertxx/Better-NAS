package filesorter

import (
	"fmt"
	"os"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	print("start mock file sorting")
	currentDir, err := os.Getwd()
	check(err)
	for 
}

func EvaluateFilename(filename string, rules []Rule) (target string, matched bool) {

}

func MoveFile(currentDir string, src string, dest string) error {
	
	

	sourcePath := filepath.Join(currentDir, src)
	destDir := filepath.Join(currentDir, dest)
	destPath := filepath.Join(destDir, src)
 
	err = os.MkdirAll(destDir, os.ModePerm) // A directory with 0755 permissions
	check(err)
	err = os.Rename(sourcePath, destPath)
	
	check(err)
	fmt.Println("File moved successfully.")
}
