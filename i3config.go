package i3config

import (
	"fmt"
	"reflect"
)

type Config struct {
	path   string
	lines  []Generator
	chords Chords

	subConfig bool
	funcs     map[string]func() error
}

type Generator interface {
	Generate() string
}
type GenerateYabaier interface {
	GenerateYabai() string
}

func New(path string) *Config {
	return &Config{
		path:      path,
		lines:     []Generator{},
		chords:    Chords{},
		subConfig: false,
		funcs:     map[string]func() error{},
	}
}

func (c *Config) newSubConfig() *Config {
	sc := New(c.path)
	sc.subConfig = true
	return sc
}

type Variable struct {
	Name  string
	Value string
}

func (v *Variable) Generate() string {
	return v.Name + " " + v.Value
}
func (c *Config) Set(variable, value string) {
	c.AddLine(&Variable{
		Name:  variable,
		Value: value,
	})
}
func (c *Config) AddLine(g Generator) {
	c.lines = append(c.lines, g)
}

func (c *Config) Generate() string {
	src := ""
	for _, b := range c.lines {
		src += b.Generate() + "\n"
	}
	return src
}

func EachKey(v interface{}, cb func(key, value string)) {

	modelReflect := reflect.ValueOf(v)

	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	modelRefType := modelReflect.Type()
	fieldsCount := modelReflect.NumField()

	for i := 0; i < fieldsCount; i++ {
		field := modelReflect.Field(i)

		value := field.Interface()
		key := modelRefType.Field(i).Tag.Get("i3")
		if gen, ok := value.(Generator); ok {
			cb(key, gen.Generate())
		} else {
			cb(key, fmt.Sprint(value))
		}
	}
}
