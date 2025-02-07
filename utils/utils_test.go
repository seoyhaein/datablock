package utils

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"testing"
)

// TestFirstCheck는 임시 폴더와 in‑memory SQLite DB를 이용해 FirstCheck 함수의 전체 동작을 검증
func TestFirstCheck(t *testing.T) {
	ctx := context.Background()

	// 임시 폴더 생성 및 더미 파일 생성
	tempDir := t.TempDir()
	dummyFilePath := filepath.Join(tempDir, "dummy.txt")
	dummyContent := []byte("Hello, world!")
	if err := os.WriteFile(dummyFilePath, dummyContent, 0644); err != nil {
		t.Fatalf("더미 파일 생성 실패: %v", err)
	}

	// in-memory SQLite DB 생성
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("in-memory DB 생성 실패: %v", err)
	}
	defer db.Close()

	// 테이블 생성 (folders, files)
	schema := []string{
		`CREATE TABLE folders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL,
			total_size INTEGER,
			file_count INTEGER,
			created_time TEXT
		);`,
		`CREATE TABLE files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			folder_id INTEGER,
			name TEXT,
			size INTEGER,
			created_time TEXT,
			FOREIGN KEY(folder_id) REFERENCES folders(id)
		);`,
	}
	for _, q := range schema {
		if _, err := db.Exec(q); err != nil {
			t.Fatalf("스키마 생성 실패: %v", err)
		}
	}

	// FirstCheck 함수 실행
	if err := FirstCheck(ctx, db, tempDir); err != nil {
		t.Fatalf("FirstCheck 실행 실패: %v", err)
	}

	// 폴더 레코드 확인
	var folderID int64
	var folderPath string
	var totalSize, fileCount int64
	var createdTime string
	row := db.QueryRow("SELECT id, path, total_size, file_count, created_time FROM folders WHERE path = ?", tempDir)
	if err := row.Scan(&folderID, &folderPath, &totalSize, &fileCount, &createdTime); err != nil {
		t.Fatalf("폴더 레코드 조회 실패: %v", err)
	}
	if folderPath != tempDir {
		t.Errorf("예상한 폴더 경로 %q, 실제 폴더 경로 %q", tempDir, folderPath)
	}

	// 파일 레코드 확인
	var fileID int64
	var dbFolderID int64
	var fileName string
	var fileSize int64
	var fileCreatedTime string
	row = db.QueryRow("SELECT id, folder_id, name, size, created_time FROM files WHERE folder_id = ?", folderID)
	if err := row.Scan(&fileID, &dbFolderID, &fileName, &fileSize, &fileCreatedTime); err != nil {
		t.Fatalf("파일 레코드 조회 실패: %v", err)
	}
	if dbFolderID != folderID {
		t.Errorf("파일의 folder_id가 %d여야 하는데, 실제는 %d", folderID, dbFolderID)
	}
	if fileName != "dummy.txt" {
		t.Errorf("예상한 파일명은 %q, 실제는 %q", "dummy.txt", fileName)
	}
	if fileSize != int64(len(dummyContent)) {
		t.Errorf("예상한 파일 크기는 %d, 실제는 %d", len(dummyContent), fileSize)
	}

	// 폴더의 total_size와 file_count 업데이트 확인
	if totalSize != fileSize {
		t.Errorf("예상한 폴더 total_size는 %d, 실제는 %d", fileSize, totalSize)
	}
	if fileCount != 1 {
		t.Errorf("예상한 폴더 file_count는 1, 실제는 %d", fileCount)
	}
}
