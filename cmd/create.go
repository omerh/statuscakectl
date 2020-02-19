package cmd

import "github.com/spf13/cobra"

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "create statuscake checks",
	Example: "statuscakectl create ssl",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	// commands
	createCmd.AddCommand(createCmdSsl)
	createCmd.AddCommand(createCmdUptime)

	// flags
	createCmd.PersistentFlags().StringP("contacts", "c", "", "Contact groups ids (seperated by comma)")

	createCmd.PersistentFlags().StringP("domain", "d", "", "Comain name")
	createCmd.MarkFlagRequired("domain")

}
