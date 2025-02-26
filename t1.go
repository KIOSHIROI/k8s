package main

import (
	"fmt"
	"strings"
)

func main() {
	// var l = []string{
	// 	"127.0.0.1:5000/nginx:latest",
	// 	"127.0.0.1:5000/testimage/redis:latest",
	// 	"127.0.0.1:5000/testimage/busybox:latest",
	// }

	var s = "127.0.0.1:5000/testimage/redis"

	fmt.Println(strings.Join(strings.Split(s, "/")[1:], "/"))
}
