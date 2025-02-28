package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	d "github.com/seoyhaein/datablock/db"
	"github.com/seoyhaein/datablock/protos"
	r "github.com/seoyhaein/datablock/rule"
	u "github.com/seoyhaein/utils"
	"os"
	"path/filepath"
)

// 내부 api 들을 여기서 wrapping 함. TODO 향후 cobra 를 사용하여 cli 로 만들어 줄 예정.

// DBApis 데이터베이스와 관련된 인터페이스
type DBApis interface {
	StoreFoldersInfo(ctx context.Context, db *sql.DB) error
	CompareFoldersAndFiles(ctx context.Context, db *sql.DB) (*bool, []d.FolderDiff, []d.FileChange, []*protos.FileBlockData, error)
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

// CompareFoldersAndFiles 폴더와 파일을 비교하고, 변경 내역을 반환
func (f *dBApisImpl) CompareFoldersAndFiles(ctx context.Context, db *sql.DB) (*bool, []d.FolderDiff, []d.FileChange, []*protos.FileBlockData, error) {
	// 폴더 비교: foldersMatch, folders, folderDiffs, err
	_, folders, folderDiffs, err := d.CompareFolders(db, f.rootDir, f.foldersExclusion, f.filesExclusions)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// 만약 폴더 차이가 있다면, 전체 결과는 동일하지 않음
	/*if !foldersMatch {
		return false, folderDiffs, nil, nil
	}*/

	var allFileChanges []d.FileChange
	var fbs []*protos.FileBlockData // 파일 블록 데이터 슬라이스

	// 각 폴더별로 파일 비교 진행
	for _, folder := range folders {
		// 파일 비교: filesMatch, files, fileChanges, err
		filesMatch, files, fileChanges, err := d.CompareFiles(db, folder.Path, f.filesExclusions)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		// 해당 폴더에 파일이 다르다면 db 업데이트 후 rule.json 을 제외하고 모든 것을 다시 만들어 줘야 함.
		if !filesMatch {
			allFileChanges = append(allFileChanges, fileChanges...)
		} else {
			// 파일과 폴더가 db 와 동일한 경우, 특수 파일 존재 여부를 확인
			special, err := SpecialFilesExist(folder.Path)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			// 특수 파일이 이미 존재하면 새로 생성할 필요가 없으므로 건너뜀
			if special != nil && *special {
				continue
			}
			// 특수 파일이 없으면, []Files 를 []string 으로 변환 여기에는 파일이름만 들어감. GenerateFileBlock 실행
			fileNames := d.ExtractFileNames(files)
			// rule.json 은 반드시 있어야 함. 없으면 에러. SpecialFilesExist 에서 검사해서 사실 필요한가 싶음.
			// *.pb 파일 생성함
			fb, err := r.GenerateFileBlock(folder.Path, fileNames)

			if err != nil {
				return nil, nil, nil, nil, err
			}
			fbs = append(fbs, fb)
		}
	}

	// 전체 동일 여부: 폴더 차이와 파일 변경 내역이 없으면 true, 아니면 false
	overallSame := len(folderDiffs) == 0 && len(allFileChanges) == 0
	if overallSame {
		return u.PTrue, nil, nil, fbs, nil
	}
	return u.PFalse, folderDiffs, allFileChanges, fbs, nil
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
