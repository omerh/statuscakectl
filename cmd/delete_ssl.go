package cmd

import (
	"fmt"
	"net/url"
	"statuscakectl/helpers"
	"statuscakectl/statuscake"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmdSSLExample = `
  # Delete by domain
  statuscakectl delete ssl --domain domain.com
  # Delete by Id
  statuscakectl delete ssl --id 11111111
`

var deleteCmdSSL = &cobra.Command{
	Use:     "ssl",
	Short:   "delete ssl checks",
	Example: deleteCmdSSLExample,
	Run: func(cmd *cobra.Command, args []string) {
		api, _ := cmd.Flags().GetString("api")
		user, _ := cmd.Flags().GetString("user")
		key, _ := cmd.Flags().GetString("key")
		id, _ := cmd.Flags().GetInt("id")
		domain, _ := cmd.Flags().GetString("domain")

		// delete with id
		if id != 0 {
			ok := statuscake.DeleteSSLCheck(strconv.Itoa(id), api, user, key)
			if !ok {
				fmt.Println("Error: Failed to delete SSL check")
				return
			}
			fmt.Printf("Sucess: SSL test %v was deleted", id)
			return
		}

		// delete with domain
		ok := helpers.DomainValidation(domain)
		if !ok {
			return
		}

		// no need for err, its parses everything
		target, _ := url.Parse(domain)
		sslCheckSlice := statuscake.ListSSLChecks(api, user, key)
		for _, r := range sslCheckSlice {
			t, _ := url.Parse(r.Domain)
			if t.Host == target.Host || t.Host == domain {
				ok := statuscake.DeleteSSLCheck(r.ID, api, user, key)
				if !ok {
					fmt.Println("Error: Failed to delete SSL check")
					return
				}
				fmt.Printf("Success: SSL test %v was deleted\n", r.ID)
			}
		}
	},
}
