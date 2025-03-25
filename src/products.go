package src

import (
	"basot/config"
	"fmt"
)

// getAvailableProducts checks which products are available at the current state
func GetAvailableProducts(state string, vendingMachineProducts []config.VendingMachineProduct) []config.VendingMachineProduct {
	var availableProducts []config.VendingMachineProduct

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
func DisplayAvailableProducts(state string, vendingMachineProducts []config.VendingMachineProduct) {
	products := GetAvailableProducts(state, vendingMachineProducts)

	if len(products) == 0 {
		fmt.Println("No products available at this state.")
		return
	}

	fmt.Println("Available Products:")
	for _, product := range products {
		fmt.Printf("  - %s (price: %d): ON\n", product.Name, product.Price)
	}
}
