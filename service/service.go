package service

import (
	"reflect"

	"github.com/plimble/micro"
)

var mapStruct map[string]method

type method struct {
	fn  reflect.Method
	req reflect.Value
	res reflect.Value
	s   reflect.Value
}

func init() {
	mapStruct = make(map[string]method)
}

func QueueSubscribe(m *micro.Micro, prefix string, v interface{}) {
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)

	if vt.Kind().String() != "ptr" {
		return
	}

	if vt.Elem().Kind().String() != "struct" {
		return
	}

	if vt.NumMethod() == 0 {
		return
	}

	for i := 0; i < vt.NumMethod(); i++ {
		var name string
		if prefix != "" {
			name = prefix + "." + vt.Method(i).Name
		} else {
			name = vt.Method(i).Name
		}
		mapStruct[name] = method{
			req: reflect.New(vt.Method(i).Type.In(1).Elem()),
			res: reflect.New(vt.Method(i).Type.In(2).Elem()),
			fn:  vt.Method(i),
			s:   vv,
		}

		m.QueueSubscribe(name, name, func(ctx *micro.Context) error {
			service := mapStruct[ctx.Subject]
			req := service.req.Interface()
			res := service.res.Interface()
			ctx.Decode(ctx.Data, req)

			fn := service.fn.Func.Call([]reflect.Value{service.s, reflect.ValueOf(req), reflect.ValueOf(res)})

			if ierr := fn[0].Interface(); ierr != nil {
				return ierr.(error)
			}

			if ctx.Reply != "" {
				ctx.Publish(ctx.Reply, res)
			}

			return nil
		})
	}
}

func Subscribe(m *micro.Micro, prefix string, v interface{}) {
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)

	if vt.Kind().String() != "ptr" {
		return
	}

	if vt.Elem().Kind().String() != "struct" {
		return
	}

	if vt.NumMethod() == 0 {
		return
	}

	for i := 0; i < vt.NumMethod(); i++ {
		var name string
		if prefix != "" {
			name = prefix + "." + vt.Method(i).Name
		} else {
			name = vt.Method(i).Name
		}
		mapStruct[name] = method{
			req: reflect.New(vt.Method(i).Type.In(1).Elem()),
			res: reflect.New(vt.Method(i).Type.In(2).Elem()),
			fn:  vt.Method(i),
			s:   vv,
		}

		m.Subscribe(name, func(ctx *micro.Context) error {
			service := mapStruct[ctx.Subject]
			req := service.req.Interface()
			res := service.res.Interface()
			ctx.Decode(ctx.Data, req)

			fn := service.fn.Func.Call([]reflect.Value{service.s, reflect.ValueOf(req), reflect.ValueOf(res)})

			if ierr := fn[0].Interface(); ierr != nil {
				return ierr.(error)
			}

			if ctx.Reply != "" {
				ctx.Publish(ctx.Reply, res)
			}

			return nil
		})
	}
}
