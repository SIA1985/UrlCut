package storage

import (
	"container/list"
	"errors"
)

/*Сохранение недавних ссылок для более быстрого доступа*/

type CacheElement struct {
	cutUrl  string
	fullUrl string
}

type Cache struct {
	data    *list.List
	maxSize int
}

func NewCache(maxSize int) (c *Cache) {
	c = &Cache{}

	c.maxSize = maxSize
	c.data = list.New()
	return c
}

func (c *Cache) Add(cutUrl string, fullUrl string) {
	if c.data.Len() >= c.maxSize {
		c.data.Remove(c.data.Front())
	}

	c.data.PushBack(CacheElement{cutUrl, fullUrl})
}

func (c *Cache) Contains(cutUrl string) bool {
	for e := c.data.Front(); e != nil; e = e.Next() {
		ce := e.Value.(CacheElement)

		if ce.cutUrl == cutUrl {
			return true
		}
	}

	return false
}

func (c *Cache) Get(cutUrl string) (string, error) {
	for e := c.data.Front(); e != nil; e = e.Next() {
		ce := e.Value.(CacheElement)

		if ce.cutUrl == cutUrl {
			return ce.fullUrl, nil
		}
	}

	return "", errors.New("Элемент по ключу '" + cutUrl + "' не найден!")
}
