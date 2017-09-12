package flags

import (
	"flag"
)

type Flags struct {
	//receiver switches
	UseStdout, UseCSV, UseCassandra, UseJSON  bool
	TwitterConsumerKey, TwitterConsumerSecret string
	URL                                       string
	//receiver specific flags
	JSONFile      string
	CSVFile       string
	CassandraHost string
}

func ParseFlags() Flags {
	var f Flags
	flag.StringVar(&f.URL, "url", "", "url where the data should be gathered from. Currently supported: twitter urls, trumptwitterarchive urls")
	flag.StringVar(&f.TwitterConsumerKey, "tkey", "", "twitter consumer key, needed when using twitter as a data source")
	flag.StringVar(&f.TwitterConsumerSecret, "tsecret", "", "twitter consumer key, needed when using twitter as a data source")
	flag.StringVar(&f.CassandraHost, "chost", "127.0.0.1", "cassandra host, only used when using cassandra as a data receiver")
	flag.StringVar(&f.JSONFile, "jsonfile", "./output.json", "determines where the json file should be stored, only used when using json as a data receiver")
	flag.StringVar(&f.CSVFile, "csvfile", "./output.csv", "determines where the csv file should be stored, only used when using csv as a data receiver")

	flag.BoolVar(&f.UseStdout, "stdout", true, "use stdout")
	flag.BoolVar(&f.UseCassandra, "cassandra", false, "use cassandra")
	flag.BoolVar(&f.UseCSV, "csv", false, "use csv")
	flag.BoolVar(&f.UseJSON, "json", false, "use json")
	flag.Parse()
	return f
}
