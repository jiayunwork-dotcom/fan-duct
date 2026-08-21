# fan-duct — 风机–风管工作点核算

fan-duct 是风机–风管工作点核算命令行工具：给定圆管风道几何与阻力参数以及风机 Q–Δp 样本曲线，解出共同工作点的流量 Q、流速 V 与压升 Δp，并按相似律把风机缩放到新转速后重新求交。纯标准库，无网络依赖，无 cgo。

## 构建 / 运行 / 测试

```text
go build ./...
go run . operate example/inline-fan.json   # CLI：核算算例并打印 Q、Δp、V 与可选新转速交点
go test ./...                              # 单元测试（duct / fan / solve）
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```

进容器后运行 `go build ./... && go test ./...`，再用 `go run . operate example/inline-fan.json` 验证 CLI 输出 Q、Δp、V。
