package dto

import "strings"

var AvailableCommands = []string{
	"/stock",
}

type Command string

func (c *Command) GetValue() string {
	parts := strings.SplitN(string(*c), "=", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func (c *Command) ToString() string {
	return string(*c)
}
