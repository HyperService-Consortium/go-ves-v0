# go-ves

ves-server启动实例

```bash
go run bin/ves-server/run_server.go
```

centered-ves-server启动实例

```bash
go run bin/centered-ves-server/run_server.go
```

启动ves-client

```bash
go run bin/ves-client/main.go
```

进度：

+ 完整地运行一次跨链转账交易
+ 完成 attestation 的上链行为
+ 完成 merkle proof 的上链行为
+ 完整地运行一次跨链合约交易
+ <del>完整地</del>完成 merkle proof 的 value 抽取功能
+ 将 transaction proof 加入 workflow

这周 TODO List:

+ 将 data proof 加入 workflow

+ 考虑设计UI界面

+ 考虑将NSB部署到docker集群

+ 第一阶段完成，优化代码，压力测试，接口测试

下周 TODO List:

+ 将ISC功能完善化

+ 加入 validate 逻辑
