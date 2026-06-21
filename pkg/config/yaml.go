package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

// LoadYAML decodes the YAML file at path into dst.
// Use in combination with LoadEnv to let environment variables override file values:
//
//	if err := config.LoadYAML("config.yaml", &cfg); err != nil { ... }
//	if err := config.LoadEnv(&cfg); err != nil { ... }
func LoadYAML(path string, dst any) error {
	dir, err := os.OpenRoot(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("open config dir: %w", err)
	}
	defer func() { _ = dir.Close() }()

	file, err := dir.Open(filepath.Base(path))
	if err != nil {
		return fmt.Errorf("open config file: %w", err)
	}
	defer func() { _ = file.Close() }()

	if err := yaml.NewDecoder(file).Decode(dst); err != nil {
		return fmt.Errorf("decode yaml: %w", err)
	}

	return nil
}
