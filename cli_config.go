package goappconfig

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type cliConfigFlag[T any] struct {
	flagSet  *flag.FlagSet
	fileArgs []string
	props    []*propertyImpl
}

func newCliConfigFlag[T any](configFileArgs []string, properties []*propertyImpl) *cliConfigFlag[T] {
	var set *flag.FlagSet
	if len(os.Args) > 0 {
		set = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	} else {
		set = flag.NewFlagSet("", flag.ContinueOnError)
	}
	if len(configFileArgs) > 0 {
		for _, arg := range configFileArgs {
			set.String(arg, "", "configuration file name")
		}
	}
	for _, v := range properties {
		for _, arg := range v.getCliArgumentNames() {
			switch v.kind {
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				set.Int(arg, 0, v.name)
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				set.Int(arg, 0, v.name)
			case reflect.Bool:
				set.Bool(arg, false, v.name)
			case reflect.Float32, reflect.Float64:
				set.Float64(arg, 0, v.name)
			default:
				set.String(arg, "", v.name)
			}
		}
	}
	return &cliConfigFlag[T]{
		flagSet:  set,
		fileArgs: configFileArgs,
		props:    properties,
	}
}

func (c *cliConfigFlag[T]) parse(cliArgs []string) (string, error) {
	err := c.flagSet.Parse(cliArgs)
	if err != nil {
		return "", err
	}
	if len(c.fileArgs) == 0 {
		return "", nil
	}
	for _, arg := range c.fileArgs {
		f := c.flagSet.Lookup(arg)
		if f != nil && f.Value.String() != "" {
			return f.Value.String(), nil
		}
	}
	return "", nil
}

func (c *cliConfigFlag[T]) applyTo(config *T) error {
	rv := reflect.ValueOf(config)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	for _, v := range c.props {
		err := c.processProperty(rv, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *cliConfigFlag[T]) processProperty(rv reflect.Value, property *propertyImpl) error {
	for _, v := range property.getCliArgumentNames() {
		ok, err := c.processPropertyArg(rv, property, v)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
	return nil
}

func (c *cliConfigFlag[T]) processPropertyArg(rv reflect.Value, property *propertyImpl, arg string) (bool, error) {
	f := c.flagSet.Lookup(arg)
	if f == nil || f.Value.String() == "" {
		return false, nil
	}
	v := f.Value.String()
	switch property.kind {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return false, fmt.Errorf("argument %s must be an integer, '%s': %s", arg, v, errors.Unwrap(err))
		}
		rv.FieldByIndex(property.path).SetInt(val)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return false, fmt.Errorf("argument %s must be a unsigned integer, '%s': %s", arg, v, errors.Unwrap(err))
		}
		rv.FieldByIndex(property.path).SetUint(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(v)
		if err != nil {
			return false, fmt.Errorf("argument %s must be a boolean, '%s': %s", arg, v, errors.Unwrap(err))
		}
		rv.FieldByIndex(property.path).SetBool(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false, fmt.Errorf("argument %s must be float number, '%s': %s", arg, v, errors.Unwrap(err))
		}
		rv.FieldByIndex(property.path).SetFloat(val)
	default:
		rv.FieldByIndex(property.path).SetString(v)
	}
	return true, nil
}
