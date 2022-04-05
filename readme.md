# todo

github.com/justinas/alice
github.com/zeromicro/go-zero/core/threading 开 goroutines
github.com/zeromicro/go-zero/core/syncx  限制请求数
断路器

# goctl

- `goctl api go -api *.api -dir ../ -style goZero`
  - `api` 通过 .api 文件生成一个 api 服务，`goctl api -h` 查看帮助文档
  - `go` 指定生成的代码类型，`goctl api go -h` 查看帮助文档
  - `-api` 指定 api 源文件
  - `-dir` 指定生成文件的路径
  - `-style` 指定生成的文件名格式，见 https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md

- `goctl rpc protoc *.proto --go_out=../ --go-grpc_out=../  --zrpc_out=../ --style=goZero`
  - 该命令可以看作 `goctl rpc $protoc_command --zrpc_out=$output` 摸板，即中间是 protoc 命令相关
  - `rpc` 通过 .proto 文件成成一个 rpc 服务
  - `protoc x.proto --go_out=xxx --go-grpc_out=xxxx` 生成 grpc 代码
    - --go_out 与 --go-grpc_out 生成的最终目录必须一致
  - `--zrpc_out` 指定 zrpc 的输出路径
    - --go_out & --go-grpc_out 和 --zrpc_out 的生成的最终目录必须不为同一目录，否则pb.go和_grpc.pb.go就与main函数同级了，这是不允许的。

- `goctl docker -go user.go -port 9009`
  - `docker` 生成 DockerFile
  - `-go` 包含了 main 函数的文件
  - `-port` 指定服务暴露的端口
  - 注意：每个服务生成的 Dockerfile 要逐个复制到根目录（有 go.mod 的目录）下，再 `docker build -t zero-demo-api:v1 .` 构建镜像！！

- `goctl kube deploy -name user-api -namespace zero-demo -image user-api:v1.0 -o user-api.yaml -port 9019`
  - `kube deploy` 生成 kubernetes 的 yaml 文件
  - `-name` 指定 deployment 的名称
  - `-namespace` 指定 deployment 的 namespace
  - `-image` 指定 deployment 使用的镜像
  - `-o` 指定生成的 yaml 文件名
  - `-port` 指定 deployment 的端口


# api

[api 语法介绍](https://go-zero.dev/cn/api-grammar.html)

定义的 type 用于绑定参数时，仅支持四种 tag：json、form、path、header，见官方文档。 

生成 api 业务代码：`goctl api go -api *.api -dir ../  -style goZero`

# model

生成 model

- `goctl model mysql ddl -src="./*.sql" -dir="./sql/model" -c` 根据 DDL 文件生成
- `goctl model mysql datasource -url="user:password@tcp(127.0.0.1:3306)/database" -table="*"  -dir="./model" -style goZero -c` 通过连接到数据库上来生成

`-c` 是指生成的代码中会优先从缓存中查询数据，然后再走数据库。

是否使用缓存生成的代码不同：

不使用缓存，代码中使用 sqlx 直接访问数据库，且 model 的构造函数仅需传入数据库连接信息

使用缓存，代码中使用 sqlc（对 sqlx 和 缓存组件包装） 查询数据，且 model 的构造函数需要传入数据库、缓存组件连接信息

goctl <= 1.3.3，当使用缓存，且表中没有唯一索引时，会有一个 bug，

如果查询数据库没有查到，会往缓存中针对该 key 插入一条有过期时间的缓存 "*" 作为占位符，这是为了防止缓存击穿！

> 再次查询同一个 key，会从缓存中读到 "*"，着对应了 errPlaceholder 的错误类型！

如果第一次查询为空，向缓存中插入记录，即使在过期时间内向数据库中插入了记录，查询也是查不到记录的。因为查到 * 以后是当作数据库中也未查到处理的，并不会再次查询数据库，见 go-zero 源码中 core/stores/cache/cachenode.go 的 doTake()。

如果表中有唯一索引（mobile）、主键索引（id），则生成的 FindOne 函数会有根据 id 和根据 mobile 两个函数，而 FindOneByMobile 函数的逻辑是，先根据 mobile 从缓存中找到对应缓存 key 的值（保存的是主键索引 id），然后根据这个 id 再去缓存中查找数据。类似于数据库的回表操作。

