# 运行原理

## 原理

Syncd通过 `git-ssh` 方式创建仓库本地副本，通过运行项目构建脚本来生成可发布的代码包，使用 `scp` 串行/并行地向生产服务器分发并解包部署。

Syncd与生产服务器通过ssh交互，因此Syncd运行用户的 ssh-key 必须加入到目标主机的用户ssh-key信任列表中。

<img class="app-image" src="assets/img/syncd-principle.png" width="650px;">

## Git

Syncd目前仅支持Git仓库部署，并且通过 `git-ssh` 方式从仓库中拉取指定tag(分支)代码。

## 构建

部署过程中，在拉取git仓库到本地后会运行项目的自定义构建脚本，编译成可上线的软件包。

在这一环节中可运行:
- 单元测试，如 `go test`, `php phpunit`
- 下载依赖，如 `go: glide install`, `php: composer install`
- 编译软件包，如 `js: npm build`, `go: go build xx.go`, `java: mvn ...`

## 分发

通过 `scp` 命令分发软件包到各机房生产服务器的临时目录, 远程执行 `pre-deploy` 配置的命令, 执行完毕后解压缩软件包到目标目录，然后执行 `post-deploy` 命令。

分发上线过程是集群之间串行执行，集群内部可串行可并行。

## SSH信任

生产服务器与部署服务器之间通过ssh-key建立信任，配置方法请参考 [秘钥配置](server.md?id=秘钥配置) 

[filename](include/footer.md ':include')