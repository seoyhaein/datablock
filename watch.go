package datablock

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"sync"
)

var (
	watchCount int
	// TODO 이걸 전역적으로 둬야 하나??? 아래 TODO 적용후 삭제
	isWatching bool
)

func AddWatch(watcher *fsnotify.Watcher, path string, maxWatchCount int, mu *sync.Mutex) error {
	mu.Lock()
	defer mu.Unlock()

	if watchCount >= maxWatchCount {
		log.Println("Warning: Maximum watch folder count reached. Cannot add more:", path)
		return nil
	}

	err := watcher.Add(path)
	if err == nil {
		watchCount++
		log.Println("Added watch:", path, "Current watch folder count:", watchCount)
	}
	return err
}

func RemoveWatch(watcher *fsnotify.Watcher, path string, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	err := watcher.Remove(path)
	if err == nil {
		watchCount--
		log.Println("Removed watch:", path, "Current watch folder count:", watchCount)
	} else {
		log.Println("Failed to remove watch:", err)
	}
}

// StartWatching TODO 테스트 필요 watchCount, isWatching 를 담고 있는 구조체로 만들자.
func StartWatching(paths []string, maxWatchCount int, mu *sync.Mutex) (*fsnotify.Watcher, error) {
	mu.Lock()
	defer mu.Unlock()

	if isWatching {
		log.Println("Already watching.")
		return nil, nil
	}
	// 새로운 watcher 생성
	var err error
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	isWatching = true
	return watcher, nil
}

// StopWatching - 감시를 중지하는 함수
func StopWatching(watcher *fsnotify.Watcher, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	if watcher == nil || !isWatching {
		log.Println("No active watcher to stop.")
		return
	}

	err := watcher.Close()
	if err != nil {
		log.Println("Failed to close watcher:", err)
		return
	}

	isWatching = false
	log.Println("Stopped watching.")
}

// WatchEvents 이벤트 처리 루프 TODO 기억용으로 넣어둠 수정 해야힘.
func WatchEvents(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("Event:", event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}
