package base

import "strings"

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

func Next(depth int) string {
	return strings.Repeat("\n", depth)
}
func Tab(depth int) string {
	return strings.Repeat("\t", depth)
}
