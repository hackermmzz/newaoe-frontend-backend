package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"newaoe/Src/config"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
	colorGray  = "\033[90m"
)

func printColor(color string, args ...interface{}) {
	fmt.Print(color)
	fmt.Println(args...)
	fmt.Print(colorReset)
}

func Debug(args ...interface{}) {
	printColor(colorGray, args...)
}

func DebugSuccess(args ...interface{}) {
	printColor(colorGreen, args...)
}

func DebugError(args ...interface{}) {
	printColor(colorRed, args...)
}

// SHA256
func SHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// 加密密码
func EncodePassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func CheckPasswordSame(encode_password string, orogin_password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encode_password), []byte(orogin_password))
	if err != nil {
		return false
	}
	return true
}

// 产生验证码
func GenerateVerifyCode(maxLen int) string {
	ret := ""
	for i := 0; i < maxLen; i += 1 {
		s := strconv.Itoa(rand.Intn(10))
		ret = ret + s
	}
	return ret
}

// 获取UTC时间
func UTC_Time() time.Time {
	return time.Now().UTC()
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 生成cookie的tooken(jwt格式)
func GenerateCookieToken(ip string, expire_time int, secretKey []byte, info map[string]interface{}) string {
	//
	expireTime := UTC_Time().Add(time.Duration(expire_time) * time.Second)
	data := jwt.MapClaims{}
	for key, value := range info {
		data[key] = value
	}
	//添加一些额外数据
	data["wlh_to_you"] = "为什么不玩原神?!"
	data["random"] = RandomString(32)
	data["ip"] = ip
	data["expire_time"] = expireTime
	data["login_time"] = UTC_Time()
	//
	encodedStd := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, err := encodedStd.SignedString(secretKey)
	//
	if err != nil {
		DebugError("GenerateCookieToken:", err)
		tokenString = ""
	}
	//
	return tokenString
}

// 判断token是否合法
func CheckTokenLegal(expire_time string) bool {
	expireTime, err := time.Parse(time.RFC3339, expire_time)
	if err != nil {
		DebugError("CheckTokenLegal:解析过期时间失败:", err)
		return false
	}
	if UTC_Time().After(expireTime) {
		return false
	}
	return true
}

// 获取tooken解析数据
func GetTokenInfo(signed_token string) map[string]interface{} {
	token, err := jwt.ParseWithClaims(signed_token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.Conf.Server.JwtSecretKey, nil
	})
	//
	if err != nil {
		DebugError("GetTokenInfo:", err)
		return nil
	}
	//
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := make(map[string]interface{})
		for key, value := range claims {
			data[key] = value
		}
		return data
	}
	//
	DebugError("GetTokenInfo:解析失败:" + signed_token)
	return nil
}

// 获取tooken解析数据
func GetCtxTookenInfo(ctx *gin.Context) map[string]interface{} {
	//注意这里我并没有采用sessionID的方式，因为我觉得没必要
	//解析token
	token, err := ctx.Cookie(config.Conf.Cookie.TokenName)
	if err != nil || token == "" {
		return nil
	}
	//解码token
	data := GetTokenInfo(token)
	if data == nil {
		DebugError("StudentInfoGet:", "为什么过得了我的过滤器却解析不了,出BUG了?")
		return nil
	}
	return data
}

// 创建文件夹
func Mkdir(path string) bool {
	//看看存不存在
	fileInfo, err := os.Stat(path)
	if err == nil {
		//存在
		if fileInfo.IsDir() { //是目录
			return true
		}
		return false //不是目录
	}
	//不存在就创建目录
	err = os.Mkdir(path, 0755)
	if err != nil {
		DebugError("Mkdir:", err)
		return false
	}
	return true
}

// 判断targetDir是否是baseDir的子目录
func IsSubDir(baseDir, targetDir string) bool {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return false
	}
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return false
	}
	// 去掉. ..和~等相对路径符号
	cleanBase := filepath.Clean(absBase)
	cleanTarget := filepath.Clean(absTarget)

	// 判断目标目录是否以基准目录为前缀（且是真正的子目录）
	baseWithSep := cleanBase + string(filepath.Separator)
	isSubDir := strings.HasPrefix(cleanTarget, baseWithSep) // 目标目录是基准目录的子目录

	return isSubDir
}

