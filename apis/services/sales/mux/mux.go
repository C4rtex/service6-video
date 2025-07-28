package mux

import (
	"github.com/ardanlabs/service/apis/services/sales/route/sys/checkapi"
	"net/http"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI() *http.ServeMux {
	mux := http.NewServeMux()

	checkapi.Routes(mux)

	return mux
}
