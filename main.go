package main

import "fmt"

func main() {

	n := newNode()
	n.Insert("/hello")
	n.Insert("/hello/make/")
	n.Insert("/hi/hello/main")
	n.Insert("/hi/fix/best/")
	n.Insert("/hi/fix/bestgroup")
	n.Insert("/hi/something/")
	n.Insert("/hi/fix/bestgrou")

	rt := NewRadixTraverse(n)
	rt.Traverse()

	fmt.Println(n.Search("/hi/fix/bestgro"))
}
