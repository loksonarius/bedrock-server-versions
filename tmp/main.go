package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	source := "https://raw.githubusercontent.com/loksonarius/bedrock-server-versions/main/versions.json"
	resp, err := http.Get(source)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	parsedResponse := struct {
		LastUpdated string
		Versions    []string
	}{}

	if err := json.Unmarshal(out, &parsedResponse); err != nil {
		panic(err)
	}

	for _, v := range parsedResponse.Versions {
		fmt.Println(v)
	}
}