package main

import (
	"github.com/seoyhaein/datablock/watch"
	"log"
)

func main() {
	watch.Init("config.json")
	defer watch.StopWatching()
	go watch.WatchEvents()
	errChan := make(chan error)
	defer close(errChan)

	go watch.ProcessEvents(errChan)
	// 에러 채널에서 에러 수신
	go func() {
		for err := range errChan {
			log.Printf("Error received: %v", err)
			// 추가적인 에러 처리 로직
		}
	}()
}
