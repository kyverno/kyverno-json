package main

import (
	"fmt"
	"os"

	"github.com/kyverno/kyverno-json/pkg/policy"
	"sigs.k8s.io/yaml"
)

func main() {
	policies, err := policy.Load("./catalog")
	if err != nil {
		panic(err)
	}
	pages, err := os.Create("website/docs/catalog/policies/.pages")
	if err != nil {
		panic(err)
	}
	defer pages.Close()
	fmt.Fprintln(pages, "nav:")
	fmt.Fprintln(pages, "- All:")
	for i, policy := range policies {
		f, err := os.Create("website/docs/catalog/policies/" + fmt.Sprintf("%d.md", i))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		bytes, err := yaml.Marshal(policy)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(f, "# "+policy.Annotations["title.catalog.kyverno.io"])
		fmt.Fprintln(f)
		fmt.Fprintln(f, policy.Annotations["description.catalog.kyverno.io"])
		fmt.Fprintln(f)
		fmt.Fprintln(f, "## Manifest")
		fmt.Fprintln(f)
		fmt.Fprintln(f, "```yaml")
		if _, err := f.Write(bytes); err != nil {
			panic(err)
		}
		fmt.Fprintln(f, "```")
		fmt.Fprintln(pages, "  - "+fmt.Sprintf("%d.md", i))
	}
}
