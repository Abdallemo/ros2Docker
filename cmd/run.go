/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helloCmd represents the hello command
var runCmd = &cobra.Command{

	Use:   "run <workspace-name>",
	Args:  cobra.ExactArgs(1),
	Short: "run a existing ros2 workspace/project",
	Long:  `run a existing ros2 workspace/project`,
	Example: `
  ros2docker run my_ws
  ros2docker run my_ws --clean
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called", args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().Bool("clean", false, "clean start of the worksapce/project. it removes all configuration done to the worksapce/project but not the actual host files")
	runCmd.Flags().Bool("shell", false, "Start the container and drop into a shell (if container is already running, attaches instead).")

}
