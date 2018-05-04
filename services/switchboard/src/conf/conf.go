package conf

import (
	"../filter"
	"../handler"
	"../header"
	"../input"
	"../output"
	"../terminus"
	"../transformer"
	"flag"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
)

/*
Note: Here there be dragons.

All ugliness in switchboard is concentrated in this one file.

Go doesn't make it easy to convert dynamic templates like yaml files into
arbitrarily structured data.  I'm sure what's here can be accomplished more
cleanly somehow (code generation? reflection? some combination of the two?)
and I'm sure there's a "standard go way" to do it, but I couldn't find one.
*/

/*
 TODO
 - treat Header section differently than the other sections
   (it should be a single map[string]interface{}, not a slice of them)
*/

var (
	confdefault string = "/etc/switchboard/conf.yml"
	conffile    string
)

func init() {
	flag.StringVar(&conffile, "conf", confdefault, "yaml configuration file")
	flag.Parse()
}

// raw yaml maps
type Stages map[string]Stage

type Stage []Service

type Service map[string]interface{}

// parsed structs
type Conf struct {
	NewHeader    func() header.Header
	Inputs       []input.Bundle
	Filters      []filter.Bundle
	Transformers []transformer.Bundle
	Outputs      []output.Bundle
	Termini      []terminus.Bundle
	Handlers     []handler.Bundle
}

func NewConf() Conf {
	var conf Conf
	if stages, err := Read(conffile); err != nil {
		log.Fatalf("Unable to load custom configuration file %s: %v\n", conffile, err)
	} else {
		conf.AddHeader(stages["header"])
		conf.AddStages(stages)
	}
	return conf
}

func Read(filename string) (Stages, error) {
	var stages Stages
	if file, err := ioutil.ReadFile(filename); err != nil {
		return stages, fmt.Errorf("Unable to find configuration file %s: %v\n", filename, err)
	} else if err = yaml.UnmarshalStrict(file, &stages); err != nil {
		return stages, fmt.Errorf("Unable to parse configuration file %s: %v\n", filename, err)
	}
	return stages, nil
}

func (conf *Conf) AddHeader(stage Stage) {
	var (
		service Service = stage[0]
		stype   string  = service["type"].(string)
	)
	switch stype {
	case "plugin":
		tmpconf := header.PluginConf{}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(tmpconf)
			log.Fatalf("Unable to parse configuration for service header/%s: %v\n", stype, err)
			return
		} else {
			conf.NewHeader = func() header.Header {
				return header.NewPlugin(tmpconf)
			}
		}
	case "unstructured":
		tmpconf := header.UnstructuredConf{}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(tmpconf)
			log.Fatalf("Unable to parse configuration for service header/%s: %v\n", stype, err)
			return
		} else {
			conf.NewHeader = func() header.Header {
				return header.NewUnstructured(tmpconf)
			}
		}
	default:
		tmpconf := header.EmptyConf{}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(tmpconf)
			log.Fatalf("Unable to parse configuration for service header/%s: %v\n", stype, err)
			return
		} else {
			conf.NewHeader = func() header.Header {
				return header.NewEmpty(tmpconf)
			}
		}
	}
}

func (conf *Conf) AddStages(stages Stages) {
	for name, stage := range stages {
		conf.AddStage(name, stage)
	}
}

func (conf *Conf) AddStage(name string, stage Stage) {
	for _, service := range stage {
		conf.AddService(name, service)
	}
}

func (conf *Conf) AddService(stage string, service Service) {
	switch stage {
	case "input":
		conf.AddInputService(service)
	case "filter":
		conf.AddFilterService(service)
	case "transformer":
		conf.AddTransformerService(service)
	case "output":
		conf.AddOutputService(service)
	case "terminus":
		conf.AddTerminusService(service)
	case "handler":
		conf.AddHandlerService(service)
	}
}

