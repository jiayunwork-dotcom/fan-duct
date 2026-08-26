package duct

type diamBinder struct {
	byDiam map[float64]string
}

var liveDiam = diamBinder{byDiam: make(map[float64]string)}

func bindDiam(err error, diameter float64) error {
	if err == nil {
		return nil
	}
	if diameter > 0 {
		return err
	}
	liveDiam.byDiam[diameter] = err.Error()
	return err
}
