package ormtool

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func UpperCamel(s string) string {
	s = strings.TrimSpace(s)

	if len(s) <= 1 {
		return strings.ToUpper(s)
	}

	split := strings.Split(s, "_")
	var tName strings.Builder
	for _, str := range split {
		tName.WriteString(strings.ToUpper(str[:1]) + str[1:])
	}
	s = tName.String()
	return s
}

// DealFilePath back save path and package name
func DealFilePath(s string, db string) (packageName, fileDir, fileName string) {
	if !strings.HasSuffix(s, ".go") {
		log.Fatalln("path error! correct example: ./models/xx.go")
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

// Write struct information to .go file
func Write(f FileInfo, data []StructInfo, oneFile bool) {
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

//format go file
func goFormat(fileName string) {
	cmd := exec.Command("gofmt", "-w", fileName)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// JsonTag deal tag by jsonType
func JsonTag(jsonType int, origin string) string {
	switch jsonType {
	//1.UserName 2.userName 3.user_name 4.user-name
	case 1:
		return UpperCamel(origin)
	case 2:
		s1 := UpperCamel(origin)
		return strings.ToLower(s1[:1]) + s1[1:]
	case 3:
		return strings.ToLower(origin)
	case 4:
		return strings.Replace(origin, "_", "-", -1)
	default:
		// 3.user_name
		return strings.ToLower(origin)
	}
}

// todo
// imports
