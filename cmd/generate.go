// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/viper"
	"gitlab.com/pawel.potrykus/go-tmpl/templates"
)

var workspace string
var outputFilename string

//workspaceConfig is a config file name and it resides inside a root directory of a specific workspace
var workspaceConfig string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate output from templates",

	RunE: func(cmd *cobra.Command, args []string) error {
		err := initConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		csvDelimiter := rune(viper.GetString("csv_delimiter")[0])
		csvFile := workspace + directories["data"] + "/" + viper.GetString("csv_data")

		data, err := templates.ReadCSV(csvFile, csvDelimiter)
		if err != nil {
			return err
		}

		var outputs []*templates.Template
		for _, d := range data {
			var filename string
			var singleFileOutput bool

			templateFilename, ok := d[viper.GetString("template_column_name")]
			if !ok {
				return fmt.Errorf("couldn't find '%s' column in data provided", viper.GetString("template_column_name"))
			}

			templatePath := workspace + directories["templates"] + "/" + templateFilename

			var tmpl *templates.Template

			if outputFilename == "" {
				filename = d[viper.GetString("output_column_name")]
			} else {
				filename = outputFilename
				singleFileOutput = true
			}

			tmpl, err = templates.New(d, templatePath, workspace+directories["output"]+"/", filename, singleFileOutput)
			if err != nil {
				return err
			}

			// Global variables defined in configuration file for a workspace goes to Template
			vars := viper.Sub("vars")
			varsName := vars.AllKeys()
			m := make(map[string]string)

			for _, name := range varsName {
				m[name] = vars.GetString(name)
			}
			tmpl.SetGlobalVars(m)

			outputs = append(outputs, tmpl)

		}

		return templates.Execute(outputs)
	},
}

func init() {
	generateCmd.Flags().StringVarP(&workspace, "workspace", "w", "", "workspace to generate files for")
	generateCmd.MarkFlagRequired("workspace")

	generateCmd.Flags().StringVarP(&workspaceConfig, "config", "c", "workspace.toml", "configuration file to use generator for")

	rootCmd.AddCommand(generateCmd)
}

func initConfig() error {
	var err error
	var file *os.File

	viper.SetConfigType("toml")

	file, err = os.Open(workspace + "/" + workspaceConfig)
	if err != nil {
		return err
	}

	err = viper.ReadConfig(file)
	if err != nil {
		return err
	}

	if viper.IsSet("csv_data") == false ||
		viper.IsSet("csv_delimiter") == false ||
		viper.IsSet("template_column_name") == false ||
		viper.IsSet("output_column_name") == false ||
		viper.IsSet("output_in_single_file") == false ||
		viper.IsSet("output_filename") == false {
		return fmt.Errorf("some configuration parametrs in config file are missing")
	}

	if viper.GetBool("output_in_single_file") != false {
		outputFilename = viper.GetString("output_filename")
	}

	return err
}

// func openFile(path string, output string) (*os.File, error) {
// 	var file *os.File

// 	switch output {
// 	case "stdout":
// 		file = os.Stdout
// 	default:
// 		_, err := os.Stat(path + output)
// 		if err == nil {
// 			os.Remove(path + output)
// 		}
// 		return os.OpenFile(path+output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
// 	}
// 	return file, nil
// }