// 判断文件是否存在
func IsFileExist(path string) bool {
	// 获取文件信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	// 检查是否是目录（防止下载目录）
	if fileInfo.IsDir() {
		return false
	}
	return true
}

// 获取随机头像
func GetRandomAvatar() string {
	cnt := len(config.Conf.User.UserDefalutAvatar)
	x := rand.Intn(cnt)
	return config.Conf.User.UserDefalutAvatar[x]
}

// 读取文件返回文本内容
func ReadFileContent(path string) string {
	content, err := os.ReadFile(path) // 返回 []byte 和 error
	if err != nil {
		DebugError("读取文件失败:", err, path)
		return ""
	}
	return string(content)
}

// 写入内容到文件
func WriteFileContent(path string, data string) bool {
	err := os.WriteFile(path, []byte(data), os.FileMode(0644))
	if err != nil {
		DebugError("WriteFileContent:", err)
		return false
	}
	return true
}

// 获取[]byte json格式数据
func JsonBytes(data []byte, addr interface{}) bool {
	err := json.Unmarshal(data, addr)
	if err != nil {
		DebugError("JsonBytes:", err)
		return false
	}
	return true
}

// 获取strinde json格式的数据
func JsonString(data string, addr interface{}) bool {
	return JsonBytes([]byte(data), addr)
}

// 获取json格式的数据
func JsonCtx(ctx *gin.Context, addr interface{}) bool {
	data_bytes, err := ctx.GetRawData()
	if err != nil {
		DebugError("JsonCtx", err)
		return false
	}
	return JsonBytes(data_bytes, addr)
}

// 获取url里面的文件路径
func GetUrlFilePath(url_ string) string {
	u, err := url.Parse(url_)
	if err != nil {
		DebugError("GetUrlFilePath:", url_, err.Error())
		return ""
	}
	return strings.TrimPrefix(u.Path, "/")
}

// 拆分bucket和文件路径
func SplitBucketAndFilePath(bucket string, path string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] != bucket {
		DebugError("SplitBucketAndFilePath:", bucket, path, parts)
		return ""
	}
	return parts[1]
}

// 获取字符串中第一个 { 到对应 } 之间的内容(包括{})
func GetBracesContent(s string) string {
	stack := make([]int, 0)
	start := -1
	for i, ch := range s {
		if ch == '{' {
			if len(stack) == 0 {
				start = i
			}
			stack = append(stack, i)
		} else if ch == '}' {
			if len(stack) == 0 {
				continue
			}
			// 出栈
			stack = stack[:len(stack)-1]
			// 栈空，说明最外层 {} 已经完整配对
			if len(stack) == 0 && start != -1 {
				return s[start : i+1]
			}
		}
	}
	return ""
}

func GetLeftTimeForOneDay() time.Duration {
	now := time.Now()
	// 明天 00:00:00
	tomorrow := time.Date(
		now.Year(),
		now.Month(),
		now.Day()+1,
		0, 0, 0, 0,
		now.Location(),
	)
	return tomorrow.Sub(now)
}

func TruncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}

// 获取文件所在的文件夹
// GetFileFolder 获取文件所在目录
func GetFileFolder(file string) string {
	return filepath.Dir(file)
}

// 生成一个指定值的数组
func NewArray[T any](n int, value T) []T {
	arr := make([]T, n)
	for i := range arr {
		arr[i] = value
	}
	return arr
}

// 获取一个结构体的所有json key
func GetJsonKeys(obj interface{}) []string {
	t := reflect.TypeOf(obj)
	// 如果传的是指针
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	keys := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		keys = append(keys, t.Field(i).Name)
	}
	return keys
}

// 获取当前时间得纳秒级字符串
func GetTimeNanoStr() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// 判断target是否为base前缀
func IsPrefix(s string, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func NewError(args ...interface{}) error {
	final_err := ""
	for i, v := range args {
		if i == 0 {
			final_err = fmt.Sprintf("%v", v)
		} else {
			final_err = fmt.Sprintf("%s %v", final_err, v)
		}
	}
	return errors.New(final_err)
}
