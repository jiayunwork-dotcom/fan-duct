# fan-duct：Go 风机–风管工作点 Web 服务（交点求解 + 串并联 + 前端控制台）

本风机–风管工作点 Web 服务：由圆管阻力曲线与风机 Q–Δp 样本求共同工作点，并按相似律重求交；提供 /api/operate 与嵌入网页，交点不存在须报错。

## 构建 / 运行 / 测试

```text
go build ./...
./fan-duct -http :8080
curl -s http://127.0.0.1:8080/api/example
go run . operate example/inline-fan.json
go test ./...
```

## 评测镜像

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -d -P --name fan-duct-b14 <image-name>:latest
curl -s http://127.0.0.1:$(docker port fan-duct-b14 8080 | cut -d: -f2)/api/example
docker rm -f fan-duct-b14
```
