package main

import (
	"log"
	"reflect"
	"time"
)

type TestStruct struct {
	age int64
}

func main() {
	// var ref TestStruct
	// Invoke(&ref, "logfunc", "123")
	var asd time.Time
	log.Println(asd.IsZero(), asd)
	log.Println(time.Now().IsZero())
}

func Invoke(any interface{}, name string, args ...interface{}) {
	inputs := make([]reflect.Value, 1)
	inputs[0] = reflect.ValueOf(args[0])
	reflect.ValueOf(&any).MethodByName("logfunc").Call([]reflect.Value{})
}

func (t *TestStruct) logfunc() {
	log.Println(time.Now())
}
