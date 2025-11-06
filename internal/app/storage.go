package app

import "context"

type Storage interface {
	Save(ctx context.Context, id, url string) error
	Get(ctx context.Context, id string) (string, bool)
}
