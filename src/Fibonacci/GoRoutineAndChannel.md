## 고루틴과 채널
_이제야 조금 고루틴과 채널에 대해서 이해할 수 있을 듯하다. 까먹기 전에 정리해보려고 한다_


고루틴은 Go 언어에서 제공하는 동시성(Concurrency)에서 등장한 개념이다. 일반적으로 병렬처리를 위해 사용하는 Thread와 다르게, 여러개의 고루틴을 하나의 OS 쓰레드에 할당하여 처리하게끔 하므로 쓰레드 생성, 삭제에 대한 비용이 적다.
 
사용법은, 함수 앞에 go 를 붙여주면 되는데

>   
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

메인함수도 하나의 고루틴이므로, 현재 익명함수 고루틴과 메인함수 고루틴 2개의 고루틴이 동시에 돌고 있다고 할 수 있다.
**함수 앞에 go를 붙이면, 그 행의 작업이 끝나는 것을 기다리지 않고 다음 작업으로 넘어간다.** 마치 c#의 async 키워드와 java의 SynchronousQueue와 비슷한 개념이다.

작업이 끝났을 때 리턴 값은 **채널(Channel)** 을 이용하여 받을 수 있는데, 이는 마치 고루틴간의 파이프(pipe)와 같은 개념이다.

![](https://cdn-images-1.medium.com/max/1200/1*GWYUFH14uOVLNHY-L1tv2w.jpeg)

Go는 Shared Memory, 공유 메모리 개념을 사용하기 때문에 고루틴들은 전부 하나의 메모리 영역을 공유해서 사용한다. 채널을 이용하여 한 루틴에서 값을 채널에 넣고, 다른 루틴은 채널에서 그 값을 받아(꺼내) 사용할 수 있다.

채널을 이용하면, 채널로 데이터가 오기 전까지는 고루틴 처리가 멈추고, 채널에 데이터가 들어오면 자동으로 고루틴이 움직이기 때문에 동기화 이슈도 간편하게 해결할 수 있다.

>

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

피보나치 수열을 나타내는 위 코드를 보면서 예시를 들어보겠다.

메인 고루틴, 익명함수 고루틴 두 개의 고루틴이 존재한다. c와 quit 라는 이름의 채널을 생성한다. c, quit는 두 고루틴간 리턴값등 값을 공유할 수 있게 해줄 파이프 역할을 할 채널이다.
이후엔 익명함수 고루틴에서 println함수에 채널에서 받은 값을 출력하려고 한다. 이때 c 채널에는 아직 값이 없기 때문에, 잠시 멈춘다. 다음 작업으로 이동하여 메인 고루틴은 Fibo함수를 호출하는데,

cur, next 변수 생성, 초기화 후에 무한반복 for 문에서 select문을 수행하여 준비된 case문을 선택한다.
이때, c 채널은 값을 받을 준비가 되어있으므로 첫번째 case문이 수행되어 c 채널에 cur 값이 들어가고, cur, next값은 case 문 내 로직에 의해 각각 next, cur+next 값으로 업데이트(수정)된다.
 
이때 잠시 중단되었던 익명함수 고루틴의 println 함수가 다시 재개되어 Fibo에서 c 채널에 넣었던 값을 받아(꺼내어) 출력해준다.

익명함수 루틴은 반복문에 의해 다시 println이 수행되고, c 채널에 이제 값이 없으므로 다시 잠시 중단된다. 메인 고루틴에서는 다시 무한반복 for 문에 의해 select 문이 수행되고, 준비된 c 채널에 cur값이 들어간다. 아까와 같은 로직이 수행된다.

또한, 이전에 언급했듯이, 중단되었던 익명함수 고루틴의 println 함수는 재개되어, c 채널의 값을 꺼내어 출력한다.

위와 같은 패턴이 반복되다가 익명함수에서 10번 println을 수행하고, 반복문을 탈출했을때 quit 채널에 0을 넣는다. Fibo 함수의 select 문에서는 c 채널에 값이 있어 cur 값을 넣지 못하므로, 즉 준비가 되어있지 않고, 두번째 case문인 <-quit의 경우에는 quit 채널에 값이 있어, 즉 준비가 되어있으므로, 두번째 case 문을 수행하여 return(반환) 이 된다.


<p color = "blue">ㅎ.ㅎ 드디어 이해 YEAH!!</p>   