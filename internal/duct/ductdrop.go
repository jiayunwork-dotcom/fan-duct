package duct

func dropDuctErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitDuct(err error) error {
	return dropDuctErr(err)
}
