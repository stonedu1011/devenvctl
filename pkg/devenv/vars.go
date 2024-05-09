package devenv

import (
	"fmt"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
	"sort"
)

type Variable struct {
	Name  string
	Value string
}

func (v Variable) String() string {
	return fmt.Sprintf(`%s="%s"`, v.Name, v.Value)
}

func Variables(p *Profile) []Variable {
	vars := make([]Variable, 0, len(p.Services)*5+2)
	vars = append(vars, VariablesGlobal(p)...)
	vars = append(vars, VariablesService(p)...)
	vars = append(vars, VariablesBuildArg(p)...)
	return vars
}

func VariablesBuildArg(p *Profile) []Variable {
	vars := make([]Variable, 0, len(p.Services)*3)
	for _, s := range p.Services {
		for arg, v := range s.BuildArgs {
			n := tmplutils.MustSprint(TemplateBuildArgName, arg)
			vars = append(vars, Variable{Name: n, Value: v})
		}
	}
	sort.SliceStable(vars, func(i, j int) bool {
		return vars[i].Name < vars[j].Name
	})
	return vars
}

func VariablesService(p *Profile) []Variable {
	vars := make([]Variable, 0, len(p.Services)*2)
	for _, s := range p.Services {
		vars = append(vars,
			Variable{Name: tmplutils.MustSprint(TemplateServiceImage, s), Value: s.Image},
			Variable{Name: tmplutils.MustSprint(TemplateServiceContainer, s), Value: s.ContainerName()},
		)
	}
	sort.SliceStable(vars, func(i, j int) bool {
		return vars[i].Name < vars[j].Name
	})
	return vars
}

func VariablesGlobal(p *Profile) []Variable {
	vars := []Variable{
		{Name: VarProjectName, Value: p.Name},
		{Name: VarLocalDataPath, Value: p.LocalDataDir()},
		{Name: VarProjectResource, Value: p.ResourceDir()},
	}
	sort.SliceStable(vars, func(i, j int) bool {
		return vars[i].Name < vars[j].Name
	})
	return vars
}
