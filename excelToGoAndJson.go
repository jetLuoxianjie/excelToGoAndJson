package main

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	excelDir      = "./excel"
	jsonOutputDir = "./jsons"      // JSON 存储目录
	goOutputDir   = "./autoConfig" // JSON 存储目录
)

type DirConf struct {
	ExcelDir      string `json:"excelDir"`
	JsonOutputDir string `json:"jsonOutputDir"`
	GoOutputDir   string `json:"goOutputDir"`
}

var dirConf *DirConf

// 读取./dirConf.json的配置
func readDirConf() {
	file, err := os.Open("./dirConf.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	dirConf = &DirConf{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dirConf)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
}

func main() {
	//0 读取配置
	readDirConf()
	currJsonOutputDir := jsonOutputDir
	currExcelDir := excelDir
	currOoOutputDir := goOutputDir
	if dirConf.JsonOutputDir != "" {
		currJsonOutputDir = dirConf.JsonOutputDir
	}
	if dirConf.ExcelDir != "" {
		currExcelDir = dirConf.ExcelDir
	}
	if dirConf.GoOutputDir != "" {
		currOoOutputDir = dirConf.GoOutputDir
	}
	// 1. 处理 Excel 生成 JSON
	processAllExcelFiles(currExcelDir, currJsonOutputDir)
	// 2. 生成 Golang 结构体
	// 遍历所有 Excel，生成结构体
	structs := ProcessExcelFiles(currExcelDir)
	// 生成 Go 结构体代码
	GenerateGoStructs(structs, currOoOutputDir)
	// 输入任意键结束
	fmt.Println("Press any key to exit...")
	var input string
	fmt.Scanln(&input)
}

// 处理所有 Excel 文件
func processAllExcelFiles(currExcelDir, currJsonOutputDir string) {
	// 清空 JSON 目录
	clearOutputDir(currJsonOutputDir)

	// 创建 JSON 目录
	if err := os.MkdirAll(currJsonOutputDir, os.ModePerm); err != nil {
		log.Fatalf("创建 JSON 目录失败: %v", err)
	}

	// 遍历当前目录下的 Excel 文件
	err := filepath.Walk(currExcelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".xlsx") || strings.HasSuffix(info.Name(), ".xls")) {
			processExcelFile(path, currJsonOutputDir)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("遍历目录时出错: %v", err)
	}
}

// 清空 JSON 目录
func clearOutputDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // 目录可能不存在，不需要清理
	}

	for _, entry := range entries {
		err := os.Remove(filepath.Join(dir, entry.Name()))
		if err != nil {
			log.Printf("删除旧 JSON 失败: %v", err)
		}
	}
}

func processExcelFile(filePath, currJsonOutputDir string) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Printf("无法打开 Excel 文件: %s, 错误: %v\n", filePath, err)
		return
	}
	defer f.Close()
	// 获取所有 sheet 名称
	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		fmt.Printf("处理 Sheet: %s\n", sheet)
		data := parseSheet(f, sheet)
		if len(data) == 0 {
			fmt.Printf("Sheet: %s 为空，跳过。\n", sheet)
			continue
		}

		// 生成 JSON 文件路径 (Sheet名_Excel名.json)
		jsonFileName := fmt.Sprintf("%s.json", sheet)
		jsonFilePath := filepath.Join(currJsonOutputDir, jsonFileName)

		// 保存 JSON
		saveJSON(jsonFilePath, data)
	}
}

// 解析 sheet 数据
func parseSheet(f *excelize.File, sheet string) []map[string]interface{} {
	rows, err := f.GetRows(sheet)
	if err != nil || len(rows) < 3 {
		log.Printf("无法读取 Sheet: %s 或数据行不足 3 行，跳过。\n", sheet)
		return nil
	}

	// 第一行: Go 类型
	typeRow := rows[0]
	// 第二行: 字段名
	fieldRow := rows[1]
	// 第三行: 注释（可忽略）

	var result []map[string]interface{}

	// 读取数据
	for i := 3; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		record := make(map[string]interface{})
		for j, field := range fieldRow {
			if j >= len(row) || field == "" {
				continue
			}

			// 解析数据类型
			val := parseValue(typeRow[j], row[j])
			record[field] = val
		}

		result = append(result, record)
	}

	return result
}

// 根据 Golang 类型解析数据
func parseValue(goType, value string) interface{} {
	switch strings.ToLower(goType) {
	case "int", "int32":
		var v int
		fmt.Sscanf(value, "%d", &v)
		return v
	case "float", "float64":
		var v float64
		fmt.Sscanf(value, "%f", &v)
		return v
	case "bool":
		return strings.ToLower(value) == "true"
	case "string":
		return value
	case "intarray":
		// 去除字符串首尾的方括号
		value = strings.Trim(value, "[]")
		var arr []int
		for _, v := range strings.Split(value, ",") {
			var num int
			fmt.Sscanf(v, "%d", &num)
			arr = append(arr, num)
		}
		return arr
	case "floatarray":
		value = strings.Trim(value, "[]")
		var arr []float64
		for _, v := range strings.Split(value, ",") {
			var num float64
			fmt.Sscanf(v, "%f", &num)
			arr = append(arr, num)
		}
		return arr
	case "stringarray":
		value = strings.Trim(value, "[]")
		var arr []string
		for _, v := range strings.Split(value, ",") {
			arr = append(arr, v)
		}
		return arr
	default:
		return value
	}
}

