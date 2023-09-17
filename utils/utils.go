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

func ProcessGitHubContent(ctx context.Context, client *github.Client, config config.ProjectConfig, owner, repo, path string) error {
	opts := &github.RepositoryContentGetOptions{}
	_, contents, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return err
	}

	for _, file := range contents {
		if file.Type != nil {
			if *file.Type == "file" {

				templateFile := *file.Name
				fmt.Println("File:", templateFile)

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

				tmpl, err := template.New(templateFile).Parse(fileContent)
				if err != nil {
					return err
				}

				createFile := strings.Join(append([]string{config.Name}, strings.Split(*file.Path, `/`)...), `/`)
				dirErr := CreateDirectoriesIfNotExist(filepath.Dir(createFile))
				if dirErr != nil {
					return dirErr
				}

				outputFile, err := os.Create(createFile)
				if err != nil {
					return err
				}
				defer outputFile.Close()
				if StringExists("templates", strings.Split(*file.Path, `\`)) {
					if err := tmpl.Execute(outputFile, config); err != nil {
						return err
					}
				} else {
					_, err = outputFile.WriteString(fileContent)
					if err != nil {
						return err
					}
				}

				fmt.Println("File '" + createFile + "' created successfully.")
			} else if *file.Type == "dir" {
				ProcessGitHubContent(ctx, client, config, owner, repo, *file.Path)
			}
		}
	}
	return nil
}
