package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	d "github.com/seoyhaein/datablock/db"
	"github.com/seoyhaein/datablock/protos"
	r "github.com/seoyhaein/datablock/rule"
	"github.com/seoyhaein/datablock/v1rpc"
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

// TODO 향후 grpc api 들어갈 예정.

// StoreFoldersInfo 폴더 정보를 DB에 저장
func (f *dBApisImpl) StoreFoldersInfo(ctx context.Context, db *sql.DB) error {
	err := d.StoreFoldersInfo(ctx, db, f.rootDir, f.foldersExclusion, f.filesExclusions)
	return err
}

// CompareFoldersAndFiles 폴더와 파일을 비교하고, 변경 내역을 반환
func (f *dBApisImpl) CompareFoldersAndFiles(ctx context.Context, db *sql.DB) (*bool, []d.FolderDiff, []d.FileChange, []*protos.FileBlockData, error) {
	// 1. 폴더 비교: 폴더 목록과 폴더 간 차이 정보를 가져옴
	_, folders, folderDiffs, err := d.CompareFolders(db, f.rootDir, f.foldersExclusion, f.filesExclusions)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var allFileChanges []d.FileChange
	var fbs []*protos.FileBlockData // 파일 블록 데이터 슬라이스

	// 2. 각 폴더에 대해 파일 비교
	for _, folder := range folders {
		// 파일 비교
		filesMatch, files, fileChanges, err := d.CompareFiles(db, folder.Path, f.filesExclusions)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		// 해당 폴더에 파일이 다르다면 변경 내역에 추가
		if !filesMatch {
			allFileChanges = append(allFileChanges, fileChanges...)
		} else {
			// 파일과 폴더가 db와 동일한 경우, 특수 파일 존재 여부를 확인

			// rule.json 파일이 없으면 에러 리턴
			ruleExists, err := FileExistsExact(folder.Path, "rule.json")
			if !ruleExists {
				return nil, nil, nil, nil, fmt.Errorf("required file rule.json does not exist in folder %s", folder.Path)
			}
			if err != nil {
				return nil, nil, nil, nil, err
			}

			// fileblock.csv 존재 여부 확인
			bfb, err := FileExistsExact(folder.Path, "fileblock.csv")
			if err != nil {
				return nil, nil, nil, nil, err
			}

			// *.pb 존재 여부 확인
			pbs, err := SearchFilesByPattern(folder.Path, "*.pb")
			if err != nil {
				return nil, nil, nil, nil, err
			}

			// 만약 pb 파일이 여러 개이면 삭제 후 빈 슬라이스로 초기화
			if len(pbs) > 1 {
				if err = DeleteFilesByPattern(folder.Path, "*.pb"); err != nil {
					return nil, nil, nil, nil, err
				}
				pbs = []string{}
			}

			// rule.json 있고, fileblock.csv 있으며, 정확히 하나의 pb 파일이 있으면 기존 파일 블록 로드
			if bfb && len(pbs) == 1 {
				pbPath := pbs[0]
				fb, err := v1rpc.LoadFileBlock(pbPath)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				fbs = append(fbs, fb)
				continue
			}

			// []Files 를 []string(파일 이름 목록)으로 변환 후, 새 파일 블록 생성
			fileNames := d.ExtractFileNames(files)
			fb, err := r.GenerateFileBlock(folder.Path, fileNames)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			fbs = append(fbs, fb)
		}
	}

	// 전체 동일 여부 결정: 폴더 차이와 파일 변경 내역이 없으면 true, 아니면 false
	overallSame := len(folderDiffs) == 0 && len(allFileChanges) == 0
	if overallSame {
		return u.PTrue, nil, nil, fbs, nil
	}
	return u.PFalse, folderDiffs, allFileChanges, fbs, nil
}

// UpdateFilesAndFolders 폴더 변경 내역과 파일 변경 내역을 DB에 반영
func UpdateFilesAndFolders(ctx context.Context, db *sql.DB, diffs []d.FolderDiff, changes []d.FileChange) error {
	// 폴더 변경 업데이트
	if err := d.UpsertFolders(ctx, db, diffs); err != nil {
		return err
	}
	// 파일 변경 업데이트
	if err := d.UpsertDelFiles(ctx, db, changes); err != nil {
		return err
	}
	return nil
}

func SaveDataBlock(inputBlocks []*protos.FileBlockData, outputFile string) error {
	dataBlock, err := v1rpc.MergeFileBlocksFromData(inputBlocks)
	if err != nil {
		return err
	}

	// DataBlock 저장
	if err := v1rpc.SaveProtoToFile(outputFile, dataBlock, os.ModePerm); err != nil {
		return fmt.Errorf("failed to save DataBlock: %w", err)
	}

	fmt.Printf("Successfully merged %d FileBlock files into %s\n", len(inputBlocks), outputFile)
	return nil
}

// FileExists 주어진 폴더 내에서 fileOrPattern 에 해당하는 파일이 존재하는지 확인
// usePattern 이 false 이면 정확한 파일명으로 확인하고, true 이면 fileOrPattern 을 glob 패턴으로 사용함
/*func FileExists(folder, fileOrPattern string, usePattern bool) (bool, error) {
	if !usePattern {
		// 정확한 파일명으로 존재 여부 체크
		path := filepath.Join(folder, fileOrPattern)
		if _, err := os.Stat(path); err == nil {
			return true, nil
		} else if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("파일 체크 실패 (%s): %w", path, err)
		}
	} else {
		// 패턴 검색: glob 사용
		fullPattern := filepath.Join(folder, fileOrPattern)
		matches, err := filepath.Glob(fullPattern)
		if err != nil {
			return false, fmt.Errorf("패턴 검색 실패 (%s): %w", fileOrPattern, err)
		}
		return len(matches) > 0, nil
	}
}*/

// FileExistsExact 주어진 폴더 내에서 정확한 파일명이 존재하는지 확인. 별도로 FileExists 가 있지만 그냥 이걸 씀.
func FileExistsExact(folder, fileName string) (bool, error) {
	path := filepath.Join(folder, fileName)
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("파일 체크 실패 (%s): %w", path, err)
	}
}

// SearchFilesByPattern는 주어진 폴더 내에서 지정한 glob 패턴에 매칭되는 파일들을 검색합니다.
// 검색 결과로 매칭된 파일 경로들의 배열을 반환합니다.
func SearchFilesByPattern(folder, pattern string) ([]string, error) {
	fullPattern := filepath.Join(folder, pattern)
	matches, err := filepath.Glob(fullPattern)
	if err != nil {
		return nil, fmt.Errorf("패턴 검색 실패 (%s): %w", pattern, err)
	}
	return matches, nil
}

// DeleteFilesByPattern는 주어진 폴더 내에서 지정한 glob 패턴에 매칭되는 파일들을 검색합니다.
// 만약 매칭된 파일이 2개 이상이면, 해당 파일들을 모두 삭제합니다.
func DeleteFilesByPattern(folder, pattern string) error {
	files, err := SearchFilesByPattern(folder, pattern)
	if err != nil {
		return fmt.Errorf("패턴 검색 실패 (%s): %w", pattern, err)
	}

	// 매칭된 파일이 여러 개인 경우에만 삭제 수행
	if len(files) > 1 {
		for _, filePath := range files {
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("파일 삭제 실패 (%s): %w", filePath, err)
			}
		}
	}
	return nil
}
