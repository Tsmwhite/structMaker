# structMaker
golang 自动生成数据库struct 用与替换 make_table_struct

```
go get github.com/Tsmwhite/structMaker
```

```GO
func main() {
	dbConfigString := "root:********@tcp(127.0.0.1:3306)/example?charset=utf8"
	db, err := sql.Open("mysql", dbConfigString)
	if err == nil {
		// @1
		//err = structMaker.Run(db, "example", structMaker.EgMySql)

		// @2
		loader := structMaker.NewMysql(db, "example")
		err = structMaker.New().SetSourceDB(loader).MakeFile()

		// @3
		err = structMaker.New().SetSourceDB(loader).SetOutput("models2", true).MakeFile()
		fmt.Println(err)
	}
}
```
