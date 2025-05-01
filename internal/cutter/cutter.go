package cutter

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type Cutter struct {
	cutSize int
}

func NewCutter(cutSize int) (c *Cutter, err error) {
	if cutSize <= 0 {
		err = fmt.Errorf("Неверный размер кэша")
		return
	}
	c = &Cutter{
		cutSize: cutSize,
	}

	return
}

func (c *Cutter) Cut(fullUrl string) (cutUrl string) {
	h := md5.New()
	h.Write([]byte(fullUrl))

	cutUrl = hex.EncodeToString(h.Sum(nil))[:c.cutSize]
	return
}
