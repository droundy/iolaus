package color

import (
	"./git"
)

const (
	resetColor = "\x1B[00m"
)

var (
	Meta = &Color{ "color.diff.meta", "blue bold" }
	Old = &Color{ "color.diff.old", "red bold" }
	New = &Color{ "color.diff.new", "green bold" }
	Plain = &Color{ "", "" }
)

type Color struct {
	name, def string
}

func String(s string, c *Color) string {
	return Get(*c)+s+resetColor
}

func Get(c Color) string {
	o,_ := git.Read("config", "--get-color", c.name, c.def)
	return string(o)
}

