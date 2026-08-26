package server

import (
	"io"
	"net/http"

	"fan-duct/internal/solve"
)

type OperateResponse struct {
	Flow       float64 `json:"flow_m3s"`
	Velocity   float64 `json:"velocity_ms"`
	Pressure   float64 `json:"pressure_pa"`
	DuctDp     float64 `json:"duct_dp_pa"`
	FanDp      float64 `json:"fan_dp_pa"`
	Residual   float64 `json:"residual_pa"`
	Iterations int     `json:"iterations"`
	Ratio      float64 `json:"speed_ratio,omitempty"`
	RespeedQ   float64 `json:"respeed_flow_m3s,omitempty"`
	RespeedDp  float64 `json:"respeed_pressure_pa,omitempty"`
	NaiveQ     float64 `json:"naive_flow_m3s,omitempty"`
}

func handleOperate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	raw, err := readBody(r, w)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	in, err := solve.ParseInput(raw)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	b, err := in.Build()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	op, err := b.OperatingPoint()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	out := OperateResponse{
		Flow:       op.Flow,
		Velocity:   op.Velocity,
		Pressure:   op.Pressure,
		DuctDp:     op.DuctDp,
		FanDp:      op.FanDp,
		Residual:   op.Residual,
		Iterations: op.Iterations,
	}
	out.Flow, out.Pressure = HoldOpAPI(out.Flow, out.Pressure)
	if in.NewSpeed != nil {
		rr, err := in.RespeedToRPM(b, in.NewSpeed.RPM)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		out.Ratio = rr.Ratio
		out.RespeedQ = rr.Respeeded.Flow
		out.RespeedDp = rr.Respeeded.Pressure
		out.NaiveQ = rr.NaiveFlow
	}
	writeJSON(w, http.StatusOK, out)
}

func readBody(r *http.Request, w http.ResponseWriter) ([]byte, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	return io.ReadAll(r.Body)
}
