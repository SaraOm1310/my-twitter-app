package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

// User represents a user in my twitter app.
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Following []string `json:"following"`
	Tweets []Tweet `json:"tweets"`
}

// Tweet represents a tweet in my twitter app.
type Tweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// App represents my twitter app.
type App struct {
	Users  map[string]User
	Tweets []Tweet
}

// NewApp creates a new instance of my twitter app.
func NewApp() *App {
	return &App{
		Users:  make(map[string]User),
		Tweets: make([]Tweet, 0),
	}
}

var app = NewApp()

func createUser(c *gin.Context) {
	// Initializing a default user.
	var newUser User
	// Ensuring JSON binding.
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure we don't have the same user id
	if _, exists := app.Users[newUser.ID]; exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "user already exists"})
		return
	}

	app.Users[newUser.ID] = newUser
	c.Status(http.StatusCreated)
}

// Looking for a user using user id.
func findUser(c *gin.Context) {
	userID := c.Query("id")
	user, exists := app.Users[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Enables a user to update his profile
func updateUser(c *gin.Context) {
	userID := c.Query("id")

	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := app.Users[userID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	app.Users[userID] = updatedUser
	c.Status(http.StatusOK)
}

// Enables a user to delete his account
func deleteUser(c *gin.Context) {
	userID := c.Query("id")
	if _, exists := app.Users[userID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	delete(app.Users, userID)
	c.Status(http.StatusOK)
}


func postTweet(c *gin.Context) {
	var newTweet Tweet
	if err := c.ShouldBindJSON(&newTweet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	user := app.Users[newTweet.UserID]
	user.Tweets = append(app.Tweets, newTweet)
	app.Tweets = append(app.Tweets, newTweet)
	c.Status(http.StatusCreated)
}

// Enables a user to follow another user
func followUser(c *gin.Context) {
	followerID := c.Query("follower_id")
	followedID := c.Query("followed_id")

	if _, exists := app.Users[followerID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follower not found"})
		return
	}

	if _, exists := app.Users[followedID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Followed user not found"})
		return
	}

	followerUser:=app.Users[followerID]
	followerUser.Following = append(followerUser.Following, followedID)
	c.Status(http.StatusOK)
}

// Enables a user to unfollow another user
func unfollowUser(c *gin.Context) {
	followerID := c.Query("follower_id")
	unfollowedID := c.Query("unfollowed_id")

	if _, exists := app.Users[followerID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follower not found"})
		return
	}

	var updatedFollowing []string
	for _, followedID := range app.Users[followerID].Following {
		if followedID != unfollowedID {
			updatedFollowing = append(updatedFollowing, followedID)
		}
	}

	followerUser:=app.Users[followerID]
	followerUser.Following = updatedFollowing

	c.Status(http.StatusOK)
}

// When a user sign in into my twitter app, he needs to see his "twitter wall".
// Usually a user sees tweets in descending order.
func userWall(c *gin.Context) {
	userID := c.Query("id")

	var wall []Tweet
	for _, tweet := range app.Tweets {
		for _, followedID := range app.Users[userID].Following {
			if tweet.UserID == followedID {
				wall = append(wall, tweet)
			}
		}
	}

	// Sort timeline by timestamp in descending order (most recent first)
	sortTweetsByTimestampDesc(wall)

	c.JSON(http.StatusOK, wall)
}

func sortTweetsByTimestampDesc(tweets []Tweet) {
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].Timestamp.After(tweets[j].Timestamp)
	})
}

func main() {
	r := gin.Default()

	// Define routes
	r.POST("/users", createUser)
	r.GET("/users", findUser)
	r.PUT("/users", updateUser)
	r.DELETE("/users", deleteUser)

	r.POST("/tweets", postTweet)

	r.POST("/follow", followUser)
	r.DELETE("/follow", unfollowUser)

	r.GET("/timeline", userWall)

	// Run the server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
