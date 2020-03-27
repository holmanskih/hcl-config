package parse

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/holmanskih/hcl-config/config"
)

var (
	DecodeHCLBlockError = errors.New("failed to decode the hcl block")
	FilterHCListError   = errors.New("only one hcl block is permitted")
	WrongHCLBlockLabel  = errors.New("wrong hcl block label name")
)

func LoadConfig(path, flag string) (*config.Config, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return LoadConfigDir(path, flag)
	}
	return LoadConfigFile(path, flag)
}

func LoadConfigDir(dir, flag string) (*config.Config, error) {
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
		config, err := LoadConfigFile(f, flag)
		if err != nil {
			return nil, errors.Wrapf(err, "error loading %s", f)
		}

		result = config
	}

	return result, nil
}

func LoadConfigFile(path, flag string) (*config.Config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the file")
	}
	return ParseConfig(string(d), flag)
}
