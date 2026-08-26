package server

import (
	"encoding/json"
	"net/http"

	"fan-duct/internal/atm"
	"fan-duct/internal/damper"
	"fan-duct/internal/network"
	"fan-duct/internal/solve"
)

type AtmosphereRequest struct {
	AltitudeM   float64 `json:"altitude_m"`
	RelHumidity float64 `json:"rel_humidity"`
}

type AtmosphereResponse struct {
	AltitudeM    float64 `json:"altitude_m"`
	TemperatureK float64 `json:"temperature_k"`
	PressurePa   float64 `json:"pressure_pa"`
	DensityKgM3  float64 `json:"density_kg_m3"`
	ViscosityPas float64 `json:"viscosity_pa_s"`
	MoistDensity float64 `json:"moist_density_kg_m3,omitempty"`
}

func handleAtmosphere(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req AtmosphereRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	st, mu, _, err := atm.PropertiesAt(req.AltitudeM)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	out := AtmosphereResponse{
		AltitudeM:    st.AltitudeM,
		TemperatureK: st.TemperatureK,
		PressurePa:   st.PressurePa,
		DensityKgM3:  st.DensityKgM3,
		ViscosityPas: mu,
	}
	if req.RelHumidity > 0 {
		moist, err := atm.MoistDensity(st.PressurePa, st.TemperatureK, req.RelHumidity)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		out.MoistDensity = moist
	}
	writeJSON(w, http.StatusOK, out)
}

type NetworkRequest struct {
	Left  json.RawMessage `json:"left"`
	Right json.RawMessage `json:"right"`
	Mode  string          `json:"mode"`
}

type NetworkResponse struct {
	Flow     float64 `json:"flow_m3s"`
	Pressure float64 `json:"pressure_pa"`
	Residual float64 `json:"residual_pa"`
	Mode     string  `json:"mode"`
}

func handleNetwork(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req NetworkRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Mode == "" {
		req.Mode = "series-duct"
	}
	left, err := solve.ParseInput(req.Left)
	if err != nil {
		writeError(w, http.StatusBadRequest, "left: "+err.Error())
		return
	}
	lb, err := left.Build()
	if err != nil {
		writeError(w, http.StatusBadRequest, "left: "+err.Error())
		return
	}
	right, err := solve.ParseInput(req.Right)
	if err != nil {
		writeError(w, http.StatusBadRequest, "right: "+err.Error())
		return
	}
	rb, err := right.Build()
	if err != nil {
		writeError(w, http.StatusBadRequest, "right: "+err.Error())
		return
	}
	var op network.Point
	switch req.Mode {
	case "series-duct":
		drop := network.SeriesDrop{Items: []network.Dropper{lb.Duct, rb.Duct}}
		op, err = network.Intersect(lb.Fan, drop)
	case "parallel-duct":
		drop := network.ParallelDrop{Items: []network.Dropper{lb.Duct, rb.Duct}}
		op, err = network.Intersect(lb.Fan, drop)
	case "series-fan":
		rise := network.SeriesRise{Items: []network.Riser{lb.Fan, rb.Fan}}
		op, err = network.Intersect(rise, lb.Duct)
	case "parallel-fan":
		rise := network.ParallelRise{Items: []network.Riser{lb.Fan, rb.Fan}}
		op, err = network.Intersect(rise, lb.Duct)
	case "damper":
		dm, derr := damper.New(damper.Config{
			Opening:  0.55,
			KOpen:    0.2,
			KShut:    8000,
			Diameter: lb.Duct.Config().Diameter,
			Density:  lb.Duct.Config().Density,
		})
		if derr != nil {
			writeError(w, http.StatusBadRequest, derr.Error())
			return
		}
		drop := network.AddedDrop{Base: lb.Duct, Extra: dm}
		op, err = network.Intersect(lb.Fan, drop)
	default:
		writeError(w, http.StatusBadRequest, "unknown mode "+req.Mode)
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	op.Pressure = HoldSeriesAPI(op.Pressure)
	writeJSON(w, http.StatusOK, NetworkResponse{
		Flow:     op.Flow,
		Pressure: op.Pressure,
		Residual: op.Residual,
		Mode:     req.Mode,
	})
}
