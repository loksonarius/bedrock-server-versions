# Bedrock Server Versions

This repo serves as a quick hack to solve an issue I've had working on
[mcsm](https://github.com/loksonarius/mcsm): how to get a list of valid versions
for Minecraft Bedrock Edition's servers.

The solution I've gone with is having a binary that scrapes Gamepedia for all
listed versions, sorts, and caches the response in a local JSON file. Now, to
consume the cached response, clients can make an HTTP request for the raw
content of the file from GitHub and parse the response into a JSON struct from
there!

## How to Use

From a Unix terminal, the following would work:

```bash
curl https://raw.githubusercontent.com/loksonarius/bedrock-server-versions/main/versions.json
```

Though realistically, any code that can make an HTTP request and parse JSON will
work. Below is a snippet for doing this in Go:

```go
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
```

## Updates

The repo is set up to scrape Gamepedia once a week for server versions and
update the `versions.json` file. The last run scrape time for a response will be
stored in the `LastUpdated` field of the struct in `versions.json`.

## Running Locally

The following can be used to run a scrape locally:

```bash
go run main.go
```
