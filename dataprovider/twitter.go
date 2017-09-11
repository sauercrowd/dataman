package dataprovider


import(
	"net/http"
	"net/url"
	b64 "encoding/base64"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"time"
	"strconv"
	"fmt"
	"regexp"
	"strings"
	df "github.com/sauercrowd/datafill/definitions"
)

type Tweet struct {
	Contributors interface{} `json:"contributors"`
	Coordinates  interface{} `json:"coordinates"`
	CreatedAt    string      `json:"created_at"`
	Entities     struct {
		Hashtags     []interface{} `json:"hashtags"`
		Symbols      []interface{} `json:"symbols"`
		Urls         []interface{} `json:"urls"`
		UserMentions []interface{} `json:"user_mentions"`
	} `json:"entities"`
	FavoriteCount        int64       `json:"favorite_count"`
	Favorited            bool        `json:"favorited"`
	Geo                  interface{} `json:"geo"`
	ID                   int64       `json:"id"`
	IDStr                string      `json:"id_str"`
	InReplyToScreenName  interface{} `json:"in_reply_to_screen_name"`
	InReplyToStatusID    interface{} `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr interface{} `json:"in_reply_to_status_id_str"`
	InReplyToUserID      interface{} `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   interface{} `json:"in_reply_to_user_id_str"`
	IsQuoteStatus        bool        `json:"is_quote_status"`
	Lang                 string      `json:"lang"`
	Place                interface{} `json:"place"`
	RetweetCount         int64       `json:"retweet_count"`
	Retweeted            bool        `json:"retweeted"`
	Source               string      `json:"source"`
	Text                 string      `json:"text"`
	Truncated            bool        `json:"truncated"`
	User                 struct {
		ContributorsEnabled bool   `json:"contributors_enabled"`
		CreatedAt           string `json:"created_at"`
		DefaultProfile      bool   `json:"default_profile"`
		DefaultProfileImage bool   `json:"default_profile_image"`
		Description         string `json:"description"`
		Entities            struct {
			Description struct {
				Urls []interface{} `json:"urls"`
			} `json:"description"`
		} `json:"entities"`
		FavouritesCount                int         `json:"favourites_count"`
		FollowRequestSent              interface{} `json:"follow_request_sent"`
		FollowersCount                 int         `json:"followers_count"`
		Following                      interface{} `json:"following"`
		FriendsCount                   int         `json:"friends_count"`
		GeoEnabled                     bool        `json:"geo_enabled"`
		HasExtendedProfile             bool        `json:"has_extended_profile"`
		ID                             int         `json:"id"`
		IDStr                          string      `json:"id_str"`
		IsTranslationEnabled           bool        `json:"is_translation_enabled"`
		IsTranslator                   bool        `json:"is_translator"`
		Lang                           string      `json:"lang"`
		ListedCount                    int         `json:"listed_count"`
		Location                       string      `json:"location"`
		Name                           string      `json:"name"`
		Notifications                  interface{} `json:"notifications"`
		ProfileBackgroundColor         string      `json:"profile_background_color"`
		ProfileBackgroundImageURL      string      `json:"profile_background_image_url"`
		ProfileBackgroundImageURLHTTPS string      `json:"profile_background_image_url_https"`
		ProfileBackgroundTile          bool        `json:"profile_background_tile"`
		ProfileBannerURL               string      `json:"profile_banner_url"`
		ProfileImageURL                string      `json:"profile_image_url"`
		ProfileImageURLHTTPS           string      `json:"profile_image_url_https"`
		ProfileLinkColor               string      `json:"profile_link_color"`
		ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color"`
		ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color"`
		ProfileTextColor               string      `json:"profile_text_color"`
		ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
		Protected                      bool        `json:"protected"`
		ScreenName                     string      `json:"screen_name"`
		StatusesCount                  int         `json:"statuses_count"`
		TimeZone                       string      `json:"time_zone"`
		TranslatorType                 string      `json:"translator_type"`
		URL                            interface{} `json:"url"`
		UtcOffset                      int         `json:"utc_offset"`
		Verified                       bool        `json:"verified"`
	} `json:"user"`
}





