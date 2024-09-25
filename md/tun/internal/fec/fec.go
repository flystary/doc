package fec

import (
	"context"
	"fmt"
	"tunnel/internal/config"
)

const (
	FEC_DISABLE = iota
	FEC_AUTO
	FEC_FORCE
)

type FEC struct {
	ctx     context.Context
	FecMode int
	FecSrc  int
	FecRed  int
}

func NewFec(ctx context.Context, config config.Config) (r FEC, err error) {
	fec_mode_str := config.GetString("fec_mode", "")
	var fec_mode, fec_src, fec_red int
	fec_src = 0
	fec_red = 0
	if len(fec_mode_str) == 0 {
		err = fmt.Errorf("fec_mode is empty")
		return
	}
	switch fec_mode_str {
	case "disable":
		fec_mode = FEC_DISABLE
	case "auto":
		fec_mode = FEC_AUTO
	case "force":
		fec_mode = FEC_FORCE
		fec_src = 5
		fec_red = 1
	default:
		err = fmt.Errorf("fec_mode is err value:%s", fec_mode_str)
		return
	}

	r = FEC{
		ctx:     ctx,
		FecSrc:  fec_src,
		FecRed:  fec_red,
		FecMode: fec_mode,
	}
	return
}
