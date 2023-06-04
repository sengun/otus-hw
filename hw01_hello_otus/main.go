package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	helloString := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(helloString)
}
