package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"bytes"
	"time"
	"os"
	"bufio"
)

const redditURL string = "https://www.reddit.com/"

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
	Permalink string
	Url string
}

func main() {
	originMenu()
}

func originMenu() {
	fmt.Println("Welcome to CLI for Reddit!")
	fmt.Println("----------Commands----------")
	fmt.Println("goTo [subReddit]	: load the posts in subReddit")
	fmt.Println("exit				: exit the program")
	fmt.Println("back				: go back to the previous level")
	
	isExit := false

	reader := bufio.NewReader(os.Stdin)
	for isExit == false {
		fmt.Printf("Command: ")
		usrIn, _ := reader.ReadString('\n')
		cmd := strings.Fields(strings.TrimRight(usrIn, "\n"))
		if (strings.ToLower(cmd[0]) == "exit") {
			//Signify exit
			isExit = 1
		} else if ((len(cmd) >= 2) & (cmd[0] == "goto")) {
			//Go to subreddit
			isExit = subreddit(cmd[1])
			
		} else {
			//Erroneous input
			fmt.Printf("Invalid input.\n")
		}
	}
}

func subreddit(subredditString string) {
	loadURL := redditURL + "r/" + subredditString + "/.json?limit=10"

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, err := http.NewRequest("GET", loadURL, nil)

	req.Header.Set("User-agent", "your bot 0.2")
	fmt.Printf("%s \n", loadURL)
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	var lst RedditResponse
	json.Unmarshal([]byte(buf.String()), &lst)

	for i, v := range lst.Data.Children {
		fmt.Printf("%d: %s \n", i+1, v.Data.Title)
	}


	
}

displaySubreddit(subredditString string) {

}

func displayComments (subredditString string, postID string) {
	loadURL := redditURL + "r/" + subredditString + "/comments/" + postID + ".json?"

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, err := http.NewRequest("GET", loadURL, nil)

	req.Header.Set("User-agent", "your bot 0.2")
	fmt.Printf("%s \n", loadURL)
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	var lst RedditResponse
	json.Unmarshal([]byte(buf.String()), &lst)

}

var licenseCookie = &http.Cookie{Name: "oraclelicense",
	Value:    "accept-securebackup-cookie",
	Expires:  time.Now().Add(356 * 24 * time.Hour),
	HttpOnly: true}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.AddCookie(licenseCookie)
	return nil
}
