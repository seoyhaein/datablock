package parsing

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

/*func main() {
	// var totalfiles [][]string
	// 루트 디렉토리에서 각 디렉토리를 recursive 하게 탐색해서 파일명을 가져와야 한다.

	path := "/tmp/testfiles"
	// 테스트로 빈파일 생성
	MakeTestFiles(path)
	files, err := ReadAllFileNames1(path)
	if err != nil { // 에러 발생 시 종료
		fmt.Println("Error reading file names:", err)
		os.Exit(1)
	}
	for _, file := range files {
		fmt.Printf("%s\n", file)
	}
	// 그룹화 했음.
	groupedFiles := GroupFilesByExtension(files)

	err = ValidateKeyCount(groupedFiles, 1) // maxKeys를 1로 설정
	if err != nil {
		// key 가 여러개일경우 해당 디렉토리에 여러 확장자의 파일이 있다 이 경우 경고로 알려주고 datablock 작업을 하지 못한다.????
		// datablock 은 같은 확장자에 한해서 가능하도록 한다. 또는 그것을 구분지어준다.
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed: map contains acceptable number of keys.")
		isfastq, err := IsFastqFiles(groupedFiles)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		// fastq 일 경우
		if isfastq {
			// 일루니마 것인지 판단
			// 일루미나 포멧이라면 각 파트를 가져와서 구조체에 넣어줌.
			// 일루미나 포멧이 아니라면.
			// 일단 파일이름 구조체에 넣더움.
		}
		// fastq 가 아닐 경우
		// 일단 파일이름 구조체에 넣더움.
	}

	ExtractFileParts(files) // return [][]string
	// 그룹핑 해줘야 하는데...
	// 0. 특정 디렉토리 중심으로 필터링 됨.
	// 0.1 파일 확장자가 모두 동일한지를 먼저 판단해야함.
	// 0.1.1 만약 다를 경우 구분지어주어야 함. 그리고 error 리턴
	// 1. 파일 확장자로 먼저 일단 필터링.
	// 1.0 fastq 인경우 illumina 파일인지 먼저 판단함.
	// smaple, range 등등 구분함.
	// 1.0.1 paired end 인지 single 인지 판단함.

}*/

// 일단 임시로 이렇게 만듬.
var (
	fileName string
	filePath string
)

// 여기에 어떻게 잘 담을지 고민해야함.
type DataBlock struct {
	Original [][]string // 원본 데이터만 저장
	RowKeys  []string   // 행 키
	ColKeys  []string   // 열 키
}

func extractParts(fileName string) []string {
	var parts []string
	var currentPart []rune

	// 문자 유형을 판별하는 함수
	getRuneType := func(r rune) string {
		switch {
		case unicode.IsLetter(r):
			return "letter"
		case unicode.IsDigit(r):
			return "digit"
		case unicode.IsSpace(r):
			return "space"
		case unicode.IsSymbol(r):
			return "symbol"
		case unicode.IsPunct(r):
			return "punctuation"
		default:
			return "other"
		}
	}

	// 이전 문자의 유형
	var previousType string

	for _, r := range fileName {
		// 현재 문자의 유형
		currentType := getRuneType(r)

		// 현재 문자와 이전 문자의 유형이 다르면 새로운 파트 시작
		if len(currentPart) > 0 && currentType != previousType {
			parts = append(parts, string(currentPart)) // 이전 파트 저장
			currentPart = []rune{}                     // 새로운 파트 시작
		}

		// 현재 문자를 현재 파트에 추가
		currentPart = append(currentPart, r)
		previousType = currentType
	}

	// 마지막 남은 파트를 추가
	if len(currentPart) > 0 {
		parts = append(parts, string(currentPart))
	}

	return parts
}

// hasCombiningCharacter 조합문자 검출. 일단 조합문자는 막는 걸로 한다.
func hasCombiningCharacter(s string) bool {
	for _, r := range s {
		if unicode.IsMark(r) {
			return true
		}
	}
	return false
}

// 일단 유전체 데이터 일반적인 유전체 파일명 기준으로 맞추고 그것이 맞지 않는다면 별도의 파일명 기준을 찾아내는 방식으로 맞추는 방식으로 진행한다.

