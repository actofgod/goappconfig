package goappconfig

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type Builder[T any] interface {
	With(option BuilderOption) Builder[T]
	Load(fileName string) error
	Build() (T, error)
	ApplyTo(config *T) error
}

type builderImpl[T any] struct {
	opts   Options
	config T
	typ    reflect.Type
	names  map[string]int
}

func NewBuilder[T any](opts ...BuilderOption) Builder[T] {
	var c T
	r := reflect.TypeOf(c)
	count := r.NumField()
	names := make(map[string]int, count)
	for i := 0; i < count; i++ {
		var field = r.Field(i)
		val, ok := field.Tag.Lookup("json")
		if ok {
			vals := strings.Split(val, ",")
			names[vals[0]] = i
		}
	}
	result := &builderImpl[T]{
		config: c,
		typ:    r,
		names:  names,
		opts:   newOptions(),
	}
	for _, option := range opts {
		result.With(option)
	}
	return result
}

func (b *builderImpl[T]) With(option BuilderOption) Builder[T] {
	b.opts = option(b.opts)
	return b
}

func (b *builderImpl[T]) Load(fileName string) error {
	s, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is directory", fileName)
	}
	if s.Size() > b.opts.maxConfigFileSize {
		return fmt.Errorf("file '%s' is to large for config file", fileName)
	}
	fd, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fd.Close()
	var config T
	err = b.getDecoder(fd).Decode(&config)
	if err != nil {
		return err
	}
	b.config = config
	return nil
}

func (b *builderImpl[T]) getDecoder(reader io.Reader) Decoder {
	if b.opts.fileDecoder != nil {
		return b.opts.fileDecoder(reader)
	}
	return json.NewDecoder(reader)
}

func (b *builderImpl[T]) Build() (T, error) {
	var config T
	err := b.ApplyTo(&config)
	return config, err
}

func (b *builderImpl[T]) ApplyTo(config *T) error {
	*config = b.config
	if len(b.opts.cliArgs) > 0 {
		err := b.parseCliArguments(config, b.opts.cliArgs, b.opts.useFlags)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *builderImpl[T]) parseCliArguments(config *T, cliArgs []string, useFlags bool) error {
	rv := reflect.ValueOf(config)
	rv = rv.Elem()
	if useFlags {
		return b.parseCliFlagArguments(rv, cliArgs)
	}
	return nil
}

func (b *builderImpl[T]) parseCliFlagArguments(rv reflect.Value, cliArgs []string) error {
	for k, v := range b.names {
		name := k
		flag.CommandLine.String(name, "", "")
		name = "p:" + k
		flag.CommandLine.String(name, "", "")
		name = "p:" + b.typ.Field(v).Name
		flag.CommandLine.String(name, "", "")
	}
	err := flag.CommandLine.Parse(cliArgs)
	if err != nil {
		return err
	}
	for k, v := range b.names {
		name := k
		f := flag.CommandLine.Lookup(name)
		if f != nil && f.Value.String() != "" {
			value := reflect.ValueOf(f.Value.String())
			rv.Field(v).Set(value)
			continue
		}
		name = "p:" + k
		f = flag.CommandLine.Lookup(name)
		if f != nil && f.Value.String() != "" {
			value := reflect.ValueOf(f.Value.String())
			rv.Field(v).Set(value)
			continue
		}
		name = "p:" + b.typ.Field(v).Name
		f = flag.CommandLine.Lookup(name)
		if f != nil && f.Value.String() != "" {
			field := rv.Field(v)
			if field.CanSet() {
				field.SetString(f.Value.String())
			} else {
				field.SetString(f.Value.String())
			}
			continue
		}
	}
	return nil
}
