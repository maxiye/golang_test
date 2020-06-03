package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	// -h 172.17.12.152 -u root -p im2NCnCwweA= -d xqdata_2345
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
	dbList := strings.Split(database, ",")
	for _, dbName := range dbList {
		fmt.Println("--------开始操作数据库：" + dbName + "--------")
		dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", user, pass, host, port, dbName, charset)
		db, err := sql.Open("mysql", dbDsn)
		if err != nil {
			fmt.Println("数据源配置不正确: " + err.Error())
			continue
		}
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(22)
		if err = db.Ping(); err != nil {
			fmt.Println("连接建立失败：" + err.Error())
			continue
		}
		fmt.Println("数据库连接成功...")
		// 统计数据库，表大小
		statSize(db, dbName)
		// 统计数据库字段
		statFields(db, dbName)
		_ = db.Close()
		fmt.Println("--------结束操作数据库：" + dbName + "--------")
	}
}

func checkFile(rewrite bool) {
	fileList := map[string]string{
		"database.csv": "数据库名称,表的数量,字段数量,占磁盘容量（TB）\r\n",
		"tables.csv":   "数据库名称,表名,表中文名,字段数量,占磁盘容量（GB）\r\n",
		"fields.csv":   "数据库名称,表名,字段名,数据类型,字段说明\r\n",
	}
	for fName, title := range fileList {
		if rewrite {
			if err := os.Remove(fName); err != nil {
				fmt.Println(fName + "删除失败：" + err.Error())
			}
		}
		_, err := os.Stat(fName)
		if err != nil && os.IsNotExist(err) {
			if file, err := os.Create(fName); err != nil {
				panic("文件创建失败" + err.Error())
			} else {
				if _, err := file.Write([]byte(title)); err != nil {
					panic("文件写入失败" + err.Error())
				}
			}
		}
	}
}

func handleStr(str string) string {
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, "\n", "")
	return str
}

func statFields(db *sql.DB, database string) {
	fmt.Println("开始统计字段信息")
	offset, limit, goSize := 0, 1000, 10
	var end bool
	var wgWrite sync.WaitGroup
	wgWrite.Add(1)
	resChan := writeDataToFile(&wgWrite, "fields.csv", func(item *map[string]string) string {
		info := *item
		return fmt.Sprintf("%s,%s,%s,%s,%s\r\n", info["dbName"], info["table"], info["column"], info["columnType"], info["columnComment"])
	})
	for {
		var wg sync.WaitGroup
		for i := 0; i < goSize; i++ {
			wg.Add(1)
			go func(offset int) {
				defer wg.Done()
				columnRes, err := db.Query("select table_name,COLUMN_NAME,COLUMN_TYPE,COLUMN_COMMENT from information_schema.COLUMNS where table_schema = ? order by table_name limit ? offset ?", database, limit, offset)
				if err != nil {
					fmt.Println("字段查询失败：" + err.Error())
					return
				}
				resCount := 0
				tmpList := make([]map[string]string, 0, limit)
				for columnRes.Next() {
					resCount++
					var table, column, columnType, columnComment string
					if err = columnRes.Scan(&table, &column, &columnType, &columnComment); err != nil {
						fmt.Println("获取数据失败", err.Error())
						continue
					}
					fields := map[string]string{
						"dbName":        database,
						"table":         handleStr(table),
						"column":        handleStr(column),
						"columnType":    handleStr(columnType),
						"columnComment": handleStr(columnComment), //byte
					}
					tmpList = append(tmpList, fields)
				}
				resChan <- &tmpList
				_ = columnRes.Close()
				if resCount < limit {
					end = true
				}
			}(offset)
			offset += limit
		}
		wg.Wait()
		if end {
			break
		}
	}
	close(resChan)
	wgWrite.Wait()
	fmt.Println("结束统计字段信息")
}

func writeDataToFile(wgOutside *sync.WaitGroup, fileName string, getFromMap func(*map[string]string) string) chan *[]map[string]string {
	out := make(chan *[]map[string]string)
	go func() {
		defer wgOutside.Done() // 否则未写完就退出了
		if dbFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend); err == nil {
			rowCount := 0
			for item := range out {
				rowCount += len(*item)
				fmt.Println("准备写入数据：", len(*item))
				for _, info := range *item {
					dbLine := getFromMap(&info)
					if _, err := dbFile.WriteString(dbLine); err != nil {
						fmt.Println(fileName + " 写入失败" + err.Error())
					}
				}
			}
			fmt.Println(fileName+"文件保存完毕，行数：", rowCount)
			dbFile.Close()
		} else {
			panic(fileName + " 打开失败：" + err.Error())
		}
	}()
	return out
}

