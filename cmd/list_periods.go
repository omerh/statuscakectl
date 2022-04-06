package cmd

import (
	"fmt"
	"log"
	"sort"
	"statuscakectl/helpers"
	"statuscakectl/statuscake"

	"github.com/spf13/cobra"
)

var listCmdPeriods = &cobra.Command{
	Use:     "periods",
	Short:   "list periods for test-id or domain (if both provided, looks for id of the domain test)",
	Example: "statuscakectl list periods --test-id 1111\nstatuscakectl list periods --domain testdomain.com",
	Run: func(cmd *cobra.Command, args []string) {
		api, _ := cmd.Flags().GetString("api")
		user, _ := cmd.Flags().GetString("user")
		key, _ := cmd.Flags().GetString("key")
		testID, _ := cmd.Flags().GetInt("test-id")
		domain, _ := cmd.Flags().GetString("domain")

		if testID < 1 && domain == "" {
			fmt.Println("Please make sure to provide a valid test-id or domain flags")
			return
		}

		if domain != "" {
			uptimeTests := statuscake.ListUptime(api, user, key)
			for _, t := range uptimeTests {
				hostname, err := helpers.GetHostnameFromUrl(t.WebsiteURL)
				if err != nil {
					log.Printf("Can't get hostname from this URL:%v\n", t.WebsiteURL)
				}
				if hostname == domain {
					testID = t.TestID
					break
				}
			}
			if testID < 1 {
				fmt.Printf("Cannot find test for domain: %v\n", domain)
				return
			}
		}

		detailedData := statuscake.GetDetailedTestData(testID, api, user, key)
		domain, err := helpers.GetHostnameFromUrl(detailedData.URI)
		if err != nil {
			log.Fatalf("Cant get hostname from URI:%v", detailedData.URI)
		}

		periods := statuscake.ListPeriods(api, user, key, testID)

		sort.Slice(periods, func(i, j int) bool {
			return periods[i].StartUnix < periods[j].StartUnix
		})

		fmt.Printf("Total %v checks in statuscake account:\n", len(periods))
		for _, r := range periods {
			fmt.Printf("%v was %v for %v from %v to %v\n", domain, r.Status, r.Period, r.Start, r.End)
		}
	},
}

func init() {
	// flags
	listCmdPeriods.Flags().Int("test-id", 0, "TestID of the test you want to get periods for")
	listCmdPeriods.Flags().StringP("domain", "d", "", "Domain name to find periods for")
}
