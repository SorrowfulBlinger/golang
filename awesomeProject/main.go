package main

import (
	"awesomeProject/server"
	"awesomeProject/template"
	"fmt"
	"log"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}
var bridge = make(chan interface{}, 100)
var closeBridge = make (chan struct{})

func producer(wg *sync.WaitGroup, bridge *chan interface{}, closeBridge *chan struct{}) {
	defer func() {
		close(*bridge)
		close(*closeBridge)
		(*wg).Done()
		fmt.Printf("Producer destroyed\n")
	}()
	for i:=0; i<20; i++ {
		if i == 3 {
			fmtStr := fmt.Sprint("Sending close bridge ", i)
			*bridge <- fmtStr
			// Send message to close channel
			*closeBridge <- struct{}{}
			break;
		} else {
			*bridge <- i
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func consumer(wg *sync.WaitGroup, bridge *chan interface{}, closeBridge *chan struct{}) {
	defer func() {
		(*wg).Done()
		fmt.Printf("Consumer destroyed\n")
	}()
	Label:
	for {
		select {
			case data := <-(*bridge) :
				fmt.Printf("Recieved %v\n", data)
			case <- (*closeBridge) :
				fmt.Printf("Close bridge recieved\n")
				break Label
		}
	}

	/*
	for data:= range bridge {
		fmt.Printf("Recieved %v\n", data)
	}
	*/

}

type myreader interface {
	read() ([]byte, error)
}

type mywriter interface {
	Write(data []byte) error
}

type writerReader interface {
	myreader
	mywriter
}

type concreteObj struct {
	data string
}

// Passing ptr updates the original object but value just is a copy
func (obj *concreteObj) read() ([]byte, error) {
	data := make([]byte, 20)
	data = []byte(obj.data)
	obj.data = "Changed"
	return data, nil
}

/*
func (obj concreteObj) read() ([]byte, error) {
	data := make([]byte, 20)
	data = []byte(obj.data)
	obj.data = "Changed"
	return data, nil
} */

type Test struct {
	Myint int
	MyString string
}

func rec2(d int){
	fmt.Println("calling rec2")
	panic("Throwing")
}

func rec1(test ...int)(error, int) {
	fmt.Println("Calling rec1")
	for m:= range test {
		fmt.Printf("%v", m)
	}

	defer func(val string) {
		fmt.Printf("Defer %v\n", val)
		if err:= recover(); err!=nil {
			fmt.Println("Catching",err)
		}
	}("killme")
	rec2(0)
	return fmt.Errorf("test"), 5
}

func testptr(test *Test, slice ...int) (*int, error) {
	a := 1
	b := 2
	if len(slice) > 0 {
		return &a, nil
	}
	// Go converts stack to heap automatically
	return &b, fmt.Errorf(" Bad adta")
}

func test(inta int , stringa string) (string, error) {
	if inta > 10 {
		return "fail", fmt.Errorf("test failed")
	}
	return "success" , nil
}

func tester(obj interface{}) {
	if _, ok := obj.(*concreteObj); !ok{
		fmt.Println("typecast failed")
	} else {
		fmt.Println("typecast success")
	}
}

func main () {
	var myint int32 = 23
	myfloat := 32.6
	fmt.Printf("%v %T %v %T\n", myint, myint, myfloat, myfloat)
	fmt.Printf("Hello go\n")

	defer log.Println ("Goodbye babes")

	// Arrays copied
	var myarray [6]int = [6]int {10, 20, 30, 50, 60, 70}
	//myarray := {1, 2, 3, 5, 6, 7}


	for k,v := range myarray {
		fmt.Printf("%v %v\n", k,v)
	}

	if myarray[1] > 2 {
		fmt.Printf("%v greater than 2\n", myarray[1])
	} else {
		fmt.Printf("%v Less than 2\n", myarray[1])
	}

	t := Test{ Myint : 20,
		MyString : "sample"}

	//Reference/Ptr
	myslice := []string {"str1", "str2"}

	for k,v := range myslice {
		fmt.Printf("%v %v\n", k,v)
	}

	// Reference /Ptrs
	myMap := map[string]string {"a":"b",
		"c" : "d"}

	for k,v := range myMap {
		fmt.Printf("%v %v\n", k,v)
	}

	fmt.Printf("%v\n",t)

	if k, err := test(12, "test"); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", k)
	}

	if k, err := testptr(&Test{ Myint : 10 , MyString: "test"}, []int {1, 2, 3}...) ; err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", *k)
	}

	rec1(2,2)

	var globalObj interface{}
	globalObj = 10

	switch globalObj.(type) {
	case string : fmt.Println("Success type inferred")
	default : fmt.Println("Failure type not inferred")
	}

	testObj := concreteObj{ data : "Hello world"}
	//var myreaderinterface myreader = &testObj
	var myreaderinterface myreader = &testObj
	dataRead, _:= myreaderinterface.read()
	fmt.Println("Using Reader interface", string(dataRead), testObj.data)

	wg.Add(2)
	go producer(&wg, &bridge, &closeBridge)
	go consumer(&wg, &bridge, &closeBridge)
	wg.Wait()

	template.ExecuteTemplate()

	/* HTTP Server */
	go server.StartHTTPServer()

	// w/o go routine for it to block main thread
	server.StartEchoSocket()
}