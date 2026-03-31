# RoomSync - 企业级轻量会议室预约系统

RoomSync 是一款为中小型企业量身定制的轻量化、高性能会议室预约系统。采用 **Golang (Gin + Gorm)** 作为后端，**Vue 3 (TypeScript + Element Plus)** 作为前端，实现了从用户鉴权到会议室管理的完整闭环。

## 🌟 项目亮点 (简历撰写参考)
- **后端 (Golang)**:
  - **高效并发检测**: 针对核心预约业务，设计了基于时间区间重叠逻辑的数据库查询算法 `(start < end_new AND end > start_new)`，并配合索引确保在高并发预订下数据的强一致性。
  - **RBAC 权限模型**: 实现了基于 JWT 的身份验证与基于中间件的角色访问控制 (Admin/Employee)，确保接口安全性。
  - **Clean Architecture**: 遵循 MVC 思想设计项目目录，分层明确（Router, Handler, Service, Repository），大幅提升了代码的可维护性。
- **前端 (Vue 3)**:
  - **现代化技术栈**: 采用 Vue 3 (Composition API) + TypeScript + Pinia + Vite，保证了极致的开发效率与运行速度。
  - **请求拦截机制**: 基于 Axios 封装了请求拦截器，实现了 Token 自动注入与全局错误消息反馈。
  - **响应式设计**: 使用 Element Plus 实现了美观的响应式管理界面，适配多种桌面显示器。

## 🛠️ 技术栈
| 模块 | 技术 |
| --- | --- |
| **后端** | Golang 1.20+, Gin, GORM, MySQL, JWT (golang-jwt), Bcrypt, Viper |
| **前端** | Vue 3, TypeScript, Vite, Element Plus, Vue Router, Pinia, Axios |
| **工具** | Git, Postman, npm |

## 📂 项目结构
```text
RoomSync/
├── backend/            # 后端工程
│   ├── api/            # 路由与控制器层 (Handler)
│   ├── config/         # 配置管理 (YAML + Viper)
│   ├── models/         # 数据库模型 (GORM)
│   ├── repository/     # 数据库连接与初始化
│   ├── utils/          # 工具函数 (JWT, Response 封装)
│   └── main.go         # 入口文件
├── frontend/           # 前端工程
│   ├── src/
│   │   ├── api/        # Axios 请求封装
│   │   ├── views/      # 页面视图 (Login, RoomList)
│   │   ├── store/      # Pinia 状态管理
│   │   └── main.ts     # 入口文件
└── README.md           # 项目文档
```

## 🚀 快速开始

### 1. 后端启动
1. 准备 MySQL 数据库，创建 `roomsync` 库。
2. 修改 `backend/config/config.yaml` 中的 `dsn` (数据库连接地址)。
3. 在 `backend/` 目录下执行：
   ```bash
   go mod tidy
   go run main.go
   ```

### 2. 前端启动
1. 在 `frontend/` 目录下安装依赖：
   ```bash
   npm install
   ```
2. 启动开发服务器：
   ```bash
   npm run dev
   ```
3. 访问 `http://localhost:5173`。

## 📸 API 验证逻辑
- **管理员**: 系统的第一个注册用户将自动获得 `admin` 权限，可以增删会议室。
- **冲突检测**: 预约时系统会自动校验所选时间段内该会议室是否已被占用。
- **状态机**: 支持 `pending` (待审核)、`approved` (已批准)、`rejected` (已拒绝)、`cancelled` (已取消) 四种状态。

## 📜 许可证
MIT License
