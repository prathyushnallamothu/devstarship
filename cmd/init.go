package cmd

import (
	"devstarship/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type ProjectConfig struct {
	Name           string
	Description    string
	Version        string
	DockerUsername string
}

var projectName string
var projectDescription string
var dockerUsername string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Run: func(cmd *cobra.Command, args []string) {
		config := ProjectConfig{
			Name:           projectName,
			Description:    projectDescription,
			Version:        "0.1.0", // You can set a default version or prompt for it.
			DockerUsername: dockerUsername,
		}
		fmt.Println(config)
		workingdir, _ := os.Getwd()
		templatesPath := "./template"

		err := filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}
			if !info.IsDir() {
				fmt.Println("File:", path)
				templateFile := strings.ReplaceAll(path, `template\`, "")
				tmpl, err := template.ParseFiles(path)
				if err != nil {
					panic(err)
				}
				createFile := strings.Join(append([]string{workingdir, config.Name}, strings.Split(templateFile, `\`)...), `\`)
				dirErr := utils.CreateDirectoriesIfNotExist(filepath.Dir(createFile))
				if dirErr != nil {
					panic(dirErr)
				}

				outputFile, err := os.Create(createFile)
				if err != nil {
					panic(err)
				}
				defer outputFile.Close()
				if !utils.StringExists("templates", strings.Split(templateFile, `\`)) {
					if err := tmpl.Execute(outputFile, config); err != nil {
						panic(err)
					}
				} else {
					err := utils.CopyFile(path, createFile)
					if err != nil {
						panic(err)
					}
				}

				println("File '" + createFile + "'created successfully.")
			}
			return nil
		})

		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	initCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name (required)")
	initCmd.Flags().StringVarP(&projectDescription, "description", "d", "", "Project description")
	initCmd.Flags().StringVarP(&dockerUsername, "docker-username", "u", "", "Docker username (required)")
	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("docker-username")
	rootCmd.AddCommand(initCmd)
}
