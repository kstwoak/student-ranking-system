package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// Student 结构体用于存储学生信息
type Student struct {
	ID        int    // 序号
	Subjects  string // 选科
	ExamID    string // 考号
	Name      string // 学生姓名
	Class     string // 班级
	TotalScore int    // 总分
	Rank      int    // 排名
	ChineseMathEnglish int // 语数外
	Physics   string // 物理
	Chemistry string // 化学
	Biology   string // 生物
	Politics  string // 政治
	History   string // 历史
	Geography string // 地理
	Flag      string // 标识
}

// SubjectCombination 定义选科组合类型
type SubjectCombination string

// 选科组合常量
const (
	PoliticsHistoryGeography SubjectCombination = "政史地"
	PhysicsChemistryBiology  SubjectCombination = "物化生"
	PhysicsChemistryGeography SubjectCombination = "物化地"
	PoliticsBiologyHistory   SubjectCombination = "政生史"
	PhysicsChemistryPolitics SubjectCombination = "物化政"
	PoliticsBiologyGeography SubjectCombination = "政生地"
	HistoryGeographyChemistry SubjectCombination = "史地化"
	PhysicsChemistryHistory  SubjectCombination = "物化史"
	PhysicsPoliticsBiology   SubjectCombination = "物政生"
)

// SortStudentsForRank 对学生列表进行排序，用于计算排名
// 规则：按总分降序排序，不考虑标识的顺序
func SortStudentsForRank(students []Student) []Student {
	// 复制切片以避免修改原数据
	sortedStudents := make([]Student, len(students))
	copy(sortedStudents, students)

	// 自定义排序规则
	sort.Slice(sortedStudents, func(i, j int) bool {
		// 首先按总分降序排序
		if sortedStudents[i].TotalScore != sortedStudents[j].TotalScore {
			return sortedStudents[i].TotalScore > sortedStudents[j].TotalScore
		}

		// 总分相同的情况下，按ID升序排序（确保稳定性）
		return sortedStudents[i].ID < sortedStudents[j].ID
	})

	return sortedStudents
}

// SortStudentsForChineseMathEnglish 对学生列表进行排序，用于语数外的赋值
// 规则：1. 首先按总分降序排序
//      2. 标识为0的学生排在最前面
//      3. 标识为"西"的学生排在0和1之间
//      4. 标识为1的学生排在中间
//      5. 标识为2的学生排在最后
func SortStudentsForChineseMathEnglish(students []Student) []Student {
	// 复制切片以避免修改原数据
	sortedStudents := make([]Student, len(students))
	copy(sortedStudents, students)

	// 自定义排序规则
	sort.Slice(sortedStudents, func(i, j int) bool {
		// 首先比较标识：按照 0 < 西 < 1 < 2 的顺序排序
		priorityI := getFlagPriority(sortedStudents[i].Flag)
		priorityJ := getFlagPriority(sortedStudents[j].Flag)
		if priorityI != priorityJ {
			return priorityI < priorityJ
		}

		// 优先级相同的情况下，按总分降序排序
		if sortedStudents[i].TotalScore != sortedStudents[j].TotalScore {
			return sortedStudents[i].TotalScore > sortedStudents[j].TotalScore
		}

		// 总分相同的情况下，按ID升序排序（确保稳定性）
		return sortedStudents[i].ID < sortedStudents[j].ID
	})

	return sortedStudents
}

// getFlagPriority 获取标识的优先级，用于排序
// 规则：0 < 西 < 1 < 2
func getFlagPriority(flag string) int {
	switch flag {
	case "0":
		return 0
	case "西":
		return 1
	case "1":
		return 2
	case "2":
		return 3
	default:
		return 4 // 其他标识值排在最后
	}
}

// CalculateRank 计算学生排名
func CalculateRank(students []Student) []Student {
	// 复制切片以避免修改原数据
	rankedStudents := make([]Student, len(students))
	copy(rankedStudents, students)

	// 计算排名
	for i := range rankedStudents {
		rankedStudents[i].Rank = i + 1
	}

	return rankedStudents
}

