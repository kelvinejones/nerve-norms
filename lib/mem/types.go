package mem

type XY struct {
	X float64
	Y float64
}

type XYZ struct {
	X float64
	Y float64
	Z float64
}

type section interface {
	Header() string
	Parser
}
