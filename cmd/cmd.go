package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alexandrst88/terraform-variables-generator/pkg/generator"
	"github.com/alexandrst88/terraform-variables-generator/pkg/utils"
)

var (
	generatorVersion string

	vars     bool
	varsFile string
)

// Execute will run main logic
func Execute(version string) {
	generatorVersion = version

	cmd := &cobra.Command{
		Use:     "generator",
		Short:   "CLI for generating terraform variables",
		Example: "  terraform-variable-generator",
		Version: generatorVersion,
		Run:     runGenerator,
	}

	cmd.PersistentFlags().BoolVar(&vars, "vars", true, "generate variables")
	cmd.PersistentFlags().StringVar(&varsFile, "vars-file", "./variables.tf", "path to generated variables file")

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func runGenerator(cmd *cobra.Command, args []string) {
	if vars {
		if utils.FileExists(varsFile) {
			utils.UserPromt(varsFile)
		}

		generator.GenerateVars(varsFile)
	}

		generator.GenerateVars(tfFiles, varsFile)
	}
}
