package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

func main() {
	var string string
	string = "f"
	num, err := strconv.Atoi(string)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(reflect.TypeOf(num))
	fmt.Println(num)
}
