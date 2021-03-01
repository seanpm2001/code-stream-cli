/*
Package cmd Copyright © 2021 Sam McGeown <smcgeown@vmware.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var state string
var exportPath string
var importPath string
var export bool
var form bool

// getPipelineCmd represents the pipeline command
var getPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()
		response, err := getPipelines(id, name, project, export, exportPath)
		if err != nil {
			fmt.Print("Unable to get Code Stream Pipelines: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			fmt.Println("No results found")
		} else if resultCount == 1 {
			if form {
				// Get the input form
				var inputs = response[0].Input
				PrettyPrint(inputs)
			} else {
				// Print the single result
				PrettyPrint(response[0])
			}
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name, c.Project})
			}
			table.Render()
		}
	},
}

// updatePipelineCmd represents the pipeline update command
var updatePipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Update a pipeline",
	Long: `A longer description that spans multiple lines
	Enable/Disable/Release:
	cs-cli update pipeline --id d0185f04-2e87-4f3c-b6d7-ee58abba3e92 --state enabled/disabled/released
	Update from YAML
	cs-cli update pipeline --importPath "/Users/sammcgeown/Desktop/pipelines/SSH Exports.yaml"
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if state != "" {
			switch strings.ToUpper(state) {
			case "ENABLED", "DISABLED", "RELEASED":
				// Valid states
				return nil
			}
			return errors.New("--state is not valid, must be ENABLED, DISABLED or RELEASED")
		}
		if export {

		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()
		if state != "" {
			response, err := patchPipeline(id, `{"state":"`+state+`"}`)
			if err != nil {
				fmt.Print("Unable to update Code Stream Pipeline: ", err)
			}
			fmt.Println("Setting pipeline " + response.Name + " to " + state)
		}

		if importPath != "" {
			if importYaml(importPath, "apply") {
				fmt.Println("Imported successfully, pipeline updated.")
			}
		}
	},
}

// createPipelineCmd represents the pipeline create command
var createPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Create a pipeline",
	Long: `Create a pipeline by importing a YAML specification.
	
	Create from YAML
	  cs-cli create pipeline --importPath "/Users/sammcgeown/Desktop/pipelines/SSH Exports.yaml"
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()
		if importPath != "" {
			if importYaml(importPath, "create") {
				fmt.Println("Imported successfully, pipeline created.")
			}
		}
	},
}

// deletePipelineCmd represents the delete pipeline command
var deletePipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Delete a Pipeline",
	Long: `Delete a Pipeline with a specific ID
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()

		response, err := deletePipeline(id)
		if err != nil {
			log.Fatalln("Delete Pipeline failed:", err)
		}
		fmt.Println("Pipeline with id " + response.ID + " deleted")

	},
}

func init() {
	// Get
	getCmd.AddCommand(getPipelineCmd)
	getPipelineCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the pipeline to list executions for")
	getPipelineCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to list")
	getPipelineCmd.Flags().StringVarP(&project, "project", "p", "", "List pipeline in project")
	getPipelineCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	getPipelineCmd.Flags().BoolVarP(&export, "export", "e", false, "Export pipeline")
	getPipelineCmd.Flags().BoolVarP(&form, "form", "f", false, "Get pipeline inputs")

	// Create
	createCmd.AddCommand(createPipelineCmd)
	createPipelineCmd.Flags().StringVarP(&importPath, "importPath", "", "", "YAML configuration file to import")
	createPipelineCmd.MarkFlagRequired("importPath")
	// Update
	updateCmd.AddCommand(updatePipelineCmd)
	updatePipelineCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the pipeline to list")
	updatePipelineCmd.Flags().StringVarP(&importPath, "importPath", "", "", "Configuration file to import")
	updatePipelineCmd.Flags().StringVarP(&state, "state", "s", "", "Set the state of the pipeline (ENABLED|DISABLED|RELEASED")
	// Delete
	deleteCmd.AddCommand(deletePipelineCmd)
	deletePipelineCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the Pipeline to delete")
	deletePipelineCmd.MarkFlagRequired("id")

}
