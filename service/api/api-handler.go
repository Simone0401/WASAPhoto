package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.wrap(rt.doLogin, false))

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

	/* Section LIKE */
	rt.router.PUT("/posts/:postid/likes/:uid", rt.wrap(rt.likePost, true))
	rt.router.DELETE("/posts/:postid/likes/:uid", rt.wrap(rt.unlikePost, true))

	/* Section COMMENT */
	rt.router.POST("/posts/:postid/comments/", rt.wrap(rt.commentPost, true))
	rt.router.DELETE("/posts/:postid/comments/:commentid", rt.wrap(rt.uncommentPost, true))

	/* ======== MYSTREAM API ========= */
	rt.router.GET("/users/:uid/mystream", rt.wrap(rt.getMyStream, true))

	/* ======== PROFILE API ========= */
	rt.router.GET("/users/:uid/profile", rt.wrap(rt.getUserProfile, true))

	/* ======== SPECIAL ROUTES ========= */
	rt.router.GET("/liveness", rt.liveness)
	return rt.router
}
