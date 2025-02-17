package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	c "github.com/seoyhaein/datablock/config"
	u "github.com/seoyhaein/datablock/db"
	"os"
)

func main() {

	//db connection
	db, err := u.ConnectDB("sqlite3", "file_monitor.db")
	if err != nil {
		os.Exit(1)
	}
	err = u.InitializeDatabase(db)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
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
	path := config.RootDir
	// 테스트로 빈파일 생성
	// 기존 파일이 생성되어 있을 경우 권한 설정을 안해줌. 버그지만 고치지 않음.
	u.MakeTestFiles(path)

	ctx := context.Background()
	err = u.FirstCheck(ctx, db, path)
	if err != nil {
		fmt.Println("FirstCheckEmbed Error")
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
