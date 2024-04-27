package allymod

type ModType int

const (
	Attack ModType = iota
	Buff
	Debuff
)

type ModTarget int

const (
	Ally ModTarget = iota
	Enemy
	Self
)

type Mod struct {
	Name   string
	Desc   string
	Type   ModType
	Target ModTarget
	Range  float64
}

func (m *Mod) IsGood() bool {
	switch m.Type {
	case Attack, Debuff:
		return m.Target == Enemy
	case Buff:
		return m.Target == Ally
	default:
		return false
	}
}
