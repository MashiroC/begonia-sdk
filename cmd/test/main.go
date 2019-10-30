// Time : 2019/10/25 13:00
// Author : MashiroC

// test
package main

import (
	"bytes"
	"fmt"
	"github.com/MashiroC/begonia-sdk"
	"net/http"
)

// testkeycenter.go something

var (
	size = 1000000

	pool = make(chan bool, size)
)

type Person struct {
	Age  int
	Name string
}

func test() (res interface{}, err error) {
	tmp := make(map[string]interface{}, 2)
	//b,_:=json.Marshal(Person{
	//	Age:  18,
	//	Name: "hhh",
	//})
	tmp["Age"] = float64(18)
	tmp["Name"] = "hhh"
	return tmp, nil
}

func main() {
	person := Person{}
	err := begonia.Result(test()).Bind(&person)
	if err!=nil{
		fmt.Println("?????")
	}
	fmt.Println(person)
	//for i:=0;i<size;i++{
	//	pool<-true
	//}
	//fmt.Println("ok")
	//wait := make(chan bool)
	//for {
	//	//<-pool
	//	go work()
	//
	//}
	//<-wait
}

func work() {
	//defer func() {
	//	pool<-true
	//}()
	//fmt.Println("咕咕你站没了")
	str := `{-"type": "num_stu", "data": null}`

	buff := bytes.NewBuffer([]byte(str))

	http.Post("http://cqupt.online/fuck/q", "application/json", buff)
	//if err!=nil{
	//	log.Fatal(err)
	//}
	//body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
}
