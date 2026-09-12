package util

import (
	"fmt"
	"math/rand"
	"net/url"
	"newaoe/config"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Debug(args ...interface{}) {
	fmt.Println(args...)
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
func GenerateCookieToken(ctx *gin.Context, id string, email string, regist_date time.Time, expire_time int, secretKey []byte) string {
	//
	expireTime := UTC_Time().Add(time.Duration(expire_time) * time.Second)
	data := jwt.MapClaims{
		"id":          id,
		"login_time":  UTC_Time(),
		"email_bind":  email,
		"regist_date": regist_date,
		"expire_time": expireTime,
		"ip":          ctx.ClientIP(),
		"wlh_to_you":  "为什么不玩原神?!",
		"random":      RandomString(32),
	}
	//
	encodedStd := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, err := encodedStd.SignedString(secretKey)
	//
	if err != nil {
		Debug("GenerateCookieToken:", err)
		tokenString = ""
	}
	//
	return tokenString
}

// 判断token是否合法
func CheckTokenLegal(token string) bool {
	data := GetTokenInfo(token)
	if data == nil {
		return false
	}
	//检查过期时间
	expireTimeStr, ok := data["expire_time"]
	if !ok {
		return false
	}
	expireTime, err := time.Parse(time.RFC3339, expireTimeStr)
	if err != nil {
		Debug("CheckTokenLegal:解析过期时间失败:", err)
		return false
	}
	if UTC_Time().After(expireTime) {
		return false
	}
	return true
}

// 获取tooken解析数据
func GetTokenInfo(signed_token string) map[string]string {
	token, err := jwt.ParseWithClaims(signed_token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.Conf.Server.JwtSecretKey, nil
	})
	//
	if err != nil {
		Debug("GetTokenInfo:", err)
		return nil
	}
	//
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := make(map[string]string)
		for key, value := range claims {
			if str, ok := value.(string); ok {
				data[key] = str
			}
		}
		return data
	}
	//
	Debug("GetTokenInfo:解析失败:" + signed_token)
	return nil
}

// 获取tooken解析数据
func GetCtxTookenInfo(ctx *gin.Context) map[string]string {
	//注意这里我并没有采用sessionID的方式，因为我觉得没必要
	//解析token
	token, err := ctx.Cookie(config.Conf.Cookie.TokenName)
	if err != nil || token == "" {
		return nil
	}
	//解码token
	data := GetTokenInfo(token)
	if data == nil {
		Debug("StudentInfoGet:", "为什么过得了我的过滤器却解析不了,出BUG了?")
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
		Debug("Mkdir:", err)
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
		Debug("读取文件失败:", err, path)
		return ""
	}
	return string(content)
}

// 写入内容到文件
func WriteFileContent(path string, data string) bool {
	err := os.WriteFile(path, []byte(data), os.FileMode(0644))
	if err != nil {
		Debug("WriteFileContent:", err)
		return false
	}
	return true
}

// 获取json格式的数据
func JsonCtx(ctx *gin.Context, addr interface{}) bool {
	err := ctx.ShouldBindBodyWithJSON(addr)
	if err != nil {
		Debug("JsonCtx:", err)
		return false
	}
	return true
}

// 获取url里面的文件路径
func GetUrlFilePath(url_ string) string {
	u, err := url.Parse(url_)
	if err != nil {
		Debug("GetUrlFilePath:", url_, err.Error())
		return ""
	}
	return strings.TrimPrefix(u.Path, "/")
}

// 拆分bucket和文件路径
func SplitBucketAndFilePath(bucket string, path string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] != bucket {
		Debug("SplitBucketAndFilePath:", bucket, path, parts)
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
