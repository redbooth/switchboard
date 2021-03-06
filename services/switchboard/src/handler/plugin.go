package handler

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"plugin"
)

type PluginConf struct {
	Filename string
	Options  map[string]interface{}
}

type Plugin struct {
	conf PluginConf
	plug *plugin.Plugin
}

func NewPlugin(conf PluginConf) Handler {
	// load plugin
	plug, err := plugin.Open(conf.Filename)
	if err != nil {
		log.Panicf("Unable to open plugin file %s: %v\n", conf.Filename, err)
	}

	// retrieve constructor function
	constructorsym, err := plug.Lookup("Constructor")
	if err != nil {
		log.Panicf("Unable to find Constructor function in %s: %v\n", conf.Filename, err)
	}

	// assert that Constructor symbol is of the desired type
	constructor, ok := constructorsym.(Constructor)
	if !ok {
		log.Panicf("Unexpected type for Constructor function in %s: %v\n", conf.Filename, err)
	}

	// retrieve plugin conf
	plugconfsym, err := plug.Lookup("Conf")
	if err != nil {
		log.Panicf("Unable to find Conf function in %s: %v\n", conf.Filename, err)
	}

	// assert that Conf symbol is of the desired type
	plugconf, ok := plugconfsym.(func() *Conf)
	if !ok {
		log.Panicf("Unexpected type for Conf function in %s: %v\n", conf.Filename, err)
	}

	// convert plugin options map to a struct
	err = mapstructure.Decode(&plugconf, conf.Options)
	if err != nil {
		log.Panicf("Unable to parse configuration for plugin %s: %v\n", conf.Filename, err)
	}

	return constructor(plugconf)
}
