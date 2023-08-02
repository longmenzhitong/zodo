```text
  __________  _____   ____
 |___  / __ \|  __ \ / __ \
    / / |  | | |  | | |  | |
   / /| |  | | |  | | |  | |
  / /_| |__| | |__| | |__| |
 /_____\____/|_____/ \____/

```

# 简介

ZODO 是一款任务管理工具，采用命令行的形式，支持跨平台的数据同步，也可以仅使用本地文件进行存储。

# 环境

主要支持 macOS，理论上也支持 Linux 和 Windows，但最新版本只在 macOS 上进行测试。

# 安装

macOS 可以通过 Homebrew 进行安装：

```shell
brew install longmenzhitong/longmenzhitong/zodo
```

或：

```shell
brew tap longmenzhitong/longmenzhitong
brew install zodo
```

# 配置

ZODO 的配置文件为 YAML 格式，地址为：

```text
~/.config/zodo/config.yml
```

ZODO 在每次运行时都会先加载配置文件，如果找不到配置文件就会进行初始化。另外大部分关键配置都有默认值，这确保了基本功能可以开箱即用。

有一个专门的子命令用来查看所有的配置：

```shell
zodo conf
```

# 功能

如果想了解所有子命令的简介、用法和参数，请查看帮助信息：

```shell
zodo -h
```

或：

```shell
zodo --help
```

各个子命令同样拥有自己的帮助信息，可以用类似的命令查看。

接下来介绍一些比较有趣的功能特性。

## 数据同步

正如简介中所介绍的，ZODO 支持多台机器间的数据同步，这一功能是通过 Redis 来实现的。ZODO 支持两种存储方式：file 和 redis，对应的配置是 storage.type。ZODO 默认使用 file 存储方式，这种方式会把所有的数据存储到本地，准确说是本地的~/.config/zodo 目录下。但是，如果你恰好有一台 Redis 服务器，你就可以将 storage.type 改为 redis，再将你的服务器信息添加到 storage.redis 下的各项配置中，ZODO 就会使用你的 Redis 服务器来存储数据。

这里值得一提的是 storeage.redis.localize 这个配置，它决定了在向 Redis 写数据时是否会也向本地文件写数据。这个配置默认是打开的，而且我建议你没有特殊的理由不要关闭，因为它可以保证你的本地数据永远是最新的，一旦你的 Redis 连接出了问题，你只需要将存储方式切换成本地文件就可以正常使用 ZODO。

## 文本颜色

目前有三个字段的文本支持设置颜色：status、deadline 和 remain，其中 remain 是根据 deadline 自动计算的，因此和 deadline 使用相同的配置。文本颜色的配置在 todo.color 下面。ZODO 目前支持以下几种颜色：

- black
- red
- green
- yellow
- blue
- magenta
- cyan
- white
- hiBlack
- hiRed
- hiGreen
- hiYellow
- hiBlue
- hiMagenta
- hiCyan
- hiWhite
