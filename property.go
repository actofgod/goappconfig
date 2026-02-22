package goappconfig

import (
	"reflect"
	"strings"
)

var tagList = []string{
	"json",
	"yaml",
	"xml",
	"ini",
	"env",
	"arg",
	"short",
}

type propertyImpl struct {
	parent *propertyImpl
	path   []int
	name   string
	tags   map[string]string
}

func newProperty(parent *propertyImpl, field reflect.StructField) *propertyImpl {
	tags := make(map[string]string)
	for _, v := range tagList {
		value, ok := field.Tag.Lookup(v)
		if ok {
			tag := strings.Split(value, ",")[0]
			if len(tag) > 0 {
				tags[v] = tag
			}
		}
	}
	var path []int
	if parent != nil {
		path = append(path, parent.path...)
	}
	path = append(path, field.Index...)
	return &propertyImpl{
		parent: parent,
		path:   path,
		name:   field.Name,
		tags:   tags,
	}
}

func (p *propertyImpl) getCliArgumentNames() []string {
	var args []string = nil
	name, ok := p.tags["arg"]
	if ok {
		args = append(args, name)
	}
	name, ok = p.tags["short"]
	if ok {
		args = append(args, name)
	}
	if len(args) == 0 {
		args = append(args, "p:"+p.name)
	}
	return args
}

func (p *propertyImpl) getEnvVariable() string {
	name, ok := p.tags["env"]
	if ok {
		return name
	}
	if p.parent != nil {
		return p.parent.getEnvVariable() + "_" + toEnvVariable(p.name)
	}
	return toEnvVariable(p.name)
}

func toEnvVariable(name string) string {
	var b strings.Builder
	var isPreviousUpper = true
	var isPreviousNumber = false
	var length = 0
	for _, c := range name {
		if c > 127 {
			continue
		}
		if c >= 'A' && c <= 'Z' {
			if !isPreviousUpper {
				length++
			}
			isPreviousUpper = true
		} else {
			isPreviousUpper = false
		}
		if c >= '0' && c <= '9' {
			if !isPreviousNumber {
				length++
			}
			isPreviousNumber = true
		} else {
			isPreviousNumber = false
		}
		length++
	}
	if length == 0 {
		return ""
	}
	b.Grow(length)
	isPreviousUpper = true
	isPreviousNumber = false
	for _, c := range name {
		if c > 127 {
			continue
		}
		if c >= '0' && c <= '9' {
			if !isPreviousNumber {
				b.WriteRune('_')
			}
			isPreviousNumber = true
		} else {
			isPreviousNumber = false
		}
		if c >= 'A' && c <= 'Z' {
			if !isPreviousUpper {
				b.WriteRune('_')
			}
			isPreviousUpper = true
		} else {
			isPreviousUpper = false
		}
		if c >= 'a' && c <= 'z' {
			c = c - 'a' + 'A'
		}
		b.WriteRune(c)
	}
	return b.String()
}
