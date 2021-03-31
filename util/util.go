package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const GAMEPEDIA_BASE_URL = "https://minecraft.gamepedia.com"

func httpGet(addr string) (io.ReadCloser, error) {
	parsed, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(parsed.String())
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func scrapeGamepediaVersions() ([]string, error) {
	var versions []string
	addr := GAMEPEDIA_BASE_URL + "/Bedrock_Edition_version_history"
	body, err := httpGet(addr)
	if err != nil {
		return versions, err
	}
	defer body.Close()

	tables := []*html.Node{}
	var crawlForVersionTables func(n *html.Node)
	crawlForVersionTables = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			for _, a := range n.Attr {
				if a.Key == "data-description" {
					if strings.Contains(a.Val, "version history") {
						tables = append(tables, n)
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawlForVersionTables(c)
		}
	}

	urlRe := regexp.MustCompile(`\/Bedrock_Edition_(?P<version>\d+\.\d+\.\d+(\.\d+)?)$`)
	var versionUrls []string
	var crawlForVersionUrls func(n *html.Node)
	crawlForVersionUrls = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && urlRe.MatchString(a.Val) {
					versionUrls = append(versionUrls, GAMEPEDIA_BASE_URL+a.Val)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawlForVersionUrls(c)
		}
	}

	vmap := make(map[string]bool)
	versionRe := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)
	var crawlForServerVersion func(n *html.Node)
	crawlForServerVersion = func(n *html.Node) {
		d := strings.TrimSpace(n.Data)
		if n.Type == html.TextNode && versionRe.MatchString(d) {
			vmap[d] = true
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawlForServerVersion(c)
		}
	}

	doc, err := html.Parse(body)
	if err != nil {
		return versions, err
	}

	crawlForVersionTables(doc)
	for _, t := range tables {
		crawlForVersionUrls(t)
	}

	for _, u := range versionUrls {
		body, err := httpGet(u)
		if err != nil {
			return versions, err
		}
		defer body.Close()

		if doc, err := html.Parse(body); err == nil {
			crawlForServerVersion(doc)
		}
	}

	for k := range vmap {
		versions = append(versions, k)
	}

	return versions, nil
}

func compareVersionStrings(s1, s2 string) bool {
	a := strings.Split(s1, ".")
	b := strings.Split(s2, ".")

	if len(b) > len(a) {
		c := a
		a = b
		b = c
	}

	for k := range a {
		if k >= len(b) {
			return true
		}

		ad, aerr := strconv.Atoi(a[k])
		bd, berr := strconv.Atoi(b[k])

		if aerr != nil || berr != nil {
			continue
		}

		if ad == bd {
			continue
		}

		return ad < bd
	}

	return false
}

func GetVersions() ([]string, error) {
	versions := []string{}

	versions, err := scrapeGamepediaVersions()
	if err != nil {
		return versions, err
	}

	if len(versions) == 0 {
		return versions, fmt.Errorf("failed to find released bedrock versions")
	}

	sort.Slice(versions, func(i, j int) bool {
		return compareVersionStrings(versions[i], versions[j])
	})

	return versions, nil
}
