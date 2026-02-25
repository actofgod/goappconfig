package goappconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

type Builder[T any] interface {
	// With method applies option to current builder config.
	With(option BuilderOption) Builder[T]

	// Load loads configuration form fileName using decoder configured via With method. JSON decoder used by default.
	Load(fileName string) error
	Build() (T, error)
	ApplyTo(config *T) error
}

type builderImpl[T any] struct {
	opts       Options
	config     T
	typ        reflect.Type
	properties []*propertyImpl
}

func NewBuilder[T any](opts ...BuilderOption) Builder[T] {
	var c T
	r := reflect.TypeOf(c)
	props := buildPropertyList(r, nil)
	result := &builderImpl[T]{
		config:     c,
		typ:        r,
		properties: props,
		opts:       newOptions(),
	}
	for _, option := range opts {
		result.With(option)
	}
	return result
}

// With method applies option to current builder config
func (b *builderImpl[T]) With(option BuilderOption) Builder[T] {
	b.opts = option(b.opts)
	return b
}

// Load loads configuration form fileName using decoder configured via With method. JSON decoder used by default.
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
	if b.opts.applyEnv {
		err := b.applyEnvironmentVariables(config)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildPropertyList(t reflect.Type, parent *propertyImpl) []*propertyImpl {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	count := t.NumField()
	props := make([]*propertyImpl, 0, count)
	for i := 0; i < count; i++ {
		var field = t.Field(i)
		switch field.Type.Kind() {
		case reflect.Array, reflect.Slice, reflect.Struct, reflect.Pointer:
			local := newProperty(parent, field)
			p := buildPropertyList(field.Type, local)
			props = append(props, p...)
		default:
			props = append(props, newProperty(parent, field))
		}
	}
	return props
}

func (b *builderImpl[T]) getDecoder(reader io.Reader) Decoder {
	if b.opts.fileDecoder != nil {
		return b.opts.fileDecoder(reader)
	}
	return json.NewDecoder(reader)
}

func (b *builderImpl[T]) applyEnvironmentVariables(config *T) error {
	rv := reflect.ValueOf(config)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	for _, v := range b.properties {
		variable := v.getEnvVariable()
		if len(variable) > 0 {
			value, ok := os.LookupEnv(variable)
			if ok {
				rv.FieldByIndex(v.path).SetString(value)
			}
		}
	}
	return nil
}

func (b *builderImpl[T]) parseCliArguments(config *T, cliArgs []string, useFlags bool) error {
	if useFlags {
		return b.parseCliFlagArguments(config, cliArgs)
	}
	return b.parseCliNaiveArguments(config, cliArgs)
}

func (b *builderImpl[T]) parseCliFlagArguments(config *T, cliArgs []string) error {
	parser := newCliConfigFlag[T](b.opts.configFileArgs, b.properties)
	fileName, err := parser.parse(cliArgs)
	if err != nil {
		return err
	}
	if len(fileName) > 0 {
		err = b.Load(fileName)
		if err != nil {
			return err
		}
		*config = b.config
	}
	return parser.applyTo(config)
}

func (b *builderImpl[T]) parseCliNaiveArguments(config *T, cliArgs []string) error {
	return fmt.Errorf("not implemented")
}
