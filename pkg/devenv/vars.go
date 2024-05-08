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

func VariablesBuildArg(p *Profile) []Variable {
	vars := make([]Variable, 0, len(p.Services) * 3)
	for _, s := range p.Services {
		for arg, v := range s.BuildArgs {
			n := tmplutils.MustSprint(BuildArgNameTemplate, arg)
			vars = append(vars, Variable{Name: n, Value: v})
		}
	}
	sort.SliceStable(vars, func(i, j int) bool {
		return vars[i].Name < vars[j].Name
	})
	return vars
}
