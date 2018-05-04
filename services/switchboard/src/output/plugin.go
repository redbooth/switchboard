package output

import (
	"../header"
	"github.com/mitchellh/mapstructure"
	"log"
	"plugin"
)

type PluginConf struct {
	Filename string
	Options  map[string]interface{}
}

func NewPlugin(conf PluginConf, errors chan<- error, h header.Header) Output {
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
		log.Panicf("Unexpected type for Constructor function in %s: %T\n", conf.Filename, constructorsym)
	}

	// retrieve plugin conf
	plugconfsym, err := plug.Lookup("Conf")
	if err != nil {
		log.Panicf("Unable to find Conf function in %s: %v\n", conf.Filename, err)
	}

	// assert that Conf symbol is of the desired type
	plugconffunc, ok := plugconfsym.(func() *Conf)
	if !ok {
		log.Panicf("Unexpected type for Conf function in %s: %T\n", conf.Filename, plugconfsym)
	}

	// convert plugin options map to a struct
	plugconf := plugconffunc()
	err = mapstructure.Decode(&plugconf, conf.Options)
	if err != nil {
		log.Panicf("Unable to parse configuration for plugin %s: %v\n", conf.Filename, err)
	}

	return constructor(plugconf, errors, h)
}
