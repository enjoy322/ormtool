package base

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// CamelCase 修改为大写开头的驼峰格式
func CamelCase(str string) string {
	var text string
	for _, p := range strings.Split(str, "_") {
		// 字段首字母大写的同时, 是否要把其他字母转换为小写
		switch len(p) {
		case 0:
		case 1:
			text += strings.ToUpper(p[0:1])
		default:
			text += strings.ToUpper(p[0:1]) + p[1:]
		}
	}
	return text
}

// Next 换行
func Next(depth int) string {
	return strings.Repeat("\n", depth)
}

// Tab 缩进
func Tab(depth int) string {
	return strings.Repeat("\t", depth)
}

// DealFilePath 处理保存路径，包名和文件名
func DealFilePath(s string, db string) (packageName, fileDir, fileName string) {
	if !strings.HasSuffix(s, ".go") {
		fmt.Println("保存路径错误，正确如./models/xx.go")
		os.Exit(0)
	}
	if len(strings.Trim(s, " ")) < 1 {
		packageName = "models"
		fileDir = "models"
		fileName = db
		return
	}
	split := strings.Split(s, "/")
	if len(split) <= 1 {
		packageName = "models"
		fileDir = "models"
		fileName = s
	} else {
		packageName = split[len(split)-2]
		fileName = split[len(split)-1]
		s2 := strings.Split(s, "/"+fileName)
		fileDir = s2[0]
	}
	return
}

//map排序
func sortMap(m map[string]string) []map[string]string {
	data := make([]map[string]string, 0)
	var ks []string
	for k, _ := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		m2 := make(map[string]string)
		m2[k] = m[k]
		data = append(data, m2)
	}
	return data
}

// Write 结构体信息写入go文件
func Write(packageName, fileDir, fileName string, content map[string]string, oneFile bool) {
	data := sortMap(content)
	err := os.MkdirAll(fileDir, 0777)
	if err != nil {
		panic(err)
	}
	if oneFile {
		fileName = fileDir + "/" + fileName
		var s strings.Builder
		s.WriteString("package " + packageName + Next(1))
		for _, datum := range data {
			for _, v := range datum {
				s.WriteString(v)
			}
		}
		writeToFile(fileName, s.String())
	} else {
		for k, v := range content {
			fileName = fileDir + "/" + k + ".go"
			var s strings.Builder
			s.WriteString("package " + packageName + Next(1))
			s.WriteString(v)
			writeToFile(fileName, s.String())
		}
	}
}

//执行写入和格式化
func writeToFile(fileName, content string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	_, err = f.WriteString(content)
	if err != nil {
		panic(err)
	}
	goFormat(fileName)
}

//格式化
func goFormat(fileName string) {
	cmd := exec.Command("gofmt", "-w", fileName)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
