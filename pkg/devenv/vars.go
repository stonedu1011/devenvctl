package devenv

import (
	"fmt"
	"github.com/stonedu1011/devenvctl/pkg/utils"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
	"sort"
)

type Variable struct {
	Name  string
	Value string
}

func (v Variable) String() string {
	return fmt.Sprintf(`%s=%s`, v.Name, v.Value)
}

type Variables struct {
	*utils.OrderedMap[string, Variable]
}

func (v Variables) Add(vars ...Variable) {
	for i := range vars {
		v.OrderedMap.Set(vars[i].Name, vars[i])
	}
}

func (v Variables) List() []Variable {
	vars := make([]Variable, 0, v.Len())
	for _, k := range v.Keys() {
		vars = append(vars, v.Get(k))
	}
	return vars
}

func (v Variables) KVMap() map[string]string {
	vars := map[string]string{}
	for _, k := range v.Keys() {
		entry := v.Get(k)
		vars[entry.Name] = entry.Value
	}
	return vars
}

func NewVariablesWithProfile(p *Profile) Variables {
	vars := Variables{
		OrderedMap: utils.NewOrderedMapWithCap[string, Variable](len(p.Services)*5+5),
	}
	vars.Add(ResolveServiceVars(p)...)
	vars.Add(ResolveBuildArgs(p)...)
	vars.Add(ResolveGlobalVars(p)...)
	return vars
}

func ResolveBuildArgs(p *Profile) []Variable {
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

func ResolveServiceVars(p *Profile) []Variable {
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

func ResolveGlobalVars(p *Profile) []Variable {
	vars := []Variable{
		{Name: VarProjectName, Value: p.Name},
	}
	sort.SliceStable(vars, func(i, j int) bool {
		return vars[i].Name < vars[j].Name
	})
	return vars
}
