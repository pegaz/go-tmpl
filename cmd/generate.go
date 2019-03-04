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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pegaz/go-tmpl/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	workspaceName      string
	workspaceConfig    string
	outputColumnName   string
	templateColumnName string
	csvFile            string
	csvDelimiter       rune
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate output from templates",

	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		err := initConfig()
		if err != nil {
			return err
		}

		data, err := templates.ReadCSV(csvFile, csvDelimiter)
		if err != nil {
			return err
		}

		// Global variables (defined in configuration file within a [vars] section
		vars := viper.Sub("vars")
		varsName := vars.AllKeys()
		globalVars := make(map[string]string)
		for _, name := range varsName {
			globalVars[name] = vars.GetString(name)
		}

		// Delete content of the output's directory
		err = pruneDirContent(workspaceName + directories["output"])
		if err != nil {
			return err
		}

		for _, d := range data {
			var outputFile io.WriteCloser

			templateFilename, ok := d[templateColumnName]
			if !ok {
				return fmt.Errorf("couldn't find '%s' column in data provided", templateColumnName)
			}

			templatePath := workspaceName + directories["templates"] + "/" + templateFilename
			if templatePath[len(templatePath)-4:] != ".tpl" {
				templatePath = templatePath + ".tpl"
			}
			templateFile, err := os.Open(templatePath)
			if err != nil {
				return err
			}
			defer templateFile.Close()

			var tmpl *templates.Template

			tmpl, err = templates.New(d, templateFilename, templateFile)
			if err != nil {
				return err
			}

			// Global variables defined in configuration file for a workspace goes to Template
			tmpl.SetGlobalVars(globalVars)

			if d[outputColumnName] != "" {
				outputFile, err = os.OpenFile(workspaceName+directories["output"]+"/"+d[outputColumnName]+".txt", os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					return err
				}
				defer outputFile.Close()
			} else {
				fmt.Println("Couldn't determine output filename for entry")
				outputFile = os.Stdout
			}

			err = tmpl.Execute(outputFile)
			if err != nil {
				fmt.Printf("error generating file from template: %s", err)
			}
		}

		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&workspaceName, "name", "n", "", "workspace to generate files for")
	generateCmd.MarkFlagRequired("workspace")

	generateCmd.Flags().StringVarP(&workspaceConfig, "config", "c", "workspace.toml", "configuration file to use generator for")

	rootCmd.AddCommand(generateCmd)
}

func initConfig() error {
	var err error
	var file *os.File

	viper.SetConfigType("toml")

	file, err = os.Open(workspaceName + "/" + workspaceConfig)
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
		viper.IsSet("output_column_name") == false {
		return fmt.Errorf("some configuration parametrs in config file are missing")
	}

	outputColumnName = viper.GetString("output_column_name")
	templateColumnName = viper.GetString("template_column_name")
	csvDelimiter = rune(viper.GetString("csv_delimiter")[0])
	csvFile = workspaceName + directories["data"] + "/" + viper.GetString("csv_data")

	return err
}

func pruneDirContent(dir string) error {
	matches, err := filepath.Glob(dir + "/*")
	if err != nil {
		return err
	}

	for _, filename := range matches {
		if strings.HasSuffix(filename, ".txt") || strings.HasSuffix(filename, ".cfg") {
			err = os.Remove(filename)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
