package droplet

type Context interface {
	Get(key string) interface{}
	GetString(key string) string
	Set(key string, value interface{})
	SetInput(interface{})
	Input() interface{}
	SetOutput(interface{})
	Output() interface{}
}

type emptyContext struct {
	dict   map[string]interface{}
	input  interface{}
	output interface{}
}

func NewContext() *emptyContext {
	c := &emptyContext{}
	c.dict = make(map[string]interface{})
	return c
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
