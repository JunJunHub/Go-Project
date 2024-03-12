# 前言
>当第一次下载后我们如何配置自己 `Git` 账号等信息(这里可以是 `Github` 也可以是国内的码云账号等等)。下面我们开始主题。
>
>这里大概有两种方式：
> - 一种是通过全局配置信息，让所有的项目都使用这个账号，
> - 另一种是在指定项目下配置单独的信息，
>
>我们第一次下载后 `Git` 后我们最好是配置到全局，个别项目我们可以考虑单独配置。

## Git查看配置信息
`git config --list`
```shell
PS E:\Project\github\Go-Project> git config --list
core.symlinks=false
core.autocrlf=true
core.fscache=true
color.diff=auto
color.status=auto
color.branch=auto
core.symlinks=false
core.autocrlf=true
core.fscache=true
color.diff=auto
color.status=auto
color.branch=auto
color.interactive=true
help.format=html
rebase.autosquash=true
http.sslcainfo=D:/Git/mingw64/ssl/certs/ca-bundle.crt
http.sslbackend=openssl
diff.astextplain.textconv=astextplain
filter.lfs.clean=git-lfs clean -- %f
filter.lfs.smudge=git-lfs smudge -- %f
filter.lfs.process=git-lfs filter-process
filter.lfs.required=true
credential.helper=manager
core.repositoryformatversion=0
core.filemode=false
core.bare=false
core.logallrefupdates=true
core.symlinks=false
core.ignorecase=true
remote.origin.url=https://github.com/JunJunHub/Go-Project.git
credential.helper=manager
core.repositoryformatversion=0
core.filemode=false
core.bare=false
core.logallrefupdates=true
core.symlinks=false
core.ignorecase=true
remote.origin.url=https://github.com/JunJunHub/Go-Project.git
remote.origin.fetch=+refs/heads/*:refs/remotes/origin/*
branch.main.remote=origin
branch.main.merge=refs/heads/main
(END)
```

## Git配置账号信息

### 全局配置
全局配置，我们首先打开终端 `Windows` 下运行 `cmd` 窗口或者 `powershell`。然后输入命令：
```shell
git config --global user.name "JunJunHub"
git config --global user.email 520familylyj@gmail.com 
```
这里的 `JunJunHub` 是你的 `GitHub` 或者码云的账户名信息，`520familylyj@gmail.com` 是你的 `GitHub` 或者码云的邮箱信息；
如果用了 `–global` 选项，那么更改的配置文件就是位于你用户主目录下的那个，以后你所有的项目都会默认使用这里配置的用户信息；

### 配置指定项目信息
如果要在某个特定的项目中使用其他名字或者电邮，只要去掉 `–global` 选项重新配置即可，新的设定保存在当前项目的 `.git/config` 文件里。
命令行：
```shell
git config --global user.name "JunJunHub"
git config --global user.email 520familylyj@gmail.com 
```
配置完成后你可以使用上面的命令分别查询你的账号信息是否配置成功去掉对应的后面的账户信息和邮箱地址即可。

也可以使用命令行查看相关配置信息，如下：
```shell
git config --list
``` 

## 行结束符`CRLF`转换配置
```shell

```