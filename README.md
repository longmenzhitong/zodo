# 简介

zodo是一款Local First的命令行任务管理工具，数据以JSON格式存储于本地，也支持跨平台的数据同步。

# 环境

主要支持macOS，理论上也支持Linux和Windows，但最新版本只在macOS上进行测试。

# 安装

macOS可以通过Homebrew进行安装：

```shell
brew install longmenzhitong/longmenzhitong/zodo
```

或：

```shell
brew tap longmenzhitong/longmenzhitong
brew install zodo
```

其他平台需要自行编译安装。

# 存储

zodo的工作目录为~/.config/zodo/，zodo的所有数据都以文件的形式存储在这个目录下：

| 文件名     | 说明     |
| ---------- | -------- |
| config.yml | 配置文件 |
| id         | 下一个id |
| todo       | 任务数据 |
| .backup    | 备份文件 |

# 配置

zodo每次运行时都会先加载配置文件，如果找不到就会进行初始化。另外大部分关键配置都有默认值，这确保了基本功能可以开箱即用。

使用`zodo conf`查看当前的所有配置。

# 功能

使用`zodo -h`或`zodo --help`查看子命令列表。各个子命令同样拥有自己的帮助信息，可以用相同的方式查看。

## 数据同步

zodo目前支持通过Redis或AWS S3同步数据，相关子命令为pull和push，相关配置为：

| 配置                | 说明                     |
| ------------------- | ------------------------ |
| sync.type           | 同步类型(redis/s3)       |
| sync.redis.address  | Redis地址(\<ip>:\<port>) |
| sync.redis.password | Redis密码                |
| sync.redis.db       | RedisDb                  |
| sync.s3.bucket      | S3的Bucket名称           |

使用Redis无须在本地额外配置，使用S3则需要在~/.aws目录下配置好region、aws_access_key_id和aws_secret_access_key，还需要给access_key所属的用户开通Bucket的操作权限（比如加入有权限的用户组），详情请参考[官方文档](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/)。

## 文本颜色

目前列表和详情所展示的字段中有以下几个支持设置颜色：

| 字段     | 说明     | 配置                |
| -------- | -------- | ------------------- |
| status   | 任务状态 | todo.color.status   |
| deadline | 任务期限 | todo.color.deadline |
| remain   | 剩余时间 | todo.color.deadline |

支持的颜色列表：

- black/hiBlack
- red/hiRed
- green/hiGreen
- yellow/hiYellow
- blue/hiBlue
- magenta/hiMagenta
- cyan/hiCyan
- white/hiWhite

## 邮件提醒

zodo的提醒功能是通过定时任务+邮件的方式来实现的。定时任务以配置的执行计划（例如每分钟一次）检查有没有需要提醒的任务，如果有就发送邮件。开启定时任务需要以Server模式运行(`zodo server`)。

提醒功能需要以下配置：

| 配置             | 说明               |
| ---------------- | ------------------ |
| reminder.enabled | 提醒功能总开关     |
| reminder.cron    | 定时任务的执行计划 |
| email.server     | 邮件服务器地址     |
| email.port       | 邮件服务器端口     |
| email.auth       | 邮件服务器密码     |
| email.from       | 发送邮箱           |
| email.to         | 目标邮箱           |

邮件服务器只测试了QQ邮箱。

## 编辑器

想要添加一个任务，最简单的方式是`zodo add <content>`。但是如果不输入content，直接使用`zodo add`的话，就会调起配置的编辑器。在编辑器中写好任务内容后保存并退出，就可以成功添加任务。这种方式适合内容文字较多或操作较复杂时使用，比如编辑场景。

以下子命令在缺少content参数时会调起编辑器：

| 子命令 | 说明         |
| ------ | ------------ |
| add    | 添加任务     |
| mod    | 编辑任务     |
| rmk    | 添加备注     |
| ddl    | 设置任务期限 |
| rmd    | 设置提醒时间 |

编辑器的配置为todo.editor，默认是vim，你可以配置成任何一款自己喜欢的编辑器。

# TODO

| 事项                                                                          | 状态   |
| ----------------------------------------------------------------------------- | ------ |
| 取消file和redis两种存储模式的切换，改为真正的Local First，redis只用来同步数据 | 已完成 |
| 添加新的数据同步方式(AWS S3)                                                  | 已完成 |
| 改stat为info，丰富输出信息，比如工作目录路径等                                | 未开始 |
| 取消urge，以及隐藏的优先级相关逻辑，改为up和down                              | 未开始 |
| 取消server，改为采用后台进程的形式，或建立单独的服务端工程                    | 未开始 |
| 添加新的任务状态Abandoned                                                     | 未开始 |
| 导出任务列表为Markdown                                                        | 未开始 |
| 导出任务列表为表格                                                            | 未开始 |
| 开发配套的移动端程序                                                          | 未开始 |
