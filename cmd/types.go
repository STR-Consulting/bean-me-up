package cmd

import (
	"context"
	"fmt"

	"github.com/STR-Consulting/bean-me-up/internal/beans"
	"github.com/STR-Consulting/bean-me-up/internal/clickup"
	"github.com/spf13/cobra"
)

var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "List available custom task types",
	Long: `Lists all custom task types (e.g., Bug, Milestone) available in your ClickUp workspaces.

Use this command to find task type IDs for configuring type_mapping
in your .beans.clickup.yml configuration.

Requires CLICKUP_TOKEN environment variable to be set.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Get ClickUp token
		token, err := getClickUpToken()
		if err != nil {
			return err
		}

		// Create client
		client := clickup.NewClient(token)

		// Fetch custom items
		items, err := client.GetCustomItems(ctx)
		if err != nil {
			return fmt.Errorf("fetching custom task types: %w", err)
		}

		if jsonOut {
			return outputTypesJSON(items)
		}

		return outputTypesText(items)
	},
}

func outputTypesJSON(items []clickup.CustomItem) error {
	return outputJSON(items)
}

func outputTypesText(items []clickup.CustomItem) error {
	if len(items) == 0 {
		fmt.Println("No custom task types found.")
		fmt.Println("\nNote: Custom task types require a ClickUp Business plan or higher.")
		return nil
	}

	fmt.Printf("Custom task types:\n\n")

	for _, item := range items {
		fmt.Printf("%s\n", item.Name)
		fmt.Printf("  ID: %d\n", item.ID)
		if item.Description != "" {
			fmt.Printf("  Description: %s\n", item.Description)
		}
		fmt.Println()
	}

	// Build a map of ClickUp types by lowercase name for matching
	clickupTypes := make(map[string]clickup.CustomItem)
	for _, item := range items {
		clickupTypes[item.Name] = item
	}

	fmt.Println("Add these to your .beans.clickup.yml to map bean types:")
	fmt.Println()
	fmt.Println("  clickup:")
	fmt.Println("    type_mapping:")

	// Show all standard bean types with suggested mappings
	for _, beanType := range beans.StandardTypes {
		if item, ok := clickupTypes[beanType]; ok {
			fmt.Printf("      %s: %d  # %s\n", beanType, item.ID, item.Name)
		} else {
			fmt.Printf("      %s: 0  # Task (default)\n", beanType)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(typesCmd)
}
