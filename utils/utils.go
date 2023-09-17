package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateDirectoriesIfNotExist(dirPath string) error {
	dirNames := filepath.SplitList(dirPath)
	for _, dirName := range dirNames {
		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
				return err
			}
			fmt.Println("Directory created:", dirName)
		} else if err != nil {
			return err
		} else {
			fmt.Println("Directory already exists:", dirName)
		}
	}

	return nil
}

func StringExists(target string, array []string) bool {
	for _, element := range array {
		if element == target {
			return true
		}
	}
	return false
}

func CopyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
