package main

import (
	"LogAgentInGo"
	"fmt"
)


func main(){
	doc := "{ \"version\" : \"1.1\", \"host\": \"example.org\", \"short_message\": \"A short message that helps you identify what is going on\", \"full_message\": \"Backtrace here\n\nmore stuff\"}"

	//if there's no defined protocol, the default is udp
	g:= LogAgentInGo.New(LogAgentInGo.Config{})

	for i:=0; i< 10000;i++{
		g.Log(doc)
	}

	fmt.Print(g.Buffer)
}
