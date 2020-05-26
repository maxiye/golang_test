package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
)

func main() {
	// -h 172.17.12.152 -u root -p im2NCnCwweA= -d xq_2345
	var host, port, user, pass, database, charset string
	var rewrite bool
	flag.StringVar(&host, "h", "", "host，ip地址")
	flag.StringVar(&port, "t", "3306", "端口号")
	flag.StringVar(&user, "u", "", "用户名")
	flag.StringVar(&pass, "p", "", "密码")
	flag.StringVar(&database, "d", "", "数据库名")
	flag.StringVar(&charset, "c", "utf8", "字符集")
	flag.BoolVar(&rewrite, "r", false, "删除已有数据")
	flag.Parse()
	if host == "" || user == "" || pass == "" || database == "" {
		fmt.Println("参数输入错误，使用 --help 确认参数列表")
		flag.Usage()
		os.Exit(0)
	}
	checkFile(rewrite)
	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", user, pass, host, port, database, charset)
	db, err := sql.Open("mysql", dbDsn)
	if err != nil {
		panic("数据源配置不正确: " + err.Error())
	}
	if _, err = db.Exec("SELECT 1"); err != nil {
		panic("连接建立失败：" + err.Error())
	}
	fmt.Println("数据库连接成功...")
	tableStat := make([]map[string]string, 0, 6000)
	var dbTables, dbColNums int
	var dbLength float64 = 0 //mb
	offset := 0
	fmt.Println("开始查询数据...")
	for {
		rows, err := db.Query("SELECT `TABLE_SCHEMA`,`TABLE_NAME`,IFNULL(`DATA_LENGTH`,0),`TABLE_COMMENT` FROM information_schema.`TABLES` WHERE `TABLE_SCHEMA` = ? limit 50 offset ?", database, offset)
		if err != nil {
			fmt.Println("表查询失败：" + err.Error())
			continue
		}
		resCount := 0
		for rows.Next() {
			resCount++
			var dbName, table, comment string
			var length uint64
			if err = rows.Scan(&dbName, &table, &length, &comment); err != nil {
				fmt.Println("获取数据失败", err.Error())
			}
			columnRes, err := db.Query("select count(*) from information_schema.COLUMNS where table_name = ? and table_schema = ?", table, dbName)
			columnCount := 0
			if err == nil && columnRes.Next() {
				if err = columnRes.Scan(&columnCount); err != nil {
					fmt.Println(table + "表 字段数查询失败：" + err.Error())
				}
				if err = columnRes.Close(); err != nil {
					fmt.Println("columnRes数据库连接关闭失败：" + err.Error())
				}
			} else if err != nil {
				fmt.Println("获取字段数量失败" + err.Error())
			}
			dbColNums += columnCount
			dbLength += float64(length) / 1024 / 1024
			tableInfo := map[string]string{
				"dbName":      dbName,
				"table":       table,
				"comment":     comment,
				"columnCount": strconv.Itoa(columnCount),
				"length":      strconv.FormatUint(length, 10), //byte
			}
			tableStat = append(tableStat, tableInfo)
		}
		dbTables += resCount
		if err = rows.Close(); err != nil {
			fmt.Println("数据库连接关闭失败：" + err.Error())
		}
		if resCount != 50 {
			break
		}
		fmt.Print(".")
		offset += 50
	}
	fmt.Println("")
	// 存储结果
	fmt.Println("数据查询完毕，开始写入文件...")
	if dbFile, err := os.OpenFile("database.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend); err == nil {
		dbLine := fmt.Sprintf("%s,%d,%d,%.8f\r\n", database, len(tableStat), dbColNums, dbLength/1024/1024)
		fmt.Println("数据库汇总结果：数据库名称、表的数量、字段数量、占磁盘容量（TB）\r\n" + dbLine)
		if _, err := dbFile.Write([]byte(dbLine)); err != nil {
			fmt.Println("database写入失败" + err.Error())
		}
		dbFile.Close()
	} else {
		fmt.Println("database.csv打开失败：" + err.Error())
	}
	if tbFile, err := os.OpenFile("tables.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend); err == nil {
		for _, v := range tableStat {
			tableLen, _ := strconv.ParseFloat(v["length"], 10)
			dbLine := fmt.Sprintf("%s,%s,%s,%s,%.8f\r\n", database, v["table"], v["comment"], v["columnCount"], tableLen/1024/1024/1024)
			if _, err := tbFile.Write([]byte(dbLine)); err != nil {
				fmt.Println("table写入失败" + err.Error())
			}
		}
		tbFile.Close()
	} else {
		fmt.Println("tables.csv打开失败：" + err.Error())
	}
	fmt.Println("写入完成")
	//fmt.Println(tableStat)
}

func checkFile(rewrite bool) {
	if rewrite {
		if err := os.Remove("database.csv"); err != nil {
			fmt.Println("database.csv删除失败：" + err.Error())
		}
		if err := os.Remove("tables.csv"); err != nil {
			fmt.Println("tables.csv删除失败：" + err.Error())
		}
	}
	_, err := os.Stat("database.csv")
	if err != nil && os.IsNotExist(err) {
		if file, err := os.Create("database.csv"); err != nil {
			panic("文件创建失败" + err.Error())
		} else {
			if _, err := file.Write([]byte("数据库名称,表的数量,字段数量,占磁盘容量（TB）\r\n")); err != nil {
				panic("文件写入失败" + err.Error())
			}
		}
	}
	_, err = os.Stat("tables.csv")
	if err != nil && os.IsNotExist(err) {
		if file, err := os.Create("tables.csv"); err != nil {
			panic("文件创建失败" + err.Error())
		} else {
			if _, err := file.Write([]byte("数据库名称,表名,表中文名,字段数量,占磁盘容量（GB）\r\n")); err != nil {
				panic("文件写入失败" + err.Error())
			}
		}
	}
}
