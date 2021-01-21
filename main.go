package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/loksonarius/bedrock-server-versions/util"
)

type Format struct {
	LastUpdated time.Time
	Versions    []string
}

func main() {
	versions, err := util.GetVersions()
	if err != nil {
		panic(err)
	}

	now := time.Now().UTC()
	format := Format{
		LastUpdated: now,
		Versions:    versions,
	}

	out, err := json.MarshalIndent(format, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
