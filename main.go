package main

import (
	"fmt"

	dp "github.com/sauercrowd/dataman/pkg/dataprovider"
	dr "github.com/sauercrowd/dataman/pkg/datareceiver"
	df "github.com/sauercrowd/dataman/pkg/definitions"
	flags "github.com/sauercrowd/dataman/pkg/flags"
)

func main() {
	f := flags.ParseFlags()
	const limit = 2000

	receivers := make([]df.DataReceiver, 0)

	//receiver switch
	switch {
	case receiverStr == "csv":
		csv := &dr.CSVProvider{}
		csv.Setup("./donaldTweets.csv", "\t")
		receivers = append(receivers, csv)
	case receiverStr == "json":
		json := &dr.JsonProvider{}
		json.Setup("./output.json")
		receivers = append(receivers, json)
	case receiverStr == "cassandra":
		cass := &dr.CassandraProvider{}
		cass.Setup("127.0.0.1", "data", "{'class': 'SimpleStrategy', 'replication_factor' : 1}")
		receivers = append(receivers, cass)
	case receiverStr == "stdout":
		r = &dr.Stdout{}
		receivers = append(receivers, r)
	}

	//provider switch
	switch {
	case dp.CheckTwitter(providerUrl):
		t := &dp.TwitterProvider{}
		t.Login(f.TwitterConsumerKey, f.TwitterConsumerSecret)
		t.GetTuples(receivers, providerUrl, limit)
		break
	case dp.CheckTTA(providerUrl):
		t := &dp.TTAProvider{}
		t.Login(f.TwitterConsumerKey, f.TwitterConsumerSecret) //still needed to get additional tweet informations
		t.GetTuples(receivers, providerUrl, limit)
		break
	default:
		log.Fatalf("No matching data provider found")
	}

	//close everything
	for _, r := receivers{
		r.Finish()
	}	
}
