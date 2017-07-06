package main

import (
	"Prac2"
)


func main(){
	doc := "{ \"version\" : \"1.1\", \"host\": \"example.org\", \"short_message\": \"A short message that helps you identify what is going on\", \"full_message\": \"Backtrace here\n\nmore stuff\"}"
	g:= prac2.New(prac2.Config{})
	g.Log(doc)
}
