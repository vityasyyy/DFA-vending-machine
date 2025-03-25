package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

// Global products available in the vending machine
var vendingMachineProducts = []VendingMachineProduct{
	{Name: "Coffee", Price: 3000, StateKey: "3000"},
	{Name: "Tea", Price: 4000, StateKey: "4000"},
	{Name: "Hot Chocolate", Price: 6000, StateKey: "6000"},
}

// Global variable to track the entire input history for diagram
var transactionHistory []string

func main() {
	// Path to your DFA configuration file
	configPath := "./dfa_config.txt"

	// Parse the DFA configuration
	dfa, err := parseDFAConfig(configPath)
	if err != nil {
		fmt.Printf("Error parsing DFA config: %v\n", err)
		return
	}
	// Print the DFA components
	fmt.Println("=== DFA Configuration ===")
	fmt.Println("States:", dfa.States)
	fmt.Println("Alphabet:", dfa.Alphabet)
	fmt.Println("Start:", dfa.StartState)
	fmt.Println("Accept States:", dfa.AcceptStates)

	fmt.Println("\n=== Transitions ===")
	for fromState, transitions := range dfa.Transitions {
		for symbol, toState := range transitions {
			fmt.Printf("%s\t%s\t%s\n", fromState, symbol, toState)
		}
	}

	// Initialize the current state with the start state
	currentState := dfa.StartState

	// Create an array to store all inputs
	var inputStream []string

	// Initialize transaction history
	transactionHistory = []string{currentState}

	// Create an endless loop for DFA evaluation
	reader := bufio.NewReader(os.Stdin)
	purchaseMade := false

	fmt.Println("\n=== Vending Machine ===")
	fmt.Println("Enter coins (positive values) or buy products (negative values):")
	fmt.Println("To purchase: -3000 (Coffee), -4000 (Tea), -6000 (Hot Chocolate)")
	fmt.Println("Type 'quit' to quit:")

	for !purchaseMade {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Check for exit command
		if strings.ToLower(input) == "exit" {
			fmt.Println("Processing input stream and exiting...")
			processInputStream(dfa, inputStream)
			break
		}

		// Check if this is a purchase command (negative number)
		if strings.HasPrefix(input, "-") {
			productPrice, err := strconv.Atoi(input[1:])
			if err == nil {
				// This is a valid purchase command
				newState, change := processPurchase(productPrice, currentState)

				// If state changed, purchase was successful
				if newState != currentState {
					currentState = newState
					// Record this purchase in our input stream
					inputStream = append(inputStream, input)
					transactionHistory = append(transactionHistory, fmt.Sprintf("Purchase: %d (Change: %d)", productPrice, change))
					transactionHistory = append(transactionHistory, currentState)

					// Set flag to exit after purchase
					purchaseMade = true

					// Process the entire input stream before exiting
					fmt.Println("\nPurchase completed! Processing transaction history and exiting...")
					processInputStream(dfa, inputStream)
					displayStateTransitionDiagram(inputStream, transactionHistory)
					break
				} else {
					// Failed purchase
					inputStream = append(inputStream, input)
					transactionHistory = append(transactionHistory, fmt.Sprintf("Failed Purchase: %d", productPrice))
					transactionHistory = append(transactionHistory, currentState)
				}
				continue
			}
		}

		// Handle regular input
		inputStream = append(inputStream, input)
		nextState := processInput(dfa, currentState, input)

		if nextState != currentState {
			currentState = nextState
			// Record the state transition
			transactionHistory = append(transactionHistory, fmt.Sprintf("Input: %s", input))
			transactionHistory = append(transactionHistory, currentState)
		}
	}
}

// getAvailableProducts checks which products are available at the current state
func getAvailableProducts(state string) []VendingMachineProduct {
	var availableProducts []VendingMachineProduct

	// Extract the numeric part from the state name
	// First, remove any non-digit characters from the state
	numericPart := ""
	for _, char := range state {
		if char >= '0' && char <= '9' {
			numericPart += string(char)
		}
	}

	// Convert to integer if possible
	stateValue := 0
	if numericPart != "" {
		fmt.Sscanf(numericPart, "%d", &stateValue)
	}

	// Check each product threshold
	for _, product := range vendingMachineProducts {
		productValue := 0
		fmt.Sscanf(product.StateKey, "%d", &productValue)

		// If current state value is equal to or greater than the product price,
		// the product is available
		if stateValue >= productValue {
			availableProducts = append(availableProducts, product)
		}
	}

	return availableProducts
}

// displayAvailableProducts shows which products can be purchased at current state
func displayAvailableProducts(state string) {
	products := getAvailableProducts(state)

	if len(products) == 0 {
		fmt.Println("No products available at this state.")
		return
	}

	fmt.Println("Available Products:")
	for _, product := range products {
		fmt.Printf("  - %s (price: %d): ON\n", product.Name, product.Price)
	}
}

