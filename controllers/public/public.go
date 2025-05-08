package public

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"Go_Backend/config"
	"Go_Backend/models"
	"Go_Backend/utils"
)

// JWT Secret Key
var jwtKey = []byte(config.LoadConfig().JwtKey)

// User Collection Reference
var userCollection *mongo.Collection = utils.GetCollection("users")

// SetupPublicRoutes adds public routes to the provided Gin engine.
func SetupPublicRoutes(router *gin.Engine) {
	publicRoutes := router.Group("/api/public")
	{
		publicRoutes.POST("/usersignup", UserSignup)
		publicRoutes.GET("/emailverify/:token", EmailVerify)
		publicRoutes.POST("/usersignin", UserSignin)
		publicRoutes.POST("/forgotpassword", ForgotPassword)
	}
}

// UserSignup handles user registration.
func UserSignup(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input"})
		return
	}

	// Check if user already exists
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": newUser.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"msg": "User already exists, please sign in"})
		return
	}

	// Hash Password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error hashing password"})
		return
	}
	newUser.Password = string(hashedPass)

	// Generate Email Token
	tokenStr := utils.GenerateRandomToken()
	newUser.UserVerifyToken.Email = &tokenStr

	// Insert User in DB
	_, err = userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
		return
	}

	// Send Email Verification
	emailData := utils.EmailData{
		From:    config.LoadConfig().Email,
		To:      newUser.Email,
		Subject: "Verification Link",
		Text:    config.LoadConfig().URL + "/api/public/emailverify/" + *newUser.UserVerifyToken.Email,
	}
	utils.SendEmail(emailData)

	c.JSON(http.StatusOK, gin.H{"msg": "You'll be registered once you verify your email! 🙌"})
}

// EmailVerify handles email verification.
func EmailVerify(c *gin.Context) {
	token := c.Param("token")

	// Find user with matching token
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"userVerifyToken.email": token}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Invalid token ❌"})
		return
	}

	// If user already verified
	if user.UserVerified.Email {
		c.JSON(http.StatusOK, gin.H{"msg": "User email already verified"})
		return
	}

	// Update user verification
	update := bson.M{
		"$set": bson.M{"userVerified.email": true},
		"$unset": bson.M{"userVerifyToken.email": ""},
	}
	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"email": user.Email}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error updating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "User email verified successfully! ✅"})
}

// UserSignin handles user login.
func UserSignin(c *gin.Context) {
	var loginUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input"})
		return
	}

	// Retrieve user
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": loginUser.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Email doesn't exist"})
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid password"})
		return
	}

	// Generate JWT
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID.Hex(),
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(),
	}).SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "User Logged in Successfully! 🙌", "token": token})
}

// ForgotPassword handles password reset.
func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input"})
		return
	}

	// Find user
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Email not found"})
		return
	}

	// Generate new password
	newPass := utils.GenerateRandomToken()
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error generating password"})
		return
	}

	// Update password
	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"email": user.Email}, bson.M{"$set": bson.M{"password": hashedPass}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error updating password"})
		return
	}

	// Send Email
	emailData := utils.EmailData{
		From:    config.LoadConfig().Email,
		To:      req.Email,
		Subject: "New Password",
		Text:    "Your new password is: " + newPass,
	}
	utils.SendEmail(emailData)

	c.JSON(http.StatusOK, gin.H{"msg": "New password sent to your email successfully!"})
}