package main

import "os"

func main() {
	env("USER")
	env("PASS")
}

func env(s string) { println(os.Getenv(s)) }
