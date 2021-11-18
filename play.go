package main

import (
	"fmt"
	"regexp"
)

func main() {
	exp := regexp.MustCompile(`Complete the (.*?)(?: job)*? on the (.*?) difficulty or above.`)

	match := exp.FindStringSubmatch("Complete the Four Stores job on the Death Wish difficulty or above.")
	for _, m := range match {
		fmt.Println(m)
	}
	fmt.Printf("map[%s]%s.", match[1], match[2])
}
