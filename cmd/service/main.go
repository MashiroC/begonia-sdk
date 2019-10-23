// Time : 2019/10/11 18:40
// Author : MashiroC

// main
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-sdk"
	"log"
	"reflect"
)

// main.go something

type HelloService struct {
}

type PersonChild struct {
	Test string
}

type Person struct {
	Name  string
	Age   int
	Child PersonChild
}

func (hs *HelloService) TestInt(a, b int) int {
	fmt.Println(a, b)
	return a + b
}

func (hs *HelloService) TestIntWithErr(a, b int) (int, error) {
	return 0, nil
}

func (hs *HelloService) TestUint(a, b uint) uint {
	fmt.Println(a, b)
	return a + b
}

func (hs *HelloService) TestInt8(a, b int8) int8 {
	fmt.Println(a, b)
	return a + b
}

func (hs *HelloService) TestUint8(a, b uint8) uint8 {
	fmt.Println(a, b)
	return a + b
}

func (hs *HelloService) TestFloat64(a, b float64) float64 {
	fmt.Println(a, b)
	return a * b
}

func (hs *HelloService) TestFloat32(a, b float32) float32 {
	fmt.Println(a, b)
	return a * b
}

func (hs *HelloService) TestStruct(per Person) Person {
	fmt.Println(per)
	per.Child = PersonChild{Test: "test"}
	return per
}

func (hs *HelloService) TestStructWithErr(per Person) (Person, error) {
	fmt.Println(per)
	return Person{}, errors.New("testErr")
}

func (hs *HelloService) TestEmpty(per Person) {
	fmt.Println(per)
}

func (hs *HelloService) TestEmptyWithErr(test string) error {
	return entity.CallError{ErrCode: "123456789", ErrMessage: "testError"}
}

//func (hs *HelloService) TestIntWithErr(a,b int)

func (hs *HelloService) World() (res string) {
	return "helloworld!"
}

func testCall(in ...interface{}) {
	b, err := json.Marshal(in)
	if err != nil {
		log.Fatal(err.Error())
	}
	realCall(b)
}

func realCall(b []byte) {
	var in []interface{}
	err := json.Unmarshal(b, &in)
	if err != nil {
		log.Fatal(err.Error())
	}
	v := reflect.ValueOf(Hello)
	param := make([]reflect.Value, len(in))
	for i, v := range in {
		param[i] = reflect.ValueOf(v)
	}
	res := v.Call(param)
	fmt.Println(res[1].IsNil())
}

func testStruct(data map[string]interface{}) {

}

func Hello(name string) (res string, err error) {
	return "hello " + name, nil
}

func main() {
	cli := begonia.Default(":4949")
	cli.Sign("Hello", &HelloService{})
	cli.KeepConnect()
}
