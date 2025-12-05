package interfaces

import "context"

type Cache interface {
	Set(
		ctx context.Context,
		key string,
		val string,
	) error
	Get(
		ctx context.Context,
		key string,
	) (string, error)
}
