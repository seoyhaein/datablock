package rule

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	u "github.com/seoyhaein/utils"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
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

// LoadRuleSetFromFile JSON 파일을 읽어 RuleSet 구조체로 디코딩. RuleSet 의 경우 값의 수정이 일어나면 안되기때문에 값으로 리턴한다.
func LoadRuleSetFromFile(filePath string) (RuleSet, error) {
	filePath, err := u.CheckPath(filePath)
	if err != nil {
		return RuleSet{}, err
	}

	// filePath 가 디렉토리인지 확인
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return RuleSet{}, fmt.Errorf("failed to access path: %w", err)
	}

	if fileInfo.IsDir() {
		// 디렉토리 내 rule.json 파일 경로 확인
		filePath = filepath.Join(filePath, "rule.json")
	}

	// rule.json 파일 존재 여부 확인
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return RuleSet{}, fmt.Errorf("rule file not found at %s", filePath)
	}

	// JSON 파일 읽기, 파일이 크지 않기 때문에 이렇게 처리 함.
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

// Helper function: 파일 이름을 JSON 규칙의 delimiter 를 기준으로 파트로 나누기
func extractParts(fileName string, delimiters []string) []string {
	for _, delim := range delimiters {
		fileName = strings.ReplaceAll(fileName, delim, " ")
	}
	return strings.Fields(fileName)
}

// FilesToMap 파일 이름을 JSON 규칙에 따라 블록화하여 맵으로 변환
func FilesToMap(fileNames []string, ruleSet RuleSet) (map[int]map[string]string, error) {
	rowMap := make(map[string]int)               // Row Key → Row Index
	rowCounter := 0                              // 행 인덱스 카운터 (0부터 시작)
	resultMap := make(map[int]map[string]string) // 결과 데이터 저장용 맵

	for _, fileName := range fileNames {
		// 파일명을 JSON 규칙의 delimiter 를 기준으로 분리
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

		// Row 에 Column Key 와 파일명 추가
		rowIdx := rowMap[rowKey]
		resultMap[rowIdx][colKey] = fileName
	}

	return resultMap, nil
}

// FilterMap 컬럼 수를 검증하고 유효/무효 행을 분리하는 메서드
func FilterMap(resultMap map[int]map[string]string, expectedColCount int) (map[int]map[string]string, []map[string]string) {
	validRows := make(map[int]map[string]string)
	var invalidRows []map[string]string
	newRowCounter := 0

	for _, row := range resultMap {
		if len(row) == expectedColCount {
			validRows[newRowCounter] = row
			newRowCounter++
		} else {
			invalidRows = append(invalidRows, row)
		}
	}

	return validRows, invalidRows
}

// WriteInvalidFiles invalidRows 의 파일명을 하나의 텍스트 파일에 기록 TODO readonly 로 하는 것 생각
func WriteInvalidFiles(invalidRows []map[string]string, outputFilePath string) (err error) {
	// invalidRows 가 비어있으면 파일을 생성하지 않고 리턴
	if len(invalidRows) == 0 {
		return nil
	}

	// outputFilePath 가 디렉토리인지 확인
	fileInfo, err := os.Stat(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to access path %s: %w", outputFilePath, err)
	}

	// 디렉토리인 경우, 날짜와 시간을 포함한 파일명을 생성
	if fileInfo.IsDir() {
		timestamp := time.Now().Format("20060102150405") // 현재 날짜와 시간 (년월일시간분초)
		outputFilePath = filepath.Join(outputFilePath, fmt.Sprintf("invalid_files_%s.txt", timestamp))
	}

	// 출력 파일 생성 (덮어쓰기)
	file, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputFilePath, err)
	}

	defer func() {
		if cErr := file.Close(); cErr != nil {
			if err == nil {
				err = fmt.Errorf("failed to close file: %w", cErr)
			} else {
				err = fmt.Errorf("%v; failed to close file: %w", err, cErr)
			}
		}
	}()

	// 파일명들을 텍스트 파일에 기록
	for _, row := range invalidRows {
		for _, fileName := range row {
			_, err := file.WriteString(fileName + "\n")
			if err != nil {
				return fmt.Errorf("failed to write to file %s: %w", outputFilePath, err)
			}
		}
	}

	return err
}

