package config

import (
	"fmt"
	"testing"
)

func TestSum(t *testing.T) {
	path := "../config.json"
	fmt.Printf("%#v \n", GetConfig(path))
}
