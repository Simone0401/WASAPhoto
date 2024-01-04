package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)

	/* ======== LOGIN API ========= */
	rt.router.POST("/session", rt.wrap(rt.doLogin, false))

	/* ======== USERNAME API ========= */
	rt.router.GET("/users/:uid/username", rt.wrap(rt.getUsername, true))
	rt.router.PUT("/users/:uid/username", rt.wrap(rt.setUsername, true))

	/* ======== FOLLOW API ========= */
	rt.router.PUT("/users/:uid/following/:fuid", rt.wrap(rt.followUser, true))
	rt.router.DELETE("/users/:uid/following/:fuid", rt.wrap(rt.unfollowUser, true))

	/* ======== MUTE API ========= */
	rt.router.PUT("/users/:uid/muted/:muteduid", rt.wrap(rt.banUser, true))
	rt.router.DELETE("/users/:uid/muted/:muteduid", rt.wrap(rt.unbanUser, true))

	/* ======== POSTS API ========= */
	rt.router.POST("/users/:uid/posts/", rt.wrap(rt.uploadPost, true))
	rt.router.GET("/images/:imageid", rt.wrap(rt.getImage, true))
	rt.router.DELETE("/users/:uid/posts/:postid", rt.wrap(rt.deletePost, true))

	/* ======== SPECIAL ROUTES ========= */
	rt.router.GET("/liveness", rt.liveness)
	return rt.router
}
