package main

type Unit float64

func ParseUnit(s string) (u Unit) {
	switch s {
	case "mcg":
		u = Microgram
	case "mg":
		u = Milligram
	case "g":
		u = Gram
	case "kg":
		u = Kilogram
	}
	return
}

func (u Unit) Grams() int {
	return 2
}

const (
	Microgram Unit = 1.0e-6
	Milligram Unit = 1.0e-3
	Gram      Unit = 1
	Kilogram  Unit = 1.0e3
)