// 保存 JSON 文件
func saveJSON(filePath string, data interface{}) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("无法创建 JSON 文件: %s, 错误: %v\n", filePath, err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 美化格式
	if err := encoder.Encode(data); err != nil {
		log.Printf("写入 JSON 失败: %s, 错误: %v\n", filePath, err)
	}
}

// ---------

// 配置结构体
type Config struct {
	Name   string
	Fields map[string]string
}

// 结构体存储
var configStructs = make(map[string]interface{})

// 遍历 Excel 目录，处理所有 Excel
func ProcessExcelFiles(dir string) []Config {
	var configs []Config

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".xlsx") {
			fmt.Println("📖 处理 Excel:", path)
			configs = append(configs, ParseExcel(path)...)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("遍历 Excel 失败: %v", err)
	}

	return configs
}

// 解析 Excel
func ParseExcel(filePath string) []Config {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("无法打开 Excel: %v", err)
	}
	defer f.Close()

	var configs []Config
	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet)
		if err != nil || len(rows) < 3 {
			continue
		}

		typeRow := rows[0] // 第一行: Golang 类型
		nameRow := rows[1] // 第二行: 字段名

		fields := make(map[string]string)
		for i, fieldName := range nameRow {
			if fieldName == "" {
				continue
			}
			fields[fieldName] = parseGoType(typeRow[i])
		}

		configs = append(configs, Config{
			Name:   toPascalCase(sheet) + "Config",
			Fields: fields,
		})
	}

	return configs
}

// 解析 Golang 类型
func parseGoType(excelType string) string {
	switch strings.ToLower(excelType) {
	case "int":
		return "int32"
	case "float", "float64":
		return "float64"
	case "string":
		return "string"
	case "intarray":
		return "[]int32"
	case "floatarray":
		return "[]float64"
	case "stringarray":
		return "[]string"
	default:
		return "interface{}"
	}
}

// 生成 Golang 结构体
func GenerateGoStructs(configs []Config, jsonOutputDir string) {
	os.MkdirAll(jsonOutputDir, os.ModePerm)
	structFile := filepath.Join(jsonOutputDir, "autoConfig.go")
	file, err := os.Create(structFile)
	if err != nil {
		log.Fatalf("无法创建文件: %v", err)
	}
	defer file.Close()

	fmt.Fprintf(file, `package config
import (
    "encoding/json"
	"io/ioutil"
	"sync"
)
`)

	str := fmt.Sprintf(`
// 调用初始化读取配置
func InitConfig() {
`)
	fmt.Fprintf(file, str)
	// 添加
	for _, cfg := range configs {
		fmt.Fprintln(file, fmt.Sprintf(`  init%v()`, cfg.Name))
	}
	str = fmt.Sprintf(`
}
`)
	fmt.Fprintf(file, str)

	for _, cfg := range configs {

		fmt.Fprintf(file, `
type %s struct {
`, cfg.Name)
		for field, fieldType := range cfg.Fields {
			fmt.Fprintf(file, "\t%s %s `json:\"%s\"`\n", toPascalCase(field), fieldType, field)
		}
		fmt.Fprintln(file, "}")

		lList := strings.Split(cfg.Name, "Config")
		name := lList[0]
		str = fmt.Sprintf(`
var %vListTable []*%v
var %vListMap sync.Map
// 解析JSON数据到结构体
func init%v(){
  fileContent, err := ioutil.ReadFile("./jsons/%v.json")
	if err != nil {
        panic("%v Error reading file:" + err.Error())
		return
	}
	err = json.Unmarshal(fileContent, &%vListTable)
	if err != nil {
		panic("%v Error unmarshalling JSON:"+ err.Error())
		return
	}
	for _, v := range %vListTable {
		%vListMap.Store(v.Id, v)
	}
}
func Get%v(id float64) *%v {
	if v, ok := %vListMap.Load(id); ok {
		return v.(*%v)
	}
    return nil
}
`, name, cfg.Name, name, cfg.Name, name, name, name, name, name, name, cfg.Name, cfg.Name, name, cfg.Name)
		fmt.Fprintln(file, str)
	}

	fmt.Println("✅ Golang 结构体生成完成！")
}

// 读取 JSON 并初始化配置
func LoadConfig(dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".json") {
			fmt.Println("📦 加载 JSON:", path)
			LoadJson(path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("遍历 JSON 失败: %v", err)
	}
}

// 读取 JSON 并存入全局变量
func LoadJson(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("无法读取 JSON: %v", err)
	}
	defer file.Close()

	var data interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatalf("解析 JSON 失败: %v", err)
	}

	name := strings.TrimSuffix(filepath.Base(filePath), ".json")
	configStructs[name] = data
	fmt.Printf("✅ %s 配置加载完成\n", name)
}

// PascalCase 转换
func toPascalCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