// AssignChineseMathEnglish 根据排序顺序对语数外字段进行赋值
// 语数外的赋值按照 1、2、3 的顺序进行，相当于设置考号
// 排序规则：考虑"西"这个标识符的影响
func AssignChineseMathEnglish(students []Student) []Student {
	// 复制切片以避免修改原数据
	assignedStudents := make([]Student, len(students))
	copy(assignedStudents, students)

	// 先对学生进行排序，按照标识和总分排序（考虑"西"的影响）
	sortedStudents := SortStudentsForChineseMathEnglish(assignedStudents)

	// 根据排序顺序对语数外字段进行赋值
	for i := range sortedStudents {
		// 语数外的值按照排序顺序赋值，从1开始递增
		sortedStudents[i].ChineseMathEnglish = i + 1
	}

	return sortedStudents
}

// 主函数
func main() {
	// 读取CSV文件
	students, err := readCSVFile("students.csv")
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return
	}

	// 首先对学生进行排序，用于计算排名（不考虑"西"的影响）
	rankSortedStudents := SortStudentsForRank(students)

	// 根据排名顺序对语数外字段进行赋值（考虑"西"的影响）
	assignedStudents := AssignChineseMathEnglish(students)

	// 对物理字段进行赋值（根据选科是否包含"物"字）
	physicsAssignedStudents := AssignPhysics(students)

	// 对化学字段进行赋值（根据选科是否包含"化"字）
	chemistryAssignedStudents := AssignChemistry(students)

	// 对生物字段进行赋值（根据选科是否包含"生"字）
	biologyAssignedStudents := AssignBiology(students)

	// 对政治字段进行赋值（根据选科是否包含"政"字）
	politicsAssignedStudents := AssignPolitics(students)

	// 对历史字段进行赋值（根据选科是否包含"史"字）
	historyAssignedStudents := AssignHistory(students)

	// 对地理字段进行赋值（根据选科是否包含"地"字）
	geographyAssignedStudents := AssignGeography(students)

	// 创建映射，用于快速查找每个学生的值
	chineseMathEnglishMap := make(map[string]int) // key: 考号
	physicsMap := make(map[string]string) // key: 考号
	chemistryMap := make(map[string]string) // key: 考号
	biologyMap := make(map[string]string) // key: 考号
	politicsMap := make(map[string]string) // key: 考号
	historyMap := make(map[string]string) // key: 考号
	geographyMap := make(map[string]string) // key: 考号

	for _, student := range assignedStudents {
		chineseMathEnglishMap[student.ExamID] = student.ChineseMathEnglish
	}
	for _, student := range physicsAssignedStudents {
		physicsMap[student.ExamID] = student.Physics
	}
	for _, student := range chemistryAssignedStudents {
		chemistryMap[student.ExamID] = student.Chemistry
	}
	for _, student := range biologyAssignedStudents {
		biologyMap[student.ExamID] = student.Biology
	}
	for _, student := range politicsAssignedStudents {
		politicsMap[student.ExamID] = student.Politics
	}
	for _, student := range historyAssignedStudents {
		historyMap[student.ExamID] = student.History
	}
	for _, student := range geographyAssignedStudents {
		geographyMap[student.ExamID] = student.Geography
	}

	// 重新计算排名
	rankedStudents := CalculateRank(rankSortedStudents)

	// 将语数外、物理、化学、生物、政治、历史、地理值赋值给排名后的学生
	for i := range rankedStudents {
		if cme, ok := chineseMathEnglishMap[rankedStudents[i].ExamID]; ok {
			rankedStudents[i].ChineseMathEnglish = cme
		}
		if physics, ok := physicsMap[rankedStudents[i].ExamID]; ok {
			rankedStudents[i].Physics = physics
		}
		if chemistry, ok := chemistryMap[rankedStudents[i].ExamID]; ok {
			rankedStudents[i].Chemistry = chemistry
		}
		if biology, ok := biologyMap[rankedStudents[i].ExamID]; ok {
			rankedStudents[i].Biology = biology
		}
		if politics, ok := politicsMap[rankedStudents[i].ExamID]; ok {
			rankedStudents[i].Politics = politics
		}
		if history, ok := historyMap[rankedStudents[i].ExamID]; ok {
			rankedStudents[i].History = history
		}
		if geography, ok := geographyMap[rankedStudents[i].ExamID]; ok {
			rankedStudents[i].Geography = geography
		}
	}

	// 打印学生信息
	for _, student := range rankedStudents {
		fmt.Printf("序号: %d, 选科: %s, 考号: %s, 学生姓名: %s, 班级: %s, 总分: %d, 排名: %d, 语数外: %d, 物理: %s, 化学: %s, 生物: %s, 政治: %s, 历史: %s, 地理: %s, 标识: %s\n",
			student.ID, student.Subjects, student.ExamID, student.Name, student.Class, student.TotalScore, student.Rank, student.ChineseMathEnglish, student.Physics, student.Chemistry, student.Biology, student.Politics, student.History, student.Geography, student.Flag)
	}

	// 分析选科组合
	analyzeSubjectCombinations(rankedStudents)

	// 分析标识字段
	analyzeFlags(rankedStudents)

	// 将结果写入文件
	if err := writeResultFile(rankedStudents, "result.csv"); err != nil {
		fmt.Printf("写入结果文件失败: %v\n", err)
		return
	}
	fmt.Println("\n结果已写入到 result.csv 文件中")
}

