package datareceiver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	df "github.com/sauercrowd/dataman/pkg/definitions"
)

type JsonProvider struct {
	tuples   []df.ProviderTuple
	filepath string
}

func jPrintln(s string) {
	fmt.Println("[JSON] ", s)
}

func (j *JsonProvider) Setup(filepath string) {
	j.filepath = filepath
	// create a slice with 64 to minimize copy processes
	j.tuples = make([]df.ProviderTuple, 0, 64)
	jPrintln("FilePath: " + filepath)
}

func (j *JsonProvider) AddTuples(tp []df.ProviderTuple) {
	j.tuples = append(j.tuples, tp...)
}

func (j *JsonProvider) AddTuple(tp df.ProviderTuple) {
	j.tuples = append(j.tuples, tp)
}

func (j *JsonProvider) Finish() {
	data, err := json.Marshal(j.tuples)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(j.filepath, data, 0644)
	if err != nil {
		panic(err)
	}
}
