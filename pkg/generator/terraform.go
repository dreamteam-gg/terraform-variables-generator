package generator

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/hcl"
	log "github.com/sirupsen/logrus"

	"github.com/alexandrst88/terraform-variables-generator/pkg/utils"
)

const (
	tfFileExt = "*.tf"
)

type terraformVars struct {
	Variables []string
}

type terraformModule struct {
	Key  string `json:"Key"`
	Dir  string `json:"Dir"`
	Root string `json:"Root"`
	// Path is full path of module
	Path string
	// Name is parsed full name of module resource
	Name string
	// Outputs is list of module outputs
	Outputs []terraformOutput `hcl:"output,block"`
}

type terraformOutput struct {
	Name  string `hcl:"name,label"`
	Value string `hcl:"value"`
}

type terraformModules struct {
	Modules []terraformModule `json:"Modules"`
}

func (t *terraformVars) matchVarPref(row, varPrefix string) {
	if strings.Contains(row, varPrefix) {
		pattern := regexp.MustCompile(`var.([a-z?A-Z?0-9?_][a-z?A-Z?0-9?_?-]*)`)
		match := pattern.FindAllStringSubmatch(row, -1)
		for _, m := range match {
			res := replacer.Replace(m[0])
			if !utils.ContainsElement(t.Variables, res) {
				t.Variables = append(t.Variables, res)
			}
		}
	}
}

func (t *terraformModule) parseOutputs() {
	// module can actually be a submodule, join root to path
	files, err := utils.GetAllFiles(path.Join(t.Path, t.Root), tfFileExt)
	utils.CheckError(err)

	if len(files) == 0 {
		log.Warnf("No terraform files to parse in module %s", t.Name)
		return
	}

	var wg sync.WaitGroup
	// messages := make(chan string)
	wg.Add(len(files))

	for _, file := range files {
		go func(file string) {
			defer wg.Done()

			fileHandle, err := os.Open(file)
			utils.CheckError(err)
			defer fileHandle.Close()

			fileContent, err := ioutil.ReadAll(fileHandle)
			utils.CheckError(err)

			var config interface{}

			err = hcl.Unmarshal(fileContent, &config)
			utils.CheckError(err)

			if outputs, ok := config.(map[string]interface{})["output"]; ok {
				for _, o := range outputs.([]map[string]interface{}) {
					// output is a map with one key
					name := utils.GetMapKeys(o)[0]

					t.Outputs = append(t.Outputs, terraformOutput{
						Name:  t.Name + "_" + name,
						Value: "${" + t.Name + "." + name + "}",
					})
				}
			}
		}(file)
	}
	wg.Wait()
	t.sortOutputs()
}

func (t *terraformVars) sortVars() {
	sort.Strings(t.Variables)
}

func (t *terraformModule) sortOutputs() {
	sort.Slice(t.Outputs, func(i, j int) bool {
		return t.Outputs[i].Name < t.Outputs[j].Name
	})
}
