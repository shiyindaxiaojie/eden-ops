# Go 目录结构

## `/cmd` 程序入口

只存放 `main` 函数，从 `/internal` 和 `/pkg` 目录导入和调用代码，除此之外没有别的东西。

示例： `/cmd/myapp/main.go`

## `/internal` 私有库代码

外部项目不可导入，由 Go 强制执行。有关更多细节，请参阅 Go 1.4 [`release notes`](https://golang.org/doc/go1.4#internalpackages) 。

示例： `/internal/api/user.go`

## `/pkg` 公有库代码

外部项目会导入这些库。关于 `pkg` 和 `internal` 目录的区别，可以查阅 Travis Jeffery 撰写的 [`I'll take pkg over internal`](https://travisjeffery.com/b/2019/11/i-ll-take-pkg-over-internal/) 博客文章。

示例： `/pkg/auth/jwt.go`

## `/vendor` 应用程序依赖项

存放项目依赖项。使用 `go mod vendor` 命令创建 `/vendor` 目录。可以手动管理或使用你喜欢的依赖项管理工具，如新的内置 [`Go Modules`](https://go.dev/wiki/Modules) 功能)。

示例：`/vendor/github.com/mattn`

> 请注意，如果未使用默认情况下处于启用状态的 Go 1.14，则可能需要在 `go build` 命令中添加 `-mod=vendor` 标志。自从 [`1.13`](https://golang.org/doc/go1.13#modules) 以后，Go 还启用了模块代理功能（默认使用 [`https://proxy.golang.org`](https://proxy.golang.org) 作为他们的模块代理服务器）。国内模块代理功能默认是被墙的，七牛云有维护专门的的[`模块代理`](https://github.com/goproxy/goproxy.cn/blob/master/README.zh-CN.md) 。

## `/api` 服务端接口描述

存放 OpenAPI/Swagger 规范和协议文件。

示例：`/vendor/github.com/mattn`

## `/web` 前端应用程序代码

存放 Web 应用程序的组件，例如静态 Web 资源、服务器端模板和 SPAs。

示例：`/web/src/App.vue`

## `/configs` 配置管理

存放配置文件和模板，例如 `confd` 或 `consul-template` 模板文件。

示例：`/configs/config.yaml`

## `/init` 自启动服务

存放系统自启服务描述文件，例如 `System init`（systemd，upstart，sysv）和 `process manager/supervisor`（runit，supervisor）配置。

示例：`/init/myapp.service`

## `/scripts` 初始化脚本

存放各种构建、安装、分析等操作的脚本。

示例：`/scripts/sql/init.sql`

## `/build` 系统构建

存放打包配置、脚本和持续集成描述文件。

示例：`/build/ci/.gitlab-ci.yml`

## `/deployments` 云原生部署

存放 IaaS、PaaS、系统和容器编排部署配置和模板，例如 `docker-compose`、`kubernetes/helm`、`mesos`、`terraform`、`bosh`。

示例：`/deployments/docker/Dockerfile`

## `/test` 测试脚本和数据

存放外部测试应用程序和测试数据，Go 会忽略以 `.` 或`_` 开头的目录或文件。

示例：`/test/data`

## `/docs` 使用文档

存放系统设计文档和用户使用文档（除了 `godoc` 生成的文档）。

### `/tools` 支持工具

存放项目的支持工具。注意，这些工具可以从 `/pkg` 和 `/internal` 目录导入代码。

### `/examples` 代码示例

存放应用程序和或公共库的使用示例。

### `/third_party` 第三方工具

存放外部辅助工具，分叉代码和其他第三方工具，例如 Swagger UI。

### `/assets` 静态资源

存放存储库一起使用的其他资源，例如图像、徽标等。

### `/website` 相关网站

存放项目的网站数据。