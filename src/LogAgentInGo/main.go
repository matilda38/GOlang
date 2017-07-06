package main

import "LogAgentInGo/LogAgentInGo"

func main(){
	doc := "{ \"version\" : \"1.1\", \"host\": \"example.org\", \"short_message\": \"A short message that helps you identify what is going on\", \"full_message\": \"Backtrace here\n\nmore stuff\"}"

	//if there's no defined protocol, the default is udp
	g:= LogAgentInGo.New(LogAgentInGo.Config{
		Protocol: "TCP",
	})

	for i:=0; i< 100;i++{
		g.Log(doc)
	}
}
