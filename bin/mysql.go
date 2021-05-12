package structMaker

import (
	"database/sql"
	"strings"
)

func NewMysql(db *sql.DB, db_name string) *mysql {
	m := &mysql{}
	m.setDB(db, db_name)
	return m
}

type mysql struct {
	db_conn *sql.DB
	db_name string
	tables  []string
	columns map[string][]column
}

//设置连接数据库
func (this *mysql) setDB(db *sql.DB, db_name string) {
	err := db.Ping()
	if err != nil {
		panic("数据库连接无效：" + err.Error())
	}
	this.db_conn = db
	this.db_name = db_name
}

//将db数据类型转化为go数据类型
func (this *mysql) getColumnType(data_type string) string {
	var res string
	switch strings.ToUpper(data_type) {
	case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "BIGINT", "INTEGER":
		res = "int64"
	case "FLOAT", "DOUBLE", "DECIMAL":
		res = "float64"
	default:
		res = "string"
	}
	return res
}

//获取数据库table及columns
func (this *mysql) getTables() []table {
	sql_str := "SELECT `COLUMN_NAME`,`DATA_TYPE`,`TABLE_NAME` FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? ORDER BY `TABLE_NAME` ASC,`ORDINAL_POSITION` ASC"
	rows, err := this.db_conn.Query(sql_str, this.db_name)
	if err != nil {
		panic("getTableColumns Query error:" + err.Error())
	}
	var columns = make(map[string][]column)
	for rows.Next() {
		var column_name, data_type, tableName string
		err := rows.Scan(&column_name, &data_type, &tableName)
		if err != nil {
			panic("GetAllTable Scan error:" + err.Error())
		}
		item := column{
			column_name,
			data_type,
		}
		//if len(columns[tableName]) < 1 {
		//	tables = append(tables, tableName)
		//}
		columns[tableName] = append(columns[tableName], item)
	}

	var tables []table
	var tmp table
	for name, col := range columns {
		tmp = table{
			name,
			col,
		}
		tables = append(tables, tmp)
	}
	return tables
}
