package docs

import (
	"log"

	"github.com/kyverno/kyverno-json/pkg/command"
	"github.com/spf13/cobra"
)

func Command(parent *cobra.Command) *cobra.Command {
	var options options
	doc := command.New(
		parent,
		command.WithDescription(
			"Generates reference documentation.",
			"The docs command generates CLI reference documentation.",
			"It can be used to generate simple markdown files or markdown to be used for the website.",
		),
		command.WithExample(
			"Generate simple markdown documentation",
			"docs -o . --autogenTag=false",
		),
		command.WithExample(
			"Generate website documentation",
			"docs -o . --website",
		),
	)
	cmd := &cobra.Command{
		Use:          "docs",
		Short:        command.Description(doc, true),
		Long:         command.Description(doc, false),
		Example:      command.Examples(doc),
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root := parent
			if root != nil {
				for {
					if !root.HasParent() {
						break
					}
					root = root.Parent()
				}
			}
			if err := options.validate(root); err != nil {
				return err
			}
			return options.execute(root)
		},
	}
	cmd.Flags().StringVarP(&options.path, "output", "o", ".", "Output path")
	cmd.Flags().BoolVar(&options.website, "website", false, "Website version")
	cmd.Flags().BoolVar(&options.autogenTag, "autogenTag", true, "Determines if the generated docs should contain a timestamp")
	if err := cmd.MarkFlagDirname("output"); err != nil {
		log.Println("WARNING", err)
	}
	if err := cmd.MarkFlagRequired("output"); err != nil {
		log.Println("WARNING", err)
	}
	return cmd
}
