package main

import (
	"encoding/xml"
	"fmt"
	"io"

	cfg "github.com/actofgod/goappconfig"
)

type XmlConfig struct {
	InputFile  string `xml:"inputFile"`
	OutputFile string `xml:"outputFile"`
}

func main() {
	builder := cfg.NewBuilder[XmlConfig]()
	builder = builder.With(cfg.ByteArrayDecoder(xml.Unmarshal))
	builder = builder.With(cfg.FileDecoder(func(reader io.Reader) cfg.Decoder {
		return xml.NewDecoder(reader)
	}))
	err := builder.Load("examples/xml/config.xml")
	if err != nil {
		panic(err)
	}
	config, err := builder.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Config is: %v\n", config)
}
