# **🚀 全球跨境电商中后台系统 \- 企业级 RBAC 权限引擎**

(Global E-commerce Admin & RBAC Engine)

基于 **Golang \+ Gin \+ GORM \+ Redis \+ MySQL** 构建的企业级中后台权限管理系统。本项目以真实的“跨境电商”业务为沙盘，落地了极其严密的 **RBAC（Role-Based Access Control）** 权限模型，并针对高并发场景下的接口鉴权进行了深度性能优化。

## **🛠️ 技术栈 (Tech Stack)**

* **核心框架**: Golang, Gin (RESTful API 设计)  
* **持久层**: MySQL 8.0, GORM (优雅处理多对多关联与预加载)  
* **缓存与高并发**: Redis (Set 集合实现 O(1) 毫秒级权限碰撞)  
* **安全机制**: UUID Token, 动态路由拦截 (c.FullPath())

## **🔥 核心技术亮点 (Technical Highlights)**

### **1\. 毫秒级防盗门：基于 Redis Set 的高并发鉴权**

传统的 RBAC 每次接口请求都需要经历 User \-\> Role \-\> Permission 的多表联查，极易导致 MySQL 磁盘 I/O 成为系统瓶颈。

* **破局方案**：在用户登录 (Login) 瞬间，利用 GORM Preload 将树形权限网打平，拼接成 Method:Path（如 DELETE:/api/v1/products/:id）的扁平化结构，并将其全量推入 Redis 的 **Set 集合** 中。  
* **极致性能**：在中间件 RBACGuard 中，使用 Redis 的 SIsMember 命令进行 O(1) 复杂度的权限碰撞，将接口的鉴权耗时压榨至 **1 毫秒内**。

### **2\. 精准打击：RESTful 动态路由拦截**

针对带有动态参数的 URL（例如前端请求 DELETE /api/v1/products/99），普通中间件无法直接与数据库中存放的模板路径进行比对。

* **优雅实现**：巧妙利用 Gin 框架底层的 c.FullPath() 方法，在中间件层面精准还原动态路由模板（识别为 /api/v1/products/:id），实现对 RESTful API 的无死角管控。

### **3\. 企业级数据沙盘：创世种子脚本 (Seed Data)**

内置 main.go 创世脚本，项目首次启动时检测空库状态，自动利用 GORM 的 Association().Append() 建立复杂的“多对多”关联，瞬间生成高度仿真的电商业务矩阵。

## **🏢 业务矩阵设计 (Business Matrix)**

系统内置了跨境电商公司中 4 个极其典型的职场角色，严格遵循**最小权限原则**：

| **员工账号** | **佩戴工牌 (Role)** | **拥有权限 (Permissions)** | **业务禁区** |

| **Admin\_Root** | 超级管理员 | 拥有系统全部 10 把钥匙 | 无 |

| **Ops\_Alice** | 商品运营总监 | 只能看/增/删 **商品** (/products) | 绝对碰不到订单和财务报表 |

| **CS\_Bob** | 一线客服专员 | 只能查看/修改 **订单** (/orders) | 绝对无权下架商品 |

| **Finance\_Charlie** | 财务审计员 | 只能查看报表和审核 **退款** (/refunds) | 碰不到业务线数据 |

## **🚀 快速启动 (Quick Start)**

### **1\. 环境准备**

请确保本地已安装并运行：

* **MySQL** (默认端口: 3306\)  
* **Redis** (默认端口: 6379\)

### **2\. 修改配置 (可选)**

检查 main.go 中的数据库和 Redis 连接字符串，确保密码与您本地一致。

默认 MySQL 账号为 root:123456，数据库名为 rbac\_db。

### **3\. 运行服务**

\# 下载依赖  
go mod tidy

\# 启动服务（首次启动会自动建表并注入电商测试数据）  
go run main.go

启动成功后，终端将输出：🌱 检测到空数据库为空，正在植入电商业务预设数据...

## **🧪 接口调用演示 (API Demo)**

### **Step 1: 登录获取 Token**

使用客服专员账号登录：

* **POST** http://127.0.0.1:8080/api/v1/login  
* **Body (JSON)**: {"name": "CS\_Bob", "password": "123"}  
* **Response**: 记录返回的 token 值。

### **Step 2: 测试合法业务（查订单）**

* **GET** http://127.0.0.1:8080/api/v1/orders  
* **Headers**: Authorization: \<刚才获取的 Token\>  
* **Response**: {"msg": "成功获取订单列表：待发货 50 单"} (安检通过✅)

### **Step 3: 测试越权拦截（企图下架商品）**

* **DELETE** http://127.0.0.1:8080/api/v1/products/1  
* **Headers**: Authorization: \<刚才获取的 Token\>  
* **Response**: **403 Forbidden** {"error": "没有 \[DELETE:/api/v1/products/:id\] 的权限"} (安检拦截⛔)

## **🔐 生产环境演进方向 (Future Work)**

本项目作为底层鉴权基建的展示，后续推向生产环境（Production）时，将增加以下模块：

1. **密码安全升级**: 引入 golang.org/x/crypto/bcrypt，告别明文密码，使用 Hash & Salt 加盐哈希存储，应对脱库风险。  
2. **主动踢出机制**: 完善 Token 黑名单机制，实现对违规账号的一键强踢。  
3. **操作审计日志**: 增加 AopLog 中间件，将所有敏感修改（如 DELETE / POST）落库或推入消息队列，实现完整的操作追溯。