# light-Y2B

一个基于 FFmpeg 的轻量级的油管转播程序

### 使用场景

此程序是为了提供简单的转播需求（不需要打码，加评论区等）设计的，通过直接转发直播流的方式节省视频编码性能，可直接在云服务器（VPS）、NAS、甚至手机终端模拟器等低性能设备上使用。
如果对转播要求比较高，建议使用传统的 OBS 转播方式或 [Alice-Liveman](https://github.com/nekoteaparty/Alice-LiveMan) (爱丽丝)，但需要强大的主机性能用于视频编码。

### 安装方法

需要预先安装 [FFmpeg](https://ffmpeg.org/download.html) 并且添加到环境变量中 (或手动在配置文件中设置路径)，安装方法可自行百度

1. 自行编译

需要安装：[golang](https://go.dev/dl/)

```bash
git clone https://github.com/ciisaichan/light-Y2B
cd light-Y2B
go mod tidy
go build -ldflags="-s -w" .
```

如果成功会在目录下生成一个可执行文件：`light-Y2B`（或 .exe）

2. 下载已编译版本

请看：[releases](https://github.com/ciisaichan/light-Y2B/releases/)。

### 使用方法

1. 在 [releases](https://github.com/ciisaichan/light-Y2B/releases/) 中下载 config.yaml 配置文件，并将其放置于程序所在目录下。
2. 用记事本或其他文本编辑器打开 config.yaml 文件，修改需要的配置项，然后保存。
3. 直接运行 `light-Y2B` 程序即可。

### 可用命令参数

```
Usage of light-Y2B:
  -c-name string
        配置文件名，不指定默认 config (文件后缀需要为.yaml) (default "config")
  -c-path string
        配置文件目录，不指定默认 当前目录 (default "./")
  -h    显示帮助信息
```

### 注意事项

1. 配置文件中设置的代理类型需要为 HTTP（SOCKS 等需要转换为 HTTP）。
2. 如果主机环境可以直接连接油管，不需要设置代理，配置文件中的 `Proxy:` 选项后面直接留空即可。
