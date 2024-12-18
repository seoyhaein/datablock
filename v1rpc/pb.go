package v1rpc

import (
	"fmt"
	pb "github.com/seoyhaein/datablock/protos"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"os"
)

func SaveProtoToFile(filePath string, message proto.Message, perm os.FileMode) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}

	err = os.WriteFile(filePath, data, perm)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

/*
// 기본 권한(0644) 사용
err := SaveProtoToFile("data.pb", message, 0644)

// 다른 권한 설정
err := SaveProtoToFile("data.pb", message, 0600) // 소유자만 읽기/쓰기 가능

// os.FileMode 상수 사용
err := SaveProtoToFile("data.pb", message, os.ModePerm) // 0777
*/

func LoadFileBlock(filePath string) (*pb.FileBlockData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	message := &pb.FileBlockData{}
	err = proto.Unmarshal(data, message)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize data: %w", err)
	}

	return message, nil
}

func LoadDataBlock(filePath string) (*pb.DataBlockData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	message := &pb.DataBlockData{}
	err = proto.Unmarshal(data, message)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize data: %w", err)
	}

	return message, nil
}

// SaveFileBlockToTextFile 함수: FileBlock 를 텍스트 포맷으로 저장
func SaveFileBlockToTextFile(filePath string, data *pb.FileBlockData) error {
	// proto 메시지를 텍스트 포맷으로 변환
	textData, err := prototext.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal to text format: %w", err)
	}

	// 텍스트 데이터를 파일에 저장
	return os.WriteFile(filePath, textData, os.ModePerm)
}

// MergeFileBlocks 여러 FileBlock 파일을 읽어 DataBlock 으로 합치는 메서드
func MergeFileBlocks(inputFiles []string, outputFile string) error {
	var blocks []*pb.FileBlockData
	// 각 입력 파일을 로드하여 blocks 에 추가
	for _, file := range inputFiles {
		block, err := LoadFileBlock(file)
		if err != nil {
			return fmt.Errorf("failed to load file %s: %w", file, err)
		}
		blocks = append(blocks, block)
	}

	// DataBlockData 생성
	dataBlockData := &pb.DataBlockData{
		Blocks: blocks,
	}

	// DataBlockData 저장
	if err := SaveProtoToFile(outputFile, dataBlockData, os.ModePerm); err != nil {
		return fmt.Errorf("failed to save DataBlock: %w", err)
	}

	fmt.Printf("Successfully merged %d FileBlock files into %s\n", len(inputFiles), outputFile)
	return nil
}

// SaveDataBlockToTextFile DataBlockData 텍스트 포맷으로 파일에 저장
func SaveDataBlockToTextFile(filePath string, data *pb.DataBlockData) error {
	// proto 메시지를 텍스트 포맷으로 변환
	textData, err := prototext.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal DataBlock to text format: %w", err)
	}

	// 텍스트 데이터를 파일에 저장
	if err := os.WriteFile(filePath, textData, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	fmt.Printf("Successfully saved DataBlock to %s\n", filePath)
	return nil
}

// GenerateRows 테스트 데이터 생성
func GenerateRows(data [][]string, headers []string) []*pb.Row {
	//rows := []*pb.Row{}
	rows := make([]*pb.Row, 0, len(data))
	for i, cells := range data {
		row := &pb.Row{
			RowNumber:   int32(i + 1), // 1부터 시작
			CellColumns: make(map[string]string, len(headers)),
		}
		for j, header := range headers {
			if j < len(cells) {
				row.CellColumns[header] = cells[j]
			}
		}
		rows = append(rows, row)
	}
	return rows
}
