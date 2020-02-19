package cmd

import "github.com/spf13/cobra"

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete statuscake checks",
	Example: "statuscakectl delete ssl -d domain.com",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	// commands
	deleteCmd.AddCommand(deleteCmdSSL)
	deleteCmd.AddCommand(deleteCmdUptime)

	//flags
	deleteCmd.PersistentFlags().StringP("domain", "d", "", "Domain to delete")
	deleteCmd.PersistentFlags().Int("id", 0, "Check id to delete")
}
