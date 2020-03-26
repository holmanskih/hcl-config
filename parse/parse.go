package parse

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
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
			return nil, errors.Wrapf(err, "error loading %q", f)
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

func ParseConfig(d, flag string) (*config.Config, error) {
	obj, err := hcl.Parse(d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the hcl file")
	}

	var result config.Config
	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, errors.Wrap(err, "file doesnt contain the root object")
	}

	// Parse hcl blocks
	if o := list.Filter("api"); len(o.Items) > 0 {
		if err := parseAPI(&result, o); err != nil {
			return nil, errors.Wrap(err, "failed to get the api hcl block")
		}
	}

	if o := list.Filter("cache"); len(o.Items) > 0 {
		if err := parseCache(&result, o); err != nil {
			return nil, errors.Wrap(err, "failed to get the cache hcl block")
		}
	}

	if o := list.Filter("rabbitmq"); len(o.Items) > 0 {
		if err := parseRabbit(&result, o, flag); err != nil {
			return nil, errors.Wrap(err, "failed to get the rabbitmq hcl block")
		}
	}

	return &result, nil
}

func parseAPI(result *config.Config, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return FilterHCListError
	}

	item := list.Items[0]

	var c config.APIConfig
	if err := hcl.DecodeObject(&c, item.Val); err != nil {
		return errors.Wrap(err, "decode hcl block err")
	}

	if result.API == (config.APIConfig{}) {
		result.API = config.APIConfig{}
	}

	if err := hcl.DecodeObject(&result.API, item.Val); err != nil {
		return DecodeHCLBlockError
	}
	return nil
}

func parseCache(result *config.Config, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return FilterHCListError
	}

	item := list.Items[0]

	var c config.CacheCfg
	if err := hcl.DecodeObject(&c, item.Val); err != nil {
		return errors.Wrap(err, "decode hcl block err")
	}

	if result.Cache == (config.CacheCfg{}) {
		result.Cache = config.CacheCfg{}
	}

	if err := hcl.DecodeObject(&result.Cache, item.Val); err != nil {
		return DecodeHCLBlockError
	}
	return nil
}

func parseRabbit(result *config.Config, list *ast.ObjectList, flag string) error {
	if len(list.Items) > 3 {
		return FilterHCListError
	}

	filteredList := list.Filter(flag)
	if len(filteredList.Items) == 0 {
		return WrongHCLBlockLabel
	}

	var m config.RabbitMQCfg
	for _, item := range filteredList.Items {

		if err := hcl.DecodeObject(&m, item); err != nil {
			return DecodeHCLBlockError
		}

		if err := hcl.DecodeObject(&result.Rabbit, item.Val); err != nil {
			return DecodeHCLBlockError
		}
	}
	return nil
}
