package main

import (
	"encoding/json"
	"fmt"
//	"os"
	"net/http"
//	"net/url"
//	"io/ioutil"
//	"strings"
	"bytes"
	"time"
)

//var myClient = &http.Client(Timeout: 10 * time.Second}
/*
type Posting struct {
	Title string
	Subreddit string
	
}*/

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
//	resp, err := http.Get("https://www.reddit.com/top/.json?count=20")

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, err := http.NewRequest("GET", "https://www.reddit.com/top/.json?count=20", nil)

	req.Header.Set("User-agent", "your bot 0.2")

	resp, err := client.Do(req)


	if err != nil {
		panic(err)
	}

//	respBytes := []byte(string(resp.Body))

//	list1 := Listing{}

//	err1 := json.Unmarshal(respBytes, &list1)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
//	newStr := buf.String()

//	fmt.Println(newStr)

	var lst RedditResponse
	json.Unmarshal([]byte(buf.String()), &lst)

	fmt.Println(lst.Kind)
	fmt.Println(lst.Data.Whitelist_status)
	fmt.Println(lst.Data.Children[0].Kind)
	fmt.Println(lst.Data.Children[0].Data.Title)
	for i := range lst.Data.Children {
		fmt.Printf("%d: %s \n", i, lst.Data.Children[i].Data.Title)
	}

//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	fmt.Println("get:\n", string(body))
}

var licenseCookie = &http.Cookie{Name: "oraclelicense",
	Value:    "accept-securebackup-cookie",
	Expires:  time.Now().Add(356 * 24 * time.Hour),
	HttpOnly: true}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.AddCookie(licenseCookie)
	return nil
}

/*
func getJson(url string, target interace{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
*/
//func getListing() Listing {
//	res, err := 
//	
//}
