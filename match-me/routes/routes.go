package routes

import (
	"context"
	"match-me/controllers"
	"match-me/graph"
	"match-me/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func graphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	return func(c *gin.Context) {

		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/graphql")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func SetupRouter(devMode bool) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", "./uploads")

	r.POST("/graphql", middleware.JWTMiddleware(), graphqlHandler())
	if devMode {
		r.GET("/graphql", playgroundHandler())
	}

	// REST endpoints
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	protected := r.Group("/")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.GET("/users/:id", controllers.GetUser)
		protected.GET("/users/:id/profile", controllers.GetUserProfile)
		protected.GET("/users/:id/bio", controllers.GetUserBio)
		protected.GET("/me", controllers.GetMe)
		protected.GET("/me/profile", controllers.GetMeProfile)
		protected.GET("/me/bio", controllers.GetMeBio)
		protected.GET("/profile", controllers.GetProfile)
		protected.PUT("/profile", controllers.UpdateProfile)
		protected.PUT("/profile/remove-picture", controllers.RemoveProfilePicture)
		protected.GET("/recommendations", controllers.GetRecommendations)
		protected.POST("/recommendations/decline", controllers.DeclineRecommendation)
		protected.POST("/connections/request", controllers.SendConnectionRequest)
		protected.GET("/connections/requests", controllers.GetConnectionRequests)
		protected.POST("/connections/accept", controllers.AcceptConnectionRequest)
		protected.POST("/connections/decline", controllers.DeclineConnectionRequest)
		protected.GET("/connections", controllers.GetConnections)
		protected.POST("/connections/disconnect", controllers.DisconnectUser)
		protected.GET("/chats/:userId", controllers.GetChatMessages)
		protected.POST("/chats/send", controllers.SendMessage)
		protected.GET("/chats", controllers.GetChats)
		protected.GET("/ws/chat", controllers.WebSocketChatHandler)
		protected.GET("/ws/chat_list", controllers.WebSocketChatListHandler)
		protected.POST("/chats/:chatId/read", controllers.MarkMessagesAsRead)
	}

	return r
}
