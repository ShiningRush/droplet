package middleware

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/shiningrush/droplet/codec"
	"github.com/shiningrush/droplet/core"
	"github.com/shiningrush/droplet/data"
)

// DefaultValidator use a single instance of Validate, it caches struct info
var DefaultValidator = validator.New()

type HttpInputOption struct {
	PathParamsFunc       func(key string) string
	InputType            reflect.Type
	IsReadFromBody       bool
	DisableUnmarshalBody bool
	Codecs               []codec.Interface
}

type HttpInputMiddleware struct {
	BaseMiddleware
	opt HttpInputOption

	req       *http.Request
	searchMap map[string][]byte
}

func NewHttpInputMiddleWare(opt HttpInputOption) *HttpInputMiddleware {
	return &HttpInputMiddleware{opt: opt}
}

func (mw *HttpInputMiddleware) Handle(ctx core.Context) error {
	httpReq := ctx.Get(KeyHttpRequest)
	if httpReq == nil {
		return fmt.Errorf("input middleware cannot get http request, please check if HttpInfoInjectorMiddleware middlle work well")
	}
	mw.req = httpReq.(*http.Request)
	if mw.opt.InputType == nil {
		return mw.BaseMiddleware.Handle(ctx)
	}

	switch mw.req.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		mw.opt.IsReadFromBody = true
	}

	pInput := reflect.New(mw.opt.InputType).Interface()
	if !mw.opt.DisableUnmarshalBody {
		if err := mw.unmarshalFieldFromBody(pInput); err != nil {
			return data.NewFormatError(err.Error())
		}
	}

	if err := mw.injectFieldFromUrlAndMap(pInput); err != nil {
		return data.NewFormatError(err.Error())
	}

	isRecovered, err := recoverPager(pInput)
	if err != nil {
		return data.NewFormatError(err.Error())
	}
	if !isRecovered {
		if err := DefaultValidator.Struct(pInput); err != nil {
			// TODO: parse err to items
			return data.NewValidateError(fmt.Sprintf("input validate failed: %s", err), nil)
		}
	}

	ctx.SetInput(pInput)
	return mw.BaseMiddleware.Handle(ctx)
}

func (mw *HttpInputMiddleware) unmarshalFieldFromBody(ptr interface{}) error {
	if !mw.opt.IsReadFromBody || mw.req.ContentLength == 0 {
		return nil
	}

	contentType := mw.req.Header.Get("Content-Type")
	var coc codec.Interface = &codec.Json{}
	for _, c := range mw.opt.Codecs {
		for _, ctt := range c.ContentType() {
			if strings.HasPrefix(contentType, ctt) {
				coc = c
			}
		}
	}

	if dir, ok := coc.(codec.Direct); ok {
		if err := dir.Unmarshal(mw.req, ptr); err != nil {
			return err
		}
	}

	if s, ok := coc.(codec.Search); ok {
		m, err := s.UnmarshalSearchMap(mw.req)
		if err != nil {
			return err
		}
		mw.searchMap = m
	}

	return nil
}

func (mw *HttpInputMiddleware) injectFieldFromUrlAndMap(ptr interface{}) error {
	elType := reflect.TypeOf(ptr).Elem()
	input := reflect.ValueOf(ptr).Elem()

	for i := 0; i < elType.NumField(); i++ {
		if input.Field(i).Kind() == reflect.Struct {
			if err := mw.injectFieldFromUrlAndMap(input.Field(i).Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		src, name := getSourceWayAndName(elType.Field(i))
		if src == "" && mw.opt.IsReadFromBody {
			if mw.searchMap != nil {
				if v, ok := mw.searchMap[name]; ok {
					if input.Field(i).Kind() == reflect.String {
						input.Field(i).Set(reflect.ValueOf(string(v)))
					} else if input.Field(i).Kind() == reflect.Slice {
						input.Field(i).Set(reflect.ValueOf(v))
					}
				}
			}
			if name == "@body" {
				if input.Field(i).Type().Implements(reflect.TypeOf((*io.ReadCloser)(nil)).Elem()) {
					input.Field(i).Set(reflect.ValueOf(mw.req.Body))
					continue
				}

				bs, err := data.CopyBody(mw.req)
				if err != nil {
					return err
				}
				input.Field(i).Set(reflect.ValueOf(bs))
			}
			continue
		}

		val := ""
		switch src {
		case "path":
			val = mw.opt.PathParamsFunc(name)
		case "header":
			val = mw.req.Header.Get(name)
		default:
			val = mw.req.FormValue(name)
		}

		tarVal, err := changeToFieldKind(val, input.Field(i).Type())
		if err != nil {
			return err
		}
		if tarVal == nil {
			continue
		}
		input.Field(i).Set(reflect.ValueOf(tarVal))
	}

	return nil
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

func changeToFieldKind(str string, t reflect.Type) (interface{}, error) {
	kind := t.Kind()
	isPtr := false
	if kind == reflect.Ptr {
		if str == "" {
			return nil, nil
		}
		isPtr = true
		kind = t.Elem().Kind()
	}

	if kind == reflect.String {
		if isPtr {
			return &str, nil
		}
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
		if isPtr {
			return &b, nil
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

		i32 := int(i)
		if isPtr {
			return &i32, nil
		}
		return i32, nil
	}

	if kind == reflect.Uint {
		if str == "" {
			return uint(0), nil
		}
		i, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			return nil, fmt.Errorf("changeToFieldKind covert to uint failed: %s", err)
		}

		ui := uint(i)
		if isPtr {
			return &ui, nil
		}
		return ui, nil
	}

	if kind == reflect.Int64 {
		if str == "" {
			return int64(0), nil
		}
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("changeToFieldKind covert to int64 failed: %s", err)
		}

		if isPtr {
			return &i, nil
		}
		return i, nil
	}

	if kind == reflect.Uint64 {
		if str == "" {
			return uint64(0), nil
		}
		i, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("changeToFieldKind covert to uint64 failed: %s", err)
		}

		if isPtr {
			return &i, nil
		}
		return i, nil
	}

	return nil, fmt.Errorf("unsupport type: %s", kind.String())
}
