package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Country struct {
	Code string
	Name string
}

var (
	_, b, _, _ = runtime.Caller(0)
	currDir    = filepath.Dir(b)
)

func Countries() []Country {
	jsonFile, err := os.ReadFile(currDir + "/countries.json")
	if err != nil {
		fmt.Println("ERROR", err)
		return nil
	}
	var countries []Country
	json.Unmarshal(jsonFile, &countries)
	return countries
}
