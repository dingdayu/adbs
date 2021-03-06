# adbs

Android Debug Bridge (adb) Remote debugging service.

Read this in other languages: [English](README.md), [简体中文](README.zh-CN.md).

## Feature List 

* ADB Server
* WEB Control
* Remote Control

## Mod

Here is the referenced third party package.

- https://github.com/gin-gonic/gin
- https://github.com/gorilla/websocket
- https://github.com/shogo82148/androidbinary
- https://github.com/kr/pty

## Roadmap

* RUN ADB shell.
* Echo `screencap`.
* Connect or disconnect device.
* Get packages list.
* Push and pull device files.
* Get dir or file stat.
* Input event to device.

## Quick start

### ADB

Please download the corresponding version, [platform-tools](https://developer.android.com/studio/releases/platform-tools).

And extract the path to the environment variable after decompressing it.

Please ensure the correct execution of `adb devices`.

### Download

```shell
git clone https://github.com/dingdayu/adbs

cd adbs

go mod tidy

go run .
```

## TODO

- [X] 写文件时的时间
- [X] 推文件写入协议优化
- [ ] 获取文件信息时的文件类型问题（目录/连接）
- [ ] 设备列表等API接口完成设备选择
- [ ] 提供编译版本
- [ ] 提供 `docker` 镜像

## Reference

The project is affected by the following projects or articles.

- [7.0上截图的问题](https://github.com/mzlogin/awesome-adb/issues/33)
- [Android之高效率截图](https://juejin.im/post/5bab409ef265da0afc2c032e)
- [Read binary stdout data from adb shell?](https://stackoverflow.com/questions/13578416/read-binary-stdout-data-from-adb-shell)
- [Go 内嵌静态资源](http://fuxiaohei.me/2016/10/1/go-binary-embed-asset.html)