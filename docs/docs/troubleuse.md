# 使用问题列表

#### 部署失败：[cmd] $ echo 'packfile empty' && exit 1 X

原因：构建脚本中未打包出tgz文件，检查是否使用 `tar -zcf` 命令进行打包。

[filename](include/footer.md ':include')