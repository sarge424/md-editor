package files

import (
	"io"
	"log"
	"os"
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

	return string(b)
}

func SaveFile(path string, text string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = f.WriteString(text)
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
