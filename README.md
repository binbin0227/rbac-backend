# **🛡️ Enterprise RBAC API (企业级权限控制引擎)**

本项目是一个基于 Go 语言构建的纯后端企业级 RBAC（Role-Based Access Control）权限管理引擎。从零设计了“用户-角色-权限”的多对多关联模型，并通过自定义 Gin 中间件实现了高安全性的细粒度接口级鉴权。

## **✨ 核心特性**

* **标准 RBAC 架构**: 彻底抛弃低效的单表权限字段，采用 users、roles、permissions 加上中间表的标准多对多（Many-to-Many）关联模型。  
* **接口级动态鉴权 (Auth Middleware)**: 独立编写 Gin 拦截器，基于每一次请求的 URL Path 和 HTTP Method 进行精准权限碰撞，坚决拦截越权访问 (返回 HTTP 403 Forbidden)。  
* **深度级联查询**: 充分利用 GORM 的 Preload 嵌套预加载技术，优雅解决复杂的 N+1 查询问题，一次性提取“用户 \-\> 角色 \-\> 权限”的完整权限关系树。  
* **数据库自动播种 (Database Seeder)**: 拥有服务启动自检能力。当检测到空数据库时，全自动植入完整的 API 路由权限与超级管理员 (Admin\_Root) 账号，实现“开箱即用”。  
* **CORS 跨域支持**: 完善的跨域资源共享配置，支持无缝对接 Vue、React 等现代前端单页应用。

## **🗄️ 核心数据架构 (ER 模型)**

系统采用经典的多对多（Many-to-Many）关联模型，由 GORM 自动维护底层的隐藏关联外键表。

erDiagram  
    %% 实体定义  
    USER {  
        uint ID PK "主键"  
        string Name "用户名"  
    }  
      
    ROLE {  
        uint ID PK "主键"  
        string Name "角色名"  
        string Description "角色描述"  
    }  
      
    PERMISSION {  
        uint ID PK "主键"  
        string Name "权限名称"  
        string Path "接口路径 (如 /api/v1/users)"  
        string Method "请求方法 (如 GET/DELETE)"  
    }

    %% 关系定义 (多对多)  
    USER }|--|{ ROLE : "属于 (对应中间表: user\_roles)"  
    ROLE }|--|{ PERMISSION : "拥有 (对应中间表: role\_permissions)"

## **⚙️ 系统启动与请求处理生命周期**

本系统在服务启动与处理 HTTP 请求时，遵循以下严密的拦截与流转逻辑：

graph TD  
    A\[1. 连接 MySQL 数据库\] \--\> B\[2. GORM 自动建表 AutoMigrate\]  
    B \--\> C{3. 检测 users 表是否为空?}  
      
    C \-- 是 (count \== 0\) \--\> D\[执行创世种子脚本: 创建基础权限 \-\> 创建超级管理员角色 \-\> 绑定权限 \-\> 创建上帝用户 \-\> 发放工牌\]  
    C \-- 否 (非空) \--\> E\[跳过种子数据初始化\]  
    D \--\> E  
      
    E \--\> F\[4. 初始化路由 SetupRouter\]  
      
    subgraph 路由与中间件拦截链路  
        F \--\> G\[加载 CORS 中间件: 允许前端跨域及 X-User-Id 请求头\]  
        G \--\> H\[挂载 V1 路由组 /api/v1\]  
        H \--\> I\[加载 RBACGuard 中间件: 校验身份及请求动作合法性\]  
        I \--\> J\[注册核心 Controller 接口\]  
        J \-.-\> K\[实体增删改查 GET / POST / DELETE\]  
        J \-.-\> L\[核心关联操作: Replace 分配角色 / 权限\]  
    end  
      
    J \--\> M\[5. 启动 Gin HTTP 服务监听 :8080\]

## **🚀 快速启动指南**

### **1\. 环境准备**

* Go 1.20 或更高版本  
* MySQL 8.0 数据库

### **2\. 克隆项目**

git clone \[https://github.com/binbin0227/rbac-backend.git\](https://github.com/binbin0227/rbac-backend.git)  
cd rbac-backend

### **3\. 配置数据库**

打开 main.go，根据你的本地环境修改 MySQL 的 DSN 连接字符串：

dsn := "root:\*\*\*\*\*\*@tcp(127.0.0.1:3306)/rbac\_db?charset=utf8mb4\&parseTime=True\&loc=Local"

*(注意：请确保 rbac\_db 数据库已提前在 MySQL 中创建，并将 \*\*\*\*\*\* 替换为真实的数据库密码)*

### **4\. 运行服务**

\# 下载依赖  
go mod tidy

\# 启动服务  
go run main.go

*注：首次启动时，终端会打印 🌱 检测到空数据库，正在植入上帝角色...，系统会自动完成建表与权限数据的播种。*

## **📝 核心 RESTful API 概览**

所有处于 /api/v1 下的接口均受到 RBACGuard 中间件保护。测试接口时，请务必在 HTTP Header 中携带身份凭证（目前演示版本为 X-User-Id: 1）。

| 请求方法 | 接口路径 | 业务描述 | 权限拦截 |
| :---- | :---- | :---- | :---- |
| GET | /api/v1/users | 获取所有用户及级联角色 | 🔒 是 |
| POST | /api/v1/roles | 创建新角色 | 🔒 是 |
| DELETE | /api/v1/perms/:id | 删除权限 | 🔒 是 |
| POST | /api/v1/users/:id/roles | **为用户分配角色 (GORM 覆盖绑定)** | 🔒 是 |
| POST | /api/v1/roles/:id/perms | **为角色分配权限 (GORM 覆盖绑定)** | 🔒 是 |

*Built with ❤️ for backend architecture study & practice.*