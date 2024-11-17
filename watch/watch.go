package watch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Config struct {
	MaxWatchCount int
	RootDir       string
}

var (
	watchCount int // 감시 디렉토리 갯수
	isWatching bool
	eventQueue []fsnotify.Event // 이벤트 큐
	queueMu    sync.Mutex       // 큐에 접근하는 뮤텍스
	once       sync.Once        // 한 번만 실행되도록 제어
	// TODO paths 의 중복확인을 해줘야 함. 물론 중복된 경우는 넘어간다고 하지만, 추가 삭제에 대한 중복확인은 해줘야 함.
	watchedPaths  = make(map[string]bool)
	maxWatchCount int
	rootDir       string
	watcher       *fsnotify.Watcher
)

func Init(path string) {

	config, err := loadConfig(path)
	if err != nil {
		log.Println("loadConfig error:", err)
		return
	}
	maxWatchCount = config.MaxWatchCount
	rootDir = config.RootDir
	watcher, err = StartWatching(nil)
	if err != nil {
		log.Println("StartWatching error:", err)
		return
	}
	log.Println("Initialization completed successfully.")
}

func AddWatch(watcher *fsnotify.Watcher, path string, mu *sync.Mutex) error {
	if watcher == nil {
		return fmt.Errorf("watcher is nil")
	}

	if mu != nil {
		mu.Lock()
		defer mu.Unlock()
	}

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

func RemoveWatch(watcher *fsnotify.Watcher, path string, mu *sync.Mutex) error {
	if watcher == nil {
		return fmt.Errorf("watcher is nil")
	}

	if mu != nil {
		mu.Lock()
		defer mu.Unlock()
	}

	err := watcher.Remove(path)
	if err == nil {
		watchCount--
		log.Println("Removed watch:", path, "Current watch folder count:", watchCount)
		return nil
	}
	log.Println("Failed to remove watch:", err)
	return err
}

// StartWatching TODO 테스트 필요 watchCount, isWatching 를 담고 있는 구조체로 만들자.
func StartWatching(paths []string) (*fsnotify.Watcher, error) {
	if isWatching {
		log.Println("Already watching.")
		return nil, errors.New("already watching")
	}
	// 새로운 watcher 생성
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	isWatching = true

	return watcher, nil
}

// StopWatching - 감시를 중지하는 함수
func StopWatching() error {
	if watcher == nil || !isWatching {
		log.Println("No active watcher to stop.")
		return fmt.Errorf("No active watcher to stop	")
	}
	err := watcher.Close()
	if err != nil {
		log.Println("Failed to close watcher:", err)
		return err
	}
	isWatching = false
	log.Println("Stopped watching.")
	return nil
}

// ListenEvents 이벤트 처리 루프
func ListenEvents(ctx context.Context, errChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping ListenEvents...")
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("Event:", event)
			addToQueue(event) // 이벤트를 큐에 추가
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
			select {
			case errChan <- err:
			case <-ctx.Done():
				return // 컨텍스트가 취소된 경우 반환
			}
		}
	}
}

/*
func ListenEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping ListenEvents...")
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("Event:", event)
			addToQueue(event) // 이벤트를 큐에 추가
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}
*/
// addToQueue 이벤트를 큐에 추가하는 함수
func addToQueue(event fsnotify.Event) {
	queueMu.Lock()
	defer queueMu.Unlock()
	eventQueue = append(eventQueue, event) // 큐에 이벤트 추가
}

