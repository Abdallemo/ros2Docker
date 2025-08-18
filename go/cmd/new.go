package cmd

import (
	"fmt"

	"github.com/Abdallemo/ros2Docker/internals/docker"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "creates a new worksapce for the ros2 project",
	Long: `[new] command creates a new worksapce for the ros2 project
	and example:
	ros2docker new <workspaceName>
	`,
	Example: `
	ros2docker run my_ws --image humble`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new ros2 Workspace")

	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().String("image", docker.Humble, "provide the ros2 version you need for this worksapce/project")
}
