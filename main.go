package main

// TODO 코드 정리 필요.
import (
	c "github.com/seoyhaein/datablock/config"
	r "github.com/seoyhaein/datablock/rule"
	u "github.com/seoyhaein/datablock/utils"
	v1rpc "github.com/seoyhaein/datablock/v1rpc"
	"log"
	"os"
)

func main() {

	//db connection
	db, err := u.ConnectDB("sqlite3", "file_monitor.db")
	if err != nil {
		os.Exit(1)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("failed to close db:", err)
			os.Exit(1) // defer 내부에서도 os.Exit 사용 가능
		}
	}()

	config, err := c.LoadConfig("config.json")
	if err != nil {
		os.Exit(1)
	}
	path := config.RootDir
	// 테스트로 빈파일 생성
	// 기존 파일이 생성되어 있을 경우 권한 설정을 안해줌. 버그지만 고치지 않음.
	u.MakeTestFiles(path)
	// 이름을 바꾸던가 내용을 좀 수정해야 할듯하다.
	// 폴더별로 rule 을 작성해줘야 하는데 이것을 사용자 친화적으로 해주는 것을 추가적으로 구현해야 한다. 하지만 이건 나중에.
	// data 는 map[int]map[string]string 형태임.
	data, err := r.ConnectProto(path)
	if err != nil { // 에러 발생 시 종료
		os.Exit(1)
	}
	headers := []string{"r1", "r2"}
	fbd := v1rpc.ConvertMapToFileBlockData(data, headers, "tester")

	err = v1rpc.SaveProtoToFile("tester.pb", fbd, 0777)
	if err != nil {
		os.Exit(1)
	}

	/*_, err = v1rpc.LoadFileBlock(path)

	if err != nil {
		os.Exit(1)
	}*/
}
