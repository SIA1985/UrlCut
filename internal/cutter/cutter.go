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

func (c *Cutter) Cut(fullUrl string) (cutUrl string, err error) {
	if len(fullUrl) == 0 {
		err = fmt.Errorf("пустой url для сокращения")
		return
	}

	h := md5.New()
	h.Write([]byte(fullUrl))

	cutUrl = hex.EncodeToString(h.Sum(nil))[:c.cutSize]
	return
}
