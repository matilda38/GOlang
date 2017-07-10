package main

import (
	"LogAgentInGo/LogAgentInGo"
)


func main(){

	message := make(chan string)
	//if there's no defined protocol, the default is udp
	g:= LogAgentInGo.New(LogAgentInGo.Config{})

	go g.InputMessage(message)
	go g.ProcessMessage(message)

	go g.Consume()

	g.Send()
}
