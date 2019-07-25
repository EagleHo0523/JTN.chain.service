package router

import (
	"net/http"

	h "../handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{Name: "DDMX", Method: "POST", Pattern: "/DDMX", HandlerFunc: h.DDMX},
	Route{Name: "ETH", Method: "POST", Pattern: "/ETH", HandlerFunc: h.ETH},
	Route{Name: "BTC", Method: "POST", Pattern: "/BTC", HandlerFunc: h.BTC},
}
