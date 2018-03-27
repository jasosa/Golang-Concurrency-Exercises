package main

import (
	"fmt"
	"strings"
)

func main() {
	c1 := gen("   hello how are you    ")
	c2 := stripwhitespace(c1)
	c3 := reverse(c2)
	fmt.Println(<-c3)
}

func gen(command string) <-chan string {
	c := make(chan string)
	go func() {
		c <- command
		close(c)
	}()
	return c
}

func stripwhitespace(cin <-chan string) <-chan string {
	cout := make(chan string)
	go func() {
		value := <-cin
		strimmed := strings.Replace(value, " ", "", -1)
		cout <- strimmed
		close(cout)
	}()
	return cout
}

func reverse(cin <-chan string) <-chan string {
	cout := make(chan string)
	go func() {
		value := <-cin
		reversed := func() string {
			var reverse string
			for i := len(value) - 1; i >= 0; i-- {
				reverse += string(value[i])
			}
			return reverse
		}()

		cout <- reversed
		close(cout)
	}()
	return cout
}
