/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getCustomIntegrationCmd represents the customintegration command
var getCustomIntegrationCmd = &cobra.Command{
	Use:   "customintegration",
	Short: "Get Custom Integrations",
	Long: `Get Code Stream Custom Integrations by name, project or by id - e.g:

Get by ID
	cs-cli get customintegration --id 6b7936d3-a19d-4298-897a-65e9dc6620c8
	
Get by Name
	cs-cli get customintegration --name my-customintegration
	
Get by Project
	cs-cli get customintegration --project production`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}
		response, err := getCustomIntegration(id, name)
		if err != nil {
			log.Errorln("Unable to get Code Stream CustomIntegrations: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Infoln("No results found")
		} else if resultCount == 1 {
			// Print the single result
			//if export {
			//exportCustomIntegration(response[0], exportFile)
			//}
			PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Status", "Description"})
			for _, c := range response {
				//if export {
				//exportCustomIntegration(c, exportFile)
				//}
				table.Append([]string{c.ID, c.Name, c.Status, c.Description})
			}
			table.Render()
		}
	},
}

// // getCustomIntegrationCmd represents the customintegration command
// var createCustomIntegrationCmd = &cobra.Command{
// 	Use:   "customintegration",
// 	Short: "A brief description of your command",
// 	Long:  ``,
// 	Run: func(cmd *cobra.Command, args []string) {
// 				if err := ensureTargetConnection(); err != nil {
// 	log.Fatalln(err)
// }

// 		if importFile != "" { // If we are importing a file
// 			customintegrations := importCustomIntegrations(importFile)
// 			for _, value := range customintegrations {
// 				if project != "" { // If the project is specified update the object
// 					value.Project = project
// 				}
// 				createResponse, err := createCustomIntegration(value.Name, value.Description, value.Type, value.Project, value.Value)
// 				if err != nil {
// 					log.Errorln("Unable to create Code Stream CustomIntegration: ", err)
// 				} else {
// 					log.Infoln("Created customintegration", createResponse.Name, "in", createResponse.Project)
// 				}
// 			}
// 		} else {
// 			createResponse, err := createCustomIntegration(name, description, typename, project, value)
// 			if err != nil {
// 				log.Errorln("Unable to create Code Stream CustomIntegration: ", err)
// 			}
// 			PrettyPrint(createResponse)
// 		}
// 	},
// }

// // updateCustomIntegrationCmd represents the customintegration command
// var updateCustomIntegrationCmd = &cobra.Command{
// 	Use:   "customintegration",
// 	Short: "A brief description of your command",
// 	Long:  ``,
// 	Run: func(cmd *cobra.Command, args []string) {
// 				if err := ensureTargetConnection(); err != nil {
// 	log.Fatalln(err)
// }

// 		if importFile != "" { // If we are importing a file
// 			customintegrations := importCustomIntegrations(importFile)
// 			for _, value := range customintegrations {
// 				exisitingCustomIntegration, err := getCustomIntegration("", value.Name, value.Project)
// 				if err != nil {
// 					log.Infoln("Update failed - unable to find existing Code Stream CustomIntegration", value.Name, "in", value.Project)
// 				} else {
// 					_, err := updateCustomIntegration(exisitingCustomIntegration[0].ID, value.Name, value.Description, value.Type, value.Value)
// 					if err != nil {
// 						log.Infoln("Unable to update Code Stream CustomIntegration: ", err)
// 					} else {
// 						log.Infoln("Updated customintegration", value.Name)
// 					}
// 				}
// 			}
// 		} else { // Else we are updating using flags
// 			updateResponse, err := updateCustomIntegration(id, name, description, typename, value)
// 			if err != nil {
// 				log.Infoln("Unable to update Code Stream CustomIntegration: ", err)
// 			}
// 			log.Infoln("Updated customintegration", updateResponse.Name)
// 		}
// 	},
// }

// // deleteCustomIntegrationCmd represents the executions command
// var deleteCustomIntegrationCmd = &cobra.Command{
// 	Use:   "customintegration",
// 	Short: "A brief description of your command",
// 	Long: `A longer description that spans multiple lines and likely contains examples
// and usage of using your command. For example:

// Cobra is a CLI library for Go that empowers applications.
// This application is a tool to generate the needed files
// to quickly create a Cobra application.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 				if err := ensureTargetConnection(); err != nil {
// 	log.Fatalln(err)
// }

// 		response, err := deleteCustomIntegration(id)
// 		if err != nil {
// 			log.Infoln("Unable to delete customintegration: ", err)
// 		}
// 		log.Infoln("CustomIntegration with id " + response.ID + " deleted")
// 	},
// }

func init() {
	// Get CustomIntegration
	getCmd.AddCommand(getCustomIntegrationCmd)
	getCustomIntegrationCmd.Flags().StringVarP(&name, "name", "n", "", "List customintegration with name")
	getCustomIntegrationCmd.Flags().StringVarP(&id, "id", "i", "", "List customintegrations by id")
	getCustomIntegrationCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	// // Create CustomIntegration
	// createCmd.AddCommand(createCustomIntegrationCmd)
	// createCustomIntegrationCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the customintegration to create")
	// createCustomIntegrationCmd.Flags().StringVarP(&typename, "type", "t", "", "The type of the customintegration to create REGULAR|SECRET|RESTRICTED")
	// createCustomIntegrationCmd.Flags().StringVarP(&project, "project", "p", "", "The project in which to create the customintegration")
	// createCustomIntegrationCmd.Flags().StringVarP(&value, "value", "v", "", "The value of the customintegration to create")
	// createCustomIntegrationCmd.Flags().StringVarP(&description, "description", "d", "", "The description of the customintegration to create")
	// createCustomIntegrationCmd.Flags().StringVarP(&importFile, "importfile", "i", "", "Path to a YAML file with the customintegrations to import")

	// // Update CustomIntegration
	// updateCmd.AddCommand(updateCustomIntegrationCmd)
	// updateCustomIntegrationCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the customintegration to update")
	// updateCustomIntegrationCmd.Flags().StringVarP(&name, "name", "n", "", "Update the name of the customintegration")
	// updateCustomIntegrationCmd.Flags().StringVarP(&typename, "type", "t", "", "Update the type of the customintegration REGULAR|SECRET|RESTRICTED")
	// updateCustomIntegrationCmd.Flags().StringVarP(&value, "value", "v", "", "Update the value of the customintegration ")
	// updateCustomIntegrationCmd.Flags().StringVarP(&description, "description", "d", "", "Update the description of the customintegration")
	// updateCustomIntegrationCmd.Flags().StringVarP(&importFile, "importfile", "", "", "Path to a YAML file with the customintegrations to import")
	// //updateCustomIntegrationCmd.MarkFlagRequired("id")
	// // Delete CustomIntegration
	// deleteCmd.AddCommand(deleteCustomIntegrationCmd)
	// deleteCustomIntegrationCmd.Flags().StringVarP(&id, "id", "i", "", "Delete customintegration by id")
	// deleteCustomIntegrationCmd.MarkFlagRequired("id")
}
