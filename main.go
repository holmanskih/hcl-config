package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Config struct {
	A string `hcl:"a"`
}

func main() {
	cfg, err := ParseFile("config.hcl")
	if err != nil {
		log.Fatalf("failed file parsing %s", err)
	}

	log.Print(cfg)
}

func ParseFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// check extension hcl or json
	extension := filepath.Ext(filename)
	if len(extension) > 0 {
		extension = extension[1:]
	}

	return parseHCL(f, filename)
}

func parseHCL(r io.Reader, filename string) (*Config, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	f, diag := hclsyntax.ParseConfig(src, filename, hcl.Pos{})

	if diag.HasErrors() {
		return nil, diag
	}

	var config Config
	diag = gohcl.DecodeBody(f.Body, nil, &config)
	if diag.HasErrors() {
		return nil, diag
	}

	return &config, nil
}
