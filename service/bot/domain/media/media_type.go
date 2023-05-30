package media

var (
	Video = Type{"video"}
	Image = Type{"image"}
)

type Type struct {
	m string
}

func (m Type) String() string {
	return m.m
}

func (h Type) IsZero() bool {
	return h == Type{}
}
