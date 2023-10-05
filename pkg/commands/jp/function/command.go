package function

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"slices"
	"strings"

	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/kyverno/kyverno-json/pkg/command"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
)

func Command(parent *cobra.Command) *cobra.Command {
	doc := command.New(
		parent,
		command.WithDescription("Provides function informations."),
		command.WithExample("List functions", "function"),
		command.WithExample("Get function infos", "function truncate"),
	)
	return &cobra.Command{
		Use:          "function [function_name]...",
		Short:        command.Description(doc, true),
		Long:         command.Description(doc, false),
		Example:      command.Examples(doc),
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			printFunctions(cmd.OutOrStdout(), args...)
		},
	}
}

func functionString(f jpfunctions.FunctionEntry) string {
	if f.Name == "" {
		return ""
	}
	var args []string
	for _, a := range f.Arguments {
		var aTypes []string
		for _, t := range a.Types {
			aTypes = append(aTypes, string(t))
		}
		args = append(args, strings.Join(aTypes, "|"))
	}
	// var returnArgs []string
	// for _, ra := range f.ReturnType {
	// 	returnArgs = append(returnArgs, string(ra))
	// }
	// output := fmt.Sprintf("%s(%s) %s", f.Name, strings.Join(args, ", "), strings.Join(returnArgs, ","))
	output := fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ", "))
	// if f.Note != "" {
	// 	output += fmt.Sprintf(" (%s)", f.Note)
	// }
	return output
}

func printFunctions(out io.Writer, names ...string) {
	funcs := template.GetFunctions(context.Background())
	slices.SortFunc(funcs, func(a, b jpfunctions.FunctionEntry) int {
		return cmp.Compare(functionString(a), functionString(b))
	})
	namesSet := sets.New(names...)
	for _, function := range funcs {
		if len(namesSet) == 0 || namesSet.Has(function.Name) {
			// note := function.Note
			// function.Note = ""
			fmt.Fprintln(out, "Name:", function.Name)
			fmt.Fprintln(out, "  Signature:", functionString(function))
			// if note != "" {
			// 	fmt.Fprintln(out, "  Note:     ", note)
			// }
			fmt.Fprintln(out)
		}
	}
}
