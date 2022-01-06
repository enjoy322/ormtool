### 数据库表转换为golang 结构体

#### 1. 获取方式

```shell
```

#### 2. 配置说明

```shell
# 保存路径
SavePath:       "./models/test.go",
# 是否生成web json tag
IsGenJsonTag:   true,
# 是否生成在同一文件
IsGenInOneFile: true,
# 生成字段信息
GenDBInfoType:  base.CodeDBInfoGorm,
# 生成web json tag 类型，配置项IsGenJsonTag需为true
JsonTagType:    base.CodeJsonTag1
```

---
*欢迎提交pr和issue*
---

参考
> https://github.com/gohouse/converter
