```text
  __________  _____   ____
 |___  / __ \|  __ \ / __ \
    / / |  | | |  | | |  | |
   / /| |  | | |  | | |  | |
  / /_| |__| | |__| | |__| |
 /_____\____/|_____/ \____/

```

# 简介

zodo 是我写给自己用的一款命令行工具，最初只是用来进行任务管理，后来我把自己在开发中用到的一些工具也整合了进去，zodo 的定位也因此变成了“程序员的百宝箱”。好吧，说 zodo 是百宝箱属实是给自己脸上贴金，杂货屋可能更加合适;-)。不管怎么样，如果你发现 zodo 对你而言有哪怕一丁点的价值，那对我来说就是无上的荣幸。欢迎通过邮件或 Issue 的方式与我交流有关 zodo 的任何话题！

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

zodo 的配置文件为 YAML 格式，地址为：

```text
~/.config/zodo/conf
```

zodo 在每次运行时都会先加载配置文件，如果找不到配置文件就会进行初始化。另外大部分关键配置都有默认值，这确保了基本功能可以开箱即用。

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

接下来介绍一些比较复杂的功能特性。

## 数据同步

是的，zodo 支持多台机器间的数据同步，这一功能是通过 Redis 来实现的。zodo 支持两种存储方式：file 和 redis，对应的配置是 storage.type。zodo 默认使用 file 存储方式，这种方式会把所有的数据存储到本地，准确说是本地的~/.config/zodo 目录下。但是，如果你恰好有一台 Redis 服务器，你就可以将 storage.type 改为 redis，再将你的服务器信息添加到 storage.redis 下的各项配置中，zodo 就会使用你的 Redis 服务器来存储数据。

这里值得一提的是 storeage.redis.localize 这个配置，它决定了在向 Redis 写数据时是否会也向本地文件写数据。这个配置默认是打开的，而且我建议你没有特殊的理由不要关闭，因为它可以保证你的本地数据永远是最新的，一旦你的 Redis 连接出了问题，你只需要将存储方式切换成本地文件就可以正常使用 zodo。如果你真的关闭了这个配置，但你又想在本地文件和 Redis 服务器间同步数据，那么 trans 子命令可以满足你的需求，虽然大多数情况下它没有什么用。

## 文本颜色

目前有三个字段的文本支持设置颜色：status、deadline 和 remain，其中 remain 是根据 deadline 自动计算的，因此和 deadline 使用相同的配置。文本颜色的配置在 todo.color 下面。zodo 目前支持以下几种颜色：

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

## 子任务

// TODO

## 服务器模式

// TODO