func statSize(db *sql.DB, database string) {
	var dbTables, dbColNums int
	var dbLength float64 = 0 //mb
	offset, limit, goSize := 0, 400, 20
	var end bool
	var wgWrite sync.WaitGroup
	wgWrite.Add(1)
	writeChan := writeDataToFile(&wgWrite, "tables.csv", func(item *map[string]string) string {
		infoMap := *item
		tableLen, _ := strconv.ParseFloat(infoMap["length"], 64)
		// 统计数据库数据
		dbTables += 1
		columnCount, _ := strconv.Atoi(infoMap["columnCount"])
		dbColNums += columnCount
		dbLength += tableLen / 1024 / 1024
		return fmt.Sprintf("%s,%s,%s,%s,%.8f\r\n", database, infoMap["table"], infoMap["comment"], infoMap["columnCount"], tableLen/1024/1024/1024)
	})
	fmt.Println("开始统计表体量字段数等数据...")
	tableFieldCount := getTableFieldCount(db, database)
	for {
		var wg sync.WaitGroup
		for i := 0; i < goSize; i++ {
			wg.Add(1)
			go func(offset int) {
				defer wg.Done()
				tableStat := make([]map[string]string, 0, limit)
				rows, err := db.Query("SELECT `TABLE_SCHEMA`,`TABLE_NAME`,IFNULL(`DATA_LENGTH`,0),IFNULL(`INDEX_LENGTH`,0),`TABLE_COMMENT` FROM information_schema.`TABLES` WHERE `TABLE_SCHEMA` = ? limit ? offset ?", database, limit, offset)
				if err != nil {
					fmt.Println("表查询失败：" + err.Error())
					return
				}
				for rows.Next() {
					var dbName, table, comment string
					var dataLen, idxLen uint64
					if err = rows.Scan(&dbName, &table, &dataLen, &idxLen, &comment); err != nil {
						fmt.Println("获取数据失败", err.Error())
						continue
					}
					columnCount := tableFieldCount[table]
					tableInfo := map[string]string{
						"dbName":      handleStr(dbName),
						"table":       handleStr(table),
						"comment":     handleStr(comment),
						"columnCount": handleStr(strconv.Itoa(columnCount)),
						"length":      handleStr(strconv.FormatUint(dataLen+idxLen, 10)), //byte
					}
					tableStat = append(tableStat, tableInfo)
				}
				resCount := len(tableStat)
				writeChan <- &tableStat
				if err = rows.Close(); err != nil {
					fmt.Println("数据库连接关闭失败：" + err.Error())
				}
				if resCount < limit {
					end = true
				}
			}(offset)
			offset += limit
		}
		wg.Wait()
		if end {
			close(writeChan)
			break
		}
	}
	wgWrite.Wait()
	// 存储结果
	fmt.Println("汇总结果写入文件：")
	if dbFile, err := os.OpenFile("database.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend); err == nil {
		dbLine := fmt.Sprintf("%s,%d,%d,%.8f\r\n", database, dbTables, dbColNums, dbLength/1024/1024)
		fmt.Println("数据库汇总结果：数据库名称、表的数量、字段数量、占磁盘容量（TB）\r\n" + dbLine)
		if _, err := dbFile.Write([]byte(dbLine)); err != nil {
			fmt.Println("database写入失败" + err.Error())
		}
		dbFile.Close()
	} else {
		fmt.Println("database.csv打开失败：" + err.Error())
	}
	fmt.Println("表体量字段数写入完成")
}

func getTableFieldCount(db *sql.DB, dbName string) map[string]int {
	tableFieldsMap := make(map[string]int, 6000)
	columnRes, err := db.Query("select table_name,count(*) as count from information_schema.COLUMNS where table_schema = ? group by table_name ", dbName)
	if err == nil {
		for columnRes.Next() {
			var (
				tableName string
				count     int
			)
			if err = columnRes.Scan(&tableName, &count); err != nil {
				fmt.Println("表 字段数查询失败：" + err.Error())
			} else {
				tableFieldsMap[tableName] = count
			}
		}
		if err = columnRes.Close(); err != nil {
			fmt.Println("columnRes数据库连接关闭失败：" + err.Error())
		}
	} else {
		fmt.Println("获取字段数量失败" + err.Error())
	}
	return tableFieldsMap
}
