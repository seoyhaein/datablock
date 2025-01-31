package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func MakeTestFiles(path string) {
	// 디렉토리 생성
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory %s: %v", path, err)
	}

	// 디렉토리 권한을 777로 설정 os.ModePerm 해줌.
	/*err = os.Chmod(path, 0777) //0o777 이 방식보다 0777 방식 사용
	if err != nil {
		log.Fatalf("Failed to set permissions for directory %s: %v", path, err)
	}*/

	// 테스트 파일 이름 목록

	fileNames := []string{
		"sample1_S1_L001_R1_001.fastq.gz",
		"sample1_S1_L001_R2_001.fastq.gz",
		"sample1_S1_L002_R1_001.fastq.gz",
		"sample1_S1_L002_R2_001.fastq.gz",
		"sample2_S2_L001_R1_001.fastq.gz",
		"sample2_S2_L001_R2_001.fastq.gz",
		"sample2_S2_L002_R1_001.fastq.gz",
		"sample2_S2_L002_R2_001.fastq.gz",
		"sample3_S3_L001_R1_001.fastq.gz",
		"sample3_S3_L001_R2_001.fastq.gz",
		"sample3_S3_L002_R1_001.fastq.gz",
		"sample3_S3_L002_R2_001.fastq.gz",
		"sample4_S4_L001_R1_001.fastq.gz",
		"sample4_S4_L001_R2_001.fastq.gz",
		"sample4_S4_L002_R1_001.fastq.gz",
		"sample4_S4_L002_R2_001.fastq.gz",
		"sample5_S5_L001_R1_001.fastq.gz",
		"sample5_S5_L001_R2_001.fastq.gz",
		"sample5_S5_L002_R1_001.fastq.gz",
		"sample5_S5_L002_R2_001.fastq.gz",
		"sample6_S6_L001_R1_001.fastq.gz",
		"sample6_S6_L001_R2_001.fastq.gz",
		"sample6_S6_L002_R1_001.fastq.gz",
		"sample6_S6_L002_R2_001.fastq.gz",
		"sample7_S7_L001_R1_001.fastq.gz",
		"sample7_S7_L001_R2_001.fastq.gz",
		"sample7_S7_L002_R1_001.fastq.gz",
		"sample7_S7_L002_R2_001.fastq.gz",
		"sample8_S8_L001_R1_001.fastq.gz",
		"sample8_S8_L001_R2_001.fastq.gz",
		"sample8_S8_L002_R1_001.fastq.gz",
		"sample8_S8_L002_R2_001.fastq.gz",
		"sample9_S9_L001_R1_001.fastq.gz",
		"sample9_S9_L001_R2_001.fastq.gz",
		"sample9_S9_L002_R1_001.fastq.gz",
		"sample9_S9_L002_R2_001.fastq.gz",
		"sample10_S10_L001_R1_001.fastq.gz",
		"sample10_S10_L001_R2_001.fastq.gz",
		"sample10_S10_L002_R1_001.fastq.gz",
		"sample10_S10_L002_R2_001.fastq.gz",
		"sample11_S11_L001_R1_001.fastq.gz",
		"sample11_S11_L001_R2_001.fastq.gz",
		"sample11_S11_L002_R1_001.fastq.gz",
		"sample11_S11_L002_R2_001.fastq.gz",
		"sample12_S12_L001_R1_001.fastq.gz",
		"sample12_S12_L001_R2_001.fastq.gz",
		"sample12_S12_L002_R1_001.fastq.gz",
		"sample12_S12_L002_R2_001.fastq.gz",
	}
	/*
		incompleteFileNames := []string{
			"sample1_S1_L001_R1_001.fastq.gz",
			"sample1_S1_L001_R2_001.fastq.gz",
			"sample13_S13_L001_R1.fastq.gz",
			"sample14_S14_L001_R2_001.fastq",
			"sample15_S15_L001_001.fastq.gz",
			"sample16_S16_L001.fastq.gz",
		}
	*/
	// 파일 생성
	for _, fileName := range fileNames {
		filePath := fmt.Sprintf("%s/%s", path, fileName)
		_, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to create file %s: %v", filePath, err)
		} else {
			log.Printf("Created file: %s", filePath)
		}
	}
}

