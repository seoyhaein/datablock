package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	d "github.com/seoyhaein/datablock/db"
	r "github.com/seoyhaein/datablock/rule"
	u "github.com/seoyhaein/utils"
	"os"
	"path/filepath"
)

// 내부 api 들을 여기서 wrapping 함. TODO 향후 cobra 를 사용하여 cli 로 만들어 줄 예정.

// DBApis 데이터베이스와 관련된 인터페이스
type DBApis interface {
	StoreFoldersInfo(ctx context.Context, db *sql.DB) error
	CompareFoldersAndFiles(ctx context.Context, db *sql.DB) (bool, error)
}

// dBApisImpl DBApis 인터페이스의 구현체
type dBApisImpl struct {
	rootDir          string // 내부적으로 사용할 폴더 경로 (예: config.RootDir)
	foldersExclusion []string
	filesExclusions  []string
}

// NewDBApis DBApis 인터페이스의 구현체를 생성하는 factory 함수
func NewDBApis(rootDir string, foldersExclusion, filesExclusions []string) DBApis {
	return &dBApisImpl{
		rootDir:          rootDir,
		foldersExclusion: foldersExclusion,
		filesExclusions:  filesExclusions,
	}
}

// StoreFoldersInfo TODO 향후 grpc api 들어갈 예정.
func (f *dBApisImpl) StoreFoldersInfo(ctx context.Context, db *sql.DB) error {
	err := d.StoreFoldersInfo(ctx, db, f.rootDir, f.foldersExclusion, f.filesExclusions)
	return err
}

// CompareFoldersAndFiles TODO 이거 생각할 거 많음. 같지 않을때 어떻게 처리할지 생각해야함. 같지 않을때는 db 를 업데이트 해야 할 거 같음.
func (f *dBApisImpl) CompareFoldersAndFiles(ctx context.Context, db *sql.DB) (bool, error) {
	// 폴더 비교: foldersMatch, folders, <unused>, err
	foldersMatch, folders, _, err := d.CompareFolders(db, f.rootDir, f.foldersExclusion, f.filesExclusions)
	if err != nil {
		return false, err
	}
	if !foldersMatch {
		return false, nil
	}

	// 각 폴더별 처리
	for _, folder := range folders {
		// 먼저, 해당 폴더 내 파일들을 비교
		filesMatch, files, _, err := d.CompareFiles(db, folder.Path, f.filesExclusions)
		if err != nil {
			return false, err
		}
		if !filesMatch {
			return false, nil
		}

		// 파일과 폴더가 동일한 경우, 특수 파일 존재 여부를 확인
		special, err := SpecialFilesExist(folder.Path)
		if err != nil {
			return false, err
		}
		// 특수 파일이 이미 존재하면 새로 생성할 필요가 없으므로 건너뜀
		if special != nil && *special {
			continue
		}

		// 특수 파일이 없으면, 파일 목록을 추출하고 GenerateFileBlock 실행
		fileNames := d.ExtractFileNames(files)
		if _, err := r.GenerateFileBlock(folder.Path, fileNames); err != nil {
			return false, err
		}
	}
	return true, nil
}

// SpecialFilesExist 파일이 하나라도 있으면 PTrue, 모두 없으면 PFalse, 에러 발생 시 nil 반환
func SpecialFilesExist(folder string) (*bool, error) {
	// 고정 파일명 체크: rule.json, fileblock.csv
	filesToCheck := []string{"rule.json", "fileblock.csv"}
	for _, fileName := range filesToCheck {
		path := filepath.Join(folder, fileName)
		if _, err := os.Stat(path); err == nil {
			return u.PTrue, nil
		} else if !os.IsNotExist(err) {
			// err != nil 이면서 os.IsNotExist(err)인 경우는, 파일이 없는 게 아니라 권한 문제 등 예상치 못한 오류가 발생한 경우
			return nil, fmt.Errorf("failed to check %s: %w", path, err)
		}
	}

	// 패턴 체크: invalid_files* 와 *.pb
	patterns := []string{"invalid_files*", "*.pb"}
	for _, pattern := range patterns {
		fullPattern := filepath.Join(folder, pattern)
		matches, err := filepath.Glob(fullPattern)
		if err != nil {
			return nil, fmt.Errorf("failed to search for files with pattern %s: %w", pattern, err)
		}
		if len(matches) > 0 {
			return u.PTrue, nil
		}
	}

	return u.PFalse, nil
}
