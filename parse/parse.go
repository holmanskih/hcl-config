package parse

import (
	"log"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/pkg/errors"

	"github.com/holmanskih/hcl-config/config"
)

type parseBlockFunc func(*config.Config, *ast.ObjectList) error

//type parseLabelBlockFunc func(*config.Config, *ast.ObjectList, string) error
//
//type Parser interface {
//	ParseBlock() error
//}
//
//func (f parseBlockFunc) ParseBlock() error {
//	return nil
//}
//
//func (f parseLabelBlockFunc) ParseBlock() error  {
//	return nil
//}
func ParseConfig(d, flag string) (*config.Config, error) {
	obj, err := hcl.Parse(d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the hcl file")
	}

	var cfg config.Config
	if err := hcl.DecodeObject(&cfg, obj); err != nil {
		return nil, err
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, errors.Wrap(err, "file doesnt contain the root object")
	}

	// Parse hcl blocks
	//err = parseHCLBlock(cfg, "api", 0, list, parseAPI)
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to get the api hcl block")
	//}

	if o := list.Filter("api"); len(o.Items) > 0 {
		if err := parseAPI(&cfg, o); err != nil {
			return nil, errors.Wrap(err, "failed to get the api hcl block")
		}
	}

	if o := list.Filter("cache"); len(o.Items) > 0 {
		if err := parseCache(&cfg, o); err != nil {
			return nil, errors.Wrap(err, "failed to get the cache hcl block")
		}
	}

	if o := list.Filter("rabbitmq"); len(o.Items) > 0 {
		if err := parseRabbit(&cfg, o, flag); err != nil {
			return nil, errors.Wrap(err, "failed to get the rabbitmq hcl block")
		}
	}

	return &cfg, nil
}

func parseHCLBlock(cfg config.Config, hclBlockName string, cfgBlocksNumber int, list *ast.ObjectList, parse parseBlockFunc) error {
	block := list.Filter(hclBlockName)
	log.Printf("block sizes are %v and %v", len(block.Items), cfgBlocksNumber)

	if len(block.Items) > cfgBlocksNumber {
		err := parse(&cfg, block)
		if err != nil {
			return errors.Wrap(err, "failed to get the cache hcl block")
		}
		return nil
	}

	return FilterHCListError
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
