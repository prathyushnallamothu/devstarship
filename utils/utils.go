package utils

import (
	"context"
	"devstarship/config"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v33/github"
)

// CreateDirectoriesIfNotExist creates the directory path if it does not exist.
func CreateDirectoriesIfNotExist(dirPath string) error {
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}
	fmt.Println("Directory ensured:", dirPath)
	return nil
}

// StringExists checks if a target string exists in an array.
func StringExists(target string, array []string) bool {
	for _, element := range array {
		if element == target {
			return true
		}
	}
	return false
}

// CopyFile copies contents from the source file to the destination file.
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
	return err
}

// ProcessGitHubContent processes GitHub repository content recursively.
func ProcessGitHubContent(ctx context.Context, client *github.Client, config config.ProjectConfig, owner, repo, path string) error {
	opts := &github.RepositoryContentGetOptions{}
	_, contents, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return err
	}

	for _, file := range contents {
		if file.Type == nil {
			continue
		}

		switch *file.Type {
		case "file":
			if err := processFile(ctx, client, config, owner, repo, file); err != nil {
				return err
			}
		case "dir":
			if err := ProcessGitHubContent(ctx, client, config, owner, repo, *file.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

// processFile downloads, parses, and writes the file content.
func processFile(ctx context.Context, client *github.Client, config config.ProjectConfig, owner, repo string, file *github.RepositoryContent) error {
	opts := &github.RepositoryContentGetOptions{}
	fileContentReader, _, err := client.Repositories.DownloadContents(ctx, owner, repo, *file.Path, opts)
	if err != nil {
		return err
	}
	defer fileContentReader.Close()

	fileContentBytes, err := io.ReadAll(fileContentReader)
	if err != nil {
		return err
	}
	fileContent := string(fileContentBytes)

	tmpl, err := template.New(*file.Name).Parse(fileContent)
	if err != nil {
		return err
	}

	createFile := filepath.Join(config.Name, *file.Path)
	if err := CreateDirectoriesIfNotExist(filepath.Dir(createFile)); err != nil {
		return err
	}

	outputFile, err := os.Create(createFile)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if StringExists("templates", strings.Split(*file.Path, `/`)) {
		_, err = outputFile.WriteString(fileContent)
	} else {
		err = tmpl.Execute(outputFile, config)
	}
	if err != nil {
		return err
	}

	fmt.Println("File created successfully:", createFile)
	return nil
}
