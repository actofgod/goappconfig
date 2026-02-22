package main

import (
	"fmt"
	"os"

	cfg "github.com/actofgod/goappconfig"
)

type PlainConfig struct {
	InputFile  string `json:"input_file" yaml:"input-file" arg:"input-file" short:"i" env:"INPUT_FILE"`
	OutputFile string `json:"output_file" yaml:"output-file" arg:"output-file" short:"o" env:"OUTPUT_FILE"`
}

func main() {
	builder := cfg.NewBuilder[PlainConfig]()
	builder = builder.With(cfg.CliArguments(os.Args[1:])).With(cfg.DisableEnv())
	err := builder.Load("config.json")
	if err != nil {
		panic(err)
	}
	config, err := builder.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Config is: %v\n", config)
}
