# newaoe-frontend-backend 项目文档

## 项目说明
newaoe-frontend-backend 是一个全栈教育管理系统，提供用户认证、代码提交、文件管理、邮件服务和管理看板等功能。系统采用前后端分离架构，前端使用Vue.js框架，后端基于Go语言开发，结合MySQL和Redis实现数据持久化与缓存。项目支持教育机构进行学生作业提交、自动化评估和教学管理。

## 技术栈
| 模块        | 技术选型          |
|-------------|-------------------|
| 前端        | Vue.js + Tailwind CSS |
| 后端        | Go + Gin          |
| 数据库      | MySQL 8.0+        |
| 缓存        | Redis 7.0+        |
| 部署        | Docker + Nginx    |

## 项目结构说明
```
├── build/                # 构建脚本与资源文件
│   ├── build.sh          # 构建脚本
│   ├── newaoe.sql        # 数据库初始化文件
│   └── *.txt             # 配置白名单文件
├── config/               # 配置管理模块
│   └── config.go         # 全局配置加载逻辑
├── dao/                  # 数据访问层
│   ├── database.go       # 数据库连接池管理
│   └── *.go              # 各业务数据操作接口
├── ServerData/           # 静态资源目录
│   └── public/           # Web静态文件（HTML/CSS/JS）
├── service/              # 核心业务逻辑
│   ├── server.go         # 主服务启动入口
│   └── *.go              # 各功能模块服务实现
├── src/                  # 前端源码目录
│   ├── App.vue           # 根组件
│   ├── main.js           # Vue入口文件
│   ├── router.js         # 路由配置
│   └── components/       # Vue组件库
└── main.go               # 程序启动入口
```

## 环境依赖
- Go 1.21+ (项目使用Go模块管理)
  - 安装指南: https://go.dev/dl/
- MySQL 8.0+ (需提前导入build/newaoe.sql)
  - 安装指南: https://dev.mysql.com/downloads/installer/
- Redis 7.0+ (用于缓存和会话管理)
  - 安装指南: https://redis.io/download/
- Node.js 18+ (前端资源构建依赖)
  - 安装指南: https://nodejs.org/en/download/

## 快速启动
```bash
# 1. 安装依赖
go mod tidy
npm install

# 2. 测试环境启动
npm run dev      # 前端开发服务器
go run main.go   # 后端API服务

# 2. 配置数据库连接
cp config/config.go{.example,}
# 修改config.go中的数据库连接字符串

# 3. 初始化数据库
mysql -u root -p < build/newaoe.sql

# 4. 构建前端
npm run build

# 5. 启动服务
go run main.go
# 服务默认监听 :8080 端口

# 6. 访问应用
open http://localhost:8080
```

## 功能模块
### 用户系统
- 注册/登录/权限控制（JWT鉴权）
- 邮件验证码与通知服务
- 密码找回与安全设置

### 教学管理
- 作业发布与提交系统
- 自动化代码评估（支持多语言）
- 学生成绩管理看板

### 文件服务
- 支持代码文件/文档上传
- 文件版本控制与历史记录
- 安全文件下载与共享

### 系统管理
- 用户权限分级（学生/教师/管理员）：
  - RBAC权限模型实现
  - 角色权限动态配置
- 操作日志审计：
  - 记录用户操作轨迹
  - 支持日志导出与分析
- 系统健康监控：
  - 实时监控服务状态
  - 系统资源使用统计
  - 异常自动告警机制

**后端模块**
1. **用户系统** - 支持注册/登录/权限控制 (dao/user.go)
2. **代码提交** - 提供代码提交与执行接口 (dao/code.go)
3. **文件管理** - 支持文件上传/下载 (service/upload.go)
4. **邮件服务** - 邮件验证码与通知 (service/EmailSend.go)
5. **管理后台** - 提供数据管理界面 (service/admin.go)

## 开发规范
### 代码规范
- 前端：遵循Vue官方风格指南，组件命名采用PascalCase
- 后端：Go代码符合Uber Go Style Guide
- API设计：RESTful规范，版本化接口（/api/v1/）

## API文档
### 认证接口
- `POST /api/v1/register` - 用户注册
  - 请求参数：用户名、密码、邮箱
  - 响应：JWT令牌
- `POST /api/v1/login` - 用户登录
  - 请求参数：用户名/邮箱、密码
  - 响应：包含用户信息的JWT
- `GET /api/v1/logout` - 用户登出
  - 需携带有效JWT
  - 清除会话信息

### 作业管理
- `GET /api/v1/assignments` - 获取作业列表
- `POST /api/v1/assignments` - 创建新作业
- `POST /api/v1/submit` - 提交作业

### 文件服务
- `POST /api/v1/upload` - 文件上传
- `GET /api/v1/files/:id` - 文件下载
- `DELETE /api/v1/files/:id` - 文件删除

### 管理接口
- `GET /api/v1/admin/users` - 用户管理
- `GET /api/v1/admin/metrics` - 系统指标
### 安全规范
- 所有密码需Bcrypt加密存储
- JWT令牌设置合理过期时间
- 敏感操作需二次验证
- 输入参数严格校验（防XSS/SQL注入）

### 日志规范
- 使用Zap日志库分级记录
- 生产环境保留ERROR级别日志7天
- 关键操作记录用户ID与IP地址
  docs: update documentation
  style: format code
  refactor: code refactoring
  test: add tests
  chore: maintenance tasks
  ```

## 贡献指南
### 开发流程
1. 创建Issue描述需求/问题
2. Fork仓库并创建feature分支
3. 实现功能并编写单元测试
4. 提交PR时关联对应Issue

## 测试指南
### 单元测试
```bash
# 前端测试
npm run test:unit

# 后端测试
go test ./... -v
```

### 接口测试
使用Postman集合进行接口验证：
1. 导入`postman_collection.json`
2. 运行集合验证所有API端点

### 测试要求
- 单元测试覆盖率需>80%
- 接口测试使用Postman集合
- 前端组件测试使用Vue Test Utils

### 代码审查
- 至少2名维护者评审
- 重大变更需更新文档
- 合并前需通过CI流水线
