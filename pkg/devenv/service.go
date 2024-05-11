package devenv

import (
	"github.com/stonedu1011/devenvctl/pkg/utils"
)

type Service struct {
	Name           string
	DisplayName    string
	DisplayVersion string
	Image          string
	Mounts         []string
	BuildArgs      map[string]string
	owner          *Profile
}

func (s Service) ContainerName() string {
	return utils.SnakeCase(s.owner.Name) + "-" + utils.SnakeCase(s.Name)
}
