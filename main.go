package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type APIConfig struct {
	Host string `hcl:"host"`
	Port int    `hcl:"port"`
}

type Config struct {
	API *APIConfig `hcl:"api"`
}

func main() {
	cfg, err := LoadConfigFile("config.hcl")
	if err != nil {
		log.Fatalf("failed file parsing %s", err)
	}

	log.Print(cfg.API)
}

func LoadConfigFile(path string) (*Config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseConfig(string(d))
}

func ParseConfig(d string) (*Config, error) {
	obj, err := hcl.Parse(d)
	if err != nil {
		return nil, err
	}

	var result Config
	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	// Parse seperate fields
	if o := list.Filter("api"); len(o.Items) > 0 {
		if err := parseAPI(&result, o); err != nil {
			return nil, err
		}
	}

	log.Print(result)

	return &result, nil
}

func parseAPI(result *Config, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'telemetry' block is permitted")
	}

	// Get our one item
	item := list.Items[0]

	log.Print(item)

	var c APIConfig
	if err := hcl.DecodeObject(&c, item.Val); err != nil {
		return err
	}

	if result.API == nil {
		result.API = &APIConfig{}
	}

	if err := hcl.DecodeObject(&result.API, item.Val); err != nil {
		return err
	}
	return nil
}
