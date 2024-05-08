package devenv

import (
	"fmt"
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/stonedu1011/devenvctl/pkg/utils"
)

type ProfileV1 struct {
	ProfileMetadata
	Services  []ServiceV1 `json:"services"`
	PreStart  []string    `json:"pre_start"`
	PostStart []string    `json:"post_start"`
	PreStop   []string    `json:"pre_stop"`
	PostStop  []string    `json:"post_stop"`
}

func (p *ProfileV1) ToProfile() *Profile {
	ret := Profile{
		ProfileMetadata: p.ProfileMetadata,
		Services:        make([]Service, len(p.Services)),
		Hooks: Hooks{
			PhasePreStart:  utils.ConvertSlice(p.PreStart, p.hookConverter(PhasePreStart)),
			PhasePostStart: utils.ConvertSlice(p.PostStart, p.hookConverter(PhasePostStart)),
			PhasePreStop:   utils.ConvertSlice(p.PreStop, p.hookConverter(PhasePreStop)),
			PhasePostStop:  utils.ConvertSlice(p.PostStop, p.hookConverter(PhasePostStop)),
		},
	}
	for i := range p.Services {
		ret.Services[i] = Service{
			Name:           p.Services[i].Name,
			DisplayName:    p.Services[i].DisplayName,
			DisplayVersion: p.Services[i].DisplayVersion,
			Image:          p.Services[i].ImageName,
			Mounts:         p.Services[i].Mounts,
			BuildArgs:      p.Services[i].BuildArgs,
		}
	}
	return &ret
}

func (p *ProfileV1) hookConverter(phase HookPhase) func(string) Hook {
	return func(value string) Hook {
		hook := Hook{
			Name:  value,
			Phase: phase,
			Type:  TypeContainer,
			Value: value,
		}
		switch phase {
		case PhasePreStart, PhasePostStop:
			hook.Type = TypeScript
		}
		return hook
	}
}

type ServiceV1 struct {
	Name           string            `json:"service"`
	DisplayName    string            `json:"display_name"`
	DisplayVersion string            `json:"display_version"`
	ImageName      string            `json:"image"`
	Mounts         []string          `json:"mounts"`
	BuildArgs      map[string]string `json:"build_args"`
}

func LoadProfileV1(meta *ProfileMetadata) (*ProfileV1, error) {
	f, e := meta.FS.Open(meta.Path)
	if e != nil {
		return nil, fmt.Errorf(`unable to open profile definition file "%s": %v`, meta.DisplayPath, e)
	}
	defer func() { _ = f.Close() }()
	p := &ProfileV1{
		ProfileMetadata: *meta,
	}
	if e := cmdutils.BindYaml(f, p); e != nil {
		return nil, fmt.Errorf(`unable to parse profile definition file "%s" as v1 format: %v`, meta.Path, e)
	}
	return p, nil
}