// ProcessEvents 큐에서 이벤트를 하나씩 처리하는 함수 - 고루틴 사용
func ProcessEvents(ctx context.Context, errChan chan<- error) {
	var mu sync.Mutex

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping ProcessEvents...")
			return
		default:
			queueMu.Lock()
			if len(eventQueue) == 0 {
				queueMu.Unlock()
				time.Sleep(100 * time.Millisecond)
				continue
			}
			event := eventQueue[0]
			eventQueue = eventQueue[1:]
			queueMu.Unlock()

			// 이벤트 처리
			log.Println("처리 중인 이벤트:", event)
			var err error
			switch {
			case event.Has(fsnotify.Create):
				log.Println("File created:", event.Name)
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					if err := AddWatch(watcher, event.Name, &mu); err != nil {
						log.Printf("Failed to add watch for directory %s: %v", event.Name, err)
					}
				}
			case event.Has(fsnotify.Remove):
				log.Println("File removed:", event.Name)
				if err := RemoveWatch(watcher, event.Name, &mu); err != nil {
					log.Printf("Failed to remove watch for %s: %v", event.Name, err)
				}
			case event.Has(fsnotify.Rename):
				log.Println("File renamed:", event.Name)
				if err := RemoveWatch(watcher, event.Name, &mu); err != nil {
					log.Printf("Failed to remove watch for renamed file %s: %v", event.Name, err)
				}
			case event.Has(fsnotify.Write):
				log.Println("File modified:", event.Name)
			case event.Has(fsnotify.Chmod):
				log.Println("File attributes changed:", event.Name)
			}

			// 에러가 발생했을 경우 에러 채널로 전달
			if err != nil {
				select {
				case errChan <- fmt.Errorf("error processing event %s: %w", event.Name, err):
				case <-ctx.Done():
					return // 컨텍스트가 취소된 경우 반환
				}
			}
		}
	}
}

/*
errChan := make(chan error)

	// 고루틴에서 ProcessEvents 함수 실행
	go ProcessEvents(errChan)

	// 에러 채널에서 에러 수신
	go func() {
		for err := range errChan {
			log.Printf("Error received: %v", err)
			// 추가적인 에러 처리 로직
		}
	}()
*/

// FirstWalk TODO 최초 디렉토리 검사 및 관련 파일 만들어 주기. 감시할때는 별도의 메서드드로 관련 파일 수정 및 만들어 주어야함.
func FirstWalk(watcher *fsnotify.Watcher) error {
	var mu sync.Mutex
	if watcher == nil {
		return fmt.Errorf("watcher is nil")
	}

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return AddWatch(watcher, path, &mu)
		}
		return nil
	})

	return err
}

// FirstWalk1 TODO once 부분 리턴에 관해서 살펴봐야 함. 익명함수 내에서 리턴도 같이 봐야 함.
// 관련파일을 작성해주는 메서드를 만들어 줘야 함.
// TODO context 넣어 주어야 함.
func FirstWalk1(watcher *fsnotify.Watcher) error {
	var mu sync.Mutex
	var err error

	once.Do(func() {
		if watcher == nil {
			err = fmt.Errorf("watcher is nil")
			return // 익명 함수에서 반환하여 `Do` 블록의 실행 중단
		}

		// 디렉토리 순회 및 감시 추가
		err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return AddWatch(watcher, path, &mu)
			}
			return nil
		})
	})

	return err // once.Do 이후 최종 반환
}

func loadConfig(filename string) (*Config, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode JSON data into a generic map
	decoder := json.NewDecoder(file)
	configData := make(map[string]interface{})
	if err := decoder.Decode(&configData); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	// Create a Config instance
	config := &Config{}

	// Process MaxWatchCount
	if value, ok := configData["MaxWatchCount"]; ok {
		if floatValue, valid := value.(float64); valid {
			config.MaxWatchCount = int(floatValue)
		} else {
			log.Printf("Invalid type for MaxWatchCount in configuration")
			return nil, fmt.Errorf("invalid type for MaxWatchCount")
		}
	} else {
		log.Printf("Missing MaxWatchCount in configuration")
		return nil, fmt.Errorf("missing MaxWatchCount in configuration")
	}

	// Process rootDir
	if value, ok := configData["rootDir"]; ok {
		if stringValue, valid := value.(string); valid {
			config.RootDir = stringValue
		} else {
			log.Printf("Invalid type for rootDir in configuration")
			return nil, fmt.Errorf("invalid type for rootDir")
		}
	} else {
		log.Printf("Missing rootDir in configuration")
		return nil, fmt.Errorf("missing rootDir in configuration")
	}

	// Return the populated config
	return config, nil
}
