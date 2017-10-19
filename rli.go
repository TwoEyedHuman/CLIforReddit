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
	"strconv"
)

const redditURL string = "https://www.reddit.com/"
const resultLimit int = 10
const charLimit int = 64

type RedditResponse struct {
	Kind string
	Data DataType
}

type DataType struct {
	Modhash string
	Whitelist_status string
	Children []RedditResponse
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
	Body string `json:"body"`
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
}

func main() {
	welcome()
	originMenu()
}

func originMenu() {
	isExit := 0  //condition that holds whether program should quit

	reader := bufio.NewReader(os.Stdin) //reads input fromm command line

	for isExit != 1 {  //value of 1 indicates quit
		fmt.Printf("<Origin> Command: ")
		usrIn, _ := reader.ReadString('\n')
		cmd := strings.Fields(strings.TrimRight(usrIn, "\n")) //trims and tokenizes user input
		if (strings.ToLower(cmd[0]) == "exit") { //signify exit
			isExit = 1
		} else if ((len(cmd) >= 2) && (cmd[0] == "goto")) { //go to subreddit
			isExit = subreddit(cmd[1])
		} else if (cmd[0] == "testing") { //go to testing module
			testing()
		} else if (cmd[0] == "help") { //display commands
			welcome()
		} else {  //erroneous or unexpected input
			fmt.Printf("Invalid input.\n")
		}
	}
}

func welcome() {
	fmt.Println("Welcome to CLI for Reddit!")
	fmt.Println("----------Commands----------")
	fmt.Println("goTo [subReddit]		: load the posts in subReddit")
	fmt.Println("exit				: exit the program")
	fmt.Println("back				: go back to the previous level")
}

func subreddit(subredditString string) int {
	loadURL := fmt.Sprintf("%s%s%s%s%d", redditURL, "r/", subredditString, "/.json?limit=", resultLimit)

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
		fmt.Printf("%d: %s \n", i+1, v.Data.Title[0:min(charLimit,len(v.Data.Title)-1)])
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
		} else if ((len(cmd) >= 2) && (cmd[0] == "goto")) {
			//Go to subreddit
			isExit = subreddit(cmd[1])
			return isExit
		} else if ((len(cmd) >= 2) && (cmd[0] == "comm")) {
			postIndex, _ := strconv.Atoi(cmd[1])
			isExit = comments(subredditString, lst.Data.Children[postIndex - 1].Data.Id)
			return isExit
		} else {
			//Erroneous input
			fmt.Printf("Invalid input.\n")
		}
	}
	return 0
}

func testing () {
	comments("nfl","77b9kt") 
}

func comments (subredditString string, postID string) int {
	loadURL := fmt.Sprintf("%s%s%s%s%s%s%d", redditURL, "r/", subredditString, "/comments/", postID, "/.json?limit=", resultLimit)

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
	result := make([]RedditResponse,0)
	json.Unmarshal([]byte(buf.String()), &result)
	fmt.Printf("Size: %d\n", len(result[1].Data.Children))
	for i, v := range result[1].Data.Children {
		if (len(v.Data.Body) > 0) {
			fmt.Printf("%d: %s\n",i+1, v.Data.Body[0:min(charLimit, len(v.Data.Body)-1)])
		}
	}
	return 0
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
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
