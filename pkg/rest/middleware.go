package rest

import (
	"log"
	"net/http"
)

type Middleware[TReqSchema any, TRespSchema any] func(r *http.Request, next HttpCommand[TReqSchema, TRespSchema]) HttpCommand[TReqSchema, TRespSchema]

func LoggingMiddleware[TReqSchema any, TRespSchema any]() Middleware[TReqSchema, TRespSchema] {
	return func(r *http.Request, next HttpCommand[TReqSchema, TRespSchema]) HttpCommand[TReqSchema, TRespSchema] {

		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		return next
	}
}
