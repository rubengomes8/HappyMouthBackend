package enums

//go:generate enumer -json -text -sql -type RecipeType -trimprefix Recipe -transform snake-upper

// RecipeType defines the list of acceptable recipe types
type RecipeType uint

const (
	Any RecipeType = iota + 1
	Soup
	Pasta
	Pizza
	Curry
	OnePot
	Burger
	OvenRoast
	Salad
)
