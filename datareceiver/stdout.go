package datareceiver

import(
	"fmt"
	df "github.com/sauercrowd/datafill/definitions"
)

type Stdout struct{}

func (s *Stdout) AddTuples(tp []df.ProviderTuple){
	for _, x := range tp{
		fmt.Printf("ID: %s AUTHOR: %s DATE: %s SCORE: %d\n", x.Id, x.Author, x.Date.Format("02.01.2006 15:04"), x.Score)
	}
}


func (s *Stdout) AddTuple(tp df.ProviderTuple){
	fmt.Printf("ID: %s AUTHOR: %s DATE: %s SCORE: %d\n", tp.Id, tp.Author, tp.Date.Format("02.01.2006 15:04"), tp.Score)
}

func (s *Stdout) Finish(){}