// ValidateRuleSet validates the given rule set for conflicts and unused parts.
func ValidateRuleSet(ruleSet RuleSet) bool {
	hasConflict := false

	usageMap := make(map[int][]string)

	// Helper for registering usage of parts
	addUsage := func(indices []int, roleName string) {
		for _, idx := range indices {
			usageMap[idx] = append(usageMap[idx], roleName)
		}
	}

	// Register usages
	addUsage(ruleSet.RowRules.MatchParts, "RowRules.MatchParts")
	addUsage(ruleSet.ColumnRules.MatchParts, "ColumnRules.MatchParts")

	// Check for conflicts - any index used in more than one role is a conflict
	for idx, roles := range usageMap {
		if len(roles) > 1 {
			log.Printf("Conflict detected: part %d is used in multiple roles: %v", idx, roles)
			hasConflict = true
		}
	}

	return !hasConflict
}

// SaveResultMapToCSV map[int]map[string]string 데이터를 CSV 파일로 저장, TODO 파일 생성날짜를 기록할지 생각, readonly 로 하는 것 생각.
func SaveResultMapToCSV(filePath string, resultMap map[int]map[string]string, headers []string) (err error) {
	filePath, err = u.CheckPath(filePath)
	if err != nil {
		return err
	}

	// filePath 가 디렉토리인지 확인
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to access path: %w", err)
	}

	if fileInfo.IsDir() {
		// 디렉토리 경로에 fileblock.csv 파일 생성
		filePath = filepath.Join(filePath, "fileblock.csv")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	//defer file.Close()
	defer func() {
		if cErr := file.Close(); cErr != nil {
			if err == nil {
				err = fmt.Errorf("failed to close file: %w", cErr)
			} else {
				err = fmt.Errorf("%v; failed to close file: %w", err, cErr)
			}
		}
	}()

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

	// 각 행 데이터를 CSV 에 추가
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

	return err
}

// GenerateMap 일단 이름 고침.
func GenerateMap(filePath string) (map[int]map[string]string, error) {
	// Load the rule set
	ruleSet, err := LoadRuleSetFromFile(filePath) // 이 메서드에서 filepath 의 검증을 해줌.
	if err != nil {
		return nil, fmt.Errorf("failed to load rule set: %w", err)
	}

	// Validate the rule set
	if !ValidateRuleSet(ruleSet) {
		return nil, fmt.Errorf("rule set has conflicts or unused parts")
	}

	// Read all file names from the directory
	// 예외 규정: rule.json, invalid_files 로 시작하는 파일, fileblock.csv
	exclusions := []string{"rule.json", "invalid_files", "fileblock.csv"}
	files, err := ReadAllFileNames(filePath, exclusions)

	if err != nil {
		return nil, fmt.Errorf("failed to read file names: %w", err)
	}

	resultMap, err := FilesToMap(files, ruleSet)
	if err != nil {
		return nil, fmt.Errorf("failed to blockify files: %w", err)
	}

	// Filter the result map into valid and invalid rows
	validRows, invalidRows := FilterMap(resultMap, len(ruleSet.Header))

	// Save valid rows to a CSV file
	if err := SaveResultMapToCSV(filePath, validRows, ruleSet.Header); err != nil {
		return nil, fmt.Errorf("failed to save result map to CSV: %w", err)
	}

	// Save invalid rows to a separate file
	if err := WriteInvalidFiles(invalidRows, filePath); err != nil {
		return nil, fmt.Errorf("failed to write invalid files: %w", err)
	}

	return validRows, nil
}

// ReadAllFileNames 디렉토리에서 파일을 읽되 예외 규정에 맞는 파일들은 제외 TODO path 나 디렉토리 관련 정규화 적용할 것.
func ReadAllFileNames(dirPath string, exclusions []string) ([]string, error) {
	// 디렉토리의 파일 목록 읽기
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	// 파일 이름을 저장할 슬라이스
	var fileNames []string

	// 파일 목록에서 제외할 파일들을 걸러내고 이름만 추출
	for _, file := range files {
		fileName := file.Name()

		// 예외 규정에 있는 파일이면 건너뛰기
		if u.ExcludeFiles(fileName, exclusions) {
			continue
		}

		// 파일 이름을 경로와 함께 추가
		fileNames = append(fileNames, fileName)
	}

	return fileNames, nil
}
