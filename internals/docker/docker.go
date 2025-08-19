package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Abdallemo/ros2Docker/internals/ui/text"
	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// These constants represent the supported ROS distribution image types.
const (
	Humble = "humble"
	Foxy   = "foxy"
)

var execConfig = container.ExecOptions{
	Cmd:          []string{"bash", "-c", "apt-get update && apt-get install -y git python3-pip"},
	AttachStdout: true,
	AttachStderr: true,
}

type Logger interface {
	StreamLogs(io.ReadCloser) error
}

type ImageType string

type DockerContainer struct {
	ID string `json:"id"`
}

// IsValid checks if the ImageType is one of the supported types.
func (i ImageType) IsValid() bool {
	return i == Humble || i == Foxy

}

type WorkspaceConfig struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Image  ImageType `json:"image"`
	Volume string    `json:"volume"`
}

func (wcfg *WorkspaceConfig) GetSaveData() (any, string) {
	return wcfg, fmt.Sprintf("workspace-%s.json", wcfg.Name)
}

type Docker struct {
	Client    *client.Client
	Container DockerContainer
	Image     ImageType
}
type ProgressDetails struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

// ImageBuildMessage represents a single message from the Docker image build stream.
type ImageBuildMessage struct {
	ID             string          `json:"id"`
	Status         string          `json:"status"`
	ProgressDetail ProgressDetails `json:"progressDetail"`
	Progress       string          `json:"progress"`
}

func StreamLogs(reader io.ReadCloser) error {
	decoder := json.NewDecoder(reader)

	for {
		var msg ImageBuildMessage
		if err := decoder.Decode(&msg); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error decoding image pull logs: %w", err)
		}
		if msg.Progress != "" {
			fmt.Printf("\rLayer %s: %s", msg.ID, msg.Progress)
		} else if msg.Status != "" {
			fmt.Printf("\r%s", msg.Status)
		}
	}
	fmt.Println("\n✅ Image pull complete")
	return nil
}

// CreateContainer pulls a base image, sets up a temporary container to
// install dependencies, commits the changes to a new image,
// and then creates the final container with a mounted volume.
func (d *Docker) CreateContainer(containerName, volume string) error {
	toolsImage := fmt.Sprintf("ros:%s-tools", d.Image)
	ctx := context.Background()

	if !d.Image.IsValid() {
		return fmt.Errorf("invalid image type: %s", d.Image)
	}

	_, inspectErr := d.Client.ImageInspect(ctx, toolsImage)
	if inspectErr == nil {
		log := text.New()
		log.Run()
		log.Append("using Cached Version", text.Info)
		finalCont, err := d.setupContainer(ctx, toolsImage, volume, containerName)
		log.Stop()
		if err != nil {
			return err
		}
		d.Container.ID = finalCont.ID
		return nil
	}
	if !errdefs.IsNotFound(inspectErr) {
		return inspectErr
	}
	reader, err := d.Client.ImagePull(ctx, fmt.Sprintf("osrf/ros:%s-desktop", d.Image), image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	StreamLogsTui(reader)

	commitRespId, err := d.installDependencies(ctx, volume, containerName, toolsImage)
	if err != nil {
		return err
	}

	finalCont, err := d.setupContainer(ctx, commitRespId, volume, containerName)
	if err != nil {
		return err
	}

	d.Container.ID = finalCont.ID

	return nil
}

// small helper to setup a container with passed meta data
func (d *Docker) setupContainer(ctx context.Context, imageId, volume, containerName string) (container.CreateResponse, error) {

	return d.Client.ContainerCreate(ctx,
		&container.Config{
			Image: imageId,
			Tty:   true,
			Cmd:   []string{"bash"},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: volume,
					Target: "/home",
				},
			},
		},
		nil,
		nil,
		containerName,
	)

}

// RemoveContainer stops and removes the Docker container.
func (d *Docker) RemoveContainer(ctx context.Context, containerId string) error {
	if err := d.Client.ContainerRemove(ctx, containerId, container.RemoveOptions{Force: true}); err != nil {
		return err
	}
	return nil

}

// NewDocker creates a new Docker instance with the provided client
func NewDocker(client *client.Client) *Docker {
	return &Docker{
		Client: client,
	}
}

// installDependencies installs exec packages into a temepory container and then commits these
// changes to a new Custome Image.
func (d *Docker) installDependencies(ctx context.Context, volume, containerName, toolsImage string) (string, error) {

	tempCont, err := d.setupContainer(ctx, fmt.Sprintf("osrf/ros:%s-desktop", d.Image), volume, containerName+"-setup")
	if err != nil {
		return "", err
	}
	if err := d.Client.ContainerStart(ctx, tempCont.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	execID, err := d.Client.ContainerExecCreate(ctx, tempCont.ID, execConfig)
	if err != nil {
		return "", err
	}
	resp, err := d.Client.ContainerExecAttach(ctx, execID.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", err
	}
	defer resp.Close()
	io.Copy(os.Stderr, resp.Reader)

	fmt.Println("cleaning up things...")
	commitResp, err := d.Client.ContainerCommit(ctx, tempCont.ID, container.CommitOptions{
		Reference: toolsImage,
	})
	if err != nil {
		return "", err
	}
	if err := d.RemoveContainer(ctx, tempCont.ID); err != nil {
		return "", err
	}
	_, err = d.Client.ImageRemove(ctx, fmt.Sprintf("osrf/ros:%s-desktop", d.Image), image.RemoveOptions{})
	if err != nil {
		return "", err
	}

	fmt.Println("Dependencies installed successfully!")

	return commitResp.ID, nil
}