// parseInt 解析整数，处理空字符串情况
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	val, _ := strconv.Atoi(s)
	return val
}

// readCSVFile 读取CSV文件并返回学生信息列表
func readCSVFile(filePath string) ([]Student, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建CSV读取器
	reader := csv.NewReader(file)

	// 读取表头
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// 验证表头
	expectedHeader := []string{"序号", "选科", "考号", "学生姓名", "班级", "总分", "排名", "语数外", "物理", "化学", "生物", "政治", "历史", "地理", "标识"}
	if len(header) != len(expectedHeader) {
		return nil, fmt.Errorf("表头长度不匹配，期望 %d 列，实际 %d 列", len(expectedHeader), len(header))
	}

	// 读取数据行
	var students []Student
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		// 解析数据
		id := parseInt(row[0])
		totalScore := parseInt(row[5])
		rank := parseInt(row[6])
		chineseMathEnglish := parseInt(row[7]) // 解析语数外列的整数值
		physics := row[8] // 直接读取物理列的字符串值
		chemistry := row[9] // 直接读取化学列的字符串值
		biology := row[10] // 直接读取生物列的字符串值
		politics := row[11] // 直接读取政治列的字符串值
		history := row[12] // 直接读取历史列的字符串值
		geography := row[13] // 直接读取地理列的字符串值
		flag := row[14] // 直接读取标识列的字符串值，包括"西"

		// 创建学生对象
		student := Student{
			ID:                  id,
			Subjects:            row[1],
			ExamID:              row[2],
			Name:                row[3],
			Class:               row[4],
			TotalScore:          totalScore,
			Rank:                rank,
			ChineseMathEnglish:  chineseMathEnglish,
			Physics:             physics,
			Chemistry:           chemistry,
			Biology:             biology,
			Politics:            politics,
			History:             history,
			Geography:           geography,
			Flag:                flag,
		}

		students = append(students, student)
	}

	return students, nil
}

// analyzeSubjectCombinations 分析选科组合
func analyzeSubjectCombinations(students []Student) {
	// 统计各选科组合的学生数量
	combinationCount := make(map[string]int)
	for _, student := range students {
		combinationCount[student.Subjects]++
	}

	// 打印统计结果
	fmt.Println("\n选科组合统计:")
	for combination, count := range combinationCount {
		fmt.Printf("%s: %d 人\n", combination, count)
	}
}

