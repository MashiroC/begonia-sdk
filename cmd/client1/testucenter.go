// Time : 2019/10/30 15:22
// Author : MashiroC

// main
package main

import (
	"fmt"
	"github.com/MashiroC/begonia-sdk"
)

// testucenter.go something

var (
	cli         *begonia.Client
	Verify      begonia.RemoteFun
	FindStudent begonia.RemoteFun
	FindRedid   begonia.RemoteFun
	Bind        begonia.RemoteFun
	UnBind      begonia.RemoteFun
)

func main() {
	res, err := Verify("2017211573", "020311")
	fmt.Println(res, err)

	//res, err = FindStudent("xbs", "testopenid")
	//fmt.Println(res, err)

	//res, err = FindRedid("xbs", "testopenid")
	//fmt.Println(res, err)

	//res,err := Bind("xbs","testtest","2017211573","02031X")
	//fmt.Println(res,err)

	//res,err:=UnBind("xbs","1")
	//fmt.Println(res,err)
}

func init() {
	cli = begonia.Default(":8080")

	ucenter := cli.Service("ucenter")

	Verify = ucenter.FunSync("Verify")
	FindStudent = ucenter.FunSync("FindStudent")
	FindRedid = ucenter.FunSync("FindRedid")
	Bind = ucenter.FunSync("Bind")
	UnBind = ucenter.FunSync("Unbind")
}
