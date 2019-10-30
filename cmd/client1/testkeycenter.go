// Time : 2019/10/12 19:49
// Author : MashiroC

// client1
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/MashiroC/begonia-sdk"
	"strings"
	"time"
)

// testkeycenter.go something

type PersonChild struct {
	Test string
}

type Person struct {
	Name  string
	Age   int
	Child PersonChild
}

func main() {
	cli := begonia.Default(":8080")
	keycenter := cli.Service("Keycenter")

	public := keycenter.FunSync("Public")

	createToken := keycenter.FunSync("CreateToken")
	hello:=keycenter.FunAsync("PublicOld")
	hello(func(res interface{}, err error) {
		fmt.Println("old",res,err)
	})
	res, err := public()

	fmt.Println(res, err)

	info := make(map[string]string, 3)
	info["test1"] = "test1"
	info["test2"] = "test2"
	info["test3"] = "test3"

	res, err = createToken(info, (time.Second * 100).Seconds(), "test")
	fmt.Println(res, err)

	token:=strings.Split(res.([]interface{})[0].(string),".")[0]
	b,_:=base64.StdEncoding.DecodeString(token)
	fmt.Println(string(b))
}
