package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var cli args

func main() {
	if flag.NFlag() < 4 {
		log.Fatal("missing flags")
	}
	bdy, err := json.Marshal(cli.data)
	if err != nil {
		log.Fatal(err)
	}
	api, err := url.JoinPath("https://api.github.com", "repos", cli.repo, "issues")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("post %s", api)
	req, err := http.NewRequest(http.MethodPost, api, bytes.NewBuffer(bdy))
	if err != nil {
		log.Fatal(err)
	}
	var b strings.Builder
	b.WriteString("Bearer ")
	b.WriteString(cli.bear)
	req.Header.Add("Authorization", b.String())
	req.Header.Add("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	rsp.Write(os.Stdout)

	p, ok := os.LookupEnv("GITHUB_OUTPUT")
	if !ok {
		log.Fatal("'GITHUB_OUTPUT' env not set")
	}
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b.Reset()
	b.WriteString("RESPONSE_STATUS=")
	b.WriteString(rsp.Status)
	b.WriteString("\n")
	_, err = f.WriteString(b.String())
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	log.SetFlags(0)
	flag.StringVar(&cli.repo, "repo", "", "")
	flag.StringVar(&cli.bear, "bear", "", "")
	flag.StringVar(&cli.data.T, "title", "", "")
	flag.StringVar(&cli.data.B, "body", "", "")
	flag.Usage = func() {
		log.Println("usage: issue [flags]")
		flag.PrintDefaults()
	}
	flag.Parse()
	log.SetPrefix("issue: ")
}

type args struct {
	repo, bear string
	data
}

type data struct {
	T string `json:"title"`
	B string `json:"body"`
}
