package src

import (
	"basot/config"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseDFAConfig reads a configuration file and returns a DFA struct
func ParseDFAConfig(filePath string) (DFA *config.DFA, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dfa := &config.DFA{
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
