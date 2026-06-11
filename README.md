# vmcontroller

vmcontroller 是一个基于 KubeVirt 的虚拟机管理 Web 应用，提供虚拟机列表、启动、停止、状态查询，以及通过 noVNC 访问 VM 控制台的能力。

## 项目简介

该项目由以下部分组成：

- 后端：Go + Gin，负责与 KubeVirt 集群交互，并提供 REST API
- 前端：Vue 3 + Vite + Vuetify，提供管理界面
- 控制台：内嵌 noVNC，支持浏览器访问虚拟机图形控制台
- 认证：JWT 登录（当前实现为示例登录逻辑）

## 主要能力

- 列出 KubeVirt 中的虚拟机
- 启动 / 停止指定虚拟机
- 查询虚拟机状态
- 通过 WebSocket / noVNC 访问虚拟机图形控制台
- 将前端资源打包进 Go 可执行文件，方便部署

## 运行环境要求

- Go 1.25+
- Node.js 20+
- Yarn
- 已配置好的 KubeVirt / Kubernetes 集群
- 可用的 kubeconfig（用于访问集群）

## 快速开始

### 1. 安装前端依赖

```sh
cd web/ui
yarn install
```

### 2. 构建前端与后端

项目根目录提供了构建脚本：

```sh
make all
```

这将执行：

- 构建前端：`web/ui` 下的 Vue 应用
- 构建后端：输出可执行文件 `./vmcontroller`

### 3. 启动应用

```sh
go run .
```

启动后，默认监听：

- http://localhost:8080

## 默认登录

当前实现中内置了一个演示用登录账号：

- 用户名：admin
- 密码：Edu@9527

> 说明：这部分认证逻辑是示例实现，生产环境建议替换为真实的用户校验与密钥管理。

## API 说明

应用通过 JWT 保护以下接口，访问时需要在请求头中携带 Authorization。

### 虚拟机管理

- GET /api/vm：获取虚拟机列表
- GET /api/vm/:vmname/start：启动虚拟机
- GET /api/vm/:vmname/stop：停止虚拟机
- GET /api/vm/:vmname：查询虚拟机状态

### 登录

- POST /login：获取 JWT Token

### VNC 控制台

- GET /vm/:vmname/vnc：用于浏览器访问虚拟机 VNC 控制台

## Docker 部署

可以直接使用仓库中的 Dockerfile 构建镜像：

```sh
docker build -t vmcontroller .
```

运行容器：

```sh
docker run -p 8080:8080 --rm vmcontroller
```

## 项目结构

```text
.
├── main.go                # 程序入口
├── vm/                    # KubeVirt 操作封装
├── web/                   # Gin 路由、JWT、VNC、前端资源
│   ├── ui/                # Vue 前端
│   └── noVNC-1.6.0/       # 嵌入的 VNC 前端资源
├── k8s-deploy/            # Kubernetes 部署示例
└── Makefile               # 构建脚本
```

## 注意事项

- 当前项目依赖本地 kubeconfig 与 KubeVirt 集群环境。
- 生产环境建议将 JWT 密钥、登录校验与网络配置改为更安全的实现。
- 如需进一步扩展，可在 `web/` 与 `vm/` 目录中继续增加更多管理能力。

