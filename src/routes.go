package main

import (
//	"net/http"
	 mux "github.com/julienschmidt/httprouter"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Handle	    mux.Handle
}

type Routes []Route

var routes = Routes{
	Route{
		"TodoCreate",
		"POST",
		"/me",
		RegisterMau,
	},
}