// analyzeFlags 分析标识字段
func analyzeFlags(students []Student) {
	// 统计标识的学生数量
	flagCount := make(map[string]int)
	for _, student := range students {
		flagCount[student.Flag]++
	}

	// 打印统计结果
	fmt.Println("\n标识统计:")
	for flag, count := range flagCount {
		fmt.Printf("标识 %s: %d 人\n", flag, count)
	}
}

// AssignSubject 根据选科是否包含指定关键字对学科字段进行赋值
// 规则：1. 选科包含关键字的学生，字段直接赋值为数字（如1, 2, 3...），超过10时继续递增
//      2. 选科不包含关键字的学生，字段赋值为自习+数字，从包含关键字的学生数量+1开始（如自习11, 自习12, 自习13...）
//      3. 排序时，先按标识优先级（0 < 1 < 2）排序，再按总分降序排序
func AssignSubject(students []Student, keyword string) map[string]string {
	// 分离包含关键字和不包含关键字的学生
	var keywordStudents, nonKeywordStudents []Student
	for _, student := range students {
		if containsString(student.Subjects, keyword) {
			keywordStudents = append(keywordStudents, student)
		} else {
			nonKeywordStudents = append(nonKeywordStudents, student)
		}
	}

	// 对包含关键字的学生排序：先按标识优先级（0 < 1 < 2），再按总分降序
	sort.Slice(keywordStudents, func(i, j int) bool {
		// 首先比较标识优先级
		flagPriorityI := getSubjectFlagPriority(keywordStudents[i].Flag)
		flagPriorityJ := getSubjectFlagPriority(keywordStudents[j].Flag)
		if flagPriorityI != flagPriorityJ {
			return flagPriorityI < flagPriorityJ
		}
		
		// 标识优先级相同的情况下，按总分降序排序
		if keywordStudents[i].TotalScore != keywordStudents[j].TotalScore {
			return keywordStudents[i].TotalScore > keywordStudents[j].TotalScore
		}
		
		// 总分相同的情况下，按ID升序排序（确保稳定性）
		return keywordStudents[i].ID < keywordStudents[j].ID
	})

	// 对不包含关键字的学生排序：先按标识优先级（0 < 1 < 2），再按总分降序
	sort.Slice(nonKeywordStudents, func(i, j int) bool {
		// 首先比较标识优先级
		flagPriorityI := getSubjectFlagPriority(nonKeywordStudents[i].Flag)
		flagPriorityJ := getSubjectFlagPriority(nonKeywordStudents[j].Flag)
		if flagPriorityI != flagPriorityJ {
			return flagPriorityI < flagPriorityJ
		}
		
		// 标识优先级相同的情况下，按总分降序排序
		if nonKeywordStudents[i].TotalScore != nonKeywordStudents[j].TotalScore {
			return nonKeywordStudents[i].TotalScore > nonKeywordStudents[j].TotalScore
		}
		
		// 总分相同的情况下，按ID升序排序（确保稳定性）
		return nonKeywordStudents[i].ID < nonKeywordStudents[j].ID
	})

	// 创建映射，用于快速查找每个学生的学科值
	subjectMap := make(map[string]string) // key: 考号

	// 为包含关键字的学生赋值，从1开始递增
	for i, student := range keywordStudents {
		score := i + 1
		subjectMap[student.ExamID] = strconv.Itoa(score)
	}

	// 为不包含关键字的学生赋值，从包含关键字的学生数量+1开始
	keywordCount := len(keywordStudents)
	for i, student := range nonKeywordStudents {
		score := keywordCount + i + 1
		subjectMap[student.ExamID] = "自习" + strconv.Itoa(score)
	}

	return subjectMap
}

// getSubjectFlagPriority 获取学科排序时标识的优先级
// 规则：西=0 < 1 < 2，其他标识值排在最后
func getSubjectFlagPriority(flag string) int {
	switch flag {
	case "0", "西":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	default:
		return 3 // 其他标识值排在最后
	}
}

