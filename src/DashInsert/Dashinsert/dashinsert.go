package Dashinsert

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
