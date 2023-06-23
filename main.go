package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Link struct {
	Type string
	Lang string
	Path string
}

func main() {
	envServices := os.Getenv("KIWIX_SERVICES")
	if envServices == "" {
		fmt.Println("No services found")
		return
	}

	envLanguages := os.Getenv("KIWIX_LANGUAGES")
	if envLanguages == "" {
		fmt.Println("No languages found")
		return
	}

	services := strings.Split(os.Getenv("KIWIX_SERVICES"), ",")

	if len(services) == 0 {
		fmt.Println("No services found")
		return
	}

	fmt.Println("Services: ", services)

	languages := strings.Split(os.Getenv("KIWIX_LANGUAGES"), ",")

	if len(languages) == 0 {
		fmt.Println("No languages found")
		return
	}

	fmt.Println("Languages: ", languages)

	downloads := make([]Link, 0)

	for _, service := range services {
		resp, err := http.Get("https://download.kiwix.org/zim/" + service + "/")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		found := make(map[string][]string)
		for _, link := range regexp.MustCompile(`href="([^"]*)"`).FindAllStringSubmatch(string(body), -1) {
			for _, language := range languages {
				if strings.HasPrefix(link[1], service+"_"+language+"_all_maxi") {
					if _, ok := found[language]; !ok {
						found[language] = make([]string, 0)
					}
					found[language] = append(found[language], link[1])
				}
			}
		}
		for fLang, fLinks := range found {
			sort.Strings(fLinks)
			downloads = append(downloads, Link{
				Type: service,
				Lang: fLang,
				Path: fLinks[len(fLinks)-1],
			})
		}
	}

	cache := make(map[string]map[string]string)

	if _, err := os.Stat("cache.json"); err == nil {
		jsonStr, err := ioutil.ReadFile("cache.json")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		cache = make(map[string]map[string]string)
		err = json.Unmarshal(jsonStr, &cache)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

	for _, download := range downloads {
		if _, ok := cache[download.Type]; ok {
			if _, ok := cache[download.Type][download.Lang]; ok {
				if cache[download.Type][download.Lang] == download.Path {
					fmt.Println("Skip: ", download)
					continue
				}
			}
		}
		fmt.Println("Download: ", download)
		//"https://download.kiwix.org/zim/"+service+"/"+download.Path
		if _, ok := cache[download.Type]; !ok {
			cache[download.Type] = make(map[string]string)
		}
		cache[download.Type][download.Lang] = download.Path
	}

	jsonStr, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	ioutil.WriteFile("cache.json", jsonStr, 0644)
}
