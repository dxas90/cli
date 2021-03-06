package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var bootableVolume bool
var createSizeGB int

var volumeCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo volume create NAME [flags]",
	Short:   "Create a new volume",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		volumeConfig := &civogo.VolumeConfig{
			Name:          args[0],
			SizeGigabytes: createSizeGB,
			Bootable:      bootableVolume,
		}

		volume, err := client.NewVolume(volumeConfig)
		if err != nil {
			utility.Error("Creating the volume failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": volume.ID, "Name": volume.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a volume called %s with ID %s\n", utility.Green(volume.Name), utility.Green(volume.ID))
		}
	},
}
