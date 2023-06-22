package i3config

import (
	"fmt"
	"os"
	"regexp"
)

type Color string

func HexColor(code string) Color {
	return Color("#" + code)
}

func (c Color) Generate() string {
	if !c.Valid() {
		fmt.Printf("invalid color \"%s\"\n", c)
		os.Exit(1)
	}
	return string(c)
}
func (c Color) GenerateYabai() string {
	if !c.Valid() {
		fmt.Printf("invalid color \"%s\"\n", c)
		os.Exit(1)
	}
	return "0xff" + string(c)[1:]
}

var colorRegExp = regexp.MustCompile("^#[0-9a-fA-F]{6}$")

func (c Color) Valid() bool {
	return colorRegExp.Match([]byte(c))
}

type ColorClass struct {
	Border      Color
	Background  Color
	Text        Color
	Indicator   Color
	ChildBorder Color
}

func ConstantColorClass(color Color) ColorClass {
	return ColorClass{
		Border:      color,
		Background:  color,
		Text:        color,
		Indicator:   color,
		ChildBorder: color,
	}
}

func (c *ColorClass) Generate() string {
	return fmt.Sprintf("%s %s %s %s %s",
		c.Border.Generate(),
		c.Background.Generate(),
		c.Text.Generate(),
		c.Indicator.Generate(),
		c.ChildBorder.Generate(),
	)
}

type ColorConfig struct {
	Focused         ColorClass
	FocusedInactive ColorClass
	Unfocused       ColorClass
	Urgent          ColorClass
	Placeholder     ColorClass
	Background      Color
}

func (c *ColorConfig) Generate() string {
	return "client.focused " + c.Focused.Generate() + "\n" +
		"client.focused_inactive " + c.FocusedInactive.Generate() + "\n" +
		"client.unfocused " + c.Unfocused.Generate() + "\n" +
		"client.urgent " + c.Urgent.Generate() + "\n" +
		"client.placeholder " + c.Placeholder.Generate() + "\n" +
		"client.background " + c.Background.Generate()

}

func (c *ColorConfig) GenerateYabai() string {
	return "yabai -m config active_window_border_color " + c.Focused.Border.GenerateYabai() + "\n" +
		"yabai -m config normal_window_border_color " + c.Unfocused.Border.GenerateYabai() + "\n"
}

func (c *Config) Colors(cc *ColorConfig) {
	c.AddLine(cc)
}
