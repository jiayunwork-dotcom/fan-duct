# fan-duct — 风机–风管工作点核算

fan-duct 是一个命令行通风核算工具：给定一段圆管风道的几何与阻力参数（长度、内径、沿程摩阻系数、局部阻力系数和、空气密度）以及一台风机的 Q–Δp 样本曲线，它解出风机与风管的共同工作点——流量 Q、断面流速 V 与压升 Δp，并按相似律把风机缩放到新转速后重新求交。能力边界是单台风机配单根圆管：管阻用 Δp=(fL/D+ΣK)·ρV²/2 在同一截面换算 Q=V·πD²/4，风机曲线支持分段线性与最小二乘二次拟合，Q 超出样本范围时按显式策略报错或外推。它不做管网平差、不做多环求解，也不做阀门水击瞬态。

## 用法

```text
go run . operate example/inline-fan.json
```

打印该算例（L=50 m、D=0.15 m、f=0.02、ΣK=3、空气 1.205 kg/m³；风机分段线性样本 0~0.6 m³/s、零流压 1450 Pa）的工作点：交点 Q≈0.265 m³/s、V≈14.99 m/s、Δp≈1308 Pa，并给出转速从 1450 提到 1595 rpm（1.1 倍）后按相似律重求交的 Q2≈0.291 m³/s、Δp2≈1583 Pa。交点流量落在风机样本点 0.15 与 0.30 m³/s 之间。

其他用法：

```text
go run . operate example/inline-fan.json --new-speed 1700  # 用命令行覆盖目标转速
go run . operate example/inline-fan.json --compact         # 单行摘要
go run . operate example/inline-fan.json --sensitivity     # 附加 +10% 风管参数灵敏度
go run . help
```

## 输入格式

`operate` 读一个 JSON 文件，字段：

```json
{
  "duct": { "length": 50, "diameter": 0.15, "friction": 0.02,
            "lossCoeff": 3, "density": 1.205, "viscosity": 1.82e-5,
            "roughness": 0 },
  "fan": { "points": [ {"q": 0, "dp": 1450}, {"q": 0.3, "dp": 1280} ],
           "efficiency": [0.6, 0.66],
           "fit": "polyline", "extrapolate": "error" },
  "speed": { "rpm": 1450 },
  "newSpeed": { "rpm": 1595 }
}
```

- `duct.friction` 缺省为 0，表示按雷诺数自动取摩阻（层流 64/Re，湍流 Blasius；`duct.roughness` > 0 时改用 Swamee–Jain）；`density` 缺省 1.205 kg/m³，`viscosity` 缺省 1.82e-5 Pa·s。
- `fan.fit` 取 `polyline`（默认）或 `quadratic`；`fan.extrapolate` 取 `error`（默认，越界报错）、`linear` 或 `quadratic`。
- `newSpeed` 与 `speed` 成对出现时输出新转速交点；也可用 `--new-speed` 覆盖。

## 关键约定

- **管阻与动压分开**：管阻 Δp=(fL/D+ΣK)·ρV²/2 是风机全压需要克服的量，动压 ρV²/2 只是其中的中间项，二者不会被混用。Q 与 V 永远换算自同一个截面积 πD²/4。
- **工作点定义**：Δp_fan(Q)=Δp_duct(Q)。用括号法在 [0, ∞) 上找变号区间后二分求根；Q=0 时管阻恒为 0，风机取曲线零流点。零流压非正、全程无交点或交点超出外推范围都会报错。
- **相似律**：转速比 r=N2/N1，缩放后的风机曲线为 Q′=Q·r、Δp′=Δp·r²（功率 ∝N³），必须用缩放后的曲线与同一管阻**重新求交**；报告中的 `naive Q1*ratio` 只是对照参考值，不是解。
- **非法输入**：直径 ≤0、长度 <0、密度 ≤0、空风机曲线、样本流量不递增、未知 JSON 字段、流量越界而外推被禁等，一律写 stderr 并以非零退出码结束。

## 构建与测试

```text
go build ./...
go test ./...
```

纯标准库，无第三方依赖，无 cgo。

## 许可

MIT，见 [LICENSE](./LICENSE)。
