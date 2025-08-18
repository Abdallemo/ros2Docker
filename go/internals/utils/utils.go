package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Location string
	Folder   string
	OS       string
}

func NewConfig(folder string) (*Config, error) {

	userCgfDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	if pathEnv, ok := os.LookupEnv("ROS2DOCKER_CONFIG_DIR"); ok {
		absPath, err := filepath.Abs(pathEnv)
		if err != nil {
			return nil, fmt.Errorf("invalid config path %s: %w", pathEnv, err)
		}
		return &Config{
			Location: absPath,
			OS:       runtime.GOOS,
		}, nil
	}
	location := filepath.Join(userCgfDir, folder)

	err = os.MkdirAll(location, 0755)
	if err != nil {
		return nil, err
	}

	return &Config{
		Location: location,
		Folder:   folder,
		OS:       runtime.GOOS,
	}, nil
}
