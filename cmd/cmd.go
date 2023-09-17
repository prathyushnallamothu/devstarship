package cmd

import (
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{Use: "devstarship"}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        println(err)
        os.Exit(1)
    }
}
