package solve

import "context"

func abortFresh() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}