/*func isIlluminaParts(parts []string) bool {
	// Illumina 주요 파츠 검증
	// sample1_S1_L001_R1_001.fastq.gz를 예상
	partsCount := len(parts)
	if partsCount < 3 {
		return false
	} // 최소 3개 이상 sample.fastq (smaple . fastq)

	// "S1": 샘플 번호는 S로 시작하고 뒤에 숫자
	sampleIndex := -1
	for i, part := range parts {
		if strings.HasPrefix(part, "S") && isNumber(part[1:]) {
			sampleIndex = i
			break
		}
	}
	if sampleIndex == -1 {
		return false // S1이 없음
	}

	// "L001": 레인 번호는 L로 시작하고 뒤에 숫자
	laneIndex := -1
	for i, part := range parts {
		if strings.HasPrefix(part, "L") && isNumber(part[1:]) {
			laneIndex = i
			break
		}
	}
	if laneIndex == -1 {
		return false // L001이 없음
	}

	// "R1" 또는 "R2": 읽기 방향
	readIndex := -1
	for i, part := range parts {
		if part == "R1" || part == "R2" {
			readIndex = i
			break
		}
	}
	if readIndex == -1 {
		return false // R1 또는 R2가 없음
	}

	// "001": 반복 번호는 숫자여야 함
	repeatIndex := -1
	for i, part := range parts {
		if isNumber(part) {
			repeatIndex = i
			break
		}
	}
	if repeatIndex == -1 {
		return false // 반복 번호가 없음
	}

	// ".fastq.gz" 또는 ".fastq"로 끝나야 함
	if !strings.HasSuffix(parts[len(parts)-1], ".fastq.gz") && !strings.HasSuffix(parts[len(parts)-1], ".fastq") {
		return false // 확장자가 올바르지 않음
	}

	return true
}*/

// 이런식으로 하면 모든 확장자를 다 확인 해야 하는 어려움이 있음.
func findSuffixFastq(parts []string) (bool, string) {
	// sample.fastq 가 최소 조합임. => "sample", ".", "fastq"
	if len(parts) < 3 {
		return false, ""
	}

	if strings.HasSuffix(parts[len(parts)-1], "fastq") {
		return true, "fastq"
	}

	if strings.HasSuffix(parts[len(parts)-1], "fq") {
		return true, "fq"
	}
	// sample.fastq.gz 일 경우. => "sample", ".", "fastq", ".", "gz"
	if strings.HasSuffix(parts[len(parts)-1], "gz") {
		if strings.HasSuffix(parts[len(parts)-3], "fastq") {
			return true, "fastq.gz"
		}

		if strings.HasSuffix(parts[len(parts)-3], "fq") {
			return true, "fq.gz"
		}
	}

	return false, ""
}

// isNumber: 문자열이 숫자인지 확인
func isNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// 먼저 확장자 별로 묶어둔다.
// 특정 디렉토리에서 확장자별로 묶는다. 만약 특정 디렉토리에 확장자가 다른 것들이 있다면 각각 묶지만 에러를 리턴한다.
// 일단 파일을 묶는 것은 proto 방식으로 파일을 묶는다.

// ReadAllFileNames 특정 디렉토리에서 모든 파일 이름을 읽어오는 함수
func ReadAllFileNames(dirPath string) ([]string, error) {
	// 디렉토리의 파일 목록 읽기
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	// 파일 이름을 저장할 슬라이스
	var fileNames []string

	// 파일 목록에서 이름만 추출
	for _, file := range files {
		// 파일 이름을 경로와 함께 추가
		fileNames = append(fileNames, filepath.Join(dirPath, file.Name())) // Append to the End
	}

	return fileNames, nil
}

// ReadAllFileNames1 일단 임시로 만듬.
func ReadAllFileNames1(dirPath string) ([]string, error) {
	// 디렉토리의 파일 목록 읽기
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	// 파일 이름을 저장할 슬라이스
	var fileNames []string

	// 파일 목록에서 이름만 추출
	for _, file := range files {
		// 파일 이름을 경로와 함께 추가
		fileNames = append(fileNames, file.Name()) // Append to the End
	}

	return fileNames, nil
}

