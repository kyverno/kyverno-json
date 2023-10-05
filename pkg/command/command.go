package command

import (
	"strings"

	"github.com/spf13/cobra"
)

type Command struct {
	parent       *cobra.Command
	experimental bool
	description  []string
	websiteUrl   string
	examples     []Example
}

func new(parent *cobra.Command, experimental bool, options ...option) Command {
	cmd := Command{
		parent:       parent,
		experimental: experimental,
	}
	for _, opt := range options {
		if opt != nil {
			opt(&cmd)
		}
	}
	return cmd
}

func New(parent *cobra.Command, options ...option) Command {
	return new(parent, false, options...)
}

func NewExperimental(parent *cobra.Command, options ...option) Command {
	return new(parent, true, options...)
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
	if c.parent != nil {
		useLine = c.parent.UseLine() + " "
	}
	var lines []string
	for _, example := range c.examples {
		lines = append(lines, "  # "+example.title)
		lines = append(lines, "  "+useLine+example.command)
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}
