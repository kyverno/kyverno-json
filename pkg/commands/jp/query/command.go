package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	"github.com/kyverno/kyverno-json/pkg/command"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func Command(parents ...string) *cobra.Command {
	doc := command.New(
		command.WithParents(parents...),
		command.WithDescription("Provides a command-line interface to JMESPath, enhanced with Kyverno specific custom functions."),
		command.WithExample("Evaluate query", "query -i object.yaml 'request.object.metadata.name | truncate(@, `9`)'"),
		command.WithExample("Evaluate query", "query -i object.yaml -q query-file"),
		command.WithExample("Evaluate multiple queries", "query -i object.yaml -q query-file-1 -q query-file-2 'request.object.metadata.name | truncate(@, `9`)'"),
	)
	var compact, unquoted bool
	var input string
	var queries []string
	cmd := &cobra.Command{
		Use:          "query [-i input] [-q query|query]...",
		Short:        command.Description(doc, true),
		Long:         command.Description(doc, false),
		Example:      command.Examples(doc),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			queries, err := loadQueries(cmd, args, queries)
			if err != nil {
				return err
			}
			input, err := loadInput(cmd, input)
			if err != nil {
				return err
			}
			if len(queries) == 0 && input == nil {
				return errors.New("at least one query or input object is required")
			}
			if len(queries) == 0 {
				query, err := readQuery(cmd)
				if err != nil {
					return err
				}
				queries = append(queries, query)
			}
			if input == nil {
				i, err := readInput(cmd)
				if err != nil {
					return err
				}
				input = i
			}
			for _, query := range queries {
				result, err := evaluate(input, query)
				if err != nil {
					return err
				}
				if err := printResult(cmd, query, result, unquoted, compact); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&compact, "compact", "c", false, "Produce compact JSON output that omits non essential whitespace")
	cmd.Flags().BoolVarP(&unquoted, "unquoted", "u", false, "If the final result is a string, it will be printed without quotes")
	cmd.Flags().StringSliceVarP(&queries, "query", "q", nil, "Read JMESPath expression from the specified file")
	cmd.Flags().StringVarP(&input, "input", "i", "", "Read input from a JSON or YAML file instead of stdin")
	return cmd
}

func readFile(reader io.Reader) ([]byte, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func loadFile(cmd *cobra.Command, file string) ([]byte, error) {
	reader, err := os.Open(filepath.Clean(file))
	if err != nil {
		return nil, fmt.Errorf("failed open file %s: %v", file, err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Error closing file: %s\n", err)
		}
	}()
	content, err := readFile(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", file, err)
	}
	return content, nil
}

func readQuery(cmd *cobra.Command) (string, error) {
	fmt.Fprintln(cmd.OutOrStdout(), "Reading from terminal input.")
	fmt.Fprintln(cmd.OutOrStdout(), "Enter a jmespath expression and hit Ctrl+D.")
	data, err := readFile(cmd.InOrStdin())
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func loadQueries(cmd *cobra.Command, args []string, files []string) ([]string, error) {
	var queries []string
	queries = append(queries, args...)
	for _, file := range files {
		query, err := loadFile(cmd, file)
		if err != nil {
			return nil, err
		}
		queries = append(queries, string(query))
	}
	return queries, nil
}

func readInput(cmd *cobra.Command) (any, error) {
	fmt.Fprintln(cmd.OutOrStdout(), "Reading from terminal input.")
	fmt.Fprintln(cmd.OutOrStdout(), "Enter input object and hit Ctrl+D.")
	data, err := readFile(cmd.InOrStdin())
	if err != nil {
		return nil, err
	}
	var input any
	if err := yaml.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("error parsing input json: %w", err)
	}
	return input, nil
}

func loadInput(cmd *cobra.Command, file string) (any, error) {
	if file == "" {
		return nil, nil
	}
	data, err := loadFile(cmd, file)
	if err != nil {
		return nil, err
	}
	var input any
	if err := yaml.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("error parsing input json: %w", err)
	}
	return input, nil
}

func evaluate(input any, query string) (any, error) {
	result, err := compilers.Execute(query, input, nil, compilers.DefaultCompiler.Jp)
	if err != nil {
		if syntaxError, ok := err.(parsing.SyntaxError); ok {
			return nil, fmt.Errorf("%s\n%s", syntaxError, syntaxError.HighlightLocation())
		}
		return nil, fmt.Errorf("error evaluating JMESPath expression: %w", err)
	}
	return result, nil
}

func printResult(cmd *cobra.Command, query string, result any, unquoted bool, compact bool) error {
	converted, isString := result.(string)
	fmt.Fprintln(cmd.OutOrStdout(), "#", query)
	if unquoted && isString {
		fmt.Fprintln(cmd.OutOrStdout(), converted)
	} else {
		var toJSON []byte
		var err error
		if compact {
			toJSON, err = json.Marshal(result)
		} else {
			toJSON, err = json.MarshalIndent(result, "", "  ")
		}
		if err != nil {
			return fmt.Errorf("error marshalling result to JSON: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(toJSON))
	}
	return nil
}
