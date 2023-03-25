package commands

import "github.com/spf13/cobra"

var address string
var mode string
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Display table of aircraft tracked by receiver running on provided port",
	Long:  `Connects to a receiver running on a provided port. Decodes messages and displays tracked aircraft in a table.`,
	Run: func(cmd *cobra.Command, args []string) {
		println("connect")
		println(address)
		println(mode)
	},
}

func init() {
	connectCmd.Flags().StringVarP(&address, "address", "a", "", "address to connect to (include port)")
	connectCmd.Flags().StringVarP(&mode, "mode", "m", "", "mode of source (currently only raw supported)")

	connectCmd.MarkFlagRequired("address")
	connectCmd.MarkFlagRequired("mode")

	rootCmd.AddCommand(connectCmd)
}