// processInputStream takes a DFA and an array of inputs and shows the state transition path
func processInputStream(dfa *DFA, inputStream []string) {
	if len(inputStream) == 0 {
		fmt.Println("No inputs to process.")
		return
	}

	currentState := dfa.StartState
	fmt.Printf("\n=== Vending Machine State Path ===\n")
	fmt.Printf("Starting State: %s\n", currentState)

	// Check for products available at the start state
	displayAvailableProducts(currentState)

	path := fmt.Sprintf("%s", currentState)

	// Track total inserted and spent
	totalInserted := 0
	totalSpent := 0
	purchasesMade := 0

	for i, input := range inputStream {
		// Check if this is a purchase command
		if strings.HasPrefix(input, "-") {
			productPrice, err := strconv.Atoi(input[1:])
			if err == nil {
				// Valid purchase command
				fmt.Printf("\nInput[%d]: '%s' - Attempting to purchase product with price %d\n",
					i, input, productPrice)

				newState, change := processPurchase(productPrice, currentState)

				if currentState != newState { // Purchase was successful
					totalSpent += productPrice
					purchasesMade++
					currentState = newState
					path += fmt.Sprintf(" -> %s (Purchase: %d)", currentState, productPrice)
					fmt.Printf("Change: %d\n", change)
				} else {
					path += fmt.Sprintf(" -> %s (Failed purchase: %d)", currentState, productPrice)
				}
				continue
			}
		}

		// Regular input processing (adding money)
		// Skip processing if input is not in alphabet
		if !isInAlphabet(dfa, input) {
			fmt.Printf("Error: Input '%s' contains symbols not in the alphabet, skipping.\n", input)
			continue
		}

		// Try to parse input as a number and track it
		if coinValue, err := strconv.Atoi(input); err == nil {
			totalInserted += coinValue
		}

		// Check if there's a transition for this state and input
		if nextState, exists := dfa.Transitions[currentState][input]; exists {
			currentState = nextState
			path += fmt.Sprintf(" -> %s", currentState)
			fmt.Printf("\nInput[%d]: '%s' transitions to state: %s\n", i, input, currentState)

			// Display available products at this state
			displayAvailableProducts(currentState)
		} else {
			fmt.Printf("Error: No transition defined for state '%s' and input '%s'.\n", currentState, input)
		}
	}

	// Check if the final state is an accept state
	isAccepted := false
	for _, acceptState := range dfa.AcceptStates {
		if acceptState == currentState {
			isAccepted = true
			break
		}
	}

	fmt.Printf("\nTransition Path: %s\n", path)
	fmt.Printf("Final state: %s\n", currentState)

	// Display transaction summary
	fmt.Printf("\n=== Transaction Summary ===\n")
	fmt.Printf("Total money inserted: %d\n", totalInserted)
	fmt.Printf("Total money spent: %d\n", totalSpent)
	fmt.Printf("Products purchased: %d\n", purchasesMade)

	// Calculate remaining balance or change
	remainingValue := totalInserted - totalSpent
	if remainingValue > 0 {
		fmt.Printf("Remaining balance/change: %d\n", remainingValue)
	}

	if isAccepted {
		fmt.Println("Result: ACCEPTED ✓ - You can make a purchase!")
	} else {
		fmt.Println("Result: REJECTED ✗ - Insufficient funds or invalid state.")
	}
}

// processInput takes a DFA, the current state, and an input string and processes it,
// displaying each transition step and returning the new state
func processInput(dfa *DFA, currentState string, input string) string {
	fmt.Printf("Starting at state: %s\n", currentState)

	// Display available products at the current state
	displayAvailableProducts(currentState)

	// Check if this is a purchase command (negative number)
	if strings.HasPrefix(input, "-") {
		productPrice, err := strconv.Atoi(input[1:])
		if err == nil {
			newState, _ := processPurchase(productPrice, currentState)
			return newState
		}
	}

	// Check if the input string is in the alphabet
	if !isInAlphabet(dfa, input) {
		fmt.Printf("Error: Input '%s' contains symbols not in the alphabet.\n", input)
		return currentState
	}

	// Check if there's a transition for this state and input
	if nextState, exists := dfa.Transitions[currentState][input]; exists {
		fmt.Printf("Transition: δ(%s, %s) = %s\n", currentState, input, nextState)
		currentState = nextState
		fmt.Println("Current state:", currentState)
	} else {
		fmt.Printf("Error: No transition defined for state '%s' and input '%s'.\n", currentState, input)
		return currentState
	}

	// Check if the final state is an accept state
	isAccepted := false
	for _, acceptState := range dfa.AcceptStates {
		if acceptState == currentState {
			isAccepted = true
			break
		}
	}

	fmt.Printf("\nFinal state: %s\n", currentState)

	// Display available products after the transition
	displayAvailableProducts(currentState)

	if isAccepted {
		fmt.Println("Result: ACCEPTED ✓")
	} else {
		fmt.Println("Result: REJECTED ✗")
	}

	return currentState
}

