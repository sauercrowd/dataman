package dataprovider

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	df "github.com/sauercrowd/dataman/pkg/definitions"
)

type TTAProvider struct {
	Token  string
	Tuples []df.ProviderTuple
}

func (t *TTAProvider) Login(key, secret string) {
	/* parsing key and secret to required string, make request and parse it into a json map */
	auth := string(b64.StdEncoding.EncodeToString([]byte(key + ":" + secret)))
	bodyStr := []byte("grant_type=client_credentials")
	u := "https://api.twitter.com/oauth2/token"
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(bodyStr))
	req.Header.Set("User-Agent", "JonasDataCollector")
	req.Header.Set("grant_type", "client_credentials")
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var dat map[string]interface{}

	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	// check if there's an in the response json
	if err, ok := dat["errors"]; ok {
		panic(err)
	}

	// check if the access token field is present
	token, ok := dat["access_token"]
	if !ok {
		panic(dat)
	}
	fmt.Println("Token: " + token.(string))
	t.Token = token.(string)
}

type TWAUser struct {
	Account      string `json:"account"`
	Display      bool   `json:"display"`
	ID           int    `json:"id"`
	Inactive     bool   `json:"inactive"`
	Linked       string `json:"linked"`
	Name         string `json:"name"`
	StartingYear string `json:"starting_year"`
	Title        string `json:"title"`
}

type TWATweet struct {
	CreatedAt string `json:"created_at"`
	IDStr     string `json:"id_str"`
	IsRetweet bool   `json:"is_retweet"`
	Source    string `json:"source"`
	Text      string `json:"text"`
}

func getFirstYear(u string) int64 {
	const url = "http://trumptwitterarchive.com/data/accounts.json"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var users []TWAUser
	if err = json.Unmarshal(body, &users); err != nil {
		panic(err)
	}

	for _, user := range users {
		if user.Account == u {
			i, err := strconv.ParseInt(user.StartingYear, 10, 64)
			if err != nil {
				panic(err)
			}
			return i
		}
	}
	return -1

}

func addTuples(u string, year int64, ch chan []df.ProviderTuple) {
	url := fmt.Sprintf("http://trumptwitterarchive.com/data/%s/%d.json", u, year)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var tweets []TWATweet
	if err = json.Unmarshal(body, &tweets); err != nil {
		panic(err)
	}

	tuples := make([]df.ProviderTuple, len(tweets))
	for i, x := range tweets {
		t, err := time.Parse("2006", fmt.Sprintf("%d", year))
		if err != nil {
			panic(err)
		}
		tuples[i] = df.ProviderTuple{t, x.IDStr, u, x.Text, -1, -1}
	}
	fmt.Printf("year: %d\n", year)
	ch <- tuples
}

func (t *TTAProvider) GetTweetDetails(r []df.DataReceiver, tp []df.ProviderTuple) {
	for i, x := range tp {
		fmt.Printf("[TWITTER] Getting Tweet %s from %s\n", x.ID, x.Author)
		URL, err := url.Parse("https://api.twitter.com/1.1/statuses/show.json")
		if err != nil {
			panic(err)
		}
		parameters := url.Values{}
		parameters.Add("id", x.ID)
		parameters.Add("trim_user", "true")
		parameters.Add("include_my_retweet", "false")

		URL.RawQuery = parameters.Encode()

		req, err := http.NewRequest("GET", URL.String(), nil)
		req.Header.Set("Authorization", "Bearer "+t.Token)
		req.Header.Set("User-Agent", "JonasDataCollector")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var tw Tweet

		if err := json.Unmarshal(body, &tw); err != nil {
			panic(err)
		}

		dt, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tw.CreatedAt)
		if err != nil {
			fmt.Println("[ERROR] while parsing date; propably tweet id not found")
			fmt.Printf("[ERROR] Response body: %s", string(body))
			fmt.Println("[ERROR] Continue...")
			time.Sleep(time.Second)
			continue
		}
		tp[i].Date = dt
		tp[i].Score = tw.FavoriteCount
		for _, receiver := range r {
			receiver.AddTuple(tp[i])
		}
		time.Sleep(time.Second)
	}
}

func (t *TTAProvider) SendTuples(r []df.DataReceiver, u string, ch chan []df.ProviderTuple, size int64) {
	for i := 0; i < int(size); i++ {
		t.GetTweetDetails(r, <-ch)
	}
}

func (t *TTAProvider) GetTuples(r []df.DataReceiver, u string, count int64) {
	path := strings.Replace(u, "http://trumptwitterarchive.com/#/archive/account/", "", 1)
	user := strings.Split(path, "/")[0]
	log.Print("user:", string(user))
	ch := make(chan []df.ProviderTuple)
	maxYear := getFirstYear(user)
	if maxYear == -1 {
		panic("Account not Found")
	}
	for i := maxYear; i <= 2017; i++ {
		go addTuples(user, i, ch)
	}
	t.SendTuples(r, u, ch, 2017-maxYear+1)
}

func CheckTTA(u string) bool {
	reg, err := regexp.Compile("http://trumptwitterarchive.com/#/archive/account/[^/]*/?")
	if err != nil {
		panic(err)
	}
	return reg.MatchString(u)
}
