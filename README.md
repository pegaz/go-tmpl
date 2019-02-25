# go-tmpl
### Textual template generator written in Go

It enables to generate text files using templates (using text/template package) and CSV file as a source for a data.

**go-tmpl** uses [template engine](https://golang.org/pkg/text/template/) from a **golang** standard library.

Standard use of this tool:
1. Create workspace:
    
    `$ go-tmpl init workspace -n workspace_name`
 
2. Apply desired changes to the *workspace* configuration file located by default in **workspace.toml** file in the main *workspace* directory
3. Put CSV file into *data/* directory
4. Prepare templates and put them into a *template/* directory
5. Generate output files with:

    `$ go-tmpl generate -w workspace_name`
    
Optionally you can create and use additional configuration files inside a main *workspace* directory (-c switch when using *generate* subcommand).

**go-tmpl** uses configuration files in [TOML format](https://github.com/toml-lang/toml)
Example of configuration file used by a **go-tmpl**:
    
    # CSV data filename, it should be placed in data directory inside of a given workspace
    csv_data = "customer.csv"
    # delimiter used in CSV file as a field separator
    csv_delimiter = ","

    template_column_name = "template_name"
    output_column_name = "hostname"
    output_in_single_file = false
    output_filename = "output.cfg"
		
    [vars]
    # customer name
    customer = "ACME"

[Table](https://github.com/toml-lang/toml#table) **[vars]** may be used to define global variables which then can be used by a templates.
