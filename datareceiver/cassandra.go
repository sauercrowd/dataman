package datareceiver

import(
	"fmt"
	"github.com/gocql/gocql"
	df "github.com/sauercrowd/datafill/definitions"

)

type CassandraProvider struct{
	session *gocql.Session
	clusterConf *gocql.ClusterConfig
	keyspace string
}

func (c *CassandraProvider) Setup(hosts string, keyspace string, optStr string){
	c.clusterConf = gocql.NewCluster(hosts)
	var err error
	c.session, err = c.clusterConf.CreateSession()
	if err != nil{
		panic(err)
	}
	err = c.session.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH replication = %s",keyspace,optStr)).Exec()
	c.session.Close()
	c.clusterConf.Keyspace=keyspace
	c.session, err = c.clusterConf.CreateSession()
	if err != nil{
		panic(err)
	}
	err = c.session.Query("CREATE TABLE IF NOT EXISTS data (id text PRIMARY KEY, date timestamp, author text, content text, score bigint, replies bigint)").Exec()
	if err != nil{
		panic(err)
	}
}

func (c *CassandraProvider) AddTuples(tp []df.ProviderTuple){
	template := "INSERT INTO data(id, date, author, content, score, replies) VALUES (?, ?, ?, ?, ?, ?)"
	for _,t := range tp{
		err := c.session.Query(template, t.Id, t.Date, t.Author, t.Content, t.Score, t.Replies).Exec()
		if err != nil{
			panic(err)
		}
	}
}

func (c *CassandraProvider) AddTuple(tp df.ProviderTuple){
	template := "INSERT INTO data(id, date, author, content, score, replies) VALUES (?, ?, ?, ?, ?, ?)"
	err := c.session.Query(template, tp.Id, tp.Date, tp.Author, tp.Content, tp.Score, tp.Replies).Exec()
	if err != nil{
		panic(err)
	}
}

func (c *CassandraProvider) Finish(){
	c.session.Close()
}
