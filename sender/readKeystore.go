package sender

import (
	"io/ioutil"
	"log"
)

func readKeystore(path string) []string {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var result []string
	for _, f := range files {
		result = append(result, "./keystore/"+f.Name())

	}
	//log.Println("readKeystore result ", result)

	return result
}
