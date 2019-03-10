// Copyright © 2019 Pawel Potrykus <pawel.potrykus@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var name string
var directories = map[string]string{
	"templates": "/templates",
	"data":      "/data",
	"output":    "/output",
}

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "init workspace directory",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := createWorkspace(name)
		if err != nil {
			return err
		}

		fmt.Printf("workspace \"%s\" successfully generated\n", name)

		return nil
	},
}

func init() {
	workspaceCmd.Flags().StringVarP(&name, "name", "n", "", "name of a workspace to be created")
	workspaceCmd.MarkFlagRequired("name")

	initCmd.AddCommand(workspaceCmd)
}

func createWorkspace(name string) error {
	var err error

	files := map[string][]byte{
		name + "/workspace.toml": []byte(`# CSV data filename, it should be placed in data directory inside of a given workspace
#csv_data = "data.csv"
# delimiter used in CSV file as a field separator
#csv_delimiter = ","

template_column_name = "router"
output_column_name = "hostname"

[vars]
# custom vars to use them inside of templates should be placed here
`),
		name + "/README.md": []byte(`## Root of a workspace, workspace.toml configurations file should be placed here
		`),
		name + "/data/README.md": []byte(`## Data in CSV format should be placed here
		`),
		name + "/templates/README.md": []byte(`## Source of all templates used by a workspace
		`),
		name + "/output/README.md": []byte(`## Place where all the generated files will be placed
		`),
	}

	_, err = os.Stat(name)
	if err == nil {
		return fmt.Errorf("given directory %s already exists in the current directory", name)
	}

	err = os.Mkdir(name, 0755)
	if err != nil {
		return err
	}

	for _, dir := range directories {
		err = os.Mkdir(name+dir, 0755)
		if err != nil {
			os.RemoveAll(name)
			return err
		}
	}

	for filename, content := range files {
		var file *os.File

		file, err = os.Create(filename)
		if err != nil {
			os.RemoveAll(name)
			return err
		}

		_, err = file.Write(content)
		if err != nil {
			os.RemoveAll(name)
			return err
		}
		file.Close()
	}

	return nil
}
