package files

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// get the filenames in a directory
func GetFilenames(path string) []string {
	var files []string

	// return an empty slice if an error is encountered.
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	//otherwise, get all the filenames
	for _, file := range entries {
		files = append(files, file.Name())
	}

	return files
}

func LoadFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := io.ReadAll(file)

	return FText2E(string(b))
}

func SaveFile(path string, text string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = f.WriteString(EText2F(text))
	if err != nil {
		log.Fatal(err)
		f.Close()
		return
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func RenameFile(oldName string, newName string, path string) bool {
	err := os.Rename(path+oldName, path+newName)
	fmt.Println("Renamed", oldName, "to", newName, "status=", err == nil)
	return err == nil
}

func EText2F(text string) string {
	// change 4x space to tabs
	return strings.ReplaceAll(text, "    ", "\t")
}

func FText2E(text string) string {
	// change 4x space to tabs
	return strings.ReplaceAll(text, "\t", "    ")
}
