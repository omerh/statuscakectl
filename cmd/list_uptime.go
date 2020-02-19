package cmd

import (
	"fmt"
	"statuscakectl/statuscake"

	"github.com/spf13/cobra"
)

var listCmdUptime = &cobra.Command{
	Use:     "uptime",
	Short:   "list uptime checks",
	Example: "statuscakectl list uptime",
	Run: func(cmd *cobra.Command, args []string) {
		api, _ := cmd.Flags().GetString("api")
		user, _ := cmd.Flags().GetString("user")
		key, _ := cmd.Flags().GetString("key")

		uptimeCheckSlice := statuscake.ListUptime(api, user, key)
		fmt.Printf("Total %v checks in statuscake account:\n", len(uptimeCheckSlice))
		for _, r := range uptimeCheckSlice {
			fmt.Printf("Uptime test %v is %v with uptime of %v%% in the last day (check id: %v)\n", r.WebsiteURL, r.Status, r.Uptime, r.TestID)
		}
	},
}
