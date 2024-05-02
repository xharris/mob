package component

type Health struct {
	Remaining int
	Total     int
}

func NewHealth() Health {
	return Health{
		Remaining: 100,
		Total:     100,
	}
}
