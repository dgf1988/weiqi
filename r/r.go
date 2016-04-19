package main

import (
	"reflect"
	"fmt"
	"time"
)

func main() {
	us := []User{
		User{1, "dgf"}, User{2, "lsl"},
	}
	TestArrayStruct(us)
}

type Ts struct {
	Id int
	Name string
	age int
	Birth time.Time
}

func TestStruct() {
	var ts = Ts{}
	vts := reflect.ValueOf(&ts)
	fmt.Println(vts.Type())
	fmt.Println(vts, vts.CanSet())
	fmt.Println(vts.Elem(), vts.Elem().CanSet())
	vts = vts.Elem()
	for i := 0 ; i < vts.NumField(); i ++ {
		fmt.Println(vts.Field(i), vts.Field(i).Type(), vts.Field(i).CanSet(), vts.Field(i).Addr(), vts.Field(i).Addr().Kind())
	}
}

func TestArray() {
	arr := []interface{}{1, 2, "a", "b"}
	v := reflect.ValueOf(arr)
	fmt.Println(v.CanSet(), v.CanAddr(), v.CanInterface())
	if v.Kind() == reflect.Slice {
		for i := 0 ; i < v.Len() ; i ++ {
			elem := v.Index(i)
			fmt.Println("index", elem.Kind(), elem.Type(), elem.CanSet(), elem.CanAddr(), elem.CanInterface())

			elem = elem.Elem()
			fmt.Println("elem", elem.Kind(), elem.Type(), elem.CanSet(), elem.CanAddr(), elem.CanInterface())
		}
	}
}

type User struct {
	Id int
	Name string
}

func TestArrayStruct(arr interface{}) {
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice:
		v := reflect.ValueOf(arr)
		for i := 0; i < v.Len(); i ++ {
			fmt.Println(i)
			e := v.Index(i)
			fmt.Println(e.Kind(), e.Type(), e.CanSet(), e.CanAddr(), e.CanInterface())
			p := e.Addr().Elem()
			fmt.Println(p.Kind(), p.Type(), p.CanSet(), p.CanAddr(), p.CanInterface())
		}
	default:
		fmt.Println("error")
	}
}