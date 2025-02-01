package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	// 폴더 내 파일 목록 탐색 **Go 1.16 이후 부터 가능.**
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
		// Go 1.16 이후 부터 os.ReadDir 함수가 반환하는 DirEntry 에는 Info()메서드가 있어, 이 메서드를 사용하면 추가적인 시스템 콜 없이 파일 정보를 가져올 수 있음.
		// fileInfo, err := os.Stat(filePath) 이 코드 대신 아래 코드로 대체 가능
		fileInfo, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info for %s: %w", filePath, err)
		}

		// 파일 이름과 크기 저장
		filesInfo[entry.Name()] = fileInfo.Size()
	}

	return filesInfo, nil
}

// executeSQLWithTx executes an SQL statement from the specified file
func executeSQLWithTx(ctx context.Context, tx *sql.Tx, filePath string, args ...interface{}) error {
	// If the provided context is nil, use the default background context.
	// 방어적 프로그램(defensive programming) 관점에서 작성하는 것이 유리함. TODO 일단 생각해보자.
	if ctx == nil {
		ctx = context.Background()
	}

	// Read the SQL file content.
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("SQL 파일 읽기 실패 (%s): %w", filePath, err)
	}

	// Trim any unnecessary whitespace from the SQL query.
	query := strings.TrimSpace(string(content))
	if query == "" {
		return fmt.Errorf("SQL 파일 (%s)이 비어 있습니다", filePath)
	}

	// Execute the SQL query within the transaction using ExecContext.
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("SQL 실행 실패 (%s): %w", filePath, err)
	}

	return nil
}

// executeSQLTx 컨텍스트 없이 호출할 때 사용
func executeSQLTx(tx *sql.Tx, filePath string, args ...interface{}) error {
	return executeSQLWithTx(context.Background(), tx, filePath, args...)
}

func executeSQLWithDB(ctx context.Context, db *sql.DB, filePath string, args ...interface{}) error {
	// If the provided context is nil, use the default background context.
	// 방어적 프로그램(defensive programming) 관점에서 작성하는 것이 유리함.
	// 예전에 nil 이면 그냥 에러 처리했는데. 이렇게 하는게 더 좋은듯. TODO 일단 생각해보자.
	if ctx == nil {
		ctx = context.Background()
	}

	// Read the SQL file content.
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file (%s): %w", filePath, err)
	}

	// Trim whitespace and check if the file is empty.
	query := strings.TrimSpace(string(content))
	if query == "" {
		return fmt.Errorf("SQL file (%s) is empty", filePath)
	}

	// Execute the SQL statement using ExecContext.
	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement from file (%s): %w", filePath, err)
	}

	return nil
}

func executeSQLDB(db *sql.DB, filePath string, args ...interface{}) error {
	return executeSQLWithDB(context.Background(), db, filePath, args...)
}

// SaveFolderAndFiles TODO 이름 바꿀 필요 있음 초기 값을 설정하는 메서드임. 한번만 실행되고 말것. folderPath 검증 해야함. context 넣을 것 생각해보자.
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

	err = executeSQLTx(tx, "queries/insert_folder.sql", folder.Path, folder.TotalSize, folder.FileCount, folder.CreatedTime)
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
		err = executeSQLTx(tx, "queries/insert_file.sql", folderID, name, size, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("파일 삽입 실패: %w", err)
		}
	}

	// 6️. 폴더의 `total_size` 및 `file_count` 업데이트
	err = executeSQLTx(tx, "queries/update_folder.sql", folderID)
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

func isDBInitialized(db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('folders', 'files')").Scan(&count)
	if err != nil {
		log.Println("Failed to check database initialization:", err)
		return false
	}
	return count == 2 // folders, files 두 테이블이 모두 있어야 true 반환
}

func InitializeDatabase(db *sql.DB) error {
	// 데이터베이스가 초기화되지 않았다면 init.sql 실행
	if !isDBInitialized(db) {
		log.Println("Running database initialization...")
		if err := executeSQLDB(db, "queries/init.sql"); err != nil {
			return fmt.Errorf("DB 초기화 실패: %w", err)
		}
		log.Println("Database initialization completed successfully.")
	} else {
		log.Println("Database already initialized. Skipping init.sql execution.")
	}
	return nil
}
