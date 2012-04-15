package configurations

import (
	"encoding/json"
	"os"
)

func Load(filename string, object interface{}) error {
	file, err := os.Open("./configurations/" + filename + ".json")
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	return dec.Decode(object)

}
