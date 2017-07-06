package main

import (
	"LogAgentInGo"
)


func main(){
	doc := "{ \"version\" : \"1.1\", \"host\": \"example.org\", \"short_message\": \"A short message that helps you identify what is going on\", \"full_message\": \"Backtrace here\n\nmore stuff\"}"
	g:= LogAgentInGo.New(LogAgentInGo.Config{})
	g.Log(doc)
}
