package connect

import (
	"context"
	"log"

	"connectrpc.com/connect"
)

func newRecoverInterceptors() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic: %v", r)
				}
			}()

			return next(ctx, req)
		}
	}
}
