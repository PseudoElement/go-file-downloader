package utils

import (
	"fmt"
	"log"
	"os"
)

func CreateTempFile(fileName string, extension string, content string) (*os.File, error) {
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

	// time.Sleep(time.Second)

	return file, nil
}

func CreateFile(fileName string, extension string, content string) (*os.File, error) {
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
