package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/pkg/errors"
)

const cfgFlag = "master"

type APIConfig struct {
	Host string `hcl:"host"`
	Port int    `hcl:"port"`
}

type NutsDBCfg struct {
	Path        string `hcl:"path"`
	SegmentSize int64  `hcl:"segment_size"`
}

type RedisConf struct {
	DevMode  bool   `hcl:"dev_mode"`
	Password string `hcl:"password"`
	Host     string `hcl:"host"`
}

type CacheCfg struct {
	Type string `hcl:"type"`

	Redis  RedisConf `hcl:"redis"`
	NutsDB NutsDBCfg `hcl:"nutsdb"`
}

type CommonCfg struct {
	Exchange     string `hcl:"exchange"`
	ExchangeType string `hcl:"exchange_type"`
}

type RabbitMQCfg struct {
	Host        string    `hcl:"host"`
	User        string    `hcl:"user"`
	Password    string    `hcl:"password"`
	ConsumerTag string    `hcl:"consumer_tag"`
	Common      CommonCfg `hcl:"common"`
}

// Root config structure
type Config struct {
	API        APIConfig   `hcl:"api"`
	EnableAuth bool        `hcl:"enable_auth"`
	Cache      CacheCfg    `hcl:"cache"`
	Rabbit     RabbitMQCfg `hcl:"rabbitmq"`
}

var (
	DecodeHCLBlockError = errors.New("failed to decode the hcl block")
	FilterHCListError   = errors.New("only one hcl block is permitted")
	WrongHCLBlockLabel  = errors.New("wrong hcl block label name")
)

func main() {
	cfg, err := LoadConfig("env", cfgFlag)
	if err != nil {
		log.Fatalf("failed file parsing: %s", err)
	}

	log.Printf("loaded config %v", cfg.Rabbit)
}

// LoadConfig loads the configuration at the given path, regardless if
// its a file or directory.
func LoadConfig(path, flag string) (*Config, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return LoadConfigDir(path, flag)
	}
	return LoadConfigFile(path, flag)
}

// LoadConfigDir loads all the configurations in the given directory
// in alphabetical order.
func LoadConfigDir(dir, flag string) (*Config, error) {
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

	var result *Config
	for _, f := range files {
		config, err := LoadConfigFile(f, flag)
		if err != nil {
			return nil, errors.Wrapf(err, "error loading %q", f)
		}

		result = config
	}

	return result, nil
}

func LoadConfigFile(path, flag string) (*Config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the file")
	}
	return ParseConfig(string(d), flag)
}

func ParseConfig(d, flag string) (*Config, error) {
	obj, err := hcl.Parse(d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the hcl file")
	}

	var result Config
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

func parseAPI(result *Config, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return FilterHCListError
	}

	item := list.Items[0]

	var c APIConfig
	if err := hcl.DecodeObject(&c, item.Val); err != nil {
		return errors.Wrap(err, "decode hcl block err")
	}

	if result.API == (APIConfig{}) {
		result.API = APIConfig{}
	}

	if err := hcl.DecodeObject(&result.API, item.Val); err != nil {
		return DecodeHCLBlockError
	}
	return nil
}

func parseCache(result *Config, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return FilterHCListError
	}

	item := list.Items[0]

	var c CacheCfg
	if err := hcl.DecodeObject(&c, item.Val); err != nil {
		return errors.Wrap(err, "decode hcl block err")
	}

	if result.Cache == (CacheCfg{}) {
		result.Cache = CacheCfg{}
	}

	if err := hcl.DecodeObject(&result.Cache, item.Val); err != nil {
		return DecodeHCLBlockError
	}
	return nil
}

func parseRabbit(result *Config, list *ast.ObjectList, flag string) error {
	if len(list.Items) > 3 {
		return FilterHCListError
	}

	filteredList := list.Filter(flag)
	if len(filteredList.Items) == 0 {
		return WrongHCLBlockLabel
	}

	var m RabbitMQCfg
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
