package api

import (
	"github.com/Simone0401/WASAPhoto/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler, auth bool) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// bearer token must be checked just for APIs that require it
		var authId int
		if auth {
			if !rt.isAuthorized(r.Header) {
				rt.baseLogger.Error("Auth Bearer Token is missing or invalid format!")
				http.Error(w, "Auth Bearer Token is missing or invalid format!", http.StatusUnauthorized)
				return
			}
			authId, err = strconv.Atoi(r.Header.Get("Authorization"))

			if err != nil {
				rt.baseLogger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			authId = 0
		}

		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
			Uid:     uint64(authId),
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}
