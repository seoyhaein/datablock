package rule

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// RuleSet 구조체 정의

type RuleSet struct {
	Version     string      `json:"version"`
	Delimiter   []string    `json:"delimiter"`
	Header      []string    `json:"header"`
	RowRules    RowRules    `json:"rowRules"`
	ColumnRules ColumnRules `json:"columnRules"`
	SizeRules   SizeRules   `json:"sizeRules"`
}

type RowRules struct {
	MatchParts []int `json:"matchParts"`
}

type ColumnRules struct {
	MatchParts []int `json:"matchParts"`
}

type SizeRules struct {
	MinSize int `json:"minSize"`
	MaxSize int `json:"maxSize"`
}

// LoadRuleSetFromFile JSON 파일을 읽어 RuleSet 구조체로 디코딩
func LoadRuleSetFromFile(filePath string) (RuleSet, error) {
	// JSON 파일 읽기
	data, err := os.ReadFile(filePath)
	if err != nil {
		return RuleSet{}, fmt.Errorf("failed to read file: %w", err)
	}

	// JSON 데이터 디코딩
	var ruleSet RuleSet
	err = json.Unmarshal(data, &ruleSet)
	if err != nil {
		return RuleSet{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return ruleSet, nil
}

// Helper function: 파일 이름을 JSON 규칙의 delimiter를 기준으로 파트로 나누기
func extractParts(fileName string, delimiters []string) []string {
	for _, delim := range delimiters {
		fileName = strings.ReplaceAll(fileName, delim, " ")
	}
	return strings.Fields(fileName)
}

// BlockifyFilesToMap 파일 이름을 JSON 규칙에 따라 블록화하여 맵으로 변환
func BlockifyFilesToMap(fileNames []string, ruleSet RuleSet) (map[int]map[string]string, error) {
	rowMap := make(map[string]int)               // Row Key → Row Index
	rowCounter := 0                              // 행 인덱스 카운터 (0부터 시작)
	resultMap := make(map[int]map[string]string) // 결과 데이터 저장용 맵

	for _, fileName := range fileNames {
		// 파일명을 JSON 규칙의 delimiter를 기준으로 분리
		parts := extractParts(fileName, ruleSet.Delimiter)

		// Row Key 생성
		var rowKeyParts []string
		for _, idx := range ruleSet.RowRules.MatchParts {
			if idx < len(parts) {
				rowKeyParts = append(rowKeyParts, parts[idx])
			}
		}
		rowKey := strings.Join(rowKeyParts, "_")

		// Row Index 확인 및 추가
		if _, exists := rowMap[rowKey]; !exists {
			rowMap[rowKey] = rowCounter
			resultMap[rowCounter] = make(map[string]string)
			rowCounter++
		}

		// Column Key 생성 (ColumnRules.MatchParts 기준)
		var colKeyParts []string
		for _, idx := range ruleSet.ColumnRules.MatchParts {
			if idx < len(parts) {
				colKeyParts = append(colKeyParts, parts[idx])
			}
		}
		colKey := strings.Join(colKeyParts, "_")

		// Row에 Column Key와 파일명 추가
		rowIdx := rowMap[rowKey]
		resultMap[rowIdx][colKey] = fileName
	}

	return resultMap, nil
}

// ValidateRuleSet validates the given rule set for conflicts and unused parts.
func ValidateRuleSet(ruleSet RuleSet) bool {
	hasConflict := false // 충돌 여부를 저장

	// Helper 함수: 충돌 감지 및 로깅
	logConflict := func(message string, idx int) {
		log.Printf(message, idx)
		hasConflict = true
	}

	// Row와 Column 규칙 매핑
	rowMatch := make(map[int]struct{})
	colMatch := make(map[int]struct{})

	for _, idx := range ruleSet.RowRules.MatchParts {
		rowMatch[idx] = struct{}{}
	}
	for _, idx := range ruleSet.ColumnRules.MatchParts {
		colMatch[idx] = struct{}{}
	}

	// Row와 Column 규칙의 MatchParts와 UnMatchParts에서 충돌 확인
	for idx := range rowMatch {
		if _, exists := colMatch[idx]; exists {
			logConflict("Conflict detected: part %d is in both RowRules.MatchParts and ColumnRules.MatchParts", idx)
		}
	}

	// 최종 결과 반환
	return !hasConflict
}

// SaveResultMapToCSV map[int]map[string]string 데이터를 CSV 파일로 저장
func SaveResultMapToCSV1(filePath string, resultMap map[int]map[string]string, headers []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 첫 번째 행에 헤더 추가
	headerRow := append([]string{"Row"}, headers...)
	if err := writer.Write(headerRow); err != nil {
		return fmt.Errorf("failed to write header row: %w", err)
	}

	// 각 행 데이터를 CSV에 추가
	for rowIdx := 0; rowIdx < len(resultMap); rowIdx++ {
		rowData := make([]string, len(headers)+1) // +1은 "Row" 열 때문
		rowData[0] = fmt.Sprintf("Row%d", rowIdx)
		// row 는 map[key]filename 임.
		if row, exists := resultMap[rowIdx]; exists {
			for colIdx, header := range headers {
				rowData[colIdx+1] = row[header]
			}
		}

		if err := writer.Write(rowData); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

// SaveResultMapToCSV map[int]map[string]string 데이터를 CSV 파일로 저장
func SaveResultMapToCSV(filePath string, resultMap map[int]map[string]string, headers []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 첫 번째 행에 헤더 추가
	headerRow := append([]string{"Row"}, headers...)
	if err := writer.Write(headerRow); err != nil {
		return fmt.Errorf("failed to write header row: %w", err)
	}

	var columnHeaders []string
	seen := make(map[string]struct{})

	// 모든 열 키를 동적으로 추출 (중복 제거)
	for _, row := range resultMap {
		for key := range row {
			if _, exists := seen[key]; !exists {
				columnHeaders = append(columnHeaders, key)
				seen[key] = struct{}{}
			}
		}
	}

	// 열 키를 정렬
	sort.Strings(columnHeaders)

	// 각 행 데이터를 CSV에 추가
	for rowIdx := 0; rowIdx < len(resultMap); rowIdx++ {
		rowData := append([]string{fmt.Sprintf("Row%d", rowIdx)}, make([]string, len(columnHeaders))...)

		if row, exists := resultMap[rowIdx]; exists {
			for i, colKey := range columnHeaders {
				if value, ok := row[colKey]; ok {
					rowData[i+1] = value
				}
			}
		}

		if err := writer.Write(rowData); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}
