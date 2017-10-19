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

const (
	NULL int = 0
	EXIT = 1 + iota
)

type RedditPost struct {
	Collection [] RedditResponse
}

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
	Subreddit_id string
	Approved_at_utc string
	Banned_by string
	Removal_reason string
	Link_id string
	Likes int
	Saved bool
	Id string
	Banned_at_utc int
	Gilded int
	Archived bool
	Report_reasons string
	Author string
	Can_mod_post bool
	Ups int
	Parent_id string
	Score int
	Approved_by string
	Downs int
	Body string
	Edited bool
	Author_flair_css_class string
	Collapsed bool
	Is_Submitter bool
	Collapsed_reason string
	Body_html string
	Stickied bool
	Can_gild bool
	Subreddit string
	Score_hidden bool
	Subreddit_type string
	Name string
	Created int
	Author_flair_text string
	Created_utc int
	Subreddit_name_prefixed string
	Controversiality int
	Depth int
	Num_reports int
	Distinguished int
	Url string
	Permalink string
	Title string
	Replies DataType
}

type Child struct {
	Kind string
	Data DataType
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
	
	isExit := 0

	reader := bufio.NewReader(os.Stdin)
	for isExit != 1 {
		fmt.Printf("<Origin> Command: ")
		usrIn, _ := reader.ReadString('\n')
		cmd := strings.Fields(strings.TrimRight(usrIn, "\n"))
		if (strings.ToLower(cmd[0]) == "exit") {
			//Signify exit
			isExit = 1
		} else if ((len(cmd) >= 2) && (cmd[0] == "goto")) {
			//Go to subreddit
			isExit = subreddit(cmd[1])
			
		} else {
			//Erroneous input
			fmt.Printf("Invalid input.\n")
		}
	}
}

func subreddit(subredditString string) int {
	loadURL := redditURL + "r/" + subredditString + "/.json?limit=10"

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, err := http.NewRequest("GET", loadURL, nil)

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
		fmt.Printf("%d: %s \n", i+1, v.Data.Title)
	}	

	reader := bufio.NewReader(os.Stdin)

	isExit := 0
	for isExit == 0 {
		fmt.Printf("<" + subredditString + "> Command: ")
		usrIn, _ := reader.ReadString('\n')
		cmd := strings.Fields(strings.TrimRight(usrIn, "\n"))
		if (strings.ToLower(cmd[0]) == "exit") {
			//Signify exit
			isExit = 1
			return 1
		} else if (strings.ToLower(cmd[0]) == "back") {
			return 0
		} else {
			//Erroneous input
			fmt.Printf("Invalid input.\n")
		}
	}
	return 0
}
/*
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
*/
var licenseCookie = &http.Cookie{Name: "oraclelicense",
	Value:    "accept-securebackup-cookie",
	Expires:  time.Now().Add(356 * 24 * time.Hour),
	HttpOnly: true}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.AddCookie(licenseCookie)
	return nil
}
