"use strict";

const $ = (id) => document.getElementById(id);

function payload() {
  return {
    duct: {
      length: Number($("length").value),
      diameter: Number($("diameter").value),
      friction: Number($("friction").value),
      lossCoeff: Number($("loss").value),
      density: Number($("density").value),
    },
    fan: {
      points: JSON.parse($("points").value),
      fit: "polyline",
      extrapolate: "error",
    },
    speed: { rpm: Number($("rpm1").value) },
    newSpeed: { rpm: Number($("rpm2").value) },
  };
}

function showError(msg) {
  $("err-text").textContent = msg;
  $("err-panel").hidden = false;
}

function hideError() {
  $("err-panel").hidden = true;
}

function fillTable(id, rows) {
  const tb = $(id);
  tb.innerHTML = "";
  for (const [k, v] of rows) {
    const tr = document.createElement("tr");
    tr.innerHTML = "<td>" + k + "</td><td>" + v + "</td>";
    tb.appendChild(tr);
  }
}

async function postJSON(path, body) {
  const resp = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const data = await resp.json();
  if (!resp.ok) {
    throw new Error(data.error || "HTTP " + resp.status);
  }
  return data;
}

async function loadExample() {
  hideError();
  $("hint").textContent = "正在读取 /api/example …";
  try {
    const resp = await fetch("/api/example");
    const data = await resp.json();
    if (!resp.ok) throw new Error(data.error || "示例加载失败");
    $("length").value = data.duct.length;
    $("diameter").value = data.duct.diameter;
    $("friction").value = data.duct.friction;
    $("loss").value = data.duct.lossCoeff;
    $("density").value = data.duct.density;
    $("rpm1").value = data.speed.rpm;
    $("rpm2").value = data.newSpeed.rpm;
    $("points").value = JSON.stringify(data.fan.points, null, 2);
    $("hint").textContent = "已加载直管算例。点“求工作点”看交点。";
  } catch (e) {
    showError(String(e));
  }
}

async function runOperate() {
  hideError();
  try {
    const out = await postJSON("/api/operate", payload());
    $("out-panel").hidden = false;
    fillTable("out-table", [
      ["流量 Q (m³/s)", out.flow_m3s.toPrecision(6)],
      ["流速 V (m/s)", out.velocity_ms.toPrecision(6)],
      ["压升 Δp (Pa)", out.pressure_pa.toPrecision(6)],
      ["残差 (Pa)", out.residual_pa.toPrecision(4)],
      ["新转速流量", out.respeed_flow_m3s ? out.respeed_flow_m3s.toPrecision(6) : "—"],
      ["新转速压升", out.respeed_pressure_pa ? out.respeed_pressure_pa.toPrecision(6) : "—"],
    ]);
    $("hint").textContent = "交点由后端求解，前端只展示返回值。";
  } catch (e) {
    showError(String(e));
  }
}

async function runAtm() {
  hideError();
  try {
    const out = await postJSON("/api/atmosphere", { altitude_m: 0, rel_humidity: 0.5 });
    $("atm-panel").hidden = false;
    fillTable("atm-table", [
      ["温度 K", out.temperature_k.toPrecision(6)],
      ["气压 Pa", out.pressure_pa.toPrecision(7)],
      ["干空气密度", out.density_kg_m3.toPrecision(5)],
      ["湿空气密度", out.moist_density_kg_m3.toPrecision(5)],
    ]);
  } catch (e) {
    showError(String(e));
  }
}

$("btn-example").addEventListener("click", loadExample);
$("btn-operate").addEventListener("click", runOperate);
$("btn-atm").addEventListener("click", runAtm);
