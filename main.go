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

	// 모든 고루틴이 종료될 때까지 대기
	log.Println("Wait...")
	wg.Wait()
	// TODO 순서를 바꾸면 에러남. 일단 테스트 진행 해야함.
	close(errChan)
	log.Println("Close Channel")
	for err := range errChan {
		log.Printf("Error received: %v", err)
	}
	log.Println("Watch service stopped gracefully.")
}
