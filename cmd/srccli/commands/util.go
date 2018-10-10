package commands

import (
	stdjson "encoding/json"
	"os"
	jww "github.com/spf13/jwalterweatherman"

)


// txFeed
type txFeed struct {
	Alias  string `json:"alias"`
	Filter string `json:"filter,omitempty"`
}


func printJSON(data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if ok != true {
		jww.ERROR.Println("invalid type assertion")
		os.Exit(1)
	}

	rawData, err := stdjson.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		jww.ERROR.Println(err)
		os.Exit(1)
	}

	jww.FEEDBACK.Println(string(rawData))
}

func printJSONList(data interface{}) {
	dataList, ok := data.([]interface{})
	if ok != true {
		jww.ERROR.Println("invalid type assertion")
		os.Exit(1)
	}

	for idx, item := range dataList {
		jww.FEEDBACK.Println(idx, ":")
		rawData, err := stdjson.MarshalIndent(item, "", "  ")
		if err != nil {
			jww.ERROR.Println(err)
			os.Exit(1)
		}

		jww.FEEDBACK.Println(string(rawData))
	}
}