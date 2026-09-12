# newaoe 后端项目

## 项目简介 📚

newaoe 是一个基于 Go 语言开发的代码运行和评估系统后端，支持用户上传代码、运行代码、提交评估等功能。

## 技术栈 🛠️

- **语言**: Go 1.25.0
- **Web 框架**: Gin v1.12.0
- **数据库**: MySQL
- **缓存**: Redis v8.11.5
- **消息队列**: RocketMQ v2.1.2
- **对象存储**: MinIO v7.0.100
- **认证**: JWT v4.5.2
- **RPC**: gRPC v1.80.0
- **配置管理**: YAML
- **其他**: CORS, xorm

## 项目结构 📁

```
newaoe-frontend-backend/
├── assets/           # 配置文件和密钥
│   ├── config.yaml   # 系统配置文件
│   └── jwt_secret.key # JWT 密钥
├── build/            # 部署相关脚本和配置
│   ├── CDN/          # CDN 部署配置
│   ├── FrontEnd/     # 前端部署配置
│   ├── MQ/           # 消息队列部署配置
│   ├── MySQL/        # 数据库部署配置
│   └── Redis/        # Redis 部署配置
├── config/           # 配置相关代码
├── dao/              # 数据访问层
│   ├── codeRun.go    # 代码运行相关数据操作
│   ├── database.go   # 数据库连接
│   ├── redis.go      # Redis 连接
│   └── oss.go        # 对象存储操作
├── service/          # 业务逻辑层
│   ├── Code/         # 代码运行相关服务
│   ├── Data/         # 数据管理服务
│   ├── Download/     # 文件下载服务
│   ├── Email/        # 邮件服务
│   ├── Grpc/         # gRPC 服务
│   ├── Home/         # 首页相关服务
│   ├── Server/       # 服务器配置
│   ├── Service/      # 服务初始化
│   ├── Upload/       # 文件上传服务
│   └── User/         # 用户相关服务
├── minio/            # MinIO 存储数据
├── main.go           # 项目入口
├── go.mod            # Go 模块依赖
├── go.sum            # 依赖校验和
└── run.sh            # 运行脚本
```

## 环境要求 🔧

- Go 1.25.0 或更高版本
- MySQL 数据库
- Redis 服务
- RocketMQ 服务
- MinIO 服务
- 网络连接

## 快速开始 🚀

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境

编辑 `assets/config.yaml` 文件，根据实际环境配置数据库、Redis、RocketMQ 等服务的连接信息。

### 3. 启动服务

```bash
# 直接运行
go run main.go

# 或使用运行脚本
chmod +x run.sh
./run.sh
```

## 配置说明 ⚙️

配置文件位于 `assets/config.yaml`，主要配置项包括：

### 数据库配置
- `mysql.username`: 数据库登录用户名
- `mysql.password`: 数据库登录密码
- `mysql.host`: 数据库连接地址
- `mysql.dbname`: 要连接的数据库名称
- `mysql.dbMaxIdleConns`: 数据库连接池最大空闲连接数
- `mysql.dbMaxOpenConns`: 数据库连接池最大打开连接数

### Redis 配置
- `redis.password`: Redis 密码
- `redis.host`: Redis 连接地址

### RocketMQ 配置
- `rocketMQ.host`: RocketMQ 服务地址

### gRPC 配置
- `grpc.port`: gRPC 服务端口

### 邮箱服务配置
- `email.senderEmail`: 发件人邮箱账号
- `email.senderAuthCode`: 邮箱授权码
- `email.emailServe`: 邮箱 SMTP 服务器地址

### 服务器基础配置
- `server.serveMode`: 服务器运行模式 (debug/release)
- `server.serverDomain`: 服务器访问域名
- `server.serverListenPort`: 服务监听端口
- `server.maxFilsSizeUserUpload`: 用户上传文件最大限制

### 代码运行配置
- `code.codeSubmitInterval`: 代码提交间隔时间
- `code.codeRunQueueMaxSize`: 代码待运行队列最大大小

### OSS 配置
- `oss.host`: OSS 服务地址
- `oss.accessKey`: OSS 访问密钥
- `oss.secretKey`: OSS 秘密密钥
- `oss.bucketName`: OSS 桶名称

## API 文档 📡

### 用户相关接口

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 注册 | POST | `/api/user/regist` | 用户注册 |
| 登录 | POST | `/api/user/login` | 用户登录 |
| 密码重置 | POST | `/api/user/password/reset` | 重置密码 |

### 代码相关接口

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 代码运行 | POST | `/api/code/run` | 运行代码 |
| 代码获取 | GET | `/api/code/get` | 获取代码 |
| 代码重运行 | POST | `/api/code/reRun` | 重新运行代码 |

### 上传相关接口

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 头像上传 | POST | `/api/upload/avatar` | 上传用户头像 |
| 代码上传 | POST | `/api/upload/code` | 上传代码文件 |
| 考核提交 | POST | `/api/upload/assessment` | 提交考核代码 |

### 下载相关接口

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 私有文件下载 | GET | `/api/download/private` | 下载私有文件 |
| 公共文件下载 | GET | `/api/download/public` | 下载公共文件 |

### 首页相关接口

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 信息获取 | GET | `/api/home/info` | 获取用户信息 |
| 历史记录 | GET | `/api/home/history` | 获取历史记录 |
| 教师信息 | GET | `/api/home/teacher` | 获取教师信息 |
| 反馈提交 | POST | `/api/home/feedback` | 提交反馈 |

## 部署 🐳

### Docker 部署

1. 确保 Docker 和 Docker Compose 已安装

2. 使用项目提供的部署脚本：

```bash
# 进入部署目录
cd build

# 启动所有服务
docker-compose up -d
```

### 服务组件

- **MySQL**: 数据库服务
- **Redis**: 缓存服务
- **RocketMQ**: 消息队列服务
- **MinIO**: 对象存储服务
- **Nginx**: 前端和 CDN 服务
- **后端服务**: 运行 newaoe 后端应用

### 部署注意事项

- 确保所有服务的端口不冲突
- 配置文件中的服务地址需要根据实际部署环境修改
- 首次部署需要初始化数据库结构

## 总结

newaoe 后端项目是一个功能完整的代码运行和评估系统，采用现代化的技术栈和架构设计。它支持用户注册登录、代码上传运行、考核提交、文件管理等功能，为前端应用提供稳定可靠的后端服务。

项目采用模块化设计，代码结构清晰，易于维护和扩展。通过合理的配置管理和部署方案，可以快速搭建生产环境，为用户提供优质的服务体验。