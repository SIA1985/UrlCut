package cutter

import (
	"crypto/sha256"
)


type Cutter struct {
	cutSize			int
}

type CutterOption func(*Cutter)

func WithCutSize(cutSize int) (CutterOption) {
	return func(c *Cutter) {
		c.cutSize = cutSize
	}
}

func NewCutter(opts... CutterOption) (c *Cutter, err error) {
	c = &Cutter{
		cutSize: 6,
	}

	for _, opt := range opts {
		opt(c)
	}

	return
}

func (c *Cutter) Cut(fullUrl string) (cutUrl string) {
	h := sha256.New()
	h.Write([]byte(fullUrl))
  
	cutUrl = string(h.Sum(nil))[:c.cutSize]
	return 
}