package cmd

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "list statuscake checks",
	Example: "statuscakectl list ssl",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	listCmd.AddCommand(listCmdSsl)
	listCmd.AddCommand(listCmdUptime)
}
