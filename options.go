package goappconfig

import (
	"io"
	"os"
	"strings"
)

const defaultMaxConfigFileSize int64 = 1024 * 1024 * 10

type Options struct {
	configFileArgs    []string
	cliArgs           []string
	useFlags          bool
	applyEnv          bool
	maxConfigFileSize int64
	fileDecoder       DecoderConstructor
}

func newOptions() Options {
	return Options{
		cliArgs:           os.Args[1:],
		useFlags:          true,
		applyEnv:          true,
		maxConfigFileSize: defaultMaxConfigFileSize,
	}
}

type BuilderOption func(Options) Options

func CliArguments(cliArguments []string) BuilderOption {
	return func(opts Options) Options {
		opts.cliArgs = cliArguments
		return opts
	}
}

func DisableCliArguments() BuilderOption {
	return func(opts Options) Options {
		opts.cliArgs = nil
		return opts
	}
}

func DisableEnv() BuilderOption {
	return func(opts Options) Options {
		opts.applyEnv = false
		return opts
	}
}

func MaxConfigFileSize(maxFileSize int64) BuilderOption {
	return func(opts Options) Options {
		opts.maxConfigFileSize = maxFileSize
		return opts
	}
}

func FileDecoder(constructor DecoderConstructor) BuilderOption {
	return func(opts Options) Options {
		opts.fileDecoder = constructor
		return opts
	}
}

func ByteArrayDecoder(decoder ByteArrayDecoderFunc) BuilderOption {
	return func(opts Options) Options {
		opts.fileDecoder = func(reader io.Reader) Decoder {
			return NewBufferedDecoder(reader, decoder)
		}
		return opts
	}
}

func ConfigFileArguments(argumentName string, argumentNames ...string) BuilderOption {
	var arguments []string
	if len(argumentNames) == 0 {
		arguments = strings.Split(argumentName, ",")
	} else {
		arguments = make([]string, 0, len(argumentNames)+1)
		if len(argumentName) > 0 {
			arguments = append(arguments, argumentName)
		}
		for _, arg := range argumentNames {
			if len(argumentName) > 0 {
				arguments = append(arguments, arg)
			}
		}
	}
	return func(opts Options) Options {
		opts.configFileArgs = arguments
		return opts
	}
}
