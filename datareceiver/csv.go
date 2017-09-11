package datareceiver

import(

	df "github.com/sauercrowd/datafill/definitions"
	"bytes"
	"io/ioutil"
	"strconv"
	"fmt"
)

type CSVProvider struct{
	buf bytes.Buffer
	filepath, s string
}

func cPrintln(s string){
	fmt.Println("[CSV] "+s)
}

func (c *CSVProvider) Setup(filepath string, seperator string){
	c.filepath = filepath
	c.s = seperator
	c.buf.WriteString("ID"+c.s+"Date"+c.s+"Author"+c.s+"Content"+c.s+"Score"+c.s+"Replies\n")
	cPrintln("FilePath: "+filepath)
}

func (c *CSVProvider) AddTuples(tp []df.ProviderTuple){
	d := `"`
	for _, x := range tp{
		c.buf.WriteString(x.Id+c.s)
		c.buf.WriteString(d+x.Date.Format("Sun Feb 05 00:48:12 +0000 2017")+d+c.s)
		c.buf.WriteString(d+x.Author+d+c.s)
		c.buf.WriteString(d+x.Content+d+c.s)
		c.buf.WriteString(strconv.FormatInt(x.Score,10)+c.s)
		c.buf.WriteString(strconv.FormatInt(x.Replies,10)+c.s)
		c.buf.WriteString("\n")
	}
}

func (c *CSVProvider) AddTuple(tp df.ProviderTuple){
	d := `"`
	c.buf.WriteString(tp.Id+c.s)
	c.buf.WriteString(d+tp.Date.Format("Sun Feb 05 00:48:12 +0000 2017")+d+c.s)
	c.buf.WriteString(d+tp.Author+d+c.s)
	c.buf.WriteString(d+tp.Content+d+c.s)
	c.buf.WriteString(strconv.FormatInt(tp.Score,10)+c.s)
	c.buf.WriteString(strconv.FormatInt(tp.Replies,10)+c.s)
	c.buf.WriteString("\n")
}


func (c *CSVProvider) Finish(){
	err := ioutil.WriteFile(c.filepath, c.buf.Bytes(),0644)
	if err != nil {
		panic(err)
	}
}
