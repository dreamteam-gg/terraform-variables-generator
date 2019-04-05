package generator

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/alexandrst88/terraform-variables-generator/pkg/utils"
)

const (
	modulesConfig = "./.terraform/modules/modules.json"
)

var outputTemplate = template.Must(template.New("output_file").Parse(`{{ range . }}
output "{{ .Name }}" {
  description = ""
  value       = "{{ .Value }}"
}
{{end}}`))

// GenerateModuleOutputs generates module outputs
func GenerateModuleOutputs(filter string, dstFile string) {
	dir, err := os.Getwd()
	utils.CheckError(err)

	modules := getModules(dir)
	if len(modules) == 0 {
		log.Warn("No terraform modules config found, run 'terraform init'. Exiting")
		return
	}

	filterRegex := regexp.MustCompile(filter)

	var outputs strings.Builder

	for _, m := range modules {
		if !matchModule(m, filterRegex) {
			continue
		}

		m.parseOutputs()
		var buf bytes.Buffer
		err = outputTemplate.Execute(&buf, m.Outputs)
		utils.CheckError(err)

		outputs.WriteString(buf.String())
		outputs.WriteString("\n\n")
	}

	f, err := os.Create(dstFile)
	utils.CheckError(err)
	defer f.Close()

	_, err = f.Write([]byte(outputs.String()))
	utils.CheckError(err)
}

func getModules(dir string) []terraformModule {
	modulesFile := path.Join(dir, modulesConfig)
	if !utils.FileExists(modulesFile) {
		return nil
	}

	var modules terraformModules

	modulesJSON, err := os.Open(modulesFile)
	utils.CheckError(err)
	defer modulesJSON.Close()

	byteValue, _ := ioutil.ReadAll(modulesJSON)
	err = json.Unmarshal(byteValue, &modules)
	utils.CheckError(err)

	for i, m := range modules.Modules {
		// remove leading 1. and split after comma
		modules.Modules[i].Name = strings.Split(m.Key[2:], ";")[0]
		modules.Modules[i].Path = path.Join(dir, m.Dir)
	}

	return modules.Modules
}

func matchModule(m terraformModule, filter *regexp.Regexp) bool {
	// submodules contain | in their name, skip them
	if strings.Contains(m.Key, "|") {
		return false
	}

	if !filter.Match([]byte(m.Name)) {
		return false
	}

	return true
}
