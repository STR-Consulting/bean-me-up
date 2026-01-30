package cmd

import (
	"context"
	"fmt"
	"sort"

	"github.com/STR-Consulting/bean-me-up/internal/beans"
	"github.com/STR-Consulting/bean-me-up/internal/clickup"
	"github.com/spf13/cobra"
)

var fixTagsDryRun bool

var fixTagsCmd = &cobra.Command{
	Use:   "fix-tags",
	Short: "Create missing space-level tags in ClickUp",
	Long: `Ensures all tags used by synced beans exist as space-level tags in ClickUp.

Tags synced to ClickUp tasks are only created at the task level, which means
they don't appear in ClickUp's tag filter/search dropdown. This command reads
all unique tags from beans that have a linked ClickUp task and creates any
missing ones at the space level.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if err := requireListID(); err != nil {
			return err
		}

		token, err := getClickUpToken()
		if err != nil {
			return err
		}

		syncStore, err := loadSyncStore()
		if err != nil {
			return fmt.Errorf("loading sync state: %w", err)
		}

		// List all beans and collect tags from those with a linked task
		beansClient := beans.NewClient(getBeansPath())
		allBeans, err := beansClient.List()
		if err != nil {
			return fmt.Errorf("listing beans: %w", err)
		}

		tagSet := make(map[string]bool)
		for _, b := range allBeans {
			if syncStore.GetTaskID(b.ID) == nil {
				continue
			}
			for _, tag := range b.Tags {
				tagSet[tag] = true
			}
		}

		if len(tagSet) == 0 {
			fmt.Println("No tags found on synced beans.")
			return nil
		}

		// Sort for deterministic output
		tags := make([]string, 0, len(tagSet))
		for tag := range tagSet {
			tags = append(tags, tag)
		}
		sort.Strings(tags)

		// Get space ID from the configured list
		client := clickup.NewClient(token)
		list, err := client.GetList(ctx, cfg.Beans.ClickUp.ListID)
		if err != nil {
			return fmt.Errorf("getting list: %w", err)
		}

		// Fetch existing space tags
		if err := client.PopulateSpaceTagCache(ctx, list.SpaceID); err != nil {
			return fmt.Errorf("fetching space tags: %w", err)
		}

		if fixTagsDryRun {
			fmt.Printf("Found %d unique tag(s) across synced beans:\n", len(tags))
			for _, tag := range tags {
				fmt.Printf("  %s\n", tag)
			}
			fmt.Println("\nDry run — no tags were created.")
			return nil
		}

		var created, existing int
		for _, tag := range tags {
			alreadyExists := client.HasSpaceTag(tag)
			if alreadyExists {
				existing++
				continue
			}
			if err := client.EnsureSpaceTag(ctx, list.SpaceID, tag); err != nil {
				return fmt.Errorf("creating space tag %q: %w", tag, err)
			}
			fmt.Printf("  created: %s\n", tag)
			created++
		}

		fmt.Printf("\nDone. %d created, %d already existed.\n", created, existing)
		return nil
	},
}

func init() {
	fixTagsCmd.Flags().BoolVar(&fixTagsDryRun, "dry-run", false, "preview which tags would be created without creating them")
	rootCmd.AddCommand(fixTagsCmd)
}
