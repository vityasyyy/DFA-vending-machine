package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type DFA struct {
	States       []string
	Alphabet     []string
	AcceptStates []string
	StartState   string
	Transitions  map[string]map[string]string
}

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

	// Create an endless loop for DFA evaluation
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n=== DFA Evaluation ===")
		fmt.Println("Enter a string to evaluate (type 'exit' to quit):")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		// Check for exit command
		if strings.ToLower(input) == "exit" {
			fmt.Println("Exiting program. Goodbye!")
			break
		}

		// Process the input through the DFA
		currentState = processInput(dfa, currentState, input)
	}
}

// processInput takes a DFA, the current state, and an input string and processes it,
// displaying each transition step
func processInput(dfa *DFA, currentState string, input string) string {
	fmt.Printf("Starting at state: %s\n", currentState)

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
