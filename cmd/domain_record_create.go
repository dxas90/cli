package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

const (
	// DNSRecordTypeA represents an A record
	DNSRecordTypeA = "a"

	// DNSRecordTypeCName represents an CNAME record
	DNSRecordTypeCName = "cname"

	// DNSRecordTypeMX represents an MX record
	DNSRecordTypeMX = "mx"

	// DNSRecordTypeTXT represents an TXT record
	DNSRecordTypeTXT = "txt"
)

var recordName, recordType, recordValue string
var recordTTL, recordPriority int

var domainRecordCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new domain record",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo domain record create DOMAIN/DOMAIN_ID [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		domain, err := client.FindDNSDomain(args[0])
		if err != nil {
			utility.Error("Unable to find the domain for your search %s", err)
			os.Exit(1)
		}

		newRecordConfig := &civogo.DNSRecordConfig{
			Name:     recordName,
			Value:    recordValue,
			TTL:      recordTTL,
			Priority: recordPriority,
		}

		if recordType == "a" || recordType == "alias" {
			newRecordConfig.Type = DNSRecordTypeA
		}

		if recordType == "cname" || recordType == "canonical" {
			newRecordConfig.Type = DNSRecordTypeCName
		}

		if recordType == "mx" || recordType == "mail" {
			newRecordConfig.Type = DNSRecordTypeMX
		}

		if recordType == "txt" || recordType == "text" {
			newRecordConfig.Type = DNSRecordTypeTXT
		}

		record, err := client.CreateDNSRecord(domain.ID, newRecordConfig)
		if err != nil {
			utility.Error("Unable to find domain for your search %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": domain.ID, "Name": domain.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created %s record %s for %s with a TTL of %s seconds and with a priority of %s with ID %s", utility.Green(string(record.Type)), utility.Green(record.Name), utility.Green(domain.Name), utility.Green(strconv.Itoa(record.TTL)), utility.Green(strconv.Itoa(record.Priority)), utility.Green(record.ID))
		}
	},
}
