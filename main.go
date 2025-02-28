package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	c "github.com/seoyhaein/datablock/config"
	d "github.com/seoyhaein/datablock/db"
	"os"
	"path"
	"path/filepath"
)

func main() {

	//_ = RemoveDBFile("file_monitor.db")
	// db connection foreign key 설정을 위해 PRAGMA foreign_keys = ON; 설정을 해줘야 함.
	db, err := d.ConnectDB("sqlite3", "file_monitor.db", true)
	if err != nil {
		os.Exit(1)
	}
	err = d.InitializeDatabase(db)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		/*if err = d.ClearDatabase(db); err != nil {
			//log.Fatal("failed to clear db:", err)
			os.Exit(1)
		}*/

		if err := db.Close(); err != nil {
			//log.Fatal("failed to close db:", err)
			os.Exit(1) // defer 내부에서도 os.Exit 사용 가능
		}
	}()

	config, err := c.LoadConfig("config.json")
	if err != nil {
		os.Exit(1)
	}

	// 테스트로 빈파일 생성
	// 기존 파일이 생성되어 있을 경우 권한 설정을 안해줌. 버그지만 고치지 않음.
	testFilePath := config.RootDir
	testFilePath = filepath.Join(testFilePath, "testFiles/")
	testFilePath = path.Clean(testFilePath)
	d.MakeTestFiles(testFilePath)
	d.MakeTestFilesA("/test/baba/")

	ctx := context.Background()
	// exclusion 은 보안상 여기다가 넣어둠. TODO 일단 생각은 해보자.
	exclusions := []string{"*.json", "invalid_files", "*.csv", "*.pb"}
	dbApis := NewDBApis(config.RootDir, nil, exclusions)
	err = dbApis.StoreFoldersInfo(ctx, db)
	if err != nil {

		os.Exit(1)
	}

	// TODO 같지 않을때 처리 해줘야 함. db 를 업데이트 해줘야 함.
	_, _, _, _, err = dbApis.CompareFoldersAndFiles(ctx, db)
	if err != nil {
		os.Exit(1)
	}

}

// RemoveDBFile 주어진 DB 파일을 삭제함.
// filePath: 삭제할 DB 파일의 경로.
func RemoveDBFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove DB file %s: %w", filePath, err)
	}
	return nil
}
