# Anon图床1.0

一个功能完整的现代化图床管理系统，基于 Vue.js 3 + Go 构建，支持POW验证、剪贴板上传等高级功能。

## 🐳 Docker 部署

### 环境要求
- Docker 20.10.0 或更高版本
- Docker Compose v2.0.0 或更高版本

### 使用 Docker Compose 部署

1. **克隆项目**
```bash
git clone <repository-url>
cd oneimgxru
```

2. **启动服务**
```bash
docker compose up -d
```

3. **访问系统**
- 前端界面：打开浏览器访问 `http://localhost`
- 后端API：`http://localhost:8080`

4. **查看日志**
```bash
# 查看所有服务日志
docker compose logs -f

# 查看特定服务日志
docker compose logs -f frontend
docker compose logs -f backend
```

5. **停止服务**
```bash
docker compose down
```

### 数据持久化
系统数据和上传的图片通过 Docker 数据卷保持持久化：
- 上传的图片存储在 `./uploads` 目录
- 数据库文件存储在 `./data` 目录

### 自定义配置
如需修改配置，可以通过环境变量或直接编辑 `docker-compose.yml` 文件：

```yaml
services:
  backend:
    environment:
      - TZ=Asia/Shanghai            # 时区设置
```

## ✨ 功能特性

### 🔐 安全认证
- POW (工作量证明) 验证登录
- Session 会话管理
- 密码加密存储
- 会话超时保护

### 📤 图片上传
- **剪贴板粘贴直接上传** - 支持 Ctrl+V 粘贴上传
- 拖拽上传支持
- 批量文件选择上传
- 支持多种图片格式 (JPEG, PNG, GIF, WebP, SVG, BMP)
- 文件大小限制和格式验证
- 上传进度显示

### 🖼️ 图片管理
- 图片预览和详情查看
- 复制链接功能
- 图片信息展示

### 📊 数据统计
- 仪表板概览
- 存储空间统计
- 实时数据更新

### 🎨 用户界面
- 现代化设计风格
- 响应式布局 (支持移动端)
- 深色/浅色主题
- 流畅的动画效果
- 直观的操作体验

### ⚙️ 系统设置
- 账户密码管理

## 🛠️ 技术栈

### 前端技术
- **Vue.js 3** - 渐进式JavaScript框架
- **Composition API** - Vue 3 组合式API
- **Vite** - 现代化构建工具
- **Vue Router** - 官方路由管理器
- **Axios** - HTTP客户端
- **SCSS** - CSS预处理器

### 后端技术
- **Go** - 高性能编程语言
- **Gin** - 轻量级Web框架
- **GORM** - Go语言ORM库
- **SQLite/MySQL** - 数据库支持
- **bcrypt** - 密码加密
- **Session** - 会话管理

### 开发工具
- **热重载** - 开发时自动刷新
- **代码分割** - 优化加载性能
- **压缩优化** - 生产环境优化
- **跨域支持** - CORS配置

## 🚀 快速开始

### 环境要求
- Go 1.19 或更高版本
- Node.js 16 或更高版本
- npm 或 yarn 包管理器

### 一键启动

#### Windows 用户
```bash
# 双击运行或在命令行执行
start.bat
```

#### Linux/Mac 用户
```bash
# 添加执行权限
chmod +x start.sh
# 运行启动脚本
./start.sh
```

## 📁 项目结构

```
oneimgxru/
├── frontend/                 # 前端项目
│   ├── src/
│   │   ├── App.vue         # 应用入口
│   │   ├── assets/         # 资源文件
│   │   ├── components/     # Vue组件
│   │   ├── views/          # 页面组件
│   │   ├── utils/          # 工具函数
│   │   └── router/         # 路由配置
│   ├── dist/               # 构建输出
│   └── package.json        # 前端依赖
├── backend/                 # 后端项目
│   ├── controllers/        # 控制器
│   ├── models/            # 数据模型
│   ├── routes/            # 路由配置
│   ├── middlewares/       # 中间件
│   ├── database/          # 数据库配置
│   ├── config/            # 配置文件
│   └── main.go            # 入口文件
├── uploads/                # 上传文件存储
├── data/                   # 数据文件
├── start.bat              # Windows启动脚本
├── start.sh               # Linux/Mac启动脚本
└── README.md              # 项目说明
```

## ⚙️ 配置说明

### 数据库配置
系统支持 SQLite 和 MySQL 两种数据库：

```go
// SQLite (默认)
cfg.IsMysql = false
cfg.SqlitePath = "./data/data.db"

// MySQL
cfg.IsMysql = true
cfg.DbHost = "localhost"
cfg.DbPort = 3306
cfg.DbUser = "root"
cfg.DbPassword = "password"
cfg.DbName = "oneimg"
```

### 上传配置
```go
cfg.UploadPath = "./uploads"     // 上传目录
cfg.MaxFileSize = 10 * 1024 * 1024  // 最大文件大小 (10MB)
```

## 📖 使用指南

### 登录系统
1. 访问 `http://localhost:8080`
2. 点击"登录"按钮
3. 完成POW验证 (人机验证)
4. 验证成功后自动登录

### 上传图片

#### 方法一：剪贴板上传
1. 复制图片到剪贴板 (截图、复制图片等)
2. 在上传区域按 `Ctrl+V` (Windows) 或 `Cmd+V` (Mac)
3. 图片自动开始上传

#### 方法二：拖拽上传
1. 将图片文件拖拽到上传区域
2. 松开鼠标完成上传

#### 方法三：点击上传
1. 点击"选择文件"按钮
2. 选择一个或多个图片文件
3. 确认上传

### 管理图片
- **查看图片**：点击图片缩略图查看大图
- **复制链接**：点击链接图标复制图片URL
- **删除图片**：点击删除按钮移除图片

### 查看统计
- 访问"统计"页面查看上传数据

## 🔧 开发指南

### 前端开发
```bash
cd frontend
pnpm run dev          # 启动开发服务器
pnpm run build        # 构建生产版本
pnpm run preview      # 预览构建结果
```

### 后端开发
```bash
cd backend
go run main.go       # 启动开发服务器
go build            # 构建可执行文件
```

### API接口

#### 认证接口
- `POST /api/login` - 用户登录
- `GET /api/user/info` - 获取用户信息
- `GET /api/user/status` - 检查登录状态
- `POST /api/logout` - 用户登出

#### 图片接口
- `POST /api/upload` - 单图上传
- `POST /api/upload/images` - 批量上传
- `GET /api/images` - 获取图片列表
- `GET /api/images/:id` - 获取图片详情
- `DELETE /api/images/:id` - 删除图片

#### 统计接口
- `GET /api/stats/dashboard` - 仪表板数据
- `GET /api/stats/images` - 图片统计数据

#### 账户接口
- `POST /api/account/change` - 修改密码
- `POST /api/sessions/clear` - 清除会话

## 🔒 安全特性

- **密码加密**：使用bcrypt加密存储
- **会话管理**：安全的session机制
- **CSRF防护**：跨站请求伪造防护
- **文件验证**：严格的文件类型检查
- **大小限制**：防止恶意大文件上传

## 🚀 性能优化

- **前端优化**：代码分割、懒加载、压缩
- **后端优化**：连接池、缓存、并发处理
- **图片优化**：自动压缩、格式转换
- **CDN支持**：静态资源CDN加速