package core

import (
	"errors"
	"io"
)

type CStream struct {
	MainStream io.ReadWriteCloser
	Hold       io.Closer
}

type CReadWriteCloser interface {
	io.ReadWriteCloser
}

func (c *CStream) Read(p []byte) (n int, err error) {
	if c.MainStream != nil {
		return c.MainStream.Read(p)
	}
	return 0, errors.New("MainStream is null")
}

func (c *CStream) Write(p []byte) (n int, err error) {
	if c.MainStream != nil {
		return c.MainStream.Write(p)
	}
	return 0, errors.New("MainStream is null")
}

func (c *CStream) Close() error {

	if c.MainStream != nil {
		c.MainStream.Close()
	}

	if c.Hold != nil {
		c.Hold.Close()
	}

	return nil
}
