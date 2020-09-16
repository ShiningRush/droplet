package droplet

import "context"

type Context interface {
	Context() context.Context
	SetContext(context.Context)
	Get(key string) interface{}
	GetString(key string) string
	Set(key string, value interface{})
	SetInput(interface{})
	Input() interface{}
	SetOutput(interface{})
	Output() interface{}
	SetPath(path string)
	Path() string
}

type emptyContext struct {
	cxt    context.Context
	dict   map[string]interface{}
	input  interface{}
	output interface{}
	path   string
}

func NewContext() *emptyContext {
	c := &emptyContext{}
	c.dict = make(map[string]interface{})
	c.cxt = context.TODO()
	return c
}

func (c *emptyContext) Context() context.Context {
	return c.cxt
}

func (c *emptyContext) SetContext(cxt context.Context) {
	c.cxt = cxt
}

func (c *emptyContext) Set(key string, value interface{}) {
	c.dict[key] = value
}

func (c *emptyContext) Get(key string) interface{} {
	rs, ok := c.dict[key]
	if !ok {
		return nil
	}

	return rs
}

func (c *emptyContext) GetString(key string) string {
	rs, ok := c.dict[key]

	s := ""
	s, ok = rs.(string)
	if !ok {
		return ""
	}

	return s
}

func (c *emptyContext) SetInput(input interface{}) {
	c.input = input
}

func (c *emptyContext) Input() interface{} {
	return c.input
}

func (c *emptyContext) SetOutput(output interface{}) {
	c.output = output
}

func (c *emptyContext) Output() interface{} {
	return c.output
}

func (c *emptyContext) SetPath(path string) {
	c.path = path
}

func (c *emptyContext) Path() string {
	return c.path
}
