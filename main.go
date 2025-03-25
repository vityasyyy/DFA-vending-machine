package main

import (
	"basot/config"
	"basot/src"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Global products available in the vending machine
var vendingMachineProducts = []config.VendingMachineProduct{
	{Name: "Coffee", Price: 3000, StateKey: "3000"},
	{Name: "Tea", Price: 4000, StateKey: "4000"},
	{Name: "Hot Chocolate", Price: 6000, StateKey: "6000"},
}

// Global variable to track the entire input history for diagram
var transactionHistory []string

func main() {
	// Path to your DFA configuration file
	configPath := "config/dfa_config.txt"

	// Parse the DFA configuration
	dfa, err := src.ParseDFAConfig(configPath)
	if err != nil {
		fmt.Printf("Error parsing DFA config: %v\n", err)
		return
	}

	// Initialize the current state with the start state
	currentState := dfa.StartState

	// Initialize transaction history
	transactionHistory = []string{currentState}

	// Create an endless loop for DFA evaluation
	reader := bufio.NewReader(os.Stdin)
	purchaseMade := false

	// Print initial message
	fmt.Println("\n=== Vending Machine ===")
	fmt.Println("Enter coins (positive values) or buy products (negative values):")
	fmt.Println("To purchase: PA (Coffee), PB (Tea), PC (Hot Chocolate)")
	fmt.Println("Type 'quit' to quit:")

	// Loop to process inputs from user, while there is no purchase made
	for !purchaseMade {
		// Read user input
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Check for exit command
		if strings.ToLower(input) == "quit" {
			// Exit command received
			fmt.Println("Processing input stream and exiting...")
			// Process the entire input stream before exiting
			displayStateTransitionDiagram(transactionHistory)
			break
		}

		// Check if this is a purchase command
		if strings.HasPrefix(input, "P") {
			// Extract the price from the input
			input = config.PriceMap[input]
			if input == "" {
				fmt.Print("Error: Product not found.\n")
				continue
			}
			productPrice, err := strconv.Atoi(input[1:])
			if err == nil {
				// This is a valid purchase command
				newState, change := src.ProcessPurchase(input, currentState, vendingMachineProducts, dfa)

				// If state changed, purchase was successful
				if newState != currentState {
					currentState = newState
					// Record this purchase in our input stream
					transactionHistory = append(transactionHistory, fmt.Sprintf("Purchase: %d (Change: %d)", productPrice, change))
					transactionHistory = append(transactionHistory, currentState)

					// Set flag to exit after purchase
					purchaseMade = true

					// Process the entire input stream before exiting
					fmt.Println("\nPurchase completed! Processing transaction history and exiting...")
					displayStateTransitionDiagram(transactionHistory)
					break
				} else {
					// Failed purchase
					transactionHistory = append(transactionHistory, fmt.Sprintf("Failed Purchase: %d", productPrice))
					transactionHistory = append(transactionHistory, currentState)
				}
				continue
			}
		}

		// Handle regular input
		nextState := src.ProcessInput(dfa, currentState, input, vendingMachineProducts)

		// Check if the machine is in a dead state
		if nextState == "DEAD" {
			// Machine is in a dead state
			fmt.Println("Error: Machine is in a dead state. Exiting...")
			break
		}

		// Check if the state has changed
		if nextState != currentState {
			currentState = nextState
			// Record the state transition
			transactionHistory = append(transactionHistory, fmt.Sprintf("Input: %s", input))
			transactionHistory = append(transactionHistory, currentState)
		}
	}
}

// displayStateTransitionDiagram shows a text-based representation of the state transition path
func displayStateTransitionDiagram(history []string) {
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

	fmt.Println("\n============================================")
	fmt.Println("Thank you for using the vending machine!")
	fmt.Println("============================================")
}
