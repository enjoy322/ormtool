### Convert database table to Golang struct

#### 1. Get

```shell
go get -u github.com/enjoy322/ormtool@main
```

#### 2. Configuration

Note: The database table column should be underlined in lowercase

```go
GenerateMySQL(
    Config{
        //[user]:[password]@tcp([host]:[port])/[database]?parseTime=true
        ConnStr: "root:qwe123@tcp(127.0.0.1:3306)/test?parseTime=true",
        // database name
        Database: "test",
        // relative path
        SavePath: "./model/model.go",
        // Generate one file or files
        IsGenInOneFile: true,
        // Generate simple database field information like: "int unsigned not null"
        // value 1:not generate; 2：simple info
        GenDBInfoType: 1,
        // json tag
        IsGenJsonTag: true,
        // json tag type. The necessary conditions：IsGenJsonTag:true.
        // 1.UserName 2.userName 3.user_name 4.user-name
        JsonTagType: 3,
        // sql of creating table in the database
        IsGenCreateSQL: true,
        // simple crud function
        IsGenFunction:true,
        // custom type relationships will be preferred. 
        // the key is the database type, The value is the golang type
        CustomType: map[string]string{
        "int":          "int",
        "int unsigned": "uint32",
        "tinyint(1)":   "bool",
        "json":         "json.RawMessage",
        },
	})


// result example
// User	用户表
type User struct {
    Id         int    `json:"id"`
    CreateTime int    `json:"create_time"` // 创建时间
    UserName   string `json:"user_name"`   // 用户名
}

func (*User) TableName() string {
    return "user"
}

var UserCol = struct {
    Id         string
    CreateTime string
    UserName   string
}{
    Id:         "id",
    CreateTime: "create_time",
    UserName:   "user_name",
}

// function

type UserModelInterface interface {
    Create(data *User) error
    Get(id int) (User, error)
    Find(condition interface{}, page, limit int) ([]User, error)
    Delete(id int) error
    DeleteUnScope(id int) error
}

type userModelService struct {
    db *gorm.DB
}

func NewUserModelService(db *gorm.DB) UserModelInterface {
    return userModelService{db: db}
}

func (s userModelService) Create(data *User) error {
    err := s.db.Create(data).Error
    if err != nil {
        return err
    }
    return nil
}
```
#### Reference
> https://github.com/gohouse/converter