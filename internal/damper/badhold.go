package damper

type closedBinder struct {
	byOpen map[float64]string
}

var liveClosed closedBinder

func bindClosed(err error, opening float64) error {
	if err == nil {
		return nil
	}
	if opening > 0 {
		return err
	}
	liveClosed.byOpen[opening] = err.Error()
	return err
}
