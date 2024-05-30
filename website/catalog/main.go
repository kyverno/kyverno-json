package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/catalog"
	"github.com/kyverno/kyverno-json/pkg/policy"
	fileinfo "github.com/kyverno/pkg/ext/file-info"
	"sigs.k8s.io/yaml"
)

const (
	path = "./catalog"
)

type pol struct {
	Path   string
	Policy *v1alpha1.ValidatingPolicy
}

func (p pol) TargetPath() string {
	base, err := filepath.Rel(path, p.Path)
	if err != nil {
		panic(err)
	}
	target := filepath.Join("website/docs/catalog/policies/", base)
	target = strings.TrimSuffix(target, filepath.Ext(target)) + ".md"
	return target
}

func (p pol) NavPath() string {
	base, err := filepath.Rel("website/docs", p.TargetPath())
	if err != nil {
		panic(err)
	}
	return base
}

func (p pol) Generate() error {
	if err := os.MkdirAll(filepath.Dir(p.TargetPath()), os.ModePerm); err != nil {
		return err
	}
	template, err := template.ParseFiles("./website/policy.gotmpl")
	if err != nil {
		return err
	}
	policy, err := os.Create(p.TargetPath())
	if err != nil {
		return err
	}
	defer policy.Close()
	if err := template.Execute(policy, p); err != nil {
		return err
	}
	return nil
}

func (p pol) Title() string {
	title := p.Policy.Annotations[catalog.AnnotationPolicyTitle]
	if title != "" {
		return title
	}
	base := filepath.Base(p.Path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func (p pol) Description() string {
	description := p.Policy.Annotations[catalog.AnnotationPolicyDescription]
	if description != "" {
		return description
	}
	return "None"
}

func (p pol) Manifest() string {
	bytes, err := yaml.Marshal(p.Policy)
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(string(bytes))
}

func (p pol) Tags() []string {
	var tags []string
	for k := range p.Policy.Labels {
		if strings.HasSuffix(k, catalog.TagsLabelSuffix) {
			tag := strings.TrimSuffix(k, catalog.TagsLabelSuffix)
			parts := strings.Split(tag, ".")
			slices.Reverse(parts)
			for i := 1; i <= len(parts); i++ {
				tags = append(tags, strings.Join(parts[:i], "/"))
			}
		}
	}
	return tags
}

func main() {
	var files []string
	err := filepath.Walk(path, func(file string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fileinfo.IsYaml(info) {
			files = append(files, file)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	var pols []pol
	for _, file := range files {
		policies, err := policy.Load(file)
		if err != nil {
			panic(err)
		}
		for _, policy := range policies {
			pols = append(pols, pol{
				Path:   file,
				Policy: policy,
			})
		}
		if err := os.RemoveAll("website/docs/catalog/policies"); err != nil {
			panic(err)
		}
		for _, pol := range pols {
			err := pol.Generate()
			if err != nil {
				panic(err)
			}
		}
		template, err := template.ParseFiles("./website/nav.gotmpl")
		if err != nil {
			panic(err)
		}
		mkdocs, err := os.Create("./website/mkdocs.yaml")
		if err != nil {
			panic(err)
		}
		defer mkdocs.Close()
		if err := template.Execute(mkdocs, map[string]any{
			"Policies": pols,
		}); err != nil {
			panic(err)
		}
	}
}
