# go-tmpl
[![Go Report Card](https://goreportcard.com/badge/github.com/pegaz/go-tmpl)](https://goreportcard.com/report/github.com/pegaz/go-tmpl) [![cover.run](https://cover.run/go/github.com/pegaz/go-tmpl.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fpegaz%2Fgo-tmpl)

**go-tmpl** is a tool written in [Go](https://golang.org) created to generate textual output with templates and CSV file as a source for a data.

**go-tmpl** uses [template engine](https://golang.org/pkg/text/template/) from a golang's standard library.

## Installation

To install `go-tmpl` run the following:

`go get github.com/pegaz/go-tmpl`

Then compile it from within the project's directory:

`go build`

## Usage

To create new _workspace_ in current directory one may use:
`go init workspace -n <workspace_name>`

This command creates `workspace directory tree` with example configuration file and several directories:
    * `output/` - directory where all generated files will be stored
    * `data/` - directory where CSV file(s) needs to be stored
    * `templates/` - directory where all templates need to be placed

To generate output files for a given workspace use:

`go generate -n <workspace_name> [-c <configuration_file>]`

If `-c <configuration_file>` parameter is omitted default configuration file `workspace.toml` will be used.

## Example

1. Create workspace:
    
    `$ go-tmpl init workspace -n workspace_name`
 
2. Apply desired changes to the *workspace* configuration file located by default in **workspace.toml** file in the main *workspace* directory
3. Put CSV file into `data/` directory
4. Prepare templates and put them into a *template/* directory
5. Generate output files with:

    `$ go-tmpl generate -n workspace_name`
    
Optionally you can create and use additional configuration files inside a main *workspace* directory (`-c` switch when using *generate* subcommand).

## Configuration file

**go-tmpl** may use one or several config files (as needed). It uses [TOML](https://github.com/toml-lang/toml) configuration file format.

Example of configuration file used by a **go-tmpl**:

    csv_data = "customer.csv"
    csv_delimiter = ","

    template_column_name = "template_name"
    output_column_name = "hostname"

    [vars]
    customer = "ACME"

**Important:** all paths in configuration file are relative to the workspace root.

`csv_data` - name of the CSV data file (first row has to be a header!).

`csv_delimiter` - delimiter sign use to separate fields in CSV file.

`template_column_name` - column name in CSV file where name of the template can be found.

`output_column_name` = column name in CSV file where output filename can be found.

`[vars]` [section](https://github.com/toml-lang/toml#table) may be used to define global variables which then can be used by a templates.

## Hints

There is no need to use full filename of a `template file` in data file. The only relevant part of it is name _without_ file extension (in this case `.tpl`). On the other hand it is assumed that all filenames located in `templates/` directory should end with `.tpl`.
