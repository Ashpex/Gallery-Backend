package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/config"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/controller"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/helper"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/middleware"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/repository"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/service"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var (
	// db is a global variable that represents the config connection
	db *gorm.DB = config.SetupDB()
	// Database repository
	userRepository      repository.UserRepository      = repository.NewUserRepository(db)
	postRepository      repository.PostRepository      = repository.NewPostRepository(db)
	topicRepository     repository.TopicRepository     = repository.NewTopicRepository(db)
	commentRepository   repository.CommentRepository   = repository.NewCommentRepository(db)
	followerRepository  repository.FollowerRepository  = repository.NewFollowerRepository(db)
	likeRepository      repository.LikeRepository      = repository.NewLikeRepository(db)
	subscribeRepository repository.SubscribeRepository = repository.NewSubscribeRepository(db)
	// jwtService is a global variable that represents the jwt service (json web token)
	jwtService helper.JWTService = helper.NewJWTService()
	// Authentication service and controller
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	// User service and controller
	userService    service.UserService       = service.NewUserService(userRepository)
	userController controller.UserController = controller.NewUserController(userService, jwtService)

	// Post service and controller
	postService    service.PostService       = service.NewPostService(postRepository)
	postController controller.PostController = controller.NewPostController(postService, jwtService, likeService, followerService, subscribeService)
	// Topic service and controller
	topicService    service.TopicService       = service.NewTopicService(topicRepository)
	topicController controller.TopicController = controller.NewTopicController(topicService, jwtService)
	// Comment service and controller
	commentService    service.CommentService       = service.NewCommentService(commentRepository)
	commentController controller.CommentController = controller.NewCommentController(commentService, jwtService)
	// Follower service and controller
	followerService    service.FollowService         = service.NewFollowService(followerRepository)
	followerController controller.FollowerController = controller.NewFollowerController(followerService, jwtService)
	// Like service and controller
	likeService    service.LikeService       = service.NewLikeService(likeRepository)
	likeController controller.LikeController = controller.NewLikeController(likeService, jwtService)
	// Subscribe service and controller
	subscribeService    service.SubscribeService       = service.NewSubscribeService(subscribeRepository)
	subscribeController controller.SubscribeController = controller.NewSubscribeController(subscribeService, jwtService, topicService)
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	defer config.CloseDB(db)
	r := gin.Default()
	r.Static("/static", "./static")
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	postRoutes := r.Group("api/posts", middleware.AuthorizeJWT(jwtService))
	{
		postRoutes.GET("/", postController.All)
		postRoutes.POST("/", postController.Insert)
		postRoutes.GET("/:id", postController.FindByID)
		postRoutes.PUT("/:id", postController.Update)
		postRoutes.DELETE("/:id", postController.Delete)
		postRoutes.GET("/topic/:id", postController.FindByTopicID)
		//postRoutes.GET("?topic=:id", postController.FindByTopicID)
		postRoutes.GET("/trending", postController.GetTrendingPosts)
		postRoutes.GET("/following/", postController.GetFollowingPosts)
		postRoutes.GET("/subscribed/", postController.GetPostsFromSubscribedTopic)
		postRoutes.GET("?search=:search", postController.SearchPosts)
	}
	topicRoutes := r.Group("api/topics")
	{
		topicRoutes.GET("/", topicController.All)
		topicRoutes.POST("/", topicController.Insert)
		topicRoutes.GET("/:id", topicController.FindByID)
	}
	commentRoutes := r.Group("api/comments", middleware.AuthorizeJWT(jwtService))
	{
		commentRoutes.GET("/", commentController.All)
		commentRoutes.POST("/:post_id", commentController.Insert)
		commentRoutes.GET("/:id", commentController.FindByID)
		commentRoutes.PUT("/:id", commentController.Update)
		commentRoutes.DELETE("/:id", commentController.Delete)
		commentRoutes.GET("/post/:id", commentController.FindCommentByPostID)
	}

	followerRoutes := r.Group("api/followers", middleware.AuthorizeJWT(jwtService))
	{
		followerRoutes.GET("/", followerController.AllFollowers)
		followerRoutes.GET("/following/:id", followerController.AllFollowing)
		followerRoutes.POST("/", followerController.Follow)
		followerRoutes.DELETE("/:id", followerController.Unfollow)
	}

	likeRoutes := r.Group("api/likes", middleware.AuthorizeJWT(jwtService))
	{
		likeRoutes.GET("/:id", likeController.AllLikes)
		likeRoutes.POST("/", likeController.Like)
		likeRoutes.DELETE("/:id", likeController.UnLike)
		likeRoutes.GET("/count/:id", likeController.CountLikes)
	}
	subscribeRoutes := r.Group("api/subscribes", middleware.AuthorizeJWT(jwtService))
	{
		subscribeRoutes.GET("/:id", subscribeController.AllSubscribes)
		subscribeRoutes.POST("/:id", subscribeController.Subscribe)
		subscribeRoutes.DELETE("/:id", subscribeController.Unsubscribe)
		subscribeRoutes.GET("/count/:id", subscribeController.CountSubscribes)
	}

	prom := ginprometheus.NewPrometheus("gin", newCustomMetrics())
	prom.Use(r)
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}
func newCustomMetrics() []*ginprometheus.Metric {
	return []*ginprometheus.Metric{
		{
			ID:          "1",
			Name:        "metric_1",
			Description: "Counter test metric",
			Type:        "counter",
		},
		{
			ID:          "2",
			Name:        "metric_2",
			Description: "Summary metric",
			Type:        "summary",
		},
		{
			ID:          "3",
			Name:        "metric_3",
			Description: "Gauge metric",
			Type:        "gauge",
		},
		{
			ID:          "4",
			Name:        "metric_4",
			Description: "Histogram test metric",
			Type:        "histogram",
		},
	}
}
