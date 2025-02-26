package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	c "github.com/seoyhaein/datablock/config"
	d "github.com/seoyhaein/datablock/db"
	r "github.com/seoyhaein/datablock/rule"
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
	err = d.InitializeDatabase(db) // TODO 향후 api 로 따로 빼놓아야 할듯. ConnectDB 는 기본으로 해주고 db 를 초기화하고 데이터를 넣어주는 api 를 따로 만들어 주어야 할듯.
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
	// TODO 메서드들을 모아주는 거 생각하자. 일단은 이렇게 해놓음. 향후 api 로 따로 빼놓아야 할듯. ConnectDB 는 기본으로 해주고 db 를 초기화하고 데이터를 넣어주는 api 를 따로 만들어 주어야 할듯.
	exclusions := []string{"*.json", "invalid_files", "*.csv", "*.pb"}

	err = d.StoreFoldersInfo(ctx, db, config.RootDir, nil, exclusions)
	if err != nil {
		fmt.Println("StoreFoldersInfo Error")
	}

	// TODO api 로 만들어 두어야 함.
	bSame, folders, _, err := d.CompareFolders(db, config.RootDir, nil, exclusions)
	if bSame {
		for _, folder := range folders {
			bSame, files, _, _ := d.CompareFiles(db, folder.Path, exclusions)
			if bSame {
				fmt.Println("Same")

				fileNames := d.ExtractFileNames(files)
				_, err := r.GenerateFileBlock(folder.Path, fileNames)

				if err != nil { // 에러 발생 시 종료
					os.Exit(1)
				}
			}

		}
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
