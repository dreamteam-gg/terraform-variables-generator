package generator

import (
	"bufio"
	"html/template"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/alexandrst88/terraform-variables-generator/pkg/utils"
)

const (
	varPrefix = "var."
)

var replacer *strings.Replacer

var varTemplate = template.Must(template.New("var_file").Parse(`{{range .}}
variable "{{ . }}" {
  description = ""
}
{{end}}`))

func init() {
	replacer = strings.NewReplacer(":", ".",
		"]", "",
		"}", "",
		"{", "",
		"\"", "",
		")", "",
		"(", "",
		"[", "",
		",", "",
		"var.", "",
		" ", "",
	)
}

// GenerateVars will write generated vars to file
func GenerateVars(dstFile string) {
	dir, err := os.Getwd()
	utils.CheckError(err)

	tfFiles, err := utils.GetAllFiles(dir, tfFileExt)
	utils.CheckError(err)
	if len(tfFiles) == 0 {
		log.Warn("No terraform files to proceed, exiting")
		return
	}

	var wg sync.WaitGroup
	messages := make(chan string)
	wg.Add(len(tfFiles))
	t := &terraformVars{}

	for _, file := range tfFiles {
		go func(file string) {
			defer wg.Done()
			fileHandle, _ := os.Open(file)
			defer fileHandle.Close()
			fileScanner := bufio.NewScanner(fileHandle)
			for fileScanner.Scan() {
				messages <- fileScanner.Text()
			}
		}(file)
	}
	go func() {
		for text := range messages {
			t.matchVarPref(text, varPrefix)
		}
	}()
	wg.Wait()
	f, err := os.Create(dstFile)
	utils.CheckError(err)

	t.sortVars()
	err = varTemplate.Execute(f, t.Variables)
	utils.CheckError(err)
	log.Infof("Variables are generated to %q file", dstFile)
}
