package base

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func DealStructName(s string) string {
	if len(s) == 1 {
		s = strings.ToUpper(s[:1])
	} else {
		split := strings.Split(s, "_")
		var tName strings.Builder
		for _, str := range split {
			tName.WriteString(strings.ToUpper(str[:1]) + str[1:])
		}
		s = tName.String()
	}
	return s
}

func predeal(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	data := reg.ReplaceAllString(s, "")
	return data
}

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

func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
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
func Write(f FileInfo, data []StructInfo, oneFile bool) {
	// data := sortMap(content)
	err := os.MkdirAll(f.FileDir, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	if oneFile {
		f.FileName = f.FileDir + "/" + f.FileName
		var s strings.Builder
		s.WriteString("package " + f.PackageName + "\n")
		for _, v := range data {
			s.WriteString(v.Note)
			s.WriteString(v.CreateSQL)
			s.WriteString(v.StructContent)
			s.WriteString("\n\n")

		}
		writeToFile(f.FileName, s.String())
	} else {
		for _, v := range data {
			fileName := f.FileDir + "/" + v.Name + ".go"
			var s strings.Builder
			s.WriteString("package " + f.PackageName + "\n")
			s.WriteString(v.Note)
			s.WriteString(v.CreateSQL)
			s.WriteString(v.StructContent)
			s.WriteString("\n\n")
			writeToFile(fileName, s.String())
		}

	}
}

//执行写入和格式化
func writeToFile(fileName, content string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(f)
	_, err = f.WriteString(content)
	if err != nil {
		log.Fatalln(err)
	}
	goFormat(fileName)
}

// todo
// imports

//格式化
func goFormat(fileName string) {
	cmd := exec.Command("gofmt", "-w", fileName)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// JsonTag 处理tag： json
func JsonTag(jsonType int, origin string) string {
	switch jsonType {
	//1.UserName 2.userName 3.user_name 4.user-name
	case 1:
		return Case2Camel(origin)
	case 2:
		s1 := Case2Camel(origin)
		return strings.ToLower(s1[:1]) + s1[1:]
	case 3:
		return strings.ToLower(origin)
	case 4:
		return strings.Replace(origin, "_", "-", -1)

	}
	panic("json tag 参数错误")
}

// GetTypeNum 获取表字段长度约束
func GetTypeNum(typeStr string) int {
	f := strings.HasSuffix(typeStr, ")")
	if f {
		//	有长度约束
		splitAfter := strings.SplitAfter(typeStr, "(")
		n := splitAfter[1][:1]
		i, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		return i
	}
	return 0
}
