package utils

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	utils "github.com/seoyhaein/utils"
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

// TODO 여기서 생각할 것이 루트 폴더를 기준으로 그 밑의 1차적인 하위 폴더만의 리스트를 먼저 가지고 있어야 한다. 그래서 아래의 메서드들을 이용해서 db 를 채워야 한다.

// FirstCheck 한번만 실행되고 말것. folderPath 검증 해야함. TODO 디렉토리 검증되는지 확인해야 함.
func FirstCheck(ctx context.Context, db *sql.DB, folderPath string) error {

	if ctx == nil {
		ctx = context.Background()
	}

	if utils.IsEmptyString(folderPath) {
		return fmt.Errorf("폴더 경로가 비어 있습니다")
	}

	// 1. 트랜잭션 시작 (Context-aware)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("트랜잭션 시작 실패: %w", err)
	}

	// 2. 폴더 정보 삽입 (Context-aware executeSQLTx 사용)
	folder := Folder{
		Path:        folderPath,
		TotalSize:   0,
		FileCount:   0,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = executeSQLWithTx(ctx, tx, "queries/insert_folder.sql", folder.Path, folder.TotalSize, folder.FileCount, folder.CreatedTime)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 삽입 실패: %w", err)
	}

	// 3. 폴더 ID 가져오기 (Context-aware QueryRow)
	var folderID int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM folders WHERE path = ?", folder.Path).Scan(&folderID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 ID 조회 실패: %w", err)
	}

	// 4. 폴더 내 파일 목록 가져오기 (여기서는 Context 사용이 필요하지 않을 수도 있음)
	filesInfo, err := GetFilesWithSize(folderPath)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 내 파일 정보 가져오기 실패: %w", err)
	}

	// 5. 파일 정보를 DB에 삽입 (Context-aware executeSQLTx 사용)
	for name, size := range filesInfo {
		err = executeSQLWithTx(ctx, tx, "queries/insert_file.sql", folderID, name, size, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("파일 삽입 실패: %w", err)
		}
	}

	// 6. 폴더의 total_size 및 file_count 업데이트 (Context-aware)
	err = executeSQLWithTx(ctx, tx, "queries/update_folder.sql", folderID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 통계 업데이트 실패: %w", err)
	}

	// 7. 트랜잭션 커밋
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

// TODO 테스트 필요.
// ----------------------------------------------------------------------------
// 아래부터는 embed 방식으로 SQL 파일을 읽어오는 새로운 함수들입니다.
// 기존 함수들은 그대로 남겨두고, 파일 경로 대신 SQL 파일명을 전달하여 embed된 파일에서 내용을 읽어옵니다.

// embed 패키지를 사용하여 쿼리 파일들을 포함합니다.
// 이 예제에서는 utils 패키지 파일 기준 상위 폴더의 queries 폴더 내의 모든 .sql 파일을 포함합니다.
//
// 주의: 실제 프로젝트 디렉토리 구조에 맞게 경로를 조정해야 합니다.
// 예: utils 폴더 안에 있을 경우 "../queries/*.sql" 처럼.

//go:embed ../queries/*.sql
var sqlFiles embed.FS

// getSQLQuery  파일명을 입력받아 해당 SQL 파일의 내용 반환
func getSQLQuery(fileName string) (string, error) {
	// "queries/" 디렉토리 내의 파일을 읽어옵니다.
	bytes, err := sqlFiles.ReadFile("queries/" + fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read SQL file %s: %w", fileName, err)
	}
	return string(bytes), nil
}

