package devenv

type Hooks map[HookPhase][]Hook

//// RawMap a template friendly raw map
//func (h Hooks) RawMap() map[string][]Hook {
//	m := map[string][]Hook{}
//	for k, v := range h {
//		m[string(k)] = v
//	}
//	return m
//}

func (h Hooks) Phase(phaseStr HookPhase) []Hook {
	hooks, _ := h[phaseStr]
	return hooks
}

type Hook struct {
	Name  string
	Phase HookPhase
	Type  HookType
	Value interface{}
}

const (
	PhasePreStart  HookPhase = "pre-start"
	PhasePostStart HookPhase = "post-start"
	PhasePreStop   HookPhase = "pre-stop"
	PhasePostStop  HookPhase = "post-stop"
)

type HookPhase string

const (
	TypeScript    HookType = "script"
	TypeContainer HookType = "container"
)

type HookType string
