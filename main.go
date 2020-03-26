package main

import (
	"io/ioutil"
	"log"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/pkg/errors"
)

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

// Root config structure
type Config struct {
	API        APIConfig `hcl:"api"`
	EnableAuth bool      `hcl:"enable_auth"`
	Cache      CacheCfg  `hcl:"cache"`
}

var (
	DecodeHCLBlockError = errors.New("failed to decode the hcl block")
	FilterHCListError   = errors.New("only one hcl block is permitted")
)

func main() {
	cfg, err := LoadConfigFile("config.hcl")
	if err != nil {
		log.Fatalf("failed file parsing %s", err)
	}

	log.Printf("loaded config %v", cfg.Cache)
}

func LoadConfigFile(path string) (*Config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the file")
	}
	return ParseConfig(string(d))
}

func ParseConfig(d string) (*Config, error) {
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
			return nil, errors.Wrap(err, "failed to get the hcl block")
		}
	}

	if o := list.Filter("cache"); len(o.Items) > 0 {
		if err := parseCache(&result, o); err != nil {
			return nil, errors.Wrap(err, "failed to get the hcl block")
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
