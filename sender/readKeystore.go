package sender

import (
	"io/ioutil"

	"github.com/golang/glog"
)

func readKeystore(path string) []string {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		glog.Fatal(err)
	}

	var result []string
	for _, f := range files {
		result = append(result, "./keystore/"+f.Name())

	}
	glog.Info("result ", result)

	return result
}
