package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)

	rt.router.POST("/session", rt.wrap(rt.doLogin, false))

	rt.router.GET("/users/:uid/username", rt.wrap(rt.getUsername, true))
	rt.router.PUT("/users/:uid/username", rt.wrap(rt.setUsername, true))

	rt.router.PUT("/users/:uid/following/:fuid", rt.wrap(rt.followUser, true))
	rt.router.DELETE("/users/:uid/following/:fuid", rt.wrap(rt.unfollowUser, true))

	rt.router.PUT("/users/:uid/muted/:muteduid", rt.wrap(rt.banUser, true))
	rt.router.DELETE("/users/:uid/muted/:muteduid", rt.wrap(rt.unbanUser, true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)
	return rt.router
}
