### 数据库表转换为golang 结构体

#### 1. 获取方式

```shell
go get github.com/enjoy322/ormtool@v1.1.0
```

#### 2. 配置说明

说明：数据库字段需为小写下划线

```shell
# 保存路径
SavePath:       "./models/model.go",
# 是否生成web json tag
IsGenJsonTag:   true,
# 是否生成在同一文件
IsGenInOneFile: true,
# 1：不生成数据库基本信息 2：生成简单的数据库字段信息
GenDBInfoType: 2,
# json tag类型，前提：IsGenJsonTag:true. 1.UserName 2.userName 3.user_name 4.user-name
JsonTagType:    1,
# 是否生成建表语句
IsGenCreateSQL: true
# 自定义对应类型，优先选择;map类型,key对应表字段类型,value对应Golang生成类型
CustomType: 
```
> 自定义类型，如下
```go
map[string]string{ "int": "int", 
	"int unsigned": "uint32",
	"tinyint(1)":   "bool",}
```

---
*欢迎提交pr和issue*
---

参考
> https://github.com/gohouse/converter