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

	modules       bool
	outputsFile   string
	modulesFilter string
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

	cmd.PersistentFlags().BoolVar(&modules, "module-outputs", true, "generate module outputs")
	cmd.PersistentFlags().StringVar(&outputsFile, "outputs-file", "./outputs.tf", "path to generated outputs file")
	cmd.PersistentFlags().StringVar(&modulesFilter, "modules-filter", "", "regexp to match modules by name")

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

	if modules {
		if utils.FileExists(outputsFile) {
			utils.UserPromt(outputsFile)
		}
		generator.GenerateModuleOutputs(modulesFilter, outputsFile)
	}
}
