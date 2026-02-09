# coupang

用于调用 Coupang Open API 的 Go 命令行工具示例，包含基础的订单与商品查询逻辑。

## 环境要求

- Go 版本：建议 Go 1.20 及以上
- 已在系统中安装 Git（仅在你需要管理代码时）

## 配置说明

程序通过 YAML 配置文件加载 Coupang 相关的访问密钥，默认读取当前目录下的 `config.yaml`。

示例配置文件：[`config_example.yaml`](file:///Users/bre/workspace/jay/coupang/config_example.yaml)

```yaml
coupang:
  api_key: your_api_key
  secret_key: your_secret_key
  vendor_id: your_vendor_id
```

使用方式：

- 将 `config_example.yaml` 复制为 `config.yaml`
- 用真实的 `api_key`、`secret_key`、`vendor_id` 替换示例占位符
- 如需使用自定义路径，可以通过 `-f` 参数指定：

```bash
./coupang -f /path/to/your_config.yaml
```

## 构建说明

项目已经提供构建脚本用于编译 Linux amd64 二进制：

```bash
chmod +x build.sh
./build.sh
```

正常执行后，会在项目根目录生成名为 `coupang` 的可执行二进制文件。

## 运行说明

### 前台运行

```bash
./coupang -f config.yaml
```

### 后台运行并输出日志到文件

当前仓库中原有的运行方式为：

```bash
nohup ./coupang >> coupang.log 2>&1 &
```

- `nohup`：在退出终端后程序仍保持运行
- `>> coupang.log`：将标准输出追加写入到 `coupang.log`
- `2>&1`：将标准错误重定向到标准输出，一并写入日志文件
- `&`：在后台运行进程

## 日志查看

后台运行时，可以使用以下命令实时查看日志：

```bash
tail -f coupang.log
```

如需只查看最后若干行：

```bash
tail -n 100 coupang.log
```

## 当前示例逻辑简介

在 [`main.go`](file:///Users/bre/workspace/jay/coupang/main.go) 中，目前示例代码主要完成以下工作：

- 通过配置文件初始化 `Coupang` 客户端
- 使用固定的订单号调用 `GetOrderByOrderId` 查询订单详情
- 遍历订单行，进一步调用 `GetSellerProduct` 查询对应商品信息

你可以根据自己的业务需求修改 `main.go` 中的示例逻辑，使用已经封装好的客户端方法进行更多 API 调用。
