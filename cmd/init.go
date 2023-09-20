package cmd

import (
	"devstarship/config"
	"devstarship/utils"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v33/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var projectName string
var projectDescription string
var dockerUsername string
var token string
// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.ProjectConfig{
			Name:           projectName,
			Description:    projectDescription,
			Version:        "0.1.0", // You can set a default version or prompt for it.
			DockerUsername: dockerUsername,
		}
		cjson,_:=json.Marshal(config)
		fmt.Println(cjson)
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(cmd.Context(), ts)
		client := github.NewClient(tc)
		owner := "prathyushnallamothu"
		repo := "go_template"
		templateDirectory := "./"
		err := utils.ProcessGitHubContent(cmd.Context(), client, config, owner, repo, templateDirectory)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	initCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name (required)")
	initCmd.Flags().StringVarP(&projectDescription, "description", "d", "", "Project description")
	initCmd.Flags().StringVarP(&dockerUsername, "docker-username", "u", "", "Docker username (required)")
	initCmd.Flags().StringVarP(&token, "github-pat", "g", "", "Github Personal Acess Token (required)")
	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("docker-username")
	initCmd.MarkFlagRequired("github-pat")
	rootCmd.AddCommand(initCmd)
}
