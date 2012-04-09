package templater

import (
	"../log"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var t *template.Template

func Get(tmpl string) (result *template.Template) {
	result = t.Lookup(tmpl)
	if result == nil {
		panic("Unable to find template \"" + tmpl + "\"")
	}
	return
}

func LoadTemplates(dir string) {
	base := filepath.Clean(dir)
	t = template.New("")
	t.Parse("")

	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.TrimLeft(filepath.Ext(info.Name()), ".")

		file, err := filepath.Rel(base, path)
		if err != nil {
			panic(err)
		}

		index := strings.Index(file, ".")
		if index == -1 {
			log.Warning("\nWarning: No file type in " + file)
			return nil
		}
		name := file[0:index]

		switch ext {
		case "tmpl":
			b, err := ioutil.ReadFile(path)
			if err != nil {
				log.Error("\nError: could not load template file:" + path)
				panic(err)
			}
			s := string(b)
			a := t.New(name)
			a.Parse(s)
		default:
			log.Warning("\nWarning: Unknown Template Type: ." + ext)
		}

		return nil
	})
}
