package main

import (
	"encoding/xml"
	"fmt"

	cfg "github.com/actofgod/goappconfig"
)

type SubConfig struct {
	MaxFileSize string `xml:"maxFileSize"`
	FileFormat  string `xml:"fileFormat"`
}

type XmlConfig struct {
	InputFile  string     `xml:"inputFile"`
	OutputFile string     `xml:"outputFile"`
	Config     *SubConfig `xml:"config"`
}

func main() {
	builder := cfg.NewBuilder[XmlConfig]()
	builder = builder.With(cfg.ByteArrayDecoder(xml.Unmarshal))
	err := builder.Load("examples/xml/config.xml")
	if err != nil {
		panic(err)
	}
	config, err := builder.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Config is: %v\n", config)
	fmt.Printf("SubConfig is: %v\n", config.Config)
}
