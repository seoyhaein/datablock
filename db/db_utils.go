package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	u "github.com/seoyhaein/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

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

// GetFilesWithSize 특정 폴더에서 파일 이름과 해당 파일 크기를 가져오는 함수, 삭제 예정.
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

func isDBInitialized(db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('folders', 'files')").Scan(&count)
	if err != nil {
		log.Println("Failed to check database initialization:", err)
		return false
	}
	return count == 2 // folders, files 두 테이블이 모두 있어야 true 반환
}

// execSQLTx 읽어온 SQL 파일을 트랜잭션 내에서 ExecContext 로 실행.
// IMPORTANT: 비 SELECT 쿼리에 사용. (결과 리턴 없음)
func execSQLTx(ctx context.Context, tx *sql.Tx, fileName string, args ...interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// "queries/" 하위의 SQL 파일을 읽어옴.
	content, err := sqlFiles.ReadFile("queries/" + fileName)
	if err != nil {
		return fmt.Errorf("failed to read SQL file (%s): %w", fileName, err)
	}

	query := strings.TrimSpace(string(content))
	if query == "" {
		return fmt.Errorf("SQL file (%s) is empty", fileName)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("SQL execution failed (%s): %w", fileName, err)
	}

	return nil
}

// execSQLTxNoCtx 컨텍스트 없이 트랜잭션 내에서 SQL 파일을 실행.
func execSQLTxNoCtx(tx *sql.Tx, fileName string, args ...interface{}) error {
	return execSQLTx(context.Background(), tx, fileName, args...)
}

// execSQL 읽어온 SQL 파일을 DB 에서 ExecContext 로 실행.
// IMPORTANT: 비 SELECT 쿼리에 사용. (결과 리턴 없음)
func execSQL(ctx context.Context, db *sql.DB, fileName string, args ...interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// "queries/" 하위의 SQL 파일을 읽어옴.
	content, err := sqlFiles.ReadFile("queries/" + fileName)
	if err != nil {
		return fmt.Errorf("failed to read SQL file (%s): %w", fileName, err)
	}

	query := strings.TrimSpace(string(content))
	if query == "" {
		return fmt.Errorf("SQL file (%s) is empty", fileName)
	}

	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("SQL execution failed (%s): %w", fileName, err)
	}

	return nil
}

// execSQLNoCtx 컨텍스트 없이 DB 에서 SQL 파일을 실행.
func execSQLNoCtx(db *sql.DB, fileName string, args ...interface{}) error {
	return execSQL(context.Background(), db, fileName, args...)
}

// querySQL 읽어온 SQL 파일을 DB 에서 QueryContext 로 실행.
// IMPORTANT: SELECT 쿼리에 사용. 결과로 *sql.Rows 를 반환하며, 호출자가 반드시 Close() 해야 함.
func querySQL(ctx context.Context, db *sql.DB, fileName string, args ...interface{}) (*sql.Rows, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// "queries/" 하위의 SQL 파일을 읽어옴.
	content, err := sqlFiles.ReadFile("queries/" + fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL file (%s): %w", fileName, err)
	}

	query := strings.TrimSpace(string(content))
	if query == "" {
		return nil, fmt.Errorf("SQL file (%s) is empty", fileName)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("SQL query failed (%s): %w", fileName, err)
	}
	return rows, nil
}

// querySQLNoCtx 컨텍스트 없이 DB 에서 SELECT 쿼리를 실행.
func querySQLNoCtx(db *sql.DB, fileName string, args ...interface{}) (*sql.Rows, error) {
	return querySQL(context.Background(), db, fileName, args...)
}

// FirstCheck 폴더 경로를 받아 폴더 내 파일 정보를 DB에 삽입하는 함수
func FirstCheck(ctx context.Context, db *sql.DB, folderPath string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if u.IsEmptyString(folderPath) {
		return fmt.Errorf("folder path is empty")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	// 폴더에서 제외 해야 할 파일들.
	exclusions := []string{"rule.json", "invalid_files", "fileblock.csv"}
	folderDetails, fileDetails, err := GetFolderDetails(folderPath, exclusions)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.Printf("rollback failed: %v", rbErr)
		}
		return fmt.Errorf("failed to get folder details: %w", err)
	}

	// DB에 폴더 정보 삽입 (insert_folder.sql)
	err = execSQLTx(ctx, tx, "insert_folder.sql",
		folderDetails.Path,
		folderDetails.TotalSize,
		folderDetails.FileCount,
		folderDetails.CreatedTime)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.Printf("rollback failed: %v", rbErr)
		}
		return fmt.Errorf("failed to insert folder: %w", err)
	}

	// 삽입된 폴더의 ID를 조회
	var folderID int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM folders WHERE path = ?", folderDetails.Path).Scan(&folderID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.Printf("rollback failed: %v", rbErr)
		}
		return fmt.Errorf("failed to query folder ID: %w", err)
	}

	// 파일 정보 삽입 (insert_file.sql)
	for _, file := range fileDetails {
		err = execSQLTx(ctx, tx, "insert_file.sql",
			folderID,
			file.Name,
			file.Size,
			file.CreatedTime)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
				log.Printf("rollback failed: %v", rbErr)
			}
			return fmt.Errorf("failed to insert file: %w", err)
		}
	}

	err = execSQLTx(ctx, tx, "update_folders.sql", folderID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.Printf("rollback failed: %v", rbErr)
		}
		return fmt.Errorf("failed to update folder statistics: %w", err)
	}

	// 트랜잭션 커밋
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// InitializeDatabase embed 된  SQL 파일(init.sql)을 사용하여 데이터베이스를 초기화
func InitializeDatabase(db *sql.DB) error {
	if !isDBInitialized(db) {
		log.Println("Running DB initialization (embed)...")
		if err := execSQLNoCtx(db, "init.sql"); err != nil {
			return fmt.Errorf("DB initialization failed: %w", err)
		}
		log.Println("DB initialization completed successfully (embed).")
	} else {
		log.Println("DB already initialized. Skipping init.sql execution.")
	}
	return nil
}

