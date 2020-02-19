package cmd

import (
	"fmt"
	"statuscakectl/statuscake"

	"github.com/spf13/cobra"
)

var listCmdSsl = &cobra.Command{
	Use:     "ssl",
	Short:   "list ssl checks",
	Example: "statuscakectl list ssl",
	Run: func(cmd *cobra.Command, args []string) {
		api, _ := cmd.Flags().GetString("api")
		user, _ := cmd.Flags().GetString("user")
		key, _ := cmd.Flags().GetString("key")

		sslCheckSlice := statuscake.ListSSLChecks(api, user, key)
		for _, r := range sslCheckSlice {
			var certStatus, issuerCn string
			if r.CertStatus == "" {
				certStatus = "PENDING"
			}
			if r.IssuerCn == "" {
				issuerCn = "PENDING"
			}
			fmt.Printf("Certificate status for domain %v is %v, certificate issues by %v and valid until %v (id: %v)\n", r.Domain, certStatus, issuerCn, r.ValidUntilUtc, r.ID)
		}

	},
}
