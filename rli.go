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
//	"io/ioutil"
	"os/exec"
)

const redditURL string = "https://www.reddit.com/"
const resultLimit int = 10
const charLimit int = 64

type RedditResponse struct {
	Kind string
	Data DataType
}

type DataType struct {
//	Modhash string
//	Whitelist_status string
	Children []RedditResponse
	After string
	Before string
//	Subreddit_id string
//	Approved_at_utc string
//	Banned_by string
//	Removal_reason string
//	Link_id string
//	Likes int
//	Saved bool
	Id string
//	Banned_at_utc int
//	Gilded int
//	Archived bool
//	Report_reasons string
//	Author string
//	Can_mod_post bool
//	Ups int
//	Parent_id string
//	Score int
//	Approved_by string
//	Downs int
	Body string `json:"body"`
//	Edited bool
//	Author_flair_css_class string
//	Collapsed bool
//	Is_Submitter bool
//	Collapsed_reason string
//	Body_html string
//	Stickied bool
//	Can_gild bool
//	Subreddit string
//	Score_hidden bool
//	Subreddit_type string
//	Name string
//	Created int
//	Author_flair_text string
//	Created_utc int
//	Subreddit_name_prefixed string
//	Controversiality int
//	Depth int
//	Num_reports int
//	Distinguished int
	Url string
//	Permalink string
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
		if (cmd[0] == "testing") { //go to testing module
			testing()
		} else if (cmd[0] == "help") { //display commands
			welcome()
		} else {
			//use one of the default switches
			isExit = defaultSwitcher(cmd)
		}
	}
}

func welcome() {
	//Display the input options of the program
	fmt.Println("Welcome to CLI for Reddit!")
	fmt.Println("----------Commands----------")
	fmt.Println("goTo [subReddit]		: load the posts in subReddit")
	fmt.Println("exit				: exit the program")
	fmt.Println("back				: go back to the previous level")
	fmt.Println("full [int]			: display the full comment or title")
}

func defaultSwitcher(cmd []string) int {
		var isExit int
		if (strings.ToLower(cmd[0]) == "exit") {
			isExit = 1 //signify exit
		} else if (strings.ToLower(cmd[0]) == "back") {
			isExit = 0 //signify back
		} else if ((len(cmd) >= 2) && (cmd[0] == "goto")) {
			isExit = subreddit(cmd[1], "") //goto subreddit
		} else {
			fmt.Printf("Invalid input.\n") //erroneous input
		}
		return isExit
}

func subreddit(subredditString string, after string) int {
	isExit := 0
	//Build the url to load the json
	var loadURL string
	if after == "" { //load the first page of the subreddit
		loadURL = fmt.Sprintf("%s%s%s%s%s%d", redditURL, "r/", subredditString, "/.json?", "limit=", resultLimit)
	} else { //load one of the following pages of the subreddit
		loadURL = fmt.Sprintf("%s%s%s%s%s%d%s%s%s%d", redditURL, "r/", subredditString, "/.json?", "limit=", resultLimit, "&after=", after, "&count=", resultLimit)
	}

	//mechanism to pull the json data
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	req, err := http.NewRequest("GET", loadURL, nil)
	req.Header.Set("User-agent", "your bot 0.2")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("There was an issue loading the URL. Please try again later.")
		return 0 //go back one level
	}

	//transform the response data into a data structure
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	var lst RedditResponse
	json.Unmarshal([]byte(buf.String()), &lst)
	f, err := os.Create("lastJson.txt")
	defer f.Close()
	f.WriteString(buf.String())
	f.Sync()

	for i, v := range lst.Data.Children {
		fmt.Printf("%d: %s \n", i+1, v.Data.Title[0:min(charLimit,len(v.Data.Title))])
	}	

	reader := bufio.NewReader(os.Stdin)

	for isExit == 0 {
		fmt.Printf("<" + subredditString + "> Command: ")
		usrIn, _ := reader.ReadString('\n')
		cmd := strings.Fields(strings.TrimRight(usrIn, "\n"))
		if ((len(cmd) >= 2) && (cmd[0] == "comm")) {
			postIndex, _ := strconv.Atoi(cmd[1])
			isExit = comments(subredditString, lst.Data.Children[postIndex - 1].Data.Id, lst.Data.Children[postIndex-1].Data.Title)
		} else if ((len(cmd) >=2) && (cmd[0] == "full")) {
			postIndex, _ := strconv.Atoi(cmd[1])
			fmt.Printf("%d: %s \n", postIndex, lst.Data.Children[postIndex-1].Data.Title)
		} else if ((len(cmd) >= 2) && (cmd[0] == "open")) {
			postIndex, _ := strconv.Atoi(cmd[1])
			cmd := exec.Command("open", lst.Data.Children[postIndex - 1].Data.Url)
			cmd.Output()
		} else if (strings.ToLower(cmd[0]) == "next") {
			//Goto next page of subreddit
			isExit = subreddit(subredditString, lst.Data.After)
		} else {
			//one of the default switches
			isExit = defaultSwitcher(cmd)
		}
	}
	return isExit
}

func testing () {
	comments("nfl","77b9kt","testing") 
}

func comments (subredditString string, postID string, postTitle string) int {
	isExit := 0
	loadURL := fmt.Sprintf("%s%s%s%s%s%s%d", redditURL, "r/", subredditString, "/comments/", postID, "/.json?")

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	req, err := http.NewRequest("GET", loadURL, nil)
	req.Header.Set("User-agent", "your bot 0.2")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("There was an issue loading the URL. Please try again later.")
		return 0 //go back one level
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	result := make([]RedditResponse,0)
	json.Unmarshal([]byte(buf.String()), &result)

	f, err := os.Create("lastJson.txt")

	defer f.Close()

	f.WriteString(buf.String())

	f.Sync()

	for i, v := range result[1].Data.Children {
		if ((len(v.Data.Body) > 0) && (i <= resultLimit + 1)){
			fmt.Printf("%d: %s\n",i+1, v.Data.Body[0:min(charLimit, len(v.Data.Body))])
		}
	}

	reader := bufio.NewReader(os.Stdin)

	for isExit == 0 {
		fmt.Printf("<" + postTitle[0:min(5, len(postTitle))] + "> Command: ")
		usrIn, _ := reader.ReadString('\n')
		cmd := strings.Fields(strings.TrimRight(usrIn, "\n"))
		if ((len(cmd) >=2) && (cmd[0] == "full")) {
			commIndex, _ := strconv.Atoi(cmd[1])
			fmt.Printf("%d: %s \n", commIndex, result[1].Data.Children[commIndex-1].Data.Body)
		} else if ((len(cmd) >=2) && (cmd[0] == "more")) {
			grabComment, _ := strconv.Atoi(cmd[1])
			fmt.Printf("Expanded: %s\n",  result[1].Data.Children[grabComment].Data.Body)
		} else {
			//one of the default switches
			isExit = defaultSwitcher(cmd)
		}
	}
	return isExit
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