// GetFolderDetails 특정 디렉토리 내의 파일들을 읽어 전체 파일 개수, 총 크기와 각 파일의 메타데이터를 수집.
// Go 1.16부터 도입된 os.ReadDir, DirEntry.Info()를 사용하여 시스템 콜을 최소화함.
func GetFolderDetails(dirPath string, exclusions []string) (Folder, []File, error) {
	var folder Folder
	var files []File

	// 디렉토리 내 파일 목록 읽기 (Go 1.16 이상에서는 os.ReadDir 사용)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return folder, nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	totalSize := int64(0)
	fileCount := int64(0)

	// 각 엔트리(파일)에 대해 처리
	for _, entry := range entries {
		if entry.IsDir() {
			continue // 하위 디렉토리는 무시
		}

		fileName := entry.Name()
		// 제외 목록에 있는 파일이면 건너뛰기
		if u.ExcludeFiles(fileName, exclusions) {
			continue
		}

		// 파일 전체 경로 생성
		filePath := filepath.Join(dirPath, fileName)

		// 파일 정보 가져오기 (os.ReadDir 가 반환하는 DirEntry 의 Info() 사용)
		info, err := entry.Info()
		if err != nil {
			return folder, nil, fmt.Errorf("failed to get file info for %s: %w", filePath, err)
		}

		size := info.Size()
		totalSize += size
		fileCount++

		// File 구조체 생성 (ID 및 FolderID는 DB 삽입 후 업데이트 가능)
		// IMPORTANT sqlite 에서 AUTOINCREMENT 로 시작하도록 하였음. 따라서 0 이 들어간것은 DB 에 들어가기 전 데이터임.
		fileRecord := File{
			ID:          0, // 아직 DB에 저장되지 않았으므로 0 또는 추후 채움
			FolderID:    0, // folder 삽입 후 업데이트
			Name:        fileName,
			Size:        size,
			CreatedTime: info.ModTime().Format("2006-01-02 15:04:05"),
		}
		files = append(files, fileRecord)
	}

	// IMPORTANT Folder 구조체 생성 (ID는 DB 삽입 후 업데이트)
	folder = Folder{
		ID:          0,
		Path:        dirPath,
		TotalSize:   totalSize,
		FileCount:   fileCount,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	return folder, files, nil
}

// GetFoldersFromDB DB의 폴더 정보를 조회하여 Folder 구조체 슬라이스로 반환함.
// IMPORTANT: 호출자가 반환된 rows를 직접 Close() 할 필요는 없음. 내부에서 모두 처리됨.
func GetFoldersFromDB(db *sql.DB) ([]Folder, error) {
	// "select_folders.sql" 파일에 정의된 SELECT 쿼리를 실행하여 폴더 정보를 조회
	rows, err := querySQLNoCtx(db, "select_folders.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to query folders: %w", err)
	}
	defer rows.Close()

	var folders []Folder
	// 각 행을 순회하면서 Folder 구조체에 스캔
	for rows.Next() {
		var f Folder
		err = rows.Scan(&f.ID, &f.Path, &f.TotalSize, &f.FileCount, &f.CreatedTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan folder: %w", err)
		}
		folders = append(folders, f)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return folders, nil
}

// GetFilesInfoFromDB DB의 파일 정보를 조회하여 File 구조체 슬라이스로 반환함.
// IMPORTANT: 호출자가 반환된 rows를 직접 Close() 할 필요는 없음. 내부에서 모두 처리됨.
func GetFilesInfoFromDB(db *sql.DB) ([]File, error) {
	// "select_files.sql" 파일에 정의된 SELECT 쿼리를 실행하여 파일 정보를 조회
	rows, err := querySQLNoCtx(db, "select_files.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []File
	// 각 행을 순회하면서 File 구조체에 스캔
	for rows.Next() {
		var f File
		err = rows.Scan(&f.ID, &f.FolderID, &f.Name, &f.Size, &f.CreatedTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, f)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return files, nil
}

// for test
func clearDatabase(db *sql.DB) error {
	// 외래 키 제약 조건이 ON DELETE CASCADE 로 설정되어 있다면, folders 테이블에서 데이터를 삭제하면 files 테이블의 데이터도 자동 삭제.
	_, err := db.Exec("DELETE FROM folders;")
	return err
}
