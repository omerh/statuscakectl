package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "statuscakectl",
	Short: "statuscakectl command line tool to control statuscake api",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

// Execute using cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(calculateSlaCmd)

	// viper
	viper.AutomaticEnv()

	// Persistent flags
	rootCmdPFlags := rootCmd.PersistentFlags()
	rootCmdPFlags.String("api", "https://app.statuscake.com", "statuscake default api url")
	rootCmdPFlags.String("user", viper.GetString("STATUSCAKE_USER"), "statuscake user")
	rootCmdPFlags.String("key", viper.GetString("STATUSCAKE_KEY"), "Your statuscake api key")

}
