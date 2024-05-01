package component

type Follow struct {
	// entity id
	Target        int
	Radius        float64
	Speed         float64
	TargetReached bool
}

func NewFollow() (f Follow) {
	f.Target = EntityNone
	return
}
