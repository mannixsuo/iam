package model

import (
	"database/sql"
	"strconv"
	"strings"
)

var mysql *sql.DB

// Init 初始化数据库连接信息
func Init() {
	db, e := sql.Open("mysql", "root:969929899@tcp(localhost:3306)/auth?loc=Local")
	if e != nil {
		panic(e)
	}
	mysql = db
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

// BuildNumberList 参数[1,2,3] 返回(1,2,3)
func BuildNumberList(idList []int) string {
	var builder strings.Builder
	builder.Write([]byte("("))
	len := len(idList)
	if len == 0 {
		return "(0)"
	}
	for i := 0; i < len-1; i++ {
		builder.Write([]byte(strconv.Itoa(idList[i])))
		builder.Write([]byte(","))
	}
	builder.Write([]byte(strconv.Itoa(idList[len-1])))
	builder.Write([]byte(")"))
	return builder.String()
}

// BuildStringList change action list [a,b,c] to (a,b,c)
func BuildStringList(stringList []string) string {
	sLen := len(stringList)
	if sLen == 0 {
		return "(\"\")"
	}
	builder := strings.Builder{}
	builder.WriteString("(")
	for i := 0; i < sLen-1; i++ {
		builder.WriteString(stringList[i])
		builder.WriteString(",")
	}
	builder.WriteString(stringList[sLen-1])
	builder.WriteString(")")
	return builder.String()
}
