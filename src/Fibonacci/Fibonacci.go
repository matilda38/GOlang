package main

func fibonacci(quit, c chan int){
	x,y := 0,1
	for{
		select{
			case c <- x:
				x,y = y, x+y
			case <-quit:
				return
		}
	}
}

func main(){
	c := make(chan int)
	quit := make(chan int)
	go func(){
		for i:=0;i<10;i++{
			println(<-c)
		}
		quit <- 0
	}()
	fibonacci(quit, c)
}