// AssignPhysics 根据选科是否包含"物"字对物理字段进行赋值
func AssignPhysics(students []Student) []Student {
	// 复制切片以避免修改原数据
	assignedStudents := make([]Student, len(students))
	copy(assignedStudents, students)

	// 获取物理值映射
	physicsMap := AssignSubject(students, "物")

	// 将物理值赋值给学生
	for i := range assignedStudents {
		if physics, ok := physicsMap[assignedStudents[i].ExamID]; ok {
			assignedStudents[i].Physics = physics
		}
	}

	return assignedStudents
}

// AssignChemistry 根据选科是否包含"化"字对化学字段进行赋值
func AssignChemistry(students []Student) []Student {
	// 复制切片以避免修改原数据
	assignedStudents := make([]Student, len(students))
	copy(assignedStudents, students)

	// 获取化学值映射
	chemistryMap := AssignSubject(students, "化")

	// 将化学值赋值给学生
	for i := range assignedStudents {
		if chemistry, ok := chemistryMap[assignedStudents[i].ExamID]; ok {
			assignedStudents[i].Chemistry = chemistry
		}
	}

	return assignedStudents
}

// AssignBiology 根据选科是否包含"生"字对生物字段进行赋值
func AssignBiology(students []Student) []Student {
	// 复制切片以避免修改原数据
	assignedStudents := make([]Student, len(students))
	copy(assignedStudents, students)

	// 获取生物值映射
	biologyMap := AssignSubject(students, "生")

	// 将生物值赋值给学生
	for i := range assignedStudents {
		if biology, ok := biologyMap[assignedStudents[i].ExamID]; ok {
			assignedStudents[i].Biology = biology
		}
	}

	return assignedStudents
}

// AssignPolitics 根据选科是否包含"政"字对政治字段进行赋值
func AssignPolitics(students []Student) []Student {
	// 复制切片以避免修改原数据
	assignedStudents := make([]Student, len(students))
	copy(assignedStudents, students)

	// 获取政治值映射
	politicsMap := AssignSubject(students, "政")

	// 将政治值赋值给学生
	for i := range assignedStudents {
		if politics, ok := politicsMap[assignedStudents[i].ExamID]; ok {
			assignedStudents[i].Politics = politics
		}
	}

	return assignedStudents
}

// AssignHistory 根据选科是否包含"史"字对历史字段进行赋值
func AssignHistory(students []Student) []Student {
	// 复制切片以避免修改原数据
	assignedStudents := make([]Student, len(students))
	copy(assignedStudents, students)

	// 获取历史值映射
	historyMap := AssignSubject(students, "史")

	// 将历史值赋值给学生
	for i := range assignedStudents {
		if history, ok := historyMap[assignedStudents[i].ExamID]; ok {
			assignedStudents[i].History = history
		}
	}

	return assignedStudents
}

// AssignGeography 根据选科是否包含"地"字对地理字段进行赋值
func AssignGeography(students []Student) []Student {
	// 复制切片以避免修改原数据
	assignedStudents := make([]Student, len(students))
	copy(assignedStudents, students)

	// 获取地理值映射
	geographyMap := AssignSubject(students, "地")

	// 将地理值赋值给学生
	for i := range assignedStudents {
		if geography, ok := geographyMap[assignedStudents[i].ExamID]; ok {
			assignedStudents[i].Geography = geography
		}
	}

	return assignedStudents
}

// containsString 检查字符串是否包含指定子串
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// writeResultFile 将结果写入CSV文件
func writeResultFile(students []Student, filePath string) error {
	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建CSV写入器
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	header := []string{"序号", "选科", "考号", "学生姓名", "班级", "总分", "排名", "语数外", "物理", "化学", "生物", "政治", "历史", "地理", "标识"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// 写入数据行
	for _, student := range students {
		row := []string{
			strconv.Itoa(student.ID),
			student.Subjects,
			student.ExamID,
			student.Name,
			student.Class,
			strconv.Itoa(student.TotalScore),
			strconv.Itoa(student.Rank),
			strconv.Itoa(student.ChineseMathEnglish),
			student.Physics,
			student.Chemistry,
			student.Biology,
			student.Politics,
			student.History,
			student.Geography,
			student.Flag,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
