package policy

import "context"

type Policy interface {
	Target(ctx context.Context) interface{}
}
