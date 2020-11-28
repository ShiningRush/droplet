package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

const (
	KeyHttpRequest = "HttpRequest"
	KeyRequestID   = "RequestID"
)

// use a single instance of Validate, it caches struct info
var vd *validator.Validate

func init() {
	vd = validator.New()
}

type HttpInputOption struct {
	ReqFunc        func() *http.Request
	PathParamsFunc func(key string) string
	InputType      reflect.Type
	IsReadFromBody bool
}

type HttpInputMiddleware struct {
	BaseMiddleware
	opt HttpInputOption

	input         interface{}
	multiPartData map[string][]byte
}

func NewHttpInputMiddleWare(opt HttpInputOption) *HttpInputMiddleware {
	return &HttpInputMiddleware{opt: opt}
}

func (mw *HttpInputMiddleware) Handle(ctx droplet.Context) error {
	mw.injectInfoToContext(ctx)
	if mw.opt.InputType == nil {
		return mw.BaseMiddleware.Handle(ctx)
	}

	switch mw.opt.ReqFunc().Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		mw.opt.IsReadFromBody = true
	}

	pInput := reflect.New(mw.opt.InputType).Interface()
	if err := mw.injectFieldFromBody(pInput); err != nil {
		return data.NewFormatError(err.Error())
	}

	if err := mw.injectFieldFromUrlAndForm(pInput); err != nil {
		return data.NewFormatError(err.Error())
	}

	isRecovered, err := recoverPager(pInput)
	if err != nil {
		return data.NewFormatError(err.Error())
	}
	if !isRecovered {
		if err := vd.Struct(pInput); err != nil {
			// TODO: parse err to items
			return data.NewValidateError(fmt.Sprintf("input validate failed: %s", err), nil)
		}
	}

	ctx.SetInput(pInput)
	return mw.BaseMiddleware.Handle(ctx)
}

func (mw *HttpInputMiddleware) injectFieldFromBody(ptr interface{}) error {
	if !mw.opt.IsReadFromBody || mw.opt.ReqFunc().ContentLength == 0 {
		return nil
	}

	contentType := mw.opt.ReqFunc().Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(contentType, "multipart/form-data"):
		return mw.readFromMultipartFormData()
	default: // application/json
		if err := mw.readJsonBody(ptr); err != nil {
			return err
		}
	}

	return nil
}

func (mw *HttpInputMiddleware) readJsonBody(ptr interface{}) error {
	dc := json.NewDecoder(mw.opt.ReqFunc().Body)
	return dc.Decode(ptr)
}

func (mw *HttpInputMiddleware) readFromMultipartFormData() error {
	reader, err := mw.opt.ReqFunc().MultipartReader()
	if err != nil {
		return fmt.Errorf("read form-data input from body failed: %s", err)
	}

	multiParts := map[string][]byte{}
	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read next part from body failed: %s", err)
		}

		bs, err := ioutil.ReadAll(p)
		if err != nil {
			return fmt.Errorf("read part body from body failed: %s", err)
		}
		if p.FileName() != "" {
			multiParts[fmt.Sprintf("%s%s", "_", p.FormName())] = []byte(p.FileName())
		}

		multiParts[p.FormName()] = bs
	}
	mw.multiPartData = multiParts
	return nil
}

func (mw *HttpInputMiddleware) injectFieldFromUrlAndForm(ptr interface{}) error {
	elType := reflect.TypeOf(ptr).Elem()
	input := reflect.ValueOf(ptr).Elem()

	for i := 0; i < elType.NumField(); i++ {
		if input.Field(i).Kind() == reflect.Struct {
			if err := mw.injectFieldFromUrlAndForm(input.Field(i).Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		src, name := getSourceWayAndName(elType.Field(i))
		if src == "" && mw.opt.IsReadFromBody {
			if mw.multiPartData != nil {
				if v, ok := mw.multiPartData[name]; ok {
					if input.Field(i).Kind() == reflect.String {
						input.Field(i).Set(reflect.ValueOf(string(v)))
					} else if input.Field(i).Kind() == reflect.Slice {
						input.Field(i).Set(reflect.ValueOf(v))
					}
				}
			}
			continue
		}

		val := ""
		switch src {
		case "path":
			val = mw.opt.PathParamsFunc(name)
		case "header":
			val = mw.opt.ReqFunc().Header.Get(name)
		default:
			val = mw.opt.ReqFunc().FormValue(name)
		}

		tarVal, err := changeToFieldKind(val, input.Field(i).Kind())
		if err != nil {
			return err
		}
		input.Field(i).Set(reflect.ValueOf(tarVal))
	}

	return nil
}

func (mw *HttpInputMiddleware) injectInfoToContext(ctx droplet.Context) {
	ctx.Set(KeyHttpRequest, mw.opt.ReqFunc())
	ctx.Set(KeyRequestID, mw.opt.ReqFunc().Header.Get(droplet.Option.HeaderKeyRequestID))
	ctx.SetPath(mw.opt.ReqFunc().URL.Path)
}

func recoverPager(pInput interface{}) (bool, error) {
	if v, ok := pInput.(data.PagerInfo); ok {
		return data.RecoverPager(v)
	}

	return false, nil
}

func getSourceWayAndName(field reflect.StructField) (src, name string) {
	src, name = "", lowerFirst(field.Name)
	tag := field.Tag.Get("auto_read")
	if tag == "" {
		return
	}

	tagArr := strings.Split(tag, ",")
	name = strings.TrimSpace(tagArr[0])
	if len(tagArr) > 1 {
		src = strings.TrimSpace(tagArr[1])
	}

	return
}

func lowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func changeToFieldKind(str string, kind reflect.Kind) (interface{}, error) {
	if kind == reflect.String {
		return str, nil
	}

	if kind == reflect.Bool {
		if str == "" {
			return false, nil
		}
		b, err := strconv.ParseBool(str)
		if err != nil {
			return nil, fmt.Errorf("changeToFieldKind covert to bool failed: %s", err)
		}

		return b, nil
	}

	if kind == reflect.Int {
		if str == "" {
			return 0, nil
		}
		i, err := strconv.ParseInt(str, 10, 0)
		if err != nil {
			return nil, fmt.Errorf("changeToFieldKind covert to int failed: %s", err)
		}
		return int(i), nil
	}

	if kind == reflect.Uint {
		if str == "" {
			return 0, nil
		}
		i, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			return nil, fmt.Errorf("changeToFieldKind covert to uint failed: %s", err)
		}
		return uint(i), nil
	}

	if kind == reflect.Int64 {
		if str == "" {
			return 0, nil
		}
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("changeToFieldKind covert to int64 failed: %s", err)
		}
		return i, nil
	}

	if kind == reflect.Uint64 {
		if str == "" {
			return 0, nil
		}
		i, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("changeToFieldKind covert to uint64 failed: %s", err)
		}
		return i, nil
	}

	return nil, nil
}
