package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesNewVersion string

var kubernetesUpgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"change", "modify"},
	Example: "civo kubernetes upgrade CLUSTER_NAME --version VERSION",
	Short:   "Upgrade/rescale a Kubernetes cluster",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("Finding the kubernetes cluster failed with %s", err)
			os.Exit(1)
		}

		configKubernetes := &civogo.KubernetesClusterConfig{
			KubernetesVersion: kubernetesNewVersion,
		}

		kubernetesCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
		if err != nil {
			utility.Error("Upgrading the kubernetes cluster failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": kubernetesCluster.ID, "Name": kubernetesCluster.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The kubernetes cluster %s was upgraded to %s\n", utility.Green(kubernetesCluster.Name), utility.Green(kubernetesNewVersion))
		}
	},
}