// Helper function to check if the input string is in the alphabet
func isInAlphabet(dfa *DFA, input string) bool {
	for _, symbol := range dfa.Alphabet {
		if symbol == input {
			return true
		}
	}
	return false
}

func parseDFAConfig(filePath string) (*DFA, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dfa := &DFA{
		States:       make([]string, 0),
		Alphabet:     make([]string, 0),
		AcceptStates: make([]string, 0),
		Transitions:  make(map[string]map[string]string),
	}

	scanner := bufio.NewScanner(file)
	section := ""

	fmt.Println("Reading configuration file:")

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Check if this is a section header
		if strings.HasSuffix(line, ":") {
			section = strings.TrimSuffix(line, ":")
			continue
		}

		// Process the line based on section
		switch section {
		case "States":
			states := strings.Split(line, ",")
			for _, state := range states {
				trimmedState := strings.TrimSpace(state)
				dfa.States = append(dfa.States, trimmedState)
				// Initialize the transitions map for this state
				dfa.Transitions[trimmedState] = make(map[string]string)
			}

		case "Alphabet":
			symbols := strings.Split(line, ",")
			for _, symbol := range symbols {
				dfa.Alphabet = append(dfa.Alphabet, strings.TrimSpace(symbol))
			}

		case "Accept":
			acceptStates := strings.Split(line, ",")
			for _, state := range acceptStates {
				dfa.AcceptStates = append(dfa.AcceptStates, strings.TrimSpace(state))
			}

		case "Start":
			startState := strings.TrimSpace(line)
			dfa.StartState = startState

		case "Transitions":
			parts := strings.Fields(line)
			if len(parts) == 3 {
				fromState := parts[0]
				symbol := parts[1]
				toState := parts[2]
				if dfa.Transitions[fromState] == nil {
					dfa.Transitions[fromState] = make(map[string]string)
				}
				dfa.Transitions[fromState][symbol] = toState
			} else {
				fmt.Printf("Warning: Invalid transition format in line: '%s'\n", line)
			}
		default:
			fmt.Printf("Warning: Unknown section '%s'\n", section)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	fmt.Println("DFA parsing complete")
	return dfa, nil
}

// ProcessPurchase handles buying a product and calculates change
func processPurchase(productPrice int, currentState string) (string, int) {
	// Extract the numeric part from the state name
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

	// Determine the new state after purchase (reset to initial state or remaining change)
	newState := "q0"
	if change > 0 {
		newState = fmt.Sprintf("q%d", change)
	}

	// After a successful purchase, add a visual indicator
	fmt.Println("\n============================================")
	fmt.Printf("🎉 Successfully purchased %s!\n", productName)
	fmt.Printf("💰 Change: %d\n", change)
	fmt.Println("============================================")

	return newState, change
}

// displayStateTransitionDiagram shows a text-based representation of the state transition path
func displayStateTransitionDiagram(inputStream []string, history []string) {
	fmt.Println("\n============================================")
	fmt.Println("        STATE TRANSITION DIAGRAM")
	fmt.Println("============================================")

	fmt.Println("\nState Transitions:")
	fmt.Println("----------------")

	// Display the full transition history with arrows
	for i := 0; i < len(history); i++ {
		if i%2 == 0 {
			// This is a state
			if i > 0 {
				fmt.Printf("\n")
			}
			fmt.Printf("State: %s", history[i])
		} else {
			// This is a transition label
			fmt.Printf("\n    |\n    | %s\n    v\n", history[i])
		}
	}

	// Check which products were purchased
	purchasedProducts := []string{}
	for _, input := range inputStream {
		if strings.HasPrefix(input, "-") {
			price, _ := strconv.Atoi(input[1:])
			for _, product := range vendingMachineProducts {
				if product.Price == price {
					purchasedProducts = append(purchasedProducts, product.Name)
				}
			}
		}
	}

	fmt.Println("\n============================================")
	fmt.Println("TRANSACTION SUMMARY")
	fmt.Println("============================================")

	// Calculate total money inserted
	totalInserted := 0
	for _, input := range inputStream {
		if !strings.HasPrefix(input, "-") {
			if value, err := strconv.Atoi(input); err == nil {
				totalInserted += value
			}
		}
	}

	// Calculate total spent
	totalSpent := 0
	for _, input := range inputStream {
		if strings.HasPrefix(input, "-") {
			if value, err := strconv.Atoi(input[1:]); err == nil {
				totalSpent += value
			}
		}
	}

	fmt.Printf("Total money inserted: %d\n", totalInserted)
	fmt.Printf("Total money spent: %d\n", totalSpent)
	fmt.Printf("Change received: %d\n", totalInserted-totalSpent)

	if len(purchasedProducts) > 0 {
		fmt.Println("\nProducts purchased:")
		for _, product := range purchasedProducts {
			fmt.Printf("  - %s\n", product)
		}
	} else {
		fmt.Println("\nNo products were purchased.")
	}

	fmt.Println("\n============================================")
	fmt.Println("Thank you for using the vending machine!")
	fmt.Println("============================================")
}
