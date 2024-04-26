package allymod

type ModType int

const (
	GOOD ModType = iota
	BAD
)

type Mod struct {
	Name string
	Desc string
	Type ModType
}
