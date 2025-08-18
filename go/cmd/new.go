package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new ros2 Workspace")
	},
}

func init() {
	st := newCmd.Flags().BoolP("start", "s", false, "start a new workspace")
	fmt.Println("value chosen ", &st)
}
