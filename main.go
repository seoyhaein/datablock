package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	c "github.com/seoyhaein/datablock/config"
	d "github.com/seoyhaein/datablock/db"
	r "github.com/seoyhaein/datablock/rule"
	"github.com/seoyhaein/datablock/v1rpc"
	"os"
	"path"
	"path/filepath"
)

func main() {

	_ = RemoveDBFile("file_monitor.db")
	//db connection foreign key 설정을 위해 PRAGMA foreign_keys = ON; 설정을 해줘야 함.
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

	// TODO 일단 주석 처리 DB 먼저 끝내고 주석 풀음.
	config, err := c.LoadConfig("config.json")
	if err != nil {
		os.Exit(1)
	}

	// 테스트를 위해서 조정해줌
	testFilePath := config.RootDir
	// 테스트로 빈파일 생성
	// 기존 파일이 생성되어 있을 경우 권한 설정을 안해줌. 버그지만 고치지 않음.
	testFilePath = filepath.Join(testFilePath, "testFiles/")
	testFilePath = path.Clean(testFilePath)
	d.MakeTestFiles(testFilePath)
	d.MakeTestFilesA("/test/baba/")

	ctx := context.Background()
	// TODO 메서드들을 모아주는 거 생각하자. 일단은 이렇게 해놓음. 향후 api 로 따로 빼놓아야 할듯. ConnectDB 는 기본으로 해주고 db 를 초기화하고 데이터를 넣어주는 api 를 따로 만들어 주어야 할듯.
	// TODO 각 폴더별로 exclusions 이 다를 수 있음. 이거 고려되어야 함. -> 동일해진거 같은데 생각해보자.
	// TODO 아래 폴더들을 생각해서 일단 작성하자.
	exclusions := []string{"*.json", "invalid_files", "*.csv", "*.pb"}

	err = d.StoreFoldersInfo(ctx, db, config.RootDir, nil, exclusions)
	if err != nil {
		fmt.Println("StoreFoldersInfo Error")
	}

	// 일단 이렇게 대략적으로 구성함.
	bSame, folders, _, err := d.CompareFolders(db, config.RootDir, nil, exclusions)
	if bSame {
		for _, folder := range folders {
			bSame, files, _, _ := d.CompareFiles(db, folder.Path, exclusions)
			if bSame {
				fmt.Println("Same")

				fileNames := d.ExtractFileNames(files)
				data, err := r.GenerateMap1(folder.Path, fileNames, exclusions)
				if err != nil { // 에러 발생 시 종료
					os.Exit(1)
				}
				// TODO ConvertMapToFileBlockData GenerateMap 에 통합 시켜도 됨. header 때문에.
				// blockID 같은 경우는 폴더명으로 한다. 고유해야 함.
				headers := []string{"r1", "r2"}
				fbd := v1rpc.ConvertMapToFileBlockData(data, headers, "tester")

				err = v1rpc.SaveProtoToFile("tester1.pb", fbd, 0777)
				if err != nil {
					os.Exit(1)
				}
			}

		}

	}

	/*
		// 이름을 바꾸던가 내용을 좀 수정해야 할듯하다.
			// 폴더별로 rule 을 작성해줘야 하는데 이것을 사용자 친화적으로 해주는 것을 추가적으로 구현해야 한다. 하지만 이건 나중에.
			// data 는 map[int]map[string]string 형태임.
			data, err := r.GenerateMap(path)
			if err != nil { // 에러 발생 시 종료
				os.Exit(1)
			}
			headers := []string{"r1", "r2"}
			fbd := v1rpc.ConvertMapToFileBlockData(data, headers, "tester")

			err = v1rpc.SaveProtoToFile("tester.pb", fbd, 0777)
			if err != nil {
				os.Exit(1)
			}
	*/
}

// RemoveDBFile 주어진 DB 파일을 삭제함.
// filePath: 삭제할 DB 파일의 경로.
func RemoveDBFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove DB file %s: %w", filePath, err)
	}
	return nil
}
