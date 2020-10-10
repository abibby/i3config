package i3config

import (
	"fmt"
	"os"
	"os/exec"
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
func (c *Config) Set(variable, value string) {
	c.raw(fmt.Sprintf("set %s %s", variable, value))
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

func (c *Config) Run() {
	if len(os.Args) > 1 && os.Args[1] == "func" {
		err := c.funcs[os.Args[2]]()
		if err != nil {
			exec.Command("notify-send", err.Error()).Run()
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	} else {
		c.chords.apply(c)
		fmt.Print(c.Generate())
	}
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
