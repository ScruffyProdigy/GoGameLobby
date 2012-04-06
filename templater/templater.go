package templater

import (
	"text/template"
	"os"
	"path/filepath"
	"strings"
	"../log"
	"io/ioutil"
)

func LoadTemplates(dir string) {
	base := filepath.Clean(dir)
	log.Debug("\nBase:" + base)
	
	filepath.Walk(base,func(path string,info os.FileInfo, err error) error{
		if err != nil {
			panic(err)
		}
		
		if(info.IsDir()){
			return nil
		}
		
		log.Debug("\n\nPath:" + path)
		
		ext := strings.TrimLeft(filepath.Ext(info.Name()),".")
		log.Debug("\nExt:" + ext)
		
		file,err := filepath.Rel(base, path)
		if(err != nil) {
			panic(err)
		}
		log.Debug("\nFile:" + file)
		
		index := strings.Index(file,".")
		if(index == -1){
			log.Warning("\nWarning: No file type in "+file)
			return nil
		}
		name := file[0:index]
		log.Debug("\nName:" + name)
		
		switch(ext){
		case "tmpl":
			b,err := ioutil.ReadFile(path)
			if(err != nil){
				log.Error("\nError: could not load template file:" + path)
				panic(err)
			}
			s := string(b)
			log.Debug("\n Parsing this file:" + s)
			template.New(name).Parse(s)
			log.Debug("\nDone Parsing!")
		default:
			log.Warning("\nWarning: Unknown Template Type: ." +ext)
		}
		
		return nil
	});
}