package cmd

import (
	"fmt"
	"log"
	"statuscakectl/helpers"
	"statuscakectl/statuscake"

	"github.com/spf13/cobra"
)

var createCmdUptime = &cobra.Command{
	Use:     "uptime",
	Short:   "create uptime check",
	Example: "statuscakectl create uptime --domain https://www.domain.com --checkrate 60 --type HTTP",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		if domain == "" {
			cmd.Usage()
			return
		}
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = domain
		}

		checkrate, _ := cmd.Flags().GetInt("checkrate")
		timeout, _ := cmd.Flags().GetInt("timeout")
		confirmation, _ := cmd.Flags().GetInt("confirmation")
		virus, _ := cmd.Flags().GetInt("virus")
		donotfind, _ := cmd.Flags().GetInt("donotfind")
		realbrowser, _ := cmd.Flags().GetInt("realbrowser")
		trigger, _ := cmd.Flags().GetInt("trigger")
		sslalert, _ := cmd.Flags().GetInt("sslalert")
		follow, _ := cmd.Flags().GetInt("follow")
		contacts, _ := cmd.Flags().GetString("contacts")
		testType, _ := cmd.Flags().GetString("type")
		findstring, _ := cmd.Flags().GetString("findstring")
		api, _ := cmd.Flags().GetString("api")
		user, _ := cmd.Flags().GetString("user")
		key, _ := cmd.Flags().GetString("key")

		domainOK := helpers.DomainValidation(domain)
		if !domainOK {
			return
		}

		testTypeValid := validateTestType(testType)
		if !testTypeValid {
			return
		}

		createCheck := statuscake.CreateUptimeCheck(name, domain, checkrate, timeout, confirmation, virus,
			donotfind, realbrowser, trigger, sslalert, follow, contacts, testType, findstring, api, user, key)
		if !createCheck {
			fmt.Println("Failed to create uptime check")
			return
		}

		log.Println("Success creating uptime check")
	},
}

func init() {
	// flags
	createCmdUptime.Flags().IntP("checkrate", "r", 86400, "Check rate in seconds default 86400 [300, 600, 1800, 3600, 86400, 2073600]")
	createCmdUptime.MarkFlagRequired("checkrate")
	createCmdUptime.Flags().Int("timeout", 10, "Timeout in an int form representing seconds.(default 10sec)")
	createCmdUptime.Flags().Int("confirmation", 1, "Confimation servers before alert (default 1)")
	createCmdUptime.Flags().Int("virus", 1, "Enable virus checking or not. default 1 = enable")
	createCmdUptime.Flags().String("findstring", "", "A string that should either be found or not found")
	createCmdUptime.Flags().String("name", "", "A name of the test")
	createCmdUptime.Flags().Int("donotfind", 0, "If the above string should be found to trigger a alert. 1 = will trigger if FindString found")
	createCmdUptime.Flags().StringP("type", "t", "HTTP", "Type of test type to use: HTTP,TCP,PING (default HTTP)")
	createCmdUptime.MarkFlagRequired("type")
	createCmdUptime.Flags().Int("realbrowser", 0, "Use 1 to TURN OFF real browser testing")
	createCmdUptime.Flags().Int("trigger", 1, "How many minutes to wait before sending an alert")
	createCmdUptime.Flags().Int("sslalert", 1, "HTTP Tests only. If enabled, tests will send warnings if the SSL certificate is about to expire. Paid users only, (default on)")
	createCmdUptime.Flags().Int("follow", 1, "HTTP Tests only. If enabled, our tests will follow redirects and logo the status of the final page")
}

func validateTestType(testType string) bool {
	testTypes := []string{
		"HTTP",
		"TCP",
		"PING",
	}
	for _, t := range testTypes {
		if t == testType {
			return true
		}
	}
	fmt.Println("Error: test type is not valid, only HTTP/TCP/PING")
	return false
}
