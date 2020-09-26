package test

import (
	"bufio"
	"golang_test/util"
	"io"
	"io/ioutil"
	"os"
	"syscall"
	"testing"
)

func TestFile(t *testing.T) {
	file, err := os.Create("tmp")
	if err != nil {
		t.Log(err)
	}
	if err = file.Truncate(10); err == nil {
		fInfo, _ := file.Stat()
		t.Log(fInfo.Size(), fInfo.Mode(), fInfo.Name(), fInfo.IsDir())
	}
	_ = file.Close() // rename tmp ../200B: The process cannot access the file because it is being used by another process.
	if err = os.Rename("tmp", "../10B"); err != nil {
		t.Log(err)
	}
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified. 只有O_APPEND在linux中会写不进去
	if file, err = os.OpenFile("../10B", os.O_RDWR|os.O_APPEND, os.ModeAppend); err == nil {
		_ = ioutil.WriteFile("../10B", []byte("abcbabcbabcbabcbabc"), os.ModePerm)
		_ = file.Sync()
		bytes, err := ioutil.ReadFile("../10B")
		t.Log(bytes, err)
		_ = file.Close()
	} else {
		t.Log(err)
	}
	if err = os.Remove("../10B"); err != nil {
		t.Log(err)
	}
}

func TestFileCheck(t *testing.T) {
	file, err := os.Open("../10B")
	t.Log(file, err, os.IsNotExist(err))
	t.Log(os.IsPermission(err))
	file, _ = os.OpenFile("tmp", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	_ = file.Chmod(0666)
	tfile, _ := os.OpenFile("tmp2", os.O_CREATE|os.O_RDWR, 0755)
	_, _ = io.Copy(tfile, file)
	_ = file.Close()
	_ = tfile.Close()
	_ = os.Remove("tmp")
	_ = os.Remove("tmp2")
}

func TestFileRW(t *testing.T) {
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified. 只有O_APPEND在linux中会写不进去
	file, _ := os.OpenFile("../tmp", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModeAppend)
	_, _ = file.Write([]byte("abcdefghijk"))
	//bytes := []byte{}//无法读取。。。
	bytes := make([]byte, 5, 200)
	_, _ = file.Read(bytes)
	t.Log(bytes) // 无法读取内容，需要seek
	_, _ = file.Seek(0, syscall.FILE_BEGIN)
	_, _ = file.Read(bytes)
	t.Log(bytes)
	_, _ = file.Read(bytes)
	t.Log(bytes)
	_, _ = file.Seek(-2, syscall.FILE_END)
	_, _ = file.Read(bytes)
	t.Log(bytes) // 前2字节写入了新内容，后边还是原数据
	_ = file.Close()
	_ = os.Remove("../tmp")
	t.Log(util.IoReadFile("tmp"))
}

func TestBuffW(t *testing.T) {
	file, _ := os.OpenFile("../tmp", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModeAppend)
	buffer := bufio.NewWriter(file) //4096
	_ = buffer.WriteByte(20)
	_, _ = buffer.WriteRune('b')
	t.Log(buffer.Buffered(), buffer.Available())
	_ = buffer.Flush()
	bufio.NewWriterSize(buffer, 1024) // 只增不减
	t.Log(buffer.Available())
}
