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

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init workspace directory",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := createWorkspace(name)
		if err != nil {
			return err
		}

		fmt.Printf("workspace \"%s\" successfully generated on \"%s\"\n", name, rootDir)

		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&name, "name", "n", "", "name of a workspace to be created")
	initCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(initCmd)
}

func createWorkspace(name string) error {
	var err error

	files := map[string][]byte{
		rootDir + "/" + name + "/workspace.toml": []byte(`# CSV data filename, it should be placed in data directory inside of a given workspace
#csv_data = "data.csv"
# delimiter used in CSV file as a field separator
#csv_delimiter = ","
# how to behave when no key is found in CSV file
# zero - nothing will be print in place of variable
# error - error will be returned when no value will be found
# invalid - '<no value>' will be print in place of variable
#missing_key = "invalid"
# by default output folder content won't be overriden
#override_output = false

template_column_name = "router"
output_column_name = "hostname"

[vars]
# custom vars to use them inside of templates should be placed here
`),
		rootDir + "/" + name + "/README.md": []byte(`## Root of a workspace, workspace.toml configurations file should be placed here
		`),
		rootDir + "/" + name + "/data/README.md": []byte(`## Data in CSV format should be placed here
		`),
		rootDir + "/" + name + "/templates/README.md": []byte(`## Source of all templates used by a workspace
		`),
		rootDir + "/" + name + "/output/README.md": []byte(`## Place where all the generated files will be placed
		`),
	}

	_, err = os.Stat(rootDir + "/" + name)
	if err == nil {
		return fmt.Errorf("given directory %s already exists in the current directory", rootDir+name)
	}

	_, err = os.Stat(rootDir)
	if err != nil {
		err = os.Mkdir(rootDir, 0755)
		if err != nil {
			return err
		}
	}

	err = os.Mkdir(rootDir+"/"+name, 0755)
	if err != nil {
		return err
	}

	for _, directory := range directories {
		err = os.Mkdir(rootDir+"/"+name+directory, 0755)
		if err != nil {
			os.RemoveAll(rootDir + "/" + name)
			return err
		}
	}

	for filename, content := range files {
		var file *os.File

		file, err = os.Create(filename)
		if err != nil {
			os.RemoveAll(rootDir + "/" + name)
			return err
		}

		_, err = file.Write(content)
		if err != nil {
			os.RemoveAll(rootDir + "/" + name)
			return err
		}
		file.Close()
	}

	return nil
}
