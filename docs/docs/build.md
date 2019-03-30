# 构建配置

## 配置说明

通过项目列表中的 `构建设置` 来编辑构建脚本。脚本需在上线单中手动触发，系统会使用 `/bin/bash -c {脚本文件}` 执行。

## 构建脚本支持的变量

- **${env_workspace}**

代码仓库本地副本目录

- **${env_pack_file}**

打包文件绝对地址，构建完成后将需要部署到线上的代码打包到此文件中，必须使用 `tar -zcf` 命令进行打包。

部署模块会将此压缩包分发到目标主机并解压缩到指定目录，请按照要求打包，否则会部署失败。

## 简单构建示例

```shell
cd ${env_workspace}
tar --exclude='.git' -zcvf ${env_pack_file} *
```

## Laravel构建示例

```shell
cd ${env_workspace}
composer install
tar --exclude='.git' -zcvf ${env_pack_file} *
```

## 前端项目构建示例

```shell
cd ${env_workspace}
yarn
yarn build
cd ./dist
tar -zcvf ${env_pack_file} *
```

## Java项目构建示例

待补充

## Syncd构建示例

```shell
cd ${env_workspace}
make
cd ./output
tar -zcvf ${env_pack_file} *
```

[filename](include/footer.md ':include')