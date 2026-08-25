基于 Go 实现的 WasteGen 项目，一款后端服务，完成垃圾焚烧发电 DCS 集控、给料燃烧汽水烟气与出灰联动控制及运行审计。

## 功能

- 焚烧炉与炉区命名空间、测温分区管理
- 炉排给料、抓斗投料、出灰与卸灰顺序控制
- 汽包压力调节、主汽门设定与汽机并网
- 烟气含氧量标定、脱硝喷氨与在线排放台账
- 蒸汽累计核算与炉排燃尽率判定
- 日运行配额与全量运行审计
- 文件型持久化与 JSON API 控制台

## 构建与运行

本地构建：

```bash
go build -mod=vendor -o wastegen.exe ./cmd/wastegen
```

启动服务：

```bash
./wastegen.exe -addr 127.0.0.1:8901 -root ./data
```

健康检查：

```bash
curl http://127.0.0.1:8901/api/health
```

Docker 构建：

```bash
bash build_benzhi_docker.sh
docker run -p 8901:8901 wastegen:latest
```

## 数据目录

运行数据保存在 `-root` 指定的目录下，包括给料记录、标定基准、排放台账、蒸汽累计、分区快照、配额与审计事件。
