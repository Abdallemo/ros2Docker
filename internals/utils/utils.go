package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Config interface {
	GetSaveData() (data any, absFile string)
}

func SaveConfig(baseDir string, s Config) error {
	data, fileName := s.GetSaveData()
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	savePath := filepath.Join(baseDir, fileName)
	if err := os.WriteFile(savePath, b, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", savePath, err)
	}

	return nil
}

type GlbalConfig struct {
	Location string
	Folder   string
	OS       string
}

func (cg GlbalConfig) GetSaveData() (any, string) {
	return cg, "config.json"
}
func NewGlobalConfig(folder string) (*GlbalConfig, error) {

	userCgfDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	if pathEnv, ok := os.LookupEnv("ROS2DOCKER_CONFIG_DIR"); ok {
		absPath, err := filepath.Abs(pathEnv)
		if err != nil {
			return nil, fmt.Errorf("invalid config path %s: %w", pathEnv, err)
		}
		return &GlbalConfig{
			Location: absPath,
			OS:       runtime.GOOS,
		}, nil
	}
	location := filepath.Join(userCgfDir, folder)

	err = os.MkdirAll(location, 0755)
	if err != nil {
		return nil, err
	}

	return &GlbalConfig{
		Location: location,
		Folder:   folder,
		OS:       runtime.GOOS,
	}, nil
}

func LoadConfigs(g Config) error {
	return nil
}
func IsValidPath(path string) error {
	if path != "" {
		fpath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("unable to resolve filepath/folder: %v\n", err)
		}

		info, err := os.Stat(fpath)
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s\n", fpath)
		}

		if err != nil {
			return fmt.Errorf("error accessing path: %v\n", err)
		}

		if !info.IsDir() {
			return fmt.Errorf("path is not a directory: %s\n", fpath)

		}

	}
	return nil

}
func GenerateMix(val string) string {
	return strings.Join([]string{val, rand.Text()}, "-")

}
