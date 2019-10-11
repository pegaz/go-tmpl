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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pegaz/go-tmpl/text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultCsvDelimiter = ','
	DefaultCsvDataFile  = "data.csv"
)

var (
	workspaceName      string
	workspaceConfig    string
	outputColumnName   string
	templateColumnName string
	csvFilename        string
	csvDelimiter       rune
	missingKey         string
	overrideOutput     bool

	fileCounter int64
	outputFiles []string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate output from templates",

	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		// set some default config parameters
		setDefaults()

		err := initConfig()
		if err != nil {
			return err
		}

		csvReader, err := os.Open(csvFilename)
		if err != nil {
			return err
		}
		defer csvReader.Close()

		data, err := text.ReadCSV(csvReader, csvDelimiter)
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

		if overrideOutput == true {
			// Delete content of the output's directory
			err = pruneDirContent(rootDir + "/" + workspaceName + directories["output"])
			if err != nil {
				return err
			}
		}

		for _, d := range data {
			var outputFile io.WriteCloser

			templateFilename, ok := d[templateColumnName]
			if !ok {
				fmt.Printf("couldn't find '%s' column in data provided", templateColumnName)
				continue
			}
			outputFilename, ok := d[outputColumnName]
			if !ok {
				fmt.Printf("couldn't find '%s' column in data provided", outputColumnName)
				continue
			}

			templatePath := rootDir + "/" + workspaceName + directories["templates"] + "/" + templateFilename
			if templatePath[len(templatePath)-4:] != ".tpl" {
				templatePath = templatePath + ".tpl"
			}
			templateReader, err := os.Open(templatePath)
			if err != nil {
				return err
			}
			defer templateReader.Close()

			var tmpl *text.Template

			tmpl, err = text.NewTemplate(d, templateFilename, templateReader)
			if err != nil {
				return err
			}

			// Global variables defined in configuration file for a workspace goes to Template
			tmpl.SetGlobalVars(globalVars)
			tmpl.SetStrict(missingKey)

			var flags int
			if contains(outputFiles, outputFilename) {
				flags = os.O_APPEND
			} else {
				// Check if given file already exist and if so don't generate output for it
				_, err = os.Stat(rootDir + "/" + workspaceName + directories["output"] + "/" + outputFilename + ".txt")
				if err == nil {
					continue
				}

				flags = os.O_CREATE
				outputFiles = append(outputFiles, outputFilename)
			}

			outputFile, err = os.OpenFile(rootDir+"/"+workspaceName+directories["output"]+"/"+outputFilename+".txt", flags, 0644)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			err = tmpl.Execute(outputFile)
			if err != nil {
				fmt.Printf("error generating file from template: %s", err)
				return err
			}
		}

		if len(outputFiles) > 0 {
			for _, outputFile := range outputFiles {
				fmt.Printf("* %s\n", outputFile)
			}
			fmt.Println()
			fmt.Printf("Succesfully generated %d output files", len(outputFiles))
		} else {
			fmt.Print("Nothing to do")
		}

		return nil
	},
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func init() {
	generateCmd.Flags().BoolVarP(&overrideOutput, "force", "f", false, "override the content of output directory")
	generateCmd.Flags().StringVarP(&workspaceName, "name", "n", "", "workspace to generate files for")
	generateCmd.MarkFlagRequired("workspace")

	generateCmd.Flags().StringVarP(&workspaceConfig, "config", "c", "workspace.toml", "configuration file to use generator for")

	rootCmd.AddCommand(generateCmd)
}

func setDefaults() {
	viper.SetDefault("csv_data", DefaultCsvDataFile)
	viper.SetDefault("csv_delimiter", DefaultCsvDelimiter)
	viper.SetDefault("missingkey", "invalid")
	viper.SetDefault("override_output", "false")
}

func initConfig() error {
	var err error
	var file *os.File

	viper.SetConfigType("toml")

	file, err = os.Open(rootDir + "/" + workspaceName + "/" + workspaceConfig)
	if err != nil {
		return err
	}

	err = viper.ReadConfig(file)
	if err != nil {
		return err
	}

	// mandatory fields
	if viper.IsSet("template_column_name") == false ||
		viper.IsSet("output_column_name") == false {
		return fmt.Errorf("some mandatory configuration parametrs in config file are missing")
	}

	outputColumnName = viper.GetString("output_column_name")
	templateColumnName = viper.GetString("template_column_name")
	csvDelimiter = rune(viper.GetString("csv_delimiter")[0])
	if viper.GetString("missing_key") == "invalid" || viper.GetString("missing_key") == "zero" || viper.GetString("missing_key") == "error" {
		missingKey = viper.GetString("missing_key")
	} else if viper.IsSet("missing_key") {
		return fmt.Errorf("invalid value for 'missing_key' value in configuration file, got: %s", viper.GetString("missing_key"))
	}
	csvFilename = rootDir + "/" + workspaceName + directories["data"] + "/" + viper.GetString("csv_data")

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