// db 관련 및 기타 부수 methods 들은 여기다가 임시적으로 넣음.

// ConnectDB is a function to connect to a database
func ConnectDB(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}

// ExecuteSQLFile is a function to execute SQL queries from a file
func ExecuteSQLFile(db *sql.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("SQL execution failed: %w", err)
	}

	return nil
}

type File struct {
	ID          int64  `db:"id"`
	FolderID    int64  `db:"folder_id"`
	Name        string `db:"name"`
	Size        int64  `db:"size"`
	CreatedTime string `db:"created_time"` // sting 으로 해도 충분
}

type Folder struct {
	ID          int64  `db:"id"`
	Path        string `db:"path"`
	TotalSize   int64  `db:"total_size"`
	FileCount   int64  `db:"file_count"`
	CreatedTime string `db:"created_time"` // sting 으로 해도 충분
}

// TODO 초기값 설정하는 메서드들 추가 수정 필요함.

// GetFilesWithSize 특정 폴더에서 파일 이름과 해당 파일 크기를 가져오는 함수
func GetFilesWithSize(directoryPath string) (map[string]int64, error) {
	filesInfo := make(map[string]int64)

	// 폴더 내 파일 목록 탐색
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", directoryPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // 폴더는 무시하고 파일만 처리
		}

		// 파일 전체 경로
		filePath := filepath.Join(directoryPath, entry.Name())

		// 파일 정보 가져오기
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get file info for %s: %w", filePath, err)
		}

		// 파일 이름과 크기 저장
		filesInfo[entry.Name()] = fileInfo.Size()
	}

	return filesInfo, nil
}

// SQL 파일을 읽어 실행하는 함수 (매개변수 개별 적용 가능)
func executeSQLFile(tx *sql.Tx, filePath string, args ...interface{}) error {
	// SQL 파일 읽기
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("SQL 파일 읽기 실패 (%s): %w", filePath, err)
	}

	// SQL 실행
	_, err = tx.Exec(string(content), args...)
	if err != nil {
		return fmt.Errorf("SQL 실행 실패 (%s): %w", filePath, err)
	}

	return nil
}

// SaveFolderAndFiles TODO 이름 바꿀 필요 있음 초기 값을 설정하는 메서드임. 한번만 실행되고 말것. folderPath 검증 해야함.
func SaveFolderAndFiles(db *sql.DB, folderPath string) error {
	// 1️.트랜잭션 시작
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("트랜잭션 시작 실패: %w", err)
	}

	// 2️.폴더 정보 삽입 TODO 수정 필요 현재 시간으로, 여기서는 file 관련 정보를 모름.
	folder := Folder{
		Path:        folderPath,
		TotalSize:   0,
		FileCount:   0,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = executeSQLFile(tx, "sql/insert_folder.sql", folder.Path, folder.TotalSize, folder.FileCount, folder.CreatedTime)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 삽입 실패: %w", err)
	}

	// 3. 폴더 ID 가져오기, TODO 향후 get_folder_id.sql 삭제 예정.
	var folderID int64
	err = tx.QueryRow("SELECT id FROM folders WHERE path = ?", folder.Path).Scan(&folderID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 ID 조회 실패: %w", err)
	}

	// 4️. 폴더 내 파일 목록 가져오기
	filesInfo, err := GetFilesWithSize(folderPath)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 내 파일 정보 가져오기 실패: %w", err)
	}

	// 5️. 파일 정보를 DB에 삽입 TODO 수정 필요 현재 시간으로
	for name, size := range filesInfo {
		err = executeSQLFile(tx, "sql/insert_file.sql", folderID, name, size, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("파일 삽입 실패: %w", err)
		}
	}

	// 6️. 폴더의 `total_size` 및 `file_count` 업데이트
	err = executeSQLFile(tx, "sql/update_folder.sql", folderID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 통계 업데이트 실패: %w", err)
	}

	// 7️. 트랜잭션 커밋
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("트랜잭션 커밋 실패: %w", err)
	}

	return nil
}
