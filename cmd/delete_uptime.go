package cmd

import (
	"fmt"
	"net/url"
	"statuscakectl/helpers"
	"statuscakectl/statuscake"

	"github.com/spf13/cobra"
)

var deleteUptimeExample = `
  # Delete by domain
  statuscakectl delete uptime --domain https://domain.com --id 11111
  # Delete by Id
  statuscakectl delete uptime --id 11111
`

var deleteCmdUptime = &cobra.Command{
	Use:     "uptime",
	Short:   "delete uptime checks",
	Example: deleteUptimeExample,
	Run: func(cmd *cobra.Command, args []string) {
		api, _ := cmd.Flags().GetString("api")
		user, _ := cmd.Flags().GetString("user")
		key, _ := cmd.Flags().GetString("key")
		id, _ := cmd.Flags().GetInt("id")
		domain, _ := cmd.Flags().GetString("domain")

		// delete with id
		if id != 0 {
			ok := statuscake.DeleteUptimeCheck(id, api, user, key)
			if !ok {
				fmt.Println("Error: Failed to delete Uptime check")
				return
			}
			fmt.Printf("Sucess: Uptime test %v was deleted\n", id)
			return
		}

		// delete with domain
		ok := helpers.DomainValidation(domain)
		if !ok {
			return
		}
		// no need for err, its parses everything
		target, _ := url.Parse(domain)
		sslCheckSlice := statuscake.ListUptime(api, user, key)
		for _, r := range sslCheckSlice {
			t, _ := url.Parse(r.WebsiteURL)
			if t.Host == target.Host || t.Host == domain {
				ok := statuscake.DeleteUptimeCheck(r.TestID, api, user, key)
				if !ok {
					fmt.Println("Error: Failed to delete Uptime check")
				}
				fmt.Printf("Success: Uptime test %v was deleted\n", r.TestID)
			}
		}

	},
}
