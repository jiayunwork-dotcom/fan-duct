# fan-duct

fan-duct 核算单台风机与单根圆管风道的共同工作点。输入圆管长度、内径、沿程摩阻（或相对粗糙度）、局部阻力系数和、空气密度，以及风机 Q–Δp 样本曲线，后端解 Δp_fan(Q)=Δp_duct(Q)，得到流量 Q、断面流速 V 与压升 Δp；可选按相似律把风机缩放到新转速后重新求交。还可以把两段风管串/并联、两台风机串/并联，或叠风阀局部阻力后再求交。它不做管网平差，也不做多环水力。

## 用法

最直接的一条命令（预置直管算例，L=50 m、D=0.15 m、f=0.02、ΣK=3）：

```text
go run . operate example/inline-fan.json
```

打印该算例的工作点：交点 Q≈0.265 m³/s、V≈14.99 m/s、Δp≈1308 Pa，以及转速从 1450 提到 1595 rpm 后重求交的 Q2 与 Δp2。其它入口：

```text
fan-duct operate <算例.json> [--new-speed rpm] [--compact] [--sensitivity]
fan-duct -http :8080
fan-duct version / help
```

浏览器打开 http://localhost:8080 可改参数并打 `/api/operate`。容器默认拉起同一服务。

## 算例格式

`example/inline-fan.json` 可直接运行：

```json
{
  "duct": { "length": 50, "diameter": 0.15, "friction": 0.02,
            "lossCoeff": 3, "density": 1.205, "viscosity": 1.82e-5 },
  "fan": { "points": [ {"q": 0, "dp": 1450}, {"q": 0.3, "dp": 1280} ],
           "fit": "polyline", "extrapolate": "error" },
  "speed": { "rpm": 1450 },
  "newSpeed": { "rpm": 1595 }
}
```

- `duct.friction` 缺省 0 时按雷诺数取摩阻（层流 64/Re，湍流 Blasius；`roughness`>0 时改 Swamee–Jain）。
- `fan.fit` 取 polyline 或 quadratic；`fan.extrapolate` 取 error、linear 或 quadratic。
- `newSpeed` 与 `speed` 成对出现时输出新转速交点。

## 关键约定

- 管阻 Δp=(fL/D+ΣK)·ρV²/2，Q 与 V 永远换算自同一截面积 πD²/4。
- 工作点是 Δp_fan(Q)=Δp_duct(Q)。零流压非正或全程无交点会报错，不会静默给数。
- 相似律 r=N2/N1 缩放整条风机曲线（Q∝N、Δp∝N²）后与同一管阻重新求交；报告里的 naive Q1·r 只是对照。
- 并联风管在同一压降下分流；并联风机在同一压降下流量相加。弱风机零流压低于系统压时该支路流量为 0，不做倒灌。
- ISA 密度随高度下降；湿空气密度低于同温同压干空气。

## API

同进程 Web 服务：

- `GET /health` 与 `GET /api/health` — 探活
- `GET /api/example` — 预置算例 JSON
- `POST /api/operate` — 工作点（请求体即算例 JSON）
- `POST /api/atmosphere` — ISA 高度处的 T、P、ρ、μ
- `POST /api/network` — 串/并联风管或风机后再求交

失败路径与 CLI 一致：非法参数返回 `{"error":"…"}` 与 400。

## 构建与测试

```text
go build ./...
go test ./...
go run . -http :8080
```

纯标准库，无第三方依赖，无 cgo。

## 许可

MIT，见 [LICENSE](./LICENSE)。
