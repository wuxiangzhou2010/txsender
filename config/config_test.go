package config

import (
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {

	fmt.Printf("%#v \n", GetConfig())
}
