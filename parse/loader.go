package parse

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/pkg/errors"

	"github.com/holmanskih/hcl-config/config"
)

func LoadConfig(path string) (*config.Config, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return LoadConfigDir(path)
	}
	return LoadConfigFile(path)
}

func LoadConfigDir(dir string) (*config.Config, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var files []string
	err = nil
	for err != io.EOF {
		var fileInfos []os.FileInfo
		fileInfos, err = f.Readdir(128)
		if err != nil && err != io.EOF {
			return nil, err
		}

		for _, fileInfo := range fileInfos {
			// Ignore directories
			if fileInfo.IsDir() {
				continue
			}

			// Filter files with .hcl extension
			name := fileInfo.Name()
			if strings.HasSuffix(name, ".hcl") {
				path := filepath.Join(dir, name)
				files = append(files, path)
			}
		}
	}

	var result *config.Config
	for _, f := range files {
		config, err := LoadConfigFile(f)
		if err != nil {
			return nil, errors.Wrapf(err, "error loading %s", f)
		}

		result = config
	}

	return result, nil
}

func LoadConfigFile(path string) (*config.Config, error) {
	var cfg config.Config
	err := hclsimple.DecodeFile(path, nil, &cfg)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	return &cfg, nil
}
