package src

import (
	"basot/config"
	"fmt"
	"strconv"
	"strings"
)

// processInput takes a DFA, the current state, and an input string and processes it,
// displaying each transition step and returning the new state
func ProcessInput(dfa *config.DFA, currentState string, input string, vendingMachineProducts []config.VendingMachineProduct) string {
	fmt.Printf("Starting at state: %s\n", currentState)

	// Display available products at the current state
	DisplayAvailableProducts(currentState, vendingMachineProducts)

	// Check if this is a purchase command
	if strings.HasPrefix(input, "P") {
		input = config.PriceMap[input]
		if input == "" {
			fmt.Printf("Error: Product '%s' not found.\n", input)
			return currentState
		}
		fmt.Print("input:", input)
		_, _ = ProcessPurchase(input, currentState, vendingMachineProducts, dfa)
	}

	// Check if the input string is in the alphabet
	if !IsInAlphabet(dfa, input) {
		fmt.Printf("Error: Input '%s' contains symbols not in the alphabet.\n", input)
		return currentState
	}

	// Check if there's a transition for this state and input
	if nextState, exists := dfa.Transitions[currentState][input]; exists {
		fmt.Printf("Transition: δ(%s, %s) = %s\n", currentState, input, nextState)
		currentState = nextState
		if currentState == "DEAD" {
			fmt.Println("Error: Machine is in a dead state. Exiting...")
			return currentState
		}
		fmt.Println("Current state:", currentState)
	} else {
		fmt.Printf("Error: No transition defined for state '%s' and input '%s'.\n", currentState, input)
		return currentState
	}

	// Display available products after the transition
	DisplayAvailableProducts(currentState, vendingMachineProducts)

	return currentState
}

// Helper function to check if the input string is in the alphabet
func IsInAlphabet(dfa *config.DFA, input string) bool {
	for _, symbol := range dfa.Alphabet {
		if symbol == input {
			return true
		}
	}
	return false
}

// ProcessPurchase handles buying a product and calculates change
func ProcessPurchase(input string, currentState string, vendingMachineProducts []config.VendingMachineProduct, dfa *config.DFA) (string, int) {
	// Extract the numeric part from the state name
	productPrice, _ := strconv.Atoi(input[1:])
	numericPart := ""
	for _, char := range currentState {
		if char >= '0' && char <= '9' {
			numericPart += string(char)
		}
	}

	// Convert to integer if possible
	currentValue := 0
	if numericPart != "" {
		fmt.Sscanf(numericPart, "%d", &currentValue)
	}

	// Check if there's enough money
	if currentValue < productPrice {
		fmt.Printf("Not enough money! Current value: %d, Product price: %d\n",
			currentValue, productPrice)
		return currentState, 0
	}

	// Calculate change
	change := currentValue - productPrice

	// Find product name
	productName := "Unknown"
	for _, product := range vendingMachineProducts {
		if product.Price == productPrice {
			productName = product.Name
			break
		}
	}

	// Determine the new state after purchase (based on remaining balance)
	newStateValue := dfa.Transitions[currentState][input]

	isAccepted := false

	for _, state := range dfa.AcceptStates {
		if newStateValue == state {
			isAccepted = true
			break
		}
	}
	// Validate if the new state exists in the DFA
	if _, exists := dfa.Transitions[currentState][input]; !exists {
		fmt.Printf("Error: Transition to state is not valid in the DFA.\n")
		return currentState, change
	}

	// After a successful purchase, add a visual indicator
	fmt.Println("\n============================================")
	fmt.Printf("🎉 Successfully purchased %s!\n", productName)
	fmt.Printf("💰 Change: %d\n", change)
	fmt.Printf("Accepted? %t\n", isAccepted)
	fmt.Println("============================================")

	return newStateValue, change
}