func ExtractFileParts(files []string) ([][]string, error) {
	var data [][]string
	var skippedFiles []string // 조합 문자가 포함된 파일 이름 목록

	for _, file := range files {
		if hasCombiningCharacter(file) {
			skippedFiles = append(skippedFiles, file) // 조합 문자 파일을 건너뜀
			continue
		}
		parts := extractParts(file)
		data = append(data, parts)
	}

	// 조합 문자 포함 파일이 있었는지 확인
	if len(skippedFiles) > 0 {
		errMsg := fmt.Sprintf("Skipped files with combining characters: %v", skippedFiles)
		return data, fmt.Errorf(errMsg)
	}

	return data, nil
}

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
		"sample2_S2_L001_R1_001.fastq.gz",
		"sample2_S2_L001_R2_001.fastq.gz",
	}

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

// GetAllExtensions returns all extensions in a file name as a slice of strings.
func GetAllExtensions(fileName string) []string {
	var extensions []string
	for {
		ext := filepath.Ext(fileName)
		if ext == "" {
			break
		}
		extensions = append([]string{ext}, extensions...) // Prepend to maintain order
		fileName = strings.TrimSuffix(fileName, ext)      // Remove the last extension
	}
	return extensions
}

func GetAllExtensionsAsString(fileName string) string {
	var extensions []string
	for {
		ext := filepath.Ext(fileName)
		if ext == "" {
			break
		}
		extensions = append([]string{ext}, extensions...) // Prepend to maintain order
		fileName = strings.TrimSuffix(fileName, ext)      // Remove the last extension
	}
	return strings.Join(extensions, "") // Concatenate all extensions
}

func GetAllExtensionsAsString1(fileName string) (baseName string, extensions string) {
	var extList []string

	for {
		ext := filepath.Ext(fileName)
		if ext == "" {
			break
		}
		extList = append([]string{ext}, extList...)  // Prepend to maintain order
		fileName = strings.TrimSuffix(fileName, ext) // Remove the last extension
	}

	baseName = fileName
	extensions = strings.Join(extList, "") // Combine extensions

	return baseName, extensions
}

func GetAllExtensionsAsStringOptimized(fileName string) string {
	var extensions string
	for {
		ext := filepath.Ext(fileName)
		if ext == "" {
			break
		}
		extensions = ext + extensions                // Prepend directly to the result string
		fileName = strings.TrimSuffix(fileName, ext) // Remove the last extension
	}
	return extensions
}

// GroupFilesByExtension 확장자별로 파일 이름을 그룹화
func GroupFilesByExtension(fileNames []string) map[string][]string {
	groups := make(map[string][]string)

	for _, fileName := range fileNames {
		// 파일 확장자를 추출
		baseName, extensions := GetAllExtensionsAsString1(fileName)

		// 확장자를 키로 그룹에 추가
		groups[extensions] = append(groups[extensions], baseName)
	}
	return groups
}

// IsFastqFiles fastq 파일인지
func IsFastqFiles(groups map[string][]string) (bool, error) {
	// 맵의 키 개수를 확인
	if len(groups) != 1 {
		return false, fmt.Errorf("map must contain exactly one key, but it has %d keys", len(groups))
	}

	// 맵의 첫 번째 키를 가져오기
	for ext := range groups {
		// fastq 또는 fq 포함 여부 반환
		return strings.Contains(ext, "fastq") || strings.Contains(ext, "fq"), nil
	}

	// 이 부분은 실행되지 않음 (안전성 유지용)
	return false, fmt.Errorf("unexpected error: map is empty after length check")
}

// ValidateKeyCount 키 개수 검증 메서드
func ValidateKeyCount(groupedFiles map[string][]string, maxKeys int) error {
	if len(groupedFiles) > maxKeys {
		return fmt.Errorf("error: map contains more than %d keys; found %d keys", maxKeys, len(groupedFiles))
	}
	return nil
}

/*
// 검증 실행
	err := ValidateKeyCount(groupedFiles, 1) // maxKeys를 1로 설정
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed: map contains acceptable number of keys.")
	}
*/
