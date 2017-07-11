package Dashinsert

import (
	"strings"
	"strconv"
	"fmt"
)
/*
func Do(input string){
	var ans string = ""
	for i, v := range input{
		if i == len(input) -1{
			ans += string(v)
			break
		}
		if int(v) % 2 == 0 && int(input[i+1])%2 ==0{
			ans += string(v) + "*"
		} else if int(v)%2 ==1 && int(input[i+1])%2 ==1{
			ans += string(v) + "-"
		} else{
			ans += string(v)
		}
	}
	print(ans)
}
*/
func Do(input string){
	exSlice := strings.Split(input,"")
	var anSlice []string
	for i, v := range exSlice{

		if i == len(exSlice) -1{
			anSlice = append(anSlice, v)
			break
		}

		n, err:= strconv.ParseInt(v,10,32)
		n2, err2 := strconv.ParseInt(exSlice[i+1], 10,32)

		if err!= nil || err2 != nil{
			return
		}

		if n % 2 == 0 && n2 % 2 == 0{
			anSlice = append(anSlice, v + "*")
		} else if n % 2 == 1 && n2 % 2 == 1{
			anSlice = append(anSlice, v + "-")
		} else {
			anSlice = append(anSlice, v)
		}
	}
	fmt.Print(strings.Join(anSlice, ""))
}