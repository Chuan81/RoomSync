# RoomSync - 企业级轻量会议室预约系统

RoomSync 是一款专为中小型企业设计的轻量化会议室预约管理系统。项目采用 **Golang (Gin + Gorm)** 作为后端核心，**Vue 3 (TypeScript + Element Plus)** 作为前端展示，实现了从预约发起、自动冲突检测、管理员审批到会议签到的全生命周期闭环管理。

## 🌟 项目核心亮点 

- **高效的时间区间冲突检测算法**:
  - 核心逻辑基于区间重叠判定：`(start_time < new_end_time AND end_time > new_start_time)`。
  - 结合数据库索引优化，确保在高并发预约场景下，同一会议室在同一时间段内绝不会出现“一房多订”的情况。
- **全生命周期状态机管理**:
  - 系统内置了严谨的状态流转逻辑：`Pending (待审)` -> `Approved (准许)` -> `Checked In (已签到)` -> `Completed (已完成)`。
  - 针对未按时签到的预约，系统会自动将其标记为 `Expired (已过期)` 并释放资源。
- **懒加载式状态自动维护 (Lazy Update)**:
  - 采用“按需更新”策略，在用户查询列表时自动触发过期与完成状态的检查与更新，相比于传统的定时任务，降低了服务器常驻负载，体现了对系统性能的权衡思考。
- **精细化的预约准入策略**:
  - **提前量限制**: 支持设置会议室至少需提前多久预约（如提前 12 小时），并在后端强制校验。
  - **活跃预约上限**: 实现了单个用户在单个房间的活跃预约数量限制（默认 3 条），有效防止资源恶意占用。
- **人性化的可视化交互**:
  - **冲突可视化**: 在预约时自动拉取并标红展示已占用时间段。
  - **智能签到控制**: 按钮动态状态切换，仅在会议开始前后 30 分钟内开放签到，并配有 Tooltip 悬停提示。

## 🛠️ 技术栈
| 模块 | 技术实现 |
| --- | --- |
| **后端 (Golang)** | Go 1.20+, Gin, GORM, MySQL 8.0, JWT, Bcrypt (加密), Viper (配置管理) |
| **前端 (Vue 3)** | TypeScript, Vite, Element Plus, Pinia (状态管理), Vue Router, Axios |
| **安全/规范** | RBAC 权限模型, JWT 鉴权中间件, RESTful API 规范, Git 规范配置 |

## 📂 核心功能模块
1. **认证与授权**: 基于 JWT 的登录注册，首位注册用户自动提权为 Admin。
2. **会议室管理**: 管理员可配置会议室容量、地点、设备、审批需求及预约限制。
3. **预约中心**:
   - 员工可查看会议室大盘及已占用时段。
   - 支持“我的预约”界面，追踪申请进度。
4. **审批工作流**: 管理员一键处理申请，支持对已通过预约的“强制撤销”功能。
5. **签到系统**: 会议现场 60 分钟窗口期内签到，支持二次确认。

## 🚀 快速启动

### 1. 数据库配置
1. 创建 MySQL 数据库 `roomsync`。
2. 修改 `backend/config/config.yaml` 中的 `dsn` 配置。

### 2. 后端服务
```bash
cd backend
go mod tidy
go run main.go
```

### 3. 前端界面
```bash
cd frontend
npm install
npm run dev
```

## 📂 目录结构
```text
RoomSync/
├── backend/            # 后端 Golang 源码
│   ├── api/            # 控制器与逻辑 (Handler)
│   ├── models/         # 数据库模型 (User, Room, Booking)
│   ├── repository/     # 数据库连接与初始化 (AutoMigrate)
│   └── utils/          # JWT 签发、响应封装、中间件
└── frontend/           # 前端 Vue 3 源码
    └── src/
        ├── api/        # Axios 封装
        ├── views/      # 页面 (Login, Register, RoomList, MyBookings, BookingAdmin)
        └── store/      # Pinia 全局状态
```

## 📜 许可证
MIT License
