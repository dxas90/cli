package cmd

import (
	"fmt"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var kubernetesListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo kubernetes ls -o custom -f "ID: Name"`,
	Short:   "List all Kubernetes clusters",
	Long: `List all Kubernetes clusters.
If you wish to use a custom format, the available fields are:

	* ID
	* Name
	* Node
	* Size
	* Status`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubernetesClusters, err := client.ListKubernetesClusters()
		if err != nil {
			utility.Error("Listing Kubernetes clusters failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, cluster := range kubernetesClusters.Items {
			ow.StartLine()

			ow.AppendData("ID", cluster.ID)
			ow.AppendData("Name", cluster.Name)
			ow.AppendData("Node", strconv.Itoa(cluster.NumTargetNode))
			ow.AppendData("Size", cluster.TargetNodeSize)
			ow.AppendData("Status", fmt.Sprintf("%s", utility.ColorStatus(cluster.Status)))

			if outputFormat == "json" || outputFormat == "custom" {
				ow.AppendData("Status", cluster.Status)
			}

		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
