package pawn

type ModType int

const (
	MOD_GOOD ModType = iota
	MOD_BAD
)

type Mod struct {
	Name string
	Desc string
	Type ModType
}
