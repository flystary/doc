package kcp

import (
	"context"
	"tunnel/internal/config"
)

type KCP struct {
	ctx    context.Context
	Enable bool
}

func NewKCP(ctx context.Context, config config.Config) (r KCP, err error) {
	return KCP{
		ctx:    ctx,
		Enable: config.GetBool("enable_kcp", false),
	}, nil
}
