package custom_utils

import (
	"fmt"
	"log"
	"os"
)

func CreateTempFile(fileName string, extension string, content string) (*os.File, error) {
	createTempDirectoryIfNotExists()
	pathDir, _ := os.Getwd()
	tempDir := fmt.Sprintf("%s/src/temp", pathDir)
	fullName := fmt.Sprintf("%s-*.%s", fileName, extension)
	file, err := os.CreateTemp(tempDir, fullName)

	if err != nil {
		return nil, err
	}
	defer file.Close()
	defer os.Remove(file.Name())

	_, err = file.WriteString(content)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func CreateFile(fileName string, extension string, content string) (*os.File, error) {
	createTempDirectoryIfNotExists()
	pathDir, _ := os.Getwd()
	fullName := fmt.Sprintf("%s/src/temp/%s.%s", pathDir, fileName, extension)
	file, e := os.Create(fullName)
	if e != nil {
		panic(e)
	}
	defer file.Close()

	n, err := file.WriteString(content)
	if err != nil {
		return nil, err
	}

	log.Println("bytes written ===> ", n)
	return file, nil
}

func CreateNewFileWithManyWords(wordsCount int) (*os.File, error) {
	var content string
	for i := range wordsCount {
		content += CreateRandomWord(15, false)
		if i < wordsCount-1 {
			content += "\n"
		} else {
			content += "."
		}
	}
	f, err := CreateTempFile("temp-pdf-file", "pdf", content)

	return f, err
}

func createTempDirectoryIfNotExists() error {
	pathDir, _ := os.Getwd()
	tempDir := fmt.Sprintf("%s/src/temp", pathDir)
	_, err := os.Stat(tempDir)

	if os.IsNotExist(err) {
		if err = os.Mkdir(tempDir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
