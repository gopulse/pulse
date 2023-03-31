package routing

import "github.com/valyala/fasthttp"

type Context struct {
	*fasthttp.RequestCtx

	router       *Router
	params       []string
	paramsValues []string
	data         map[string]interface{}
	handlerIndex int
	handlers     []Handler
}
