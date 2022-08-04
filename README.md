# droplet
[![Go Report Card](https://goreportcard.com/badge/github.com/shiningrush/droplet)](https://goreportcard.com/report/github.com/shiningrush/droplet)
[![codecov](https://codecov.io/gh/ShiningRush/droplet/branch/master/graph/badge.svg?token=YL8PzEOyD7)](https://codecov.io/gh/ShiningRush/droplet)

Decouple the service access layer so that the business only cares about input and output.

## Background

When you write a http web server, you will see code like below everywhere.
Such as gin:
```go
func consumerCreate(c *gin.Context) {
	// read header
    requestId, _ := c.Get("X-Request-Id")
    param, exist := c.Get("requestBody")
    
    u4 := uuid.NewV4()
    
    if !exist || len(param.([]byte)) < 1 {
        err := errno.New(errno.InvalidParam)
        logger.WithField(conf.RequestId, requestId).Error(err.ErrorDetail())
        
        // write response
        c.AbortWithStatusJSON(http.StatusBadRequest, err.Response())
        return
    }
    
    if err := service.ConsumerCreate(param, u4.String()); err != nil {
    	// handler error
        handleServiceError(c, requestId, err)
        return
    }
    
    // write response
    c.JSON(http.StatusOK, errno.Succeed())
}
```
Wow, All of those codes is not related with business logic, right?
So is here anyway could let us just care what application need and what application should return?
That is why droplet was born.

Let We have a look at the code after used droplet:
```go
type JsonInput struct {
    ID    string   `auto_read:"id,path" json:"id"`
    User  string   `auto_read:"user,header" json:"user"`
    IPs   []string `json:"ips"`
    Count int      `json:"count"`
    Body  []byte   `auto_read:"@body"`
}

func JsonInputDo(ctx droplet.Context) (interface{}, error) {
    input := ctx.Input().(*JsonInput)
    return input, nil
}
```
Here are things droplet help you to do:
- Automatically read parameter from request by you input's tag
- Check returns result and error, check error and convert it to a pre-define data structure
- Write result to response

Droplet is not only that, please keep reading to find more.

## Concept

Droplet is designed to work between the service access layer and the application layer,
and is an intermediate layer framework. It provides a pipeline mechanism based on the reified Struct,
and provides a lot of convenient middleware based on this.




## Intall

```
go get github.com/shiningrush/droplet
```
