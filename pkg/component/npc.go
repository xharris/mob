package component

type ModType int

const (
	Attack ModType = iota
	Buff
	Debuff
)

type ModTarget int

const (
	TargetAlly ModTarget = iota
	TargetEnemy
	TargetSelf
)

type Mod struct {
	Name   string
	Desc   string
	Type   ModType
	Target ModTarget
	Range  NPCRange
}

func (m *Mod) IsGood() bool {
	switch m.Type {
	case Attack, Debuff:
		return m.Target == TargetEnemy
	case Buff:
		return m.Target == TargetAlly || m.Target == TargetSelf
	default:
		return false
	}
}

type NPCType int

const (
	Ally NPCType = iota
	Enemy
)

type NPCStrategy int

const (
	Closest NPCStrategy = iota
	Farthest
	HighAttack
	LowAttack
	HighHealth
	LowHealth
	Random
)

type NPCRange float64

const (
	Melee  NPCRange = 30
	Ranged NPCRange = 90
)

type NPC struct {
	Mods            []Mod
	Name            string
	Type            NPCType
	HealthRemaining int
	HealthTotal     int
	// targeting strategy
	Strategy NPCStrategy
	// entity id
	Target int
	// -1 for infinite range
	AggroRange float64
}

func NewNPC(options ...NPCOption) (n NPC) {
	n.Target = EntityNone
	n.AggroRange = -1
	for _, opt := range options {
		opt(&n)
	}
	return
}

type NPCOption func(*NPC)

func WithMod(m Mod) NPCOption {
	return func(n *NPC) {
		n.AddMod(m)
	}
}

func WithMods(mods []Mod) NPCOption {
	return func(n *NPC) {
		for _, m := range mods {
			n.AddMod(m)
		}
	}
}

func WithName(name string) NPCOption {
	return func(n *NPC) {
		n.Name = name
	}
}

func WithType(t NPCType) NPCOption {
	return func(n *NPC) {
		n.Type = t
	}
}

func WithHealth(total int, remaining int) NPCOption {
	return func(n *NPC) {
		n.HealthTotal = total
		n.HealthRemaining = remaining
	}
}

func (n *NPC) AddMod(m Mod) {
	var mods []Mod
	found := false
	for _, am := range n.Mods {
		if am.Name == m.Name {
			found = true
		} else {
			mods = append(mods, am)
		}
	}
	if !found {
		mods = append(mods, m)
	}
	n.Mods = mods
}

func (n *NPC) RemoveMod(m Mod) {
	var mods []Mod
	for _, am := range n.Mods {
		if am.Name != m.Name {
			mods = append(mods, m)
		}
	}
	n.Mods = mods
}

func (n *NPC) GetAttackRange() (r float64) {
	for _, m := range n.Mods {
		if float64(m.Range) > r {
			r = float64(m.Range)
		}
	}
	return
}
