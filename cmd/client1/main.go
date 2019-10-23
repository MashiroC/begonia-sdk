// Time : 2019/10/12 19:49
// Author : MashiroC

// client1
package main

import (
	"fmt"
	"github.com/MashiroC/begonia-sdk"
)

// main.go something

type PersonChild struct {
	Test string
}

type Person struct {
	Name  string
	Age   int
	Child PersonChild
}

func main(){
	cli:=begonia.Default(":4949")
	helloService:=cli.Service("Hello")

	testIntWithErr := helloService.FunSync("TestIntWithErr")
	res,err:=testIntWithErr(1,1)
	fmt.Println("res",res,err)


	testEmpty :=helloService.FunSync("TestEmpty")
	res,err=testEmpty(Person{
		Name:  "aaaaa",
		Age:   123,
		Child: PersonChild{Test:"test"},
	})
	fmt.Println("res",res,err)
	//
	testEmptyWithErr :=helloService.FunSync("TestEmptyWithErr")
	res,err=testEmptyWithErr("hello")
	fmt.Println("res",res,err)

	world :=helloService.FunSync("World")
	res,err=world()
	fmt.Println("res",res,err)
}