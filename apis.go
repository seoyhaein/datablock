package main

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	d "github.com/seoyhaein/datablock/db"
)

// 내부 api 들을 여기서 wrapping 함. TODO 향후 cobra 를 사용하여 cli 로 만들어 줄 예정.

// DBApis 데이터베이스와 관련된 인터페이스
type DBApis interface {
	StoreFoldersInfo(ctx context.Context, db *sql.DB) error
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