type TwitterProvider struct{
	Token string
	Tuples []df.ProviderTuple
}

func (t *TwitterProvider) Login(key, secret string){
	/* parsing key and secret to required string, make request and parse it into a json map */
	auth := string(b64.StdEncoding.EncodeToString([]byte(key+":"+secret)))
	bodyStr := []byte("grant_type=client_credentials")
	u := "https://api.twitter.com/oauth2/token"
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(bodyStr))
	req.Header.Set("User-Agent","JonasDataCollector")
	req.Header.Set("grant_type","client_credentials")
	req.Header.Set("Authorization","Basic "+auth)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var dat map[string]interface{}

	if err := json.Unmarshal(body, &dat); err != nil{
		panic(err)
	}

	// check if there's an in the response json
	if err, ok := dat["errors"]; ok{
		panic(err)
	}

	// check if the access token field is present
	token, ok := dat["access_token"]
	if !ok {
		panic(dat)
	}
	fmt.Println("Token: "+token.(string))
	t.Token = token.(string)
}


func (t *TwitterProvider) addTuples(r df.DataReceiver, offset int64, indexOffset int64, u string, countLeft int64, sleep int64){
	const COUNT_MAX=200
	fmt.Printf("Offset: %d\n",offset)
	Url, err := url.Parse("https://api.twitter.com/1.1/statuses/user_timeline.json")
	if err != nil{
		panic(err)
	}
	parameters := url.Values{}
	if offset != -1{
		parameters.Add("max_id", strconv.FormatInt(offset-1, 10))
	}
	parameters.Add("screen_name",u)
	parameters.Add("include_rts","true")
	parameters.Add("exclude_replies","false")
	if countLeft == -1{
		parameters.Add("count", strconv.FormatInt(COUNT_MAX, 10) )
	}else{
		parameters.Add("count", strconv.FormatInt(countLeft, 10) )
	}
	Url.RawQuery = parameters.Encode()

	req, err := http.NewRequest("GET", Url.String(), nil)
	req.Header.Set("Authorization","Bearer "+t.Token)
	req.Header.Set("User-Agent","JonasDataCollector")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var tweets []Tweet

	if err := json.Unmarshal(body, &tweets); err != nil{
		panic(err)
	}

	tuples := make([]df.ProviderTuple, len(tweets))
	var idOffset int64
	idOffset=-1
	for i, obj := range tweets {
		dt, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006",obj.CreatedAt)
		if err != nil{
			fmt.Println(`"`+obj.CreatedAt+`"`);
			fmt.Println("couldnt parse date")
			panic(err)
		}
		tuples[i] = df.ProviderTuple{dt, obj.IDStr, obj.User.ScreenName, obj.Text, -1, obj.FavoriteCount}
		if idOffset == -1 || obj.ID < idOffset{
			idOffset = obj.ID
		}
	}
	fmt.Printf("Tweets size: %d\n",len(tuples))
	r.AddTuples(tuples)

	// 20 or more tweets got requested but less were returned, so there arent more
	//if len(tweets) < COUNT_MAX && (countLeft == -1 || countLeft >= COUNT_MAX) || (countLeft != -1 && countLeft-int64(len(tweets)) < 1){
	if (countLeft == -1 && len(tweets) == 0) || (countLeft != -1 && countLeft-int64(len(tweets)) <= 0){
		return
	}
	// 0.6 seconds to not exceed limit
	time.Sleep(time.Duration(sleep)*time.Millisecond)
	if countLeft != -1 {
		countLeft-=int64(len(tweets))
	}
	t.addTuples(r, idOffset, indexOffset+int64(len(tweets)), u, countLeft, sleep)
}

func (t *TwitterProvider) GetTuples(r df.DataReceiver, u string, count int64){
	path := strings.Replace(u, "https://twitter.com/","", 1)
	user := strings.Split(path,"/")[0]
	fmt.Printf("user: ",string(user))
	t.addTuples(r, -1, 0, user, count, 600)
}

func CheckTwitter(u string) bool{
	reg, err := regexp.Compile("https://twitter.com/[^/]*/?")
	if err != nil{
		panic(err)
	}
	return reg.MatchString(u)
}
