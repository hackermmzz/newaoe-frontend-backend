package config

import (
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mysql                    MySqlConfig                    `yaml:"mysql"`
	Email                    EmailConfig                    `yaml:"email"`
	Server                   ServerConfig                   `yaml:"server"`
	Code                     CodeConfig                     `yaml:"code"`
	RegistVerifyCode         RegistVerifyCodeConfig         `yaml:"registverifycode"`
	PasswordForgetVerifyCode PasswordForgetVerifyCodeConfig `yaml:"passwordforgetverifycode"`
	User                     UserConfig                     `yaml:"user"`
	Cookie                   CookieConfig                   `yaml:"cookie"`
	Redis                    RedisConfig                    `yaml:"redis"`
	OSS                      OSSConfig                      `yaml:"oss"`
	CDN                      CDNConfig                      `yaml:"cdn"`
	GRPC                     GRPCConfig                     `yaml:"grpc"`
	RocketMQ                 RocketMQConfig                 `yaml:"rocketMQ"`
	Other                    OtherConfig                    `yaml:"other"`
}
type RocketMQConfig struct {
	Host string `yaml:"host"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}
type CDNConfig struct {
	Host string `yaml:"host"` //CDN的地址
}
type MySqlConfig struct {
	Username          string `yaml:"username"`          // 用户名
	Password          string `yaml:"password"`          // 密码
	Host              string `yaml:"host"`              //地址
	DBName            string `yaml:"dbname"`            // 数据库名称
	DBMaxIdleConns    int    `yaml:"dbMaxIdleConns"`    // 数据库连接池最大空闲连接数
	DBMaxOpenConns    int    `yaml:"dbMaxOpenConns"`    // 数据库连接池最大打开连接数
	DBConnMaxLifetime int    `yaml:"dbConnMaxLifetime"` // 数据库连接最大生命周期（秒）
	DBConnMaxIdleTime int    `yaml:"dbConnMaxIdleTime"` // 数据库连接最大空闲时间（秒）

}

type RedisConfig struct {
	Password string `yaml:"password"` // 密码
	Host     string `yaml:"host"`     // 地址
}
type EmailConfig struct {
	SenderEmail       string `yaml:"senderEmail"`       // 邮箱账号
	SenderAuthCode    string `yaml:"senderAuthCode"`    // 邮箱授权码
	EmailServe        string `yaml:"emailServe"`        // 邮箱服务器
	EmailServePort    int    `yaml:"emailServePort"`    // 邮箱服务器端口
	EmailMQTopic      string `yaml:"emailMQTopic"`      // 邮箱发送topic
	EmailMQGroupCount int    `yaml:"emailMQGroupCount"` // 邮箱发送消费者组数量
}

type ServerConfig struct {
	ServeMode             string   `yaml:"serveMode"`             // 服务器运行模式
	ServerDomain          string   `yaml:"serverDomain"`          // 服务器域名
	ServerListenPort      int      `yaml:"serverListenPort"`      // 监听端口
	MaxFilsSizeUserUpload int64    `yaml:"maxFilsSizeUserUpload"` // 用户上传文件的最大限制
	TrustedProxy          []string `yaml:"TrustedProxy"`          //受信任的可以访问后端的一些地址
	JwtSecretKeyFile      string   `yaml:"jwtSecretKeyFile"`      // 存储JWT密钥的文件路径
	JwtSecretKey          []byte   `yaml:"-"`                     // JWT密钥(从文件中读取后存储在这里)
}

type CodeConfig struct {
	CodeSubmitInterval                    int    `yaml:"codeSubmitInterval"`                    // 代码提交间隔时间，单位为秒
	CodeRunQueueMaxSize                   int    `yaml:"codeRunQueueMaxSize"`                   // 代码待运行队列最大大小
	CodeRunStatusTopic                    string `yaml:"codeRunStatusTopic"`                    // 代码运行状态topic
	CodeRunStatusGroupCount               int    `yaml:"codeRunStatusGroupCount"`               // 消费者组数量
	CodeRunStatusUpdateQueueMaxSize       int    `yaml:"codeRunStatusUpdateQueueMaxSize"`       // 代码运行状态更新队列最大大小
	CodeRunStatusUpdateQueueSizeThreshold int    `yaml:"codeRunStatusUpdateQueueSizeThreshold"` // 队列触发更新大小阈值
	CodeRunStatusUpdateInterval           int    `yaml:"codeRunStatusUpdateInterval"`           // 队列更新间隔 单位：秒
	CodeWaitForRunRedisQueueTopic         string `yaml:"codeWaitForRunRedisQueueTopic"`         //代码待运行topic
	CodeSubmitTimesPerDay                 int    `yaml:"codeSubmitTimesPerDay"`                 //一天可以提交的次数
	CodeAssessmentTimes                   int    `yaml:"codeAssessmentTimes"`                   //考核提交可以提交的次数
}

type RegistVerifyCodeConfig struct {
	ExpireTime int `yaml:"expireTime"` // 过期时间
	ResendTime int `yaml:"resendTime"` // 重发时间
	Length     int `yaml:"length"`     // 验证码长度
}

type PasswordForgetVerifyCodeConfig struct {
	ExpireTime int `yaml:"expireTime"` // 过期时间
	ResendTime int `yaml:"resendTime"` // 重发时间
	Length     int `yaml:"length"`     // 验证码长度
}

type CookieConfig struct {
	ExpireTime int    `yaml:"expireTime"` // 过期时间
	TokenName  string `yaml:"tokenName"`  // Token名称，定位sessionID
}

type UserConfig struct {
	UserDefalutAvatar []string `yaml:"userDefalutAvatar"` // 默认头像(随机的)
	UserAvatarFolder  string   `yaml:"userAvatarFolder"`  // 头像存放目录
	UserCodeFolder    string   `yaml:"userCodeFolder"`    // 代码存放目录
}

type OSSConfig struct {
	Host              string `yaml:"host"`              // TOS服务地址
	AccessKey         string `yaml:"accessKey"`         // TOS访问密钥
	SecretKey         string `yaml:"secretKey"`         // TOS秘密密钥
	BucketName        string `yaml:"bucketName"`        // TOS桶名称
	PublicBaseFolder  string `yaml:"PublicBaseFolder"`  // 公共数据存储目录
	PrivateBaseFolder string `yaml:"PrivateBaseFolder"` // 私有用户数据目录
	MaxRetry          int    `yaml:"maxRetry"`          //最大尝试次数
}

type OtherConfig struct {
	AssessmentCodeUploadCategory      string   `yaml:"assessmentCodeUploadCategory"`      //考核代码上传的分类名称
	AssessmentCodeFileNameMask        string   `yaml:"assessmentCodeFileNameMask"`        //用于区分普通代码文件和考核代码文件的mask
	AvatarUploadUrlExpireTime         int      `yaml:"avatarUploadUrlExpireTime"`         //头像上传链接过期时间
	CodeCommonUploadUrlExpireTime     int      `yaml:"codeCommonUploadUrlExpireTime"`     //普通代码文件上传链接过期时间
	AssessmentCodeUploadUrlExpireTime int      `yaml:"assessmentCodeUploadUrlExpireTime"` //考核代码文件上传链接过期时间
	PrivateFileDownloadUrlExpireTime  int      `yaml:"privateFileDownloadUrlExpireTime"`  //私有文件下载链接过期时间
	PublicFileDownloadUrlExpireTime   int      `yaml:"publicFileDownloadUrlExpireTime"`   //公有文件下载链接过期时间
	SuperUser                         []string `yaml:"superUser"`                         //超级用户
	SuperUserPassword                 string   `yaml:"superUserPassword"`                 //超级用户密码
	SuperUserEmail                    string   `yaml:"superUserEmail"`                    //超级用户邮箱
	FeedbackNeedSendToEmailTopic      string   `yaml:"feedbackNeedSendToEmailTopic"`      //反馈发送邮箱队列
	FeedbackSendToEmail               string   `yaml:"feedbackSendToEmail"`               //反馈发送的邮箱
}

var (
	Conf Config
)

// 初始化一些配置
func ConfigInit() {
	//读取YAML配置文件
	WorkDir, _ := os.Getwd()
	configFilePath := path.Join(WorkDir, "assets/config.yaml")
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		panic("ConfigInit:读取配置文件失败 " + err.Error())
	}

	err = yaml.Unmarshal(configFile, &Conf)
	if err != nil {
		panic("ConfigInit:解析配置文件失败 " + err.Error())
	}
	//读取jwt密钥
	Conf.Server.JwtSecretKey, err = os.ReadFile(path.Join(WorkDir, Conf.Server.JwtSecretKeyFile))
	if err != nil {
		panic("ConfigInit:读取JWT密钥文件失败 " + err.Error())
	}
	//设置基础时区为UTC
	time.LoadLocation(time.UTC.String())
}
