# datablock

## dependencies
- fsnotify
- install : go get -u github.com/fsnotify/fsnotify
- golang.org/x/sys 이것도 자동으로 설치됨.
- fsnotify v1.8.0, x/sys v0.26.0

## 개발 시작
- 먼저 사용자 입력 값은 최소한으로 받도록 한다. main 으로 일단 접근해서 향후 필요에 따라서 붙이는 방향으로 간다.  

## TODO
- 테스트 작성해야 함.  
- Qodana 로 테스트 결과 해결 해야 함.  
- firstwalk 적용 해야 함.  
- 문서 정리해야 함. rule 에 대한 상세 설명문서도 만들어야 함.  
- main 에서 부터 이제 어떻게 다시 시나리오를 만들어 갈지 구상 해야함.  
- 검색 기능 넣고, grpc 연동 진행.
- 기초 grpc 넣어두고 grpc 프로젝트 만들고 고도화 함.
- 로그 정리 현재, 좀 불명확한 로그들이 존재함. 
- 사용자 편의성 생각할 것. exit 을 넣으면 종료 되는데 이게 로그가 올라오면 사라짐.

## 수정할 것들.
~~- 일단 굴러가게만 하자.~~  
~~- invalid 해줘야 함.~~  
~~- 특정파일에 동일 규칙의 파일들이 있는지 검사해야 할까?~~  
~~- 하나의 디렉토리에 하나의 블럭만 존재하게 해야 하는가? 복수개도 존재할 수 있도록 해야 하지 않은가?~~  
~~- 디렉토리안에서 여러개의 블럭이 존재할 수 있다고 보는데. 이건 정신이 맑아지면 컨디션 좋을대 살펴보자.~~  
- 디렉토리는 이름이 unique 함 따라서 이것은 파일명을 잡는데 중요한 기준점임.  
- 정신이 없어서 테스트 코드 살펴봐야 한다.  
~~- 필터링하는 거 해줘야 함.~~

## 생각해봐야 할 것
- 조금더 사용자 친화적으로 할 수 있는 방법이 없는지 생각해보자.  
- 디렉토리 관련해서 개발을 시작해야 할 것 같다.    
## update 
- 로그 관련 표준 정하자.  


````txt
	//path := "/tmp/testfiles"
	// 테스트로 빈파일 생성
	//MakeTestFiles(path)

	//err := rule.ApplyRule(path)
	//if err != nil { // 에러 발생 시 종료
	//	os.Exit(1)
	//}
````

- 일단 watch 를 분리할 예정임.
- 
### maing.go backup
```go
package main

import (
	"bufio"
	"context"
	"github.com/seoyhaein/datablock/watch"
	"log"
	"os"
	"sync"
)

func main() {

	// Context 와 WaitGroup 정의
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	defer cancel()

	// buffer 둘것
	errChan := make(chan error, 100) // 버퍼 크기 추가
	//defer close(errChan)

	// 초기화 및 이벤트 처리
	watch.Init("config.json")
	defer watch.StopWatching()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		watch.ListenEvents(ctx, errChan)
	}(ctx)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		watch.ProcessEvents(ctx, errChan)
	}(ctx)

	// 에러 채널에서 에러 수신
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping error listener...")
				return
				/*case err, ok := <-errChan:
				if !ok {
					log.Println("Error channel closed.")
					return
				}
				// 모든 에러 출력
				log.Printf("Error received: %v", err)*/
			}
		}
	}(ctx)

	// 사용자 입력 대기
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Type 'exit' to stop the watch service.")
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading input: %v", err)
				continue
			}
			if input == "exit\n" {
				log.Println("Exiting watch service...")
				cancel() // Context 를 취소하여 모든 고루틴 종료
				return
			}
		}
	}()
	// TODO 당분가 삭제 금지. 이렇게 한 것을 아래 와 같이 넣었다. 기억용으로 일단 남겨둔다.
	/*
		wg.Add(1)
		go func(errChan <-chan error, ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case err, ok := <-errChan:
					if !ok {
						log.Println("Error channel closed.")
						return
					}
					log.Printf("Error received: %v", err)
				case <-ctx.Done():
					log.Println("Stopping error listener...")
					// 남은 에러 처리
					for err := range errChan {
						log.Printf("Error received: %v", err)
					}
					log.Println("All errors processed. Exiting listener.")
					return
				}
			}
		}(errChan, ctx)
	*/

	// TODO 모든 고루틴이 종료될 때까지 대기 일단 주석 처리함. 테스트 끝나면 삭제.
	// log.Println("Wait...")
	wg.Wait()
	// TODO 순서를 바꾸면 에러남. 일단 테스트 진행 해야함.
	close(errChan)
	log.Println("Close Channel")
	for err := range errChan {
		log.Printf("Error received: %v", err)
	}
	log.Println("Watch service stopped gracefully.")
}

```