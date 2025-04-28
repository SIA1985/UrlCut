package cutter

import (
	"crypto/sha256"
)


type Cutter struct {
	cutSize			int
	cutFunc			func(fullUrl string)(cutUrl string)
}

type CutterOption func(*Cutter)

func WithCutFunc(cutFunc func(fullUrl string)(cutUrl string)) (CutterOption) {
	return func(c *Cutter) {
		c.cutFunc = cutFunc
	}
}

func WithCutSize(cutSize int) (CutterOption) {
	return func(c *Cutter) {
		c.cutSize = cutSize
	}
}

func NewCutter(opts... CutterOption) (c *Cutter, err error) {
	c = &Cutter{
		cutSize: 6,
		cutFunc: func(fullUrl string) (cutUrl string) {
			h := sha256.New()
			h.Write([]byte(fullUrl))

			cutUrl = string(h.Sum(nil))
			return 
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return
}

func (c *Cutter) Cut(fullUrl string) (cutUrl string) {
	return c.cutFunc(fullUrl)[:c.cutSize]
}