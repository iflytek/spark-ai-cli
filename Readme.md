<h2 align="center">

<picture>
  <img src="./docs/img/logo.png" />
</picture>

讯飞星火智能终端 (spark-ai-cli)
</h2>

## 项目地址

* Github: [https://github.com/iflytek/spark-ai-cli](https://github.com/iflytek/spark-ai-cli)
  欢迎点赞，star

## 前言

感谢开源的力量，希望讯飞开源越做越好，星火大模型效果越来越好！。

## 近期规划新特性

- [x] config模块接入
- [x] 执行结果优化
- [x] 获取终端环境，进行个性化
- [ ] agents
- [ ] plugins
- [ ] 知识库

## 快速开始

**一键安装**

* 安装:

**linux和Mac**

`sudo bash -c "$(curl -s -L https://521github.com/iflytek/spark-ai-cli/releases/download/latest/install.sh)" `

**windows**

使用Windows PowerShell管理员权限执行如下脚本

```shell
Invoke-Expression (Invoke-RestMethod 'https://521github.com/iflytek/spark-ai-cli/releases/download/latest/install.ps1')
```

* 升级最新stable版本:

`aispark update`


* 查看当前版本:
  `aispark version`

```bash
aispark cli version: v0.0.20
Git Commit Hash: 69ad2242e775d58299c62bd57477ccd2eab43ae6
Build TimeStamp: Wed, 24 Apr 2024 11:39:26 +0000
GoLang Version: 1.21

     _     ___  ____   ____    _     ____   _  __
    / \   |_ _|/ ___| |  _ \  / \   |  _ \ | |/ /
   / _ \   | | \___ \ | |_) |/ _ \  | |_) || ' /
  / ___ \  | |  ___) ||  __// ___ \ |  _ < | . \
 /_/   \_\|___||____/ |_|  /_/   \_\|_| \_\|_|\_\


讯飞云提供计算服务
https://xinghuo.xfyun.cn/sparkapi
```
## 账号
### 账号申请
前往[讯飞开放平台](https://passport.xfyun.cn/register)免费注册
### 账号登录
```shell
aispark login
```

### 账户注销
```shell
aispark logout
```


## 如何使用

### shell相关的问题
```shell
aispark awk如何截取字符
aispark nginx重启
#或者q
aispark q nginx重启
```

### 知识问答问题
```shell
aispark c 今天天气怎么样
```

### fuck模式
支持平台：

- [x] windows powershell
- [x] macos bash
- [x] linux bash



#### 使用方式
```shell
aispark fuck
aispark fuck aptget
```

#### 配置

**Windows**

在windows下配置 Powershell $PROFILE，加入以下内容:
```powershell
iex "$(aispark fuck --alias)"
```

这个文件是一个脚本文件，当 PowerShell 启动时自动执行里面的内容。如果你想知道你的 $PROFILE 文件的具体位置，可以在 PowerShell 窗口中运行以下命令：

```powershell
echo $PROFILE
```

您的目录下可能没有这个文件，若没有该文件，新建此文件并配置上述命令即可

**Linux**

将以下命令放在 .bash_profile，.bashrc,.zshrc 或其他启动脚本中：

```shell
eval $(aispark fuck --alias)
```
例如：
```shell
echo 'eval "$(aispark fuck --alias)"' >> ~/.bashrc
```

更改仅在新的 shell 会话中可用。要立即进行更改，请运行 source ~/.bashrc （或 shell 配置文件，如 .zshrc ）。




### mode切换
```shell
aispark -l nginx重启      #可交互模式
aispark -v nginx重启      #啰嗦模式，解释脚本内容
```

### 配置

aispark现已开箱即用，默认的appid为公用，若想体验，可登录https://www.xfyun.cn/ 免费注册个人账户，独享个人账户，高效稳定

```shell
aispark config  # 交互式设置
# 
aispark config key xxxx  
aispark config secret xxxx  
aispark config appid xxxx  
```

## 欢迎贡献

扫码加入交流群

![用户交流群](./docs/img/wetchat.jpg "Shiprock")]

## 已知问题

* 项目目前处于开发阶段，部分思想借鉴开源实现


## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=iflytek/spark-ai-cli&type=Date)](https://star-history.com/#iflytek/spark-ai-cli&Date)

## 致谢
