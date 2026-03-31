# RoomSync - 会议室预约系统

这是一个企业级轻量会议室预约系统 MVP 版。采用前后端分离架构。

## 技术栈
- **后端:** Golang, Gin, Gorm, MySQL, JWT
- **前端:** Vue 3, Vite, Element Plus, TypeScript (预留结构)

## 核心功能
1. 用户认证与授权 (基于 JWT, 划分 Admin/Employee 角色)
2. 会议室资源的 CRUD
3. 会议室预约，包括严格的时间冲突检测逻辑
4. 审批工作流 (特定会议室需管理员审核)
5. 取消预约功能

## 如何运行后端

### 1. 环境准备
- 确保安装了 `Go 1.20+`
- 确保安装并运行 `MySQL`

### 2. 数据库配置
在 `MySQL` 中创建一个数据库:
```sql
CREATE DATABASE roomsync CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

修改 `backend/config/config.yaml` 中的 `dsn` 配置:
```yaml
database:
  dsn: "root:你的密码@tcp(127.0.0.1:3306)/roomsync?charset=utf8mb4&parseTime=True&loc=Local"
```

### 3. 运行项目
```bash
cd backend
go mod tidy
go run main.go
```
启动后，系统将自动执行 Gorm 的 AutoMigrate 构建表结构。

### 4. 测试 API
> 注意：注册的第一个用户将自动成为 `admin` 角色。
可以使用 Postman 导入或直接发送 HTTP 请求到 `http://localhost:8080/api/...` 进行接口测试。

## 项目亮点 (简历撰写参考)
- **架构设计**: 遵循 MVC 和 Clean Architecture 思想，使路由、中间件、逻辑层和数据访问层分离。
- **并发控制**: 借助 SQL 的条件约束及事务锁机制，有效预防了会议室高并发预约产生的时间冲突。
- **安全机制**: 实现了基于 JWT 的身份验证和基于中间件的 RBAC (Role-Based Access Control) 角色控制。
- **规范编码**: 遵循统一的 HTTP 状态码与响应格式封装，提高了前后端对接效率。

## 前端 (待完善)
前端框架已预留在 `frontend` 目录，建议使用 `npm create vite@latest frontend -- --template vue-ts` 初始化 Vue 3 并在其中使用 Element Plus 快速构建 UI。
