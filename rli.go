package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"bytes"
	"time"
)

type RedditResponse struct {
	Kind string
	Data DataType
}

type DataType struct {
	Modhash string
	Whitelist_status string
	Children []Child
	After string
	Before string
}

type Child struct {
	Kind string
	Data ChildDataType
}

type ChildDataType struct {
	Subreddit string
	Title string
	Author string
}

func main() {
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, err := http.NewRequest("GET", "https://www.reddit.com/top/.json?count=20", nil)

	req.Header.Set("User-agent", "your bot 0.2")

	resp, err := client.Do(req)


	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	var lst RedditResponse
	json.Unmarshal([]byte(buf.String()), &lst)

	for i, v := range lst.Data.Children {
		fmt.Printf("%d: %s \n", i, v.Data.Title)
	}
}

var licenseCookie = &http.Cookie{Name: "oraclelicense",
	Value:    "accept-securebackup-cookie",
	Expires:  time.Now().Add(356 * 24 * time.Hour),
	HttpOnly: true}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.AddCookie(licenseCookie)
	return nil
}
