package config

// DFA represents a Deterministic Finite Automaton
type DFA struct {
	States       []string
	Alphabet     []string
	AcceptStates []string
	StartState   string
	Transitions  map[string]map[string]string
}

// VendingMachineProduct represents a product in the vending machine
type VendingMachineProduct struct {
	Name     string
	Price    int
	StateKey string
}

// price mapping for the product
var PriceMap = map[string]string{
	"PA": "-3000",
	"PB": "-4000",
	"PC": "-6000",
}
