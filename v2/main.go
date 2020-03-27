package main

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/holmanskih/hcl-config/config"
)

func main() {
	var cfg config.Config
	err := hclsimple.DecodeFile("env/common.config.hcl", nil, &cfg)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	log.Printf("Configuration is %#v", cfg)
}

func parseHCL(r io.Reader, filename string) (*config.Config, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	f, diag := hclsyntax.ParseConfig(src, filename, hcl.Pos{})

	if diag.HasErrors() {
		return nil, diag
	}

	var config config.Config
	diag = gohcl.DecodeBody(f.Body, nil, &config)
	if diag.HasErrors() {
		return nil, diag
	}

	return &config, nil
}
