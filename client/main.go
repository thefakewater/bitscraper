package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const GET = "GET"

var amount string

type File struct {
	Name string
	Url  string
}
type Manifest struct {
	Author          string
	Name            string
	Software        string
	Source          string
	ManifestVersion int
	Version         string
	Files           []File
}

func GetManifest() Manifest {
	var manifest Manifest
	dat, err := os.ReadFile("manifest.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(dat, &manifest)
	return manifest
}

func GetAmountFiles(manifest Manifest) int {
	amount := 0
	for i := range manifest.Files {
		amount = i + 1
	}
	return amount
}
func ResolveSource(manifest Manifest) string {
	if manifest.Source == "BitMidi" {
		return "https://bitmidi.com"
	}
	log.Fatal("Could not find the source " + manifest.Source)
	return ""
}
func Download(manifest Manifest) {
	err := os.Mkdir(manifest.Name, 0755)
	if err != nil {
		log.Fatal(err)
	}
	URL := ResolveSource(manifest)
	for i, s := range manifest.Files {
		req, _ := http.NewRequest(GET, URL+s.Url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		f, err := os.OpenFile(manifest.Name+"/"+s.Name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		io.Copy(io.MultiWriter(f), resp.Body)
		fmt.Println("[" + strconv.Itoa(i+1) + "/" + amount + "] Downloaded " + s.Name)
	}
}
func main() {
	manifest := GetManifest()
	amount = strconv.Itoa(GetAmountFiles(manifest))
	fmt.Println("Manifest Info")
	fmt.Println("----------------------------")
	fmt.Printf("Name : %s", manifest.Name+"\n")
	fmt.Printf("Version : %s", manifest.Version+"\n")
	fmt.Printf("Author : %s", manifest.Author+"\n")
	fmt.Printf("Files : %s", amount+"\n")
	fmt.Printf("Source : %s", manifest.Source+"\n")
	fmt.Println("----------------------------")
	fmt.Println("Are you sure to download this pack? [PRESS ENTER]")
	fmt.Scanln()
	start := time.Now()
	Download(manifest)
	elapsed := time.Since(start)
	fmt.Println("Downloaded all files in " + elapsed.String())
}
