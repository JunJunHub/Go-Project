
# 项目仓库
- http://172.16.6.108/svn/MSP/trunk/mspv1/msp_vob/40-servers/mpuaps
- https://szgitlab.kedacom.com/goframe/mpuaps

---
# 项目介绍
- `./document/mpuaps梳理.md`

---
# 编码规范
### 命名规范
- 文件名称、数据库名称、表名称统一采用下划线命名方式
- 变量名、函数名、数据结构命名统一采用驼峰命名方式
- 命名必须有具体含义，今天使用简短的全拼单次，不可随意缩写

### 注释规范
- httpApi需添加swg文档注释
- 封装业务接口需添加注释说明
- 业务功能模块需添加注释说明
- 数据库表定义需添加注释说明

### 代码风格
- 必须校验接口返回值
- 必须使用【defer panic recover】机制封装业务逻辑, 防止访问空指针导致异常
- 统一使用框架日志组件打印日志
- todo


### 注意事项
开发新功能必须考虑兼容以下几点
- 级联
- 冗余(机箱冗余、主控冗余)
- 权限(功能权限)
- 国际化(支持多语言)

---
# 开源框架
### GoFrame Project
- https://goframe.org

### 使用的GoFrame版本
- gf-cli工具下载`https://github.com/gogf/gf-cli/releases` `\\10.0.0.11\u盘中转站\@lyj\gf-cli`
- windows版本gf-cli工具安装:`./gf_windows_amd64.exe install`
- linux版本gf-cli工具安装:`./gf_linux_amd64 install`

```
PS F:\project\MSP\40-servers\mpuaps> gf -v
GoFrame CLI Tool v1.17.0, https://goframe.org
GoFrame Version: v1.16.6 in current go.mod
CLI Built Detail:
  Go Version:  go1.15.14
  GF Version:  v1.16.4
  Git Commit:  971ed46f0b9d4dfebd1907cd3ed851cf9e1a5503
  Build Time:  2021-08-10 02:22:48
```

### 参考项目
- https://gitee.com/tiger1103/gfast
- https://github.com/1307super/gea
- https://github.com/gogf/focus-mono
- https://github.com/gogf/gf-demos

### 已知问题
- 应避免在http服务启动之前,执行sql语句(如果数据库配置错误,执行sql(异步也不行)会导致服务启动阻塞)
- gdb.SetConfig()接口会覆盖数据库默认配置,慎用.
- 在使用 gconv.Struct 进行输数据结构转换时，当结构体内部有重名的字段时(或嵌套数据结构中有重名字段)数据转换有bug


---
# 工程构建
- goland打开工程目录mpuaps
- idea配置`GOPROXY=https://goproxy.io`
- 依赖同步: 进入项目根目录`/40-servers/mpuaps`执行`go mod tidy`
- 为避免多个golang工程依赖的开源组件相互冲突,可以每个项目独立设置自己的GOPATH路径
- 下载所有依赖到本地 go mod vendor, 方便在内网编译机环境中编译代码
- 配置文件修改(./config/config.toml、./config/config_dev.toml)
- 执行`\10-common\include\proto\protobuf\gen_mspgrpclib_proto.bat` 生成 `.\library\mpurpc\mspgrpclib.pb.go`
- dao、model代码生成，源码跟目录执行 `gf gen dao`，代码生成配置 `./config/config.toml`
- `./config/config.toml` 配置编译选项 `gf build` 编译(默认编译win_amd64、linux_amd64、linux_arm64版本)
- 打包Docker镜像 `gf docker -t kedacom/mpuaps:v1.0` 或 `docker build -t kedacom/mpuaps:v1.0 .`

- 多平台协作CRLF、LF换行符问题，win平台配置
  - git config --global core.autocrlf input  #commit 时去掉CRLF, checkout时什么都不做
  - git config --global core.safecrlf warn   #发现CRLF时警告并继续

---
# 接口文档
- 生成swagger文档`gf swagger --pack`
- 启动程序，访问: http://127.0.0.1:9180/swagger 
  账户密码：admin admin123
- 文档规范：https://openapi.apifox.cn/

---
# 开发工具
- goland IDEA已上传到服务器`\\10.0.0.11\u盘中转站\@lyj\JetBrains`
- vscode