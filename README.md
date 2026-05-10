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
        string Path "接口路径"  
        string Method "请求方法"  
    }

    USER }|--|{ ROLE : "拥有角色 (user\_roles)"  
    ROLE }|--|{ PERMISSION : "绑定权限 (role\_permissions)"

## **⚙️ 系统启动与请求流转生命周期 (时序图)**

相比于传统的流程图，本系统采用严格的中间件洋葱模型拦截。以下是服务启动及 API 请求流转的完整生命周期：

sequenceDiagram  
    autonumber  
    actor Client as 前端 / 客户端  
    participant Gin as Gin 路由引擎  
    participant CORS as CORS 跨域中间件  
    participant Auth as RBAC 鉴权拦截器  
    participant Ctrl as 业务逻辑层 (Controller)  
    participant DB as MySQL 数据库

    Note over Gin, DB: 【第一阶段】系统启动与数据预热  
    Gin-\>\>DB: 连接数据库 & AutoMigrate 自动建表  
    DB--\>\>Gin: 表结构同步完成  
    Gin-\>\>DB: 检查 users 表是否为空？  
    alt 数据库为空 (首次启动)  
        Gin-\>\>DB: 植入上帝角色 (Admin) 与基础 API 权限 (Seeder)  
    end  
    Gin-\>\>Gin: 挂载全局中间件与 API 路由，监听 :8080

    Note over Client, DB: 【第二阶段】API 接口请求流转 (以删除用户为例)  
    Client-\>\>Gin: 发起请求 (DELETE /api/v1/users/2)\<br\>携带 X-User-Id: 1  
      
    Gin-\>\>CORS: 1\. 跨域安全检查  
    CORS--\>\>Gin: 允许放行  
      
    Gin-\>\>Auth: 2\. 身份与 RBAC 权限校验  
    Auth-\>\>DB: 级联查询该请求者的 \[用户 \-\> 角色 \-\> 权限\] 树  
    DB--\>\>Auth: 返回当前用户持有的所有权限钥匙  
      
    alt 权限不足 / 未传凭证  
        Auth--\>\>Client: 拦截请求 (HTTP 403 / 401 拒绝访问)  
    else 权限校验通过 (拥有 DELETE /api/v1/users/:id 钥匙)  
        Auth-\>\>Ctrl: 3\. 安检通过，进入具体业务逻辑  
        Ctrl-\>\>DB: 执行对应记录的软删除操作  
        DB--\>\>Ctrl: 返回执行结果  
        Ctrl--\>\>Client: 响应操作成功 (HTTP 200 OK)  
    end

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