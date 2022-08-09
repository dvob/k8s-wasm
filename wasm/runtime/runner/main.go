package runner

import "context"

type RawRunner interface {
	Run(ctx context.Context, fnName string, input []byte) (output []byte, err error)
}
