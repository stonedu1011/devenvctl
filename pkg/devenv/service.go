package devenv

type Service struct {
	Name           string
	DisplayName    string
	DisplayVersion string
	Image          string
	Mounts         []string
	BuildArgs      map[string]string
}
