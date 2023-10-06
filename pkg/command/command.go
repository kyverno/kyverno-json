package command

import (
	"strings"
)

type Command struct {
	parents      []string
	experimental bool
	description  []string
	websiteUrl   string
	examples     []Example
}

func New(options ...option) Command {
	var cmd Command
	for _, opt := range options {
		if opt != nil {
			opt(&cmd)
		}
	}
	return cmd
}

func Description(c Command, short bool) string {
	if len(c.description) == 0 {
		return ""
	}
	var lines []string
	lines = append(lines, c.description[0])
	if !short {
		lines = append(lines, "")
		lines = append(lines, c.description[1:]...)
		if c.experimental {
			lines = append(lines, "", "NOTE: This is an experimental command.")
		}
		if c.websiteUrl != "" {
			lines = append(lines, "", "For more information visit "+c.websiteUrl)
		}
	}
	return strings.Join(lines, "\n")
}

func Examples(c Command) string {
	if len(c.examples) == 0 {
		return ""
	}
	var useLine string
	if len(c.parents) != 0 {
		useLine = strings.Join(c.parents, " ") + " "
	}
	var lines []string
	for _, example := range c.examples {
		lines = append(lines, "  # "+example.title)
		lines = append(lines, "  "+useLine+example.command)
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}
