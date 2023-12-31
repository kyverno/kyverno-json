package parse

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	"github.com/kyverno/kyverno-json/pkg/command"
	"github.com/spf13/cobra"
)

func Command(parents ...string) *cobra.Command {
	doc := command.New(
		command.WithParents(parents...),
		command.WithDescription("Parses jmespath expression and prints corresponding AST."),
		command.WithExample("Parse expression", "parse 'request.object.metadata.name | truncate(@, `9`)'"),
		command.WithExample("Parse expression from a file", "parse -f my-file"),
		command.WithExample("Parse expression from stdin", "parse"),
		command.WithExample("Parse multiple expressionxs", "parse -f my-file1 -f my-file-2 'request.object.metadata.name | truncate(@, `9`)'"),
	)
	var files []string
	cmd := &cobra.Command{
		Use:          "parse [-f file|expression]...",
		Short:        command.Description(doc, true),
		Long:         command.Description(doc, false),
		Example:      command.Examples(doc),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			expressions, err := loadExpressions(cmd, args, files)
			if err != nil {
				return err
			}
			for _, expression := range expressions {
				if err := printAst(cmd.OutOrStdout(), expression); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVarP(&files, "file", "f", nil, "Read input from a JSON or YAML file instead of stdin")
	return cmd
}

func readFile(reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func loadFile(cmd *cobra.Command, file string) (string, error) {
	reader, err := os.Open(filepath.Clean(file))
	if err != nil {
		return "", fmt.Errorf("failed open file %s: %v", file, err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Error closing file: %s\n", err)
		}
	}()
	content, err := readFile(reader)
	if err != nil {
		return "", fmt.Errorf("failed read file %s: %v", file, err)
	}
	return content, nil
}

func loadExpressions(cmd *cobra.Command, args []string, files []string) ([]string, error) {
	var expressions []string
	expressions = append(expressions, args...)
	for _, file := range files {
		expression, err := loadFile(cmd, file)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expression)
	}
	if len(expressions) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Reading from terminal input.")
		fmt.Fprintln(cmd.OutOrStdout(), "Enter a jmespath expression and hit Ctrl+D.")
		data, err := readFile(cmd.InOrStdin())
		if err != nil {
			return nil, fmt.Errorf("failed to read file STDIN: %v", err)
		}
		expressions = append(expressions, data)
	}
	return expressions, nil
}

// The following function has been adapted from
// https://github.com/jmespath/jp/blob/54882e03bd277fc4475a677fab1d35eaa478b839/jp.go
func printAst(out io.Writer, expression string) error {
	parser := parsing.NewParser()
	parsed, err := parser.Parse(expression)
	if err != nil {
		if syntaxError, ok := err.(parsing.SyntaxError); ok {
			return fmt.Errorf("%w\n%s", syntaxError, syntaxError.HighlightLocation())
		}
		return err
	}
	fmt.Fprintln(out, "#", expression)
	fmt.Fprintln(out, parsed)
	return nil
}
