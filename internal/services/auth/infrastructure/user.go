package infrastructure

import "context"

type User interface {
	CheckNamePassword(
		ctx context.Context,
		name, password string,
	) (int32, error)
}
