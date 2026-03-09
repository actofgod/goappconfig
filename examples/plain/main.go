package main

import (
	"fmt"

	cfg "github.com/actofgod/goappconfig"
)

type PlainConfig struct {
	InputFile  string `json:"input_file" yaml:"input-file" arg:"input-file" short:"i" env:"INPUT_FILE"`
	OutputFile string `json:"output_file" yaml:"output-file" arg:"output-file" short:"o" env:"OUTPUT_FILE"`
	// IntValue parameter description
	IntValue int `json:"int_value"`
	// UIntValue parameter description
	UIntValue uint16 `json:"uint_value" arg:"uint16"`
	// BoolValue parameter description
	BoolValue bool `json:"bool_value"`
}

func main() {
	builder := cfg.NewBuilder[PlainConfig]()
	builder = builder.With(cfg.ConfigFileArguments("config,c"))
	config, err := builder.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Config is: %v\n", config)
}
