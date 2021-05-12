package structMaker

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

const default_output_path = "./models/"
const default_output_single = false
const default_output_filename = "models.go"

const (
	EgMySql = iota
	EgSqlServer
)

type table struct {
	name    string
	columns []column
}

type column struct {
	name      string
	data_type string
}

type output struct {
	outputPath       string
	outputSingleFile bool //全部输出在一个文件
}

type DBLoader interface {
	setDB(db *sql.DB, dbName string)
	getTables() []table
	getColumnType(string) string
}

type Maker interface {
	makeStruct(table) string
	isOutputSingleFile() bool
	SetOutput(outputPath string, outputSingleFile bool) *defaultMaker
	SetSourceDB(DBLoader) *defaultMaker
	MakeFile() error
}

type defaultMaker struct {
	source DBLoader
	output
}

func (this *defaultMaker) SetSourceDB(loader DBLoader) *defaultMaker {
	this.source = loader
	return this
}

func (this *defaultMaker) MakeFile() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error: ", err)
		}
	}()
	if !CheckFileIsExist(this.outputPath) {
		err := os.MkdirAll(this.outputPath, os.ModePerm)
		if err != nil {
			panic("MkdirAll Error:" + err.Error())
		}
	}
	tables := this.source.getTables()
	if this.isOutputSingleFile() {
		//输出到一个文件
		structStr := "package model\n"
		for _, _table := range tables {
			structStr += this.makeStruct(_table)
		}
		file, err := OpenFile(this.outputPath +"/"+ default_output_filename)
		defer file.Close()
		if err != nil {
			panic("OpenFile Error:" + err.Error())
		}
		_, err = io.WriteString(file, structStr)
		if err != nil {
			panic("WriteString Error:" + err.Error())
		}

	} else {
		wait := sync.WaitGroup{}
		for _, _table := range tables {
			wait.Add(1)
			go func(tb table) {
				var file *os.File
				var err error
				structStr := "package model\n"
				structStr += this.makeStruct(tb)
				filename := this.outputPath + "/" + tb.name + ".go"
				file, err = OpenFile(filename)
				defer file.Close()
				if err != nil {
					panic("OpenFile Error:" + err.Error())
				}
				_, err = io.WriteString(file, structStr)
				if err != nil {
					panic("WriteString Error:" + err.Error())
				}
				wait.Done()
			}(_table)
		}
		wait.Wait()
	}
	return nil
}

func (this *defaultMaker) makeStruct(_table table) string {
	var type_str, structContent string
	structContent += "\ntype " + HumpFormat(_table.name) + " struct {"
	for _, col := range _table.columns {
		type_str = this.source.getColumnType(col.data_type)
		structContent += "\n" + "	" + HumpFormat(col.name) + " " + type_str + " `json:\"" + col.name + "\"`"
	}
	structContent += "\n}"
	return structContent
}

func (this *defaultMaker) SetOutput(outputPath string, outputSingleFile bool) *defaultMaker {
	this.outputPath = outputPath
	this.outputSingleFile = outputSingleFile
	return this
}

func (this *defaultMaker) isOutputSingleFile() bool {
	return this.outputSingleFile
}

func (this *defaultMaker) Run() error {
	if this.source == nil {
		return errors.New("缺少DB配置")
	}
	this.MakeFile()
	return nil
}

func New() *defaultMaker {
	return &defaultMaker{
		source: nil,
		output: output{
			outputPath:       default_output_path,
			outputSingleFile: false,
		},
	}
}

func Run(db *sql.DB, database string, engine int) error {
	maker := New()
	var loader DBLoader
	switch engine {
	case EgMySql:
		loader = NewMysql(db, database)
	case EgSqlServer:
		return errors.New("暂不支持 sqlserver")
	}
	maker.SetSourceDB(loader)
	return maker.MakeFile()
}
