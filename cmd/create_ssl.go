package cmd

import (
	"fmt"
	"log"
	"net/url"
	"statuscakectl/statuscake"
	"strings"

	"github.com/spf13/cobra"
)

var createCmdSsl = &cobra.Command{
	Use:     "ssl",
	Short:   "create ssl check",
	Example: "statuscakectl create ssl -d domain.com -r 86400",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		if domain == "" {
			cmd.Usage()
			return
		}
		checkrate, _ := cmd.Flags().GetInt("checkrate")
		contacts, _ := cmd.Flags().GetString("contacts")
		api, _ := cmd.Flags().GetString("api")
		user, _ := cmd.Flags().GetString("user")
		key, _ := cmd.Flags().GetString("key")

		if !strings.Contains(domain, ".") {
			// url.Parse parses anything so checking for a . in the string
			fmt.Println("Please make sure to enter a valid domain (e.g domain.com)")
			return
		}

		sslCheckSlice := statuscake.ListSSLChecks(api, user, key)
		target, err := url.Parse(domain)
		if err != nil {
			fmt.Println("Please make sure to enter a valid domain (e.g domain.com)")
			return
		}

		for _, r := range sslCheckSlice {
			d, _ := url.Parse(r.Domain)
			if d.Host == target.String() || d.String() == target.String() {
				fmt.Printf("The domain %v is already is setup in statuscake\n", domain)
				return
			}
		}
		ok := statuscake.CreateSSLStatuscakeCheck(domain, checkrate, contacts, api, user, key)
		if !ok {
			log.Printf("There was a problem creating ssl check for domain %v", domain)
			return
		}
		log.Println("Success creating check")
	},
}

func init() {
	// flags
	createCmdSsl.Flags().IntP("checkrate", "r", 86400, "Check rate in seconds default 86400 [300, 600, 1800, 3600, 86400, 2073600]")
}