// executeSQLWithTxEmbed embed된 SQL 파일을 읽어와 트랜잭션 내에서 실행합니다.
func executeSQLWithTxEmbed(ctx context.Context, tx *sql.Tx, fileName string, args ...interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// embed 파일 시스템에서 "queries/" 하위의 fileName 파일을 읽어옵니다.
	// (embeddedSQLFiles 포함된 경로는 컴파일 시점의 상대 경로에 따라 달라집니다.)
	content, err := sqlFiles.ReadFile("queries/" + fileName)
	if err != nil {
		return fmt.Errorf("SQL 파일 읽기 실패 (%s): %w", fileName, err)
	}

	query := strings.TrimSpace(string(content))
	if query == "" {
		return fmt.Errorf("SQL 파일 (%s)이 비어 있습니다", fileName)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("SQL 실행 실패 (%s): %w", fileName, err)
	}

	return nil
}

// executeSQLTxEmbed 컨텍스트 없이 embed SQL 실행 시 사용할 수 있습니다.
func executeSQLTxEmbed(tx *sql.Tx, fileName string, args ...interface{}) error {
	return executeSQLWithTxEmbed(context.Background(), tx, fileName, args...)
}

// executeSQLWithDBEmbed embed된 SQL 파일을 읽어 DB에서 실행합니다.
func executeSQLWithDBEmbed(ctx context.Context, db *sql.DB, fileName string, args ...interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}

	content, err := sqlFiles.ReadFile("queries/" + fileName)
	if err != nil {
		return fmt.Errorf("SQL 파일 읽기 실패 (%s): %w", fileName, err)
	}

	query := strings.TrimSpace(string(content))
	if query == "" {
		return fmt.Errorf("SQL 파일 (%s)이 비어 있습니다", fileName)
	}

	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("SQL 실행 실패 (%s): %w", fileName, err)
	}

	return nil
}

// executeSQLDBEmbed 컨텍스트 없이 embed된 SQL 파일을 DB에서 실행할 때 사용합니다.
func executeSQLDBEmbed(db *sql.DB, fileName string, args ...interface{}) error {
	return executeSQLWithDBEmbed(context.Background(), db, fileName, args...)
}

// FirstCheckEmbed embed 방식으로 SQL 파일을 읽어와 폴더 및 파일 정보를 DB에 삽입하는 함수입니다.
func FirstCheckEmbed(ctx context.Context, db *sql.DB, folderPath string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if utils.IsEmptyString(folderPath) {
		return fmt.Errorf("폴더 경로가 비어 있습니다")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("트랜잭션 시작 실패: %w", err)
	}

	folder := Folder{
		Path:        folderPath,
		TotalSize:   0,
		FileCount:   0,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	// embed SQL 파일 (insert_folder.sql)을 사용합니다.
	err = executeSQLWithTxEmbed(ctx, tx, "insert_folder.sql", folder.Path, folder.TotalSize, folder.FileCount, folder.CreatedTime)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 삽입 실패: %w", err)
	}

	var folderID int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM folders WHERE path = ?", folder.Path).Scan(&folderID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 ID 조회 실패: %w", err)
	}

	filesInfo, err := GetFilesWithSize(folderPath)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 내 파일 정보 가져오기 실패: %w", err)
	}

	for name, size := range filesInfo {
		err = executeSQLWithTxEmbed(ctx, tx, "insert_file.sql", folderID, name, size, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("파일 삽입 실패: %w", err)
		}
	}

	err = executeSQLWithTxEmbed(ctx, tx, "update_folder.sql", folderID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("폴더 통계 업데이트 실패: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("트랜잭션 커밋 실패: %w", err)
	}

	return nil
}

// InitializeDatabaseEmbed embed된 SQL 파일(init.sql)을 사용하여 데이터베이스를 초기화합니다.
func InitializeDatabaseEmbed(db *sql.DB) error {
	if !isDBInitialized(db) {
		log.Println("Running database initialization (embed)...")
		if err := executeSQLDBEmbed(db, "init.sql"); err != nil {
			return fmt.Errorf("DB 초기화 실패: %w", err)
		}
		log.Println("Database initialization completed successfully (embed).")
	} else {
		log.Println("Database already initialized. Skipping init.sql execution.")
	}
	return nil
}