func (conf *Conf) AddInputService(service Service) {
	var (
		stype string = service["type"].(string)
		sconf input.Conf
		scons input.Constructor
	)
	switch stype {
	case "file":
		sconf = input.FileConf{}
		scons = func(c input.Conf, e chan<- error, r chan<- io.ReadCloser) input.Input {
			ic, ok := c.(input.FileConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return input.NewFile(ic, e, r)
		}
	case "http":
		sconf = input.HTTPConf{}
		scons = func(c input.Conf, e chan<- error, r chan<- io.ReadCloser) input.Input {
			ic, ok := c.(input.HTTPConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return input.NewHTTP(ic, e, r)
		}
	case "plugin":
		sconf = input.PluginConf{}
		scons = func(c input.Conf, e chan<- error, r chan<- io.ReadCloser) input.Input {
			ic, ok := c.(input.PluginConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return input.NewPlugin(ic, e, r)
		}
	case "s3":
		sconf = input.S3Conf{}
		scons = func(c input.Conf, e chan<- error, r chan<- io.ReadCloser) input.Input {
			ic, ok := c.(input.S3Conf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return input.NewS3(ic, e, r)
		}
	case "stdin":
		sconf = input.StdinConf{}
		scons = func(c input.Conf, e chan<- error, r chan<- io.ReadCloser) input.Input {
			ic, ok := c.(input.StdinConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return input.NewStdin(ic, e, r)
		}
	case "tcp":
		tmpconf := input.TCPConf{}
		sconf = input.TCPConf{}
		scons = func(c input.Conf, e chan<- error, r chan<- io.ReadCloser) input.Input {
			ic, ok := c.(input.TCPConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return input.NewTCP(ic, e, r)
		}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(sconf)
			log.Fatalf("Unable to parse configuration for service input/%s: %v\n", stype, err)
			return
		} else {
			conf.AddInput(tmpconf, scons)
		}
	case "unixsocket":
		sconf = input.UnixSocketConf{}
		scons = func(c input.Conf, e chan<- error, r chan<- io.ReadCloser) input.Input {
			ic, ok := c.(input.UnixSocketConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return input.NewUnixSocket(ic, e, r)
		}
	case "websocket":
		tmpconf := input.WebsocketConf{}
		sconf = input.WebsocketConf{}
		scons = func(c input.Conf, e chan<- error, r chan<- io.ReadCloser) input.Input {
			ic, ok := c.(input.WebsocketConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return input.NewWebSocket(ic, e, r)
		}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(sconf)
			log.Fatalf("Unable to parse configuration for service input/%s: %v\n", stype, err)
			return
		} else {
			conf.AddInput(tmpconf, scons)
		}
	default:
		log.Fatalf("Unrecognized service input/%s\n", stype)
		return
	}
}

func (conf *Conf) AddFilterService(service Service) {
	var (
		stype string = service["type"].(string)
		sconf filter.Conf
		scons filter.Constructor
	)
	switch stype {
	case "allow":
		sconf = filter.AllowConf{}
		scons = func(c filter.Conf, e chan<- error) filter.Filter {
			ic, ok := c.(filter.AllowConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return filter.NewAllow(ic, e)
		}
	case "deny":
		sconf = filter.DenyConf{}
		scons = func(c filter.Conf, e chan<- error) filter.Filter {
			ic, ok := c.(filter.DenyConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return filter.NewDeny(ic, e)
		}
	case "exec":
		sconf = filter.ExecConf{}
		scons = func(c filter.Conf, e chan<- error) filter.Filter {
			ic, ok := c.(filter.ExecConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return filter.NewExec(ic, e)
		}
	case "http":
		sconf = filter.HTTPConf{}
		scons = func(c filter.Conf, e chan<- error) filter.Filter {
			ic, ok := c.(filter.HTTPConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return filter.NewHTTP(ic, e)
		}
	case "plugin":
		sconf = filter.PluginConf{}
		scons = func(c filter.Conf, e chan<- error) filter.Filter {
			ic, ok := c.(filter.PluginConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return filter.NewPlugin(ic, e)
		}
	case "tcp":
		sconf = filter.TCPConf{}
		scons = func(c filter.Conf, e chan<- error) filter.Filter {
			ic, ok := c.(filter.TCPConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return filter.NewTCP(ic, e)
		}
	case "unixsocket":
		sconf = filter.UnixSocketConf{}
		scons = func(c filter.Conf, e chan<- error) filter.Filter {
			ic, ok := c.(filter.UnixSocketConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return filter.NewUnixSocket(ic, e)
		}
	default:
		log.Fatalf("Unrecognized service filter/%s\n", stype)
		return
	}
	err := mapstructure.Decode(service, &sconf)
	if err != nil {
		log.Fatalf("Unable to parse configuration for service filter/%s: %v\n", stype, err)
		return
	} else {
		conf.AddFilter(sconf, scons)
	}
}

func (conf *Conf) AddTransformerService(service Service) {
	var (
		stype string = service["type"].(string)
		sconf transformer.Conf
		scons transformer.Constructor
	)
	switch stype {
	case "exec":
		sconf = transformer.ExecConf{}
		scons = func(c transformer.Conf, e chan<- error) transformer.Transformer {
			ic, ok := c.(transformer.ExecConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return transformer.NewExec(ic, e)
		}
	case "identity":
		sconf = transformer.IdentityConf{}
		scons = func(c transformer.Conf, e chan<- error) transformer.Transformer {
			ic, ok := c.(transformer.IdentityConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return transformer.NewIdentity(ic, e)
		}
	case "plugin":
		sconf = transformer.PluginConf{}
		scons = func(c transformer.Conf, e chan<- error) transformer.Transformer {
			ic, ok := c.(transformer.PluginConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return transformer.NewPlugin(ic, e)
		}
	case "tcp":
		sconf = transformer.TCPConf{}
		scons = func(c transformer.Conf, e chan<- error) transformer.Transformer {
			ic, ok := c.(transformer.TCPConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return transformer.NewTCP(ic, e)
		}
	case "unixsocket":
		sconf = transformer.UnixSocketConf{}
		scons = func(c transformer.Conf, e chan<- error) transformer.Transformer {
			ic, ok := c.(transformer.UnixSocketConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return transformer.NewUnixSocket(ic, e)
		}
	default:
		log.Fatalf("Unrecognized service transformer/%s\n", stype)
		return
	}
	err := mapstructure.Decode(service, &sconf)
	if err != nil {
		log.Fatalf("Unable to parse configuration for service transformer/%s: %v\n", stype, err)
		return
	} else {
		conf.AddTransformer(sconf, scons)
	}
}

func (conf *Conf) AddOutputService(service Service) {
	var (
		stype string = service["type"].(string)
		sconf output.Conf
		scons output.Constructor
	)
	switch stype {
	case "file":
		tmpconf := output.FileConf{}
		sconf = output.FileConf{}
		scons = func(c output.Conf, e chan<- error, h header.Header) output.Output {
			ic, ok := c.(output.FileConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return output.NewFile(ic, e, h)
		}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(sconf)
			log.Fatalf("Unable to parse configuration for service output/%s: %v\n", stype, err)
			return
		} else {
			conf.AddOutput(tmpconf, scons)
		}
	case "http":
		sconf = output.HTTPConf{}
		scons = func(c output.Conf, e chan<- error, h header.Header) output.Output {
			ic, ok := c.(output.HTTPConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return output.NewHTTP(ic, e, h)
		}
	case "plugin":
		sconf = output.PluginConf{}
		scons = func(c output.Conf, e chan<- error, h header.Header) output.Output {
			ic, ok := c.(output.PluginConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return output.NewPlugin(ic, e, h)
		}
	case "s3":
		tmpconf := output.S3Conf{}
		sconf = output.S3Conf{}
		scons = func(c output.Conf, e chan<- error, h header.Header) output.Output {
			ic, ok := c.(output.S3Conf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return output.NewS3(ic, e, h)
		}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(sconf)
			log.Fatalf("Unable to parse configuration for service output/%s: %v\n", stype, err)
			return
		} else {
			conf.AddOutput(tmpconf, scons)
		}
	case "stdout":
		tmpconf := output.StdoutConf{}
		sconf = output.StdoutConf{}
		scons = func(c output.Conf, e chan<- error, h header.Header) output.Output {
			ic, ok := c.(output.StdoutConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return output.NewStdout(ic, e, h)
		}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(sconf)
			log.Fatalf("Unable to parse configuration for service output/%s: %v\n", stype, err)
			return
		} else {
			conf.AddOutput(tmpconf, scons)
		}
	case "tcp":
		tmpconf := output.TCPConf{}
		sconf = output.TCPConf{}
		scons = func(c output.Conf, e chan<- error, h header.Header) output.Output {
			ic, ok := c.(output.TCPConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return output.NewTCP(ic, e, h)
		}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(sconf)
			log.Fatalf("Unable to parse configuration for service output/%s: %v\n", stype, err)
			return
		} else {
			conf.AddOutput(tmpconf, scons)
		}
	case "unixsocket":
		sconf = output.UnixSocketConf{}
		scons = func(c output.Conf, e chan<- error, h header.Header) output.Output {
			ic, ok := c.(output.UnixSocketConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return output.NewUnixSocket(ic, e, h)
		}
	case "websocket":
		sconf = output.WebsocketConf{}
		scons = func(c output.Conf, e chan<- error, h header.Header) output.Output {
			ic, ok := c.(output.WebsocketConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return output.NewWebSocket(ic, e, h)
		}
	default:
		log.Fatalf("Unrecognized service output/%s\n", stype)
		return
	}
}

func (conf *Conf) AddTerminusService(service Service) {
	var (
		stype string = service["type"].(string)
		sconf terminus.Conf
		scons terminus.Constructor
	)
	switch stype {
	case "file":
		tmpconf := terminus.FileConf{}
		sconf = terminus.FileConf{}
		scons = func(c terminus.Conf, e chan<- error) terminus.Terminus {
			ic, ok := c.(terminus.FileConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return terminus.NewFile(ic, e)
		}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(sconf)
			log.Fatalf("Unable to parse configuration for service terminus/%s: %v\n", stype, err)
			return
		} else {
			conf.AddTerminus(tmpconf, scons)
		}
	default:
		log.Fatalf("Unrecognized service output/%s\n", stype)
		return
	}
}

func (conf *Conf) AddHandlerService(service Service) {
	var (
		stype string = service["type"].(string)
		sconf handler.Conf
		scons handler.Constructor
	)
	switch stype {
	case "file":
		sconf = handler.FileConf{}
		scons = func(c handler.Conf) handler.Handler {
			ic, ok := c.(handler.FileConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return handler.NewFile(ic)
		}
	case "http":
		sconf = handler.HTTPConf{}
		scons = func(c handler.Conf) handler.Handler {
			ic, ok := c.(handler.HTTPConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return handler.NewHTTP(ic)
		}
	case "log":
		tmpconf := handler.LogConf{}
		sconf = handler.LogConf{}
		scons = func(c handler.Conf) handler.Handler {
			ic, ok := c.(handler.LogConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return handler.NewLog(ic)
		}
		err := mapstructure.Decode(service, &tmpconf)
		if err != nil {
			log.Println(sconf)
			log.Fatalf("Unable to parse configuration for service handler/%s: %v\n", stype, err)
			return
		} else {
			conf.AddHandler(tmpconf, scons)
		}
	case "plugin":
		sconf = handler.PluginConf{}
		scons = func(c handler.Conf) handler.Handler {
			ic, ok := c.(handler.PluginConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return handler.NewPlugin(ic)
		}
	case "stderr":
		sconf = handler.StderrConf{}
		scons = func(c handler.Conf) handler.Handler {
			ic, ok := c.(handler.StderrConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return handler.NewStderr(ic)
		}
	case "tcp":
		sconf = handler.TCPConf{}
		scons = func(c handler.Conf) handler.Handler {
			ic, ok := c.(handler.TCPConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return handler.NewTCP(ic)
		}
	case "unixsocket":
		sconf = handler.UnixSocket{}
		scons = func(c handler.Conf) handler.Handler {
			ic, ok := c.(handler.UnixSocketConf)
			if !ok {
				log.Panicf("Unexpected type for Configuration: %v", c)
			}
			return handler.NewUnixSocket(ic)
		}
	default:
		log.Fatalf("Unrecognized service handler/%s\n", stype)
		return
	}
}

func (c *Conf) AddInput(conf input.Conf, constructor input.Constructor) {
	c.Inputs = append(c.Inputs, input.Bundle{conf, constructor})
}

func (c *Conf) AddFilter(conf filter.Conf, constructor filter.Constructor) {
	c.Filters = append(c.Filters, filter.Bundle{conf, constructor})
}

func (c *Conf) AddTransformer(conf transformer.Conf, constructor transformer.Constructor) {
	c.Transformers = append(c.Transformers, transformer.Bundle{conf, constructor})
}

func (c *Conf) AddOutput(conf output.Conf, constructor output.Constructor) {
	c.Outputs = append(c.Outputs, output.Bundle{conf, constructor})
}

func (c *Conf) AddTerminus(conf terminus.Conf, constructor terminus.Constructor) {
	c.Termini = append(c.Termini, terminus.Bundle{conf, constructor})
}

func (c *Conf) AddHandler(conf handler.Conf, constructor handler.Constructor) {
	c.Handlers = append(c.Handlers, handler.Bundle{conf, constructor})
}
