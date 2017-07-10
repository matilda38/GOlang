package main

func Fibo(quit, c chan int){
	cur,next := 0,1
	for{
		select {
			case c <- cur:
				cur, next = next, cur+next
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

	Fibo(quit, c)
}