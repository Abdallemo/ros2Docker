package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

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
	Name   string    `json:"name"`
	Image  ImageType `json:"image"`
	Volume string    `json:"volume"`
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

// StreamLogs decodes and prints the logs from an image build or pull stream.
func (msg *ImageBuildMessage) StreamLogs(reader io.ReadCloser) error {
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
func (d *Docker) CreateContainer(containerName, volume string, img ImageType) error {
	m := ImageBuildMessage{}
	ctx := context.Background()

	if !img.IsValid() {
		return fmt.Errorf("invalid image type: %s", img)
	}

	reader, err := d.Client.ImagePull(ctx, fmt.Sprintf("osrf/ros:%s-desktop", img), image.PullOptions{})
	if err != nil {
		return err
	}
	m.StreamLogs(reader)

	tempCont, err := d.Client.ContainerCreate(ctx,
		&container.Config{
			Image: fmt.Sprintf("osrf/ros:%s-desktop", img),
			Tty:   true,
			Cmd:   []string{"bash"},
		},
		nil, nil, nil, containerName+"-setup")
	if err != nil {
		return err
	}

	if err := d.Client.ContainerStart(ctx, tempCont.ID, container.StartOptions{}); err != nil {
		return err
	}

	execConfig := container.ExecOptions{
		Cmd:          []string{"bash", "-c", "apt-get update && apt-get install -y git python3-pip"},
		AttachStdout: true,
		AttachStderr: true,
	}
	execID, err := d.Client.ContainerExecCreate(ctx, tempCont.ID, execConfig)
	if err != nil {
		return err
	}
	resp, err := d.Client.ContainerExecAttach(ctx, execID.ID, container.ExecAttachOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()
	io.Copy(os.Stdout, resp.Reader)

	commitResp, err := d.Client.ContainerCommit(ctx, tempCont.ID, container.CommitOptions{
		Reference: fmt.Sprintf("ros:%s-tools", img),
	})
	if err != nil {
		return err
	}

	if err := d.Client.ContainerRemove(ctx, tempCont.ID, container.RemoveOptions{Force: true}); err != nil {
		return err
	}

	finalCont, err := d.Client.ContainerCreate(ctx,
		&container.Config{
			Image: commitResp.ID,
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
	if err != nil {
		return err
	}

	d.Container.ID = finalCont.ID
	return nil
}

// RemoveContainer stops and removes the Docker container.
func (d *Docker) RemoveContainer() {}

// NewDocker creates a new Docker instance with the provided client
func NewDocker(client *client.Client) *Docker {
	return &Docker{
		Client: client,
	}
}
