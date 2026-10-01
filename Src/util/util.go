package util

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/uuid"
)

func UUID() string {
	return uuid.New().String()
}

func ReadFile(file io.Reader) []byte {
	b, e := io.ReadAll(file)
	if e != nil {
		return nil
	}
	return b
}

func ReadAnyFile(path string) []byte {
	b, e := os.ReadFile(path)
	if e != nil {
		return nil
	}
	return b
}
func ReadAnyText(path string) string {
	b, e := os.ReadFile(path)
	if e != nil {
		return ""
	}
	return string(b)
}

type byteReader struct {
	b []byte
	i int
}

func (r *byteReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}
func bytesReader(b []byte) io.Reader { return &byteReader{b: b} }
func FileExist(path string) bool     { _, e := os.Stat(path); return e == nil }
func GetFolerRandomFileIfExist(path string) string {
	entries, e := os.ReadDir(path)
	if e != nil || len(entries) == 0 {
		return ""
	}
	return filepath.Join(path, entries[0].Name())
}
func GetRandomStr() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%s_f%06d", time.Now().Format("2006_01_02_15_04_05.000000000"), n.Int64())
}

// 返回拼接好的绝对路径
func JoinPath(path ...string) string {
	ret := filepath.Join(path...)
	abs, err := filepath.Abs(ret)
	if err != nil {
		fmt.Printf("JoinPath error:%v %v\n", err, ret)
		return ""
	}
	return filepath.ToSlash(abs)
}

func FormatFile(filepath string, data map[string]interface{}) (string, error) {
	content := ReadAnyText(filepath)

	tpl, err := template.New(filepath).Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func JsonToMap(s string) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// OpenFiles 打开文件
func OpenFiles(files []string, flag int, mode os.FileMode) ([]*os.File, error) {
	var ret []*os.File
	var err error
	var f *os.File
	for _, file := range files {
		file, _ = filepath.Abs(file)
		if f, err = os.OpenFile(file, flag, mode); err != nil {
			//关闭已打开的文件
			for _, f := range ret {
				f.Close()
			}
			return nil, err
		}
		ret = append(ret, f)
	}
	return ret, nil
}

// CloseFiles 关闭文件
func CloseFiles(files []*os.File) {
	for _, f := range files {
		f.Close()
	}
}

// 从当前位置往前读文件直到遇到新的一行或者到头
func ReadLineBack(file *os.File) (string, error) {
	pos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}

	if pos == 0 {
		return "", io.EOF
	}

	var result []byte
	var one [1]byte

	// 如果当前位置在换行符后，先跳过换行符
	for pos > 0 {
		_, err = file.ReadAt(one[:], pos-1)
		if err != nil {
			return "", err
		}

		if one[0] == '\n' {
			pos--
			break
		}

		result = append(result, one[0])
		pos--
	}

	// 反转字节顺序
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	// 文件指针移动到当前行的开头
	if _, err := file.Seek(pos, io.SeekStart); err != nil {
		return "", err
	}

	return string(result), nil
}

// 序列化map为字符串
func JsonMap(m map[string]interface{}) string {
	b, _ := json.Marshal(m)
	return string(b)
}

// 读取最大指定字节的文件数据
func ReadFileLimitFromCurrentOffset(f *os.File, limit int64) ([]byte, error) {
	curOff, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	endOff := stat.Size()
	if curOff >= endOff {
		// 当前已经在末尾，没有新内容
		return nil, nil
	}
	want := min(limit, endOff-curOff)
	buf := make([]byte, want)
	n, err := io.ReadFull(f, buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}
