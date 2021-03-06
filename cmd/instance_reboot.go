package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceRebootCmd = &cobra.Command{
	Use:     "reboot",
	Example: "civo instance reboot ID/HOSTNAME",
	Aliases: []string{"hard-reboot"},
	Short:   "Hard reboot an instance",
	Long: `Pull the power and restart the specified instance by part of its ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Finding instance %s", err)
			os.Exit(1)
		}

		_, err = client.RebootInstance(instance.ID)
		if err != nil {
			utility.Error("Rebooting instance %s", err)
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being rebooted\n", utility.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
