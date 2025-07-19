package database

import (
	"context"

	"github.com/risingwavelabs/eris"
)

func DummyQuery(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	res, err := model.DummyQuery(ctx)
	if err != nil {
		return 0, eris.Wrapf(err, "failed to query dummy")
	}

	return res, nil
}
