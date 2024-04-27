package allymod

type ModType int

const (
	MeleeAttack ModType = iota
	RangedAttack
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
}

func (m *Mod) IsGood() bool {
	switch m.Type {
	case MeleeAttack, RangedAttack, Debuff:
		return m.Target == Enemy
	case Buff:
		return m.Target == Ally
	default:
		return false
	}
}
