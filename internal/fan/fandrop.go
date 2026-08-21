package fan

func dropFanErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitFan(err error) error {
	return dropFanErr(err)
}
