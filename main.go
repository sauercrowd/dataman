package main

import(
	dp "github.com/sauercrowd/datafill/dataprovider"
	dr "github.com/sauercrowd/datafill/datareceiver"
	df "github.com/sauercrowd/datafill/definitions"
	tools "github.com/sauercrowd/datafill/tools"
	"fmt"
)

func main(){
	const receiverStr = "cassandra"
	const providerUrl = "https://twitter.com/realDonaldTrump"
	//const providerUrl = "http://trumptwitterarchive.com/#/archive/account/realdonaldtrump"
	const limit       = 200

	var r df.DataReceiver

	config := tools.GetConfigFromFile("config.props")

	//receiver switch
	switch{
	case receiverStr == "csv":
		csv := &dr.CSVProvider{}
		csv.Setup("./donaldTweets.csv","\t")
		r = csv
		break
	case receiverStr == "json":
		json := &dr.JsonProvider{}
		json.Setup("./output.json")
		r = json
		break
	case receiverStr == "cassandra":
		cass := &dr.CassandraProvider{}
		cass.Setup("127.0.0.1","data","{'class': 'SimpleStrategy', 'replication_factor' : 1}")
		r = cass
		break
	case receiverStr == "stdout":
		r = &dr.Stdout{}
		break
	default:
		fmt.Println("No matching data receiver found")
		return
	}

	//provider switch
	switch{
	case dp.CheckTwitter(providerUrl):
		t := &dp.TwitterProvider{}
		t.Login(config["TWITTER_CONSUMER_KEY"], config["TWITTER_CONSUMER_SECRET"])
		t.GetTuples(r, providerUrl, limit)
		break
	case dp.CheckTTA(providerUrl):
		t := &dp.TTAProvider{}
		t.Login(config["TWITTER_CONSUMER_KEY"], config["TWITTER_CONSUMER_SECRET"])
		t.GetTuples(r, providerUrl, limit)
		break
	default:
		fmt.Println("No matching data provider found")
		break
	}
	r.Finish()
}
