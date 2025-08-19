package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Abdallemo/ros2Docker/internals/docker"
	"github.com/Abdallemo/ros2Docker/internals/ui/text"
	"github.com/Abdallemo/ros2Docker/internals/utils"
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
		log := text.New()
		log.Run()
		if len(args) < 1 {
			log.Append("please provide a workspace name", text.Error)
			return
		}
		workspaceName := args[0]
		image, _ := cmd.Flags().GetString("image")
		path, _ := cmd.Flags().GetString("path")
		currentDir, err := os.Getwd()
		if err != nil {
			log.Append("unable to locate the current working direcorty", text.Error)
			return
		}
		dockerClient.Image = docker.ImageType(image)

		if ok := dockerClient.Image.IsValid(); !ok {
			log.Append(fmt.Sprintf("unsupported image type: %s", image), text.Error)
			return
		}
		if err := utils.IsValidPath(path); err != nil {
			log.Append(fmt.Sprintf("%s", err.Error()), text.Error)
			return
		}
		log.Append(fmt.Sprintf("Creating new ros2 workspace: %s with image %s\n", workspaceName, image), text.Ok)
		containerName := utils.GenerateMix(workspaceName)

		var wscfg docker.WorkspaceConfig
		if path != "" {
			if err := dockerClient.CreateContainer(containerName, path); err != nil {
				fmt.Println(err)
				return
			}
			wscfg = docker.WorkspaceConfig{
				ID:     dockerClient.Container.ID,
				Name:   workspaceName,
				Image:  dockerClient.Image,
				Volume: path,
			}
		} else {
			if err := dockerClient.CreateContainer(containerName, currentDir); err != nil {
				fmt.Println(err)
				return
			}
			wscfg = docker.WorkspaceConfig{
				ID:     dockerClient.Container.ID,
				Name:   workspaceName,
				Image:  dockerClient.Image,
				Volume: currentDir,
			}
		}
		workspaceDir := filepath.Join(globalCfg.Location, workspaceName)
		if err := os.MkdirAll(workspaceDir, 0755); err != nil {
			fmt.Printf("failed to create workspace config folder: %v\n", err)
			return
		}

		if err := utils.SaveConfig(workspaceDir, &wscfg); err != nil {
			fmt.Printf("failed to save workspace config: %v\n", err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().String("image", docker.Humble, "provide the ros2 version you need for this worksapce/project")
	newCmd.Flags().String("path", "", "provide the absolute worksapce/project path. if not specifed it default to the current reltive folder")
}
