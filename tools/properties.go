package tools

import(
	"strings"
	"io/ioutil"
)

func GetConfigFromFile(filepath string)map[string]string{
	dat, err := ioutil.ReadFile(filepath)
	if err != nil{
		panic(err)
	}
	m := make(map[string]string)

	for _, line := range strings.Split(string(dat), "\n"){
		splitted := strings.Split(line, "=")
		if len(splitted) < 2{
			continue
		}
		key   := strings.Trim(splitted[0], " ")
		value := strings.Trim(splitted[1], " ")
		m[key] = value
	}
	return m
}
