package entry

import (
	"context"
)

type Entry interface {
	Run() error
	Stop(ctx context.Context) error
}
