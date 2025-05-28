package public

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"Go_Backend/config"
	"Go_Backend/models"
	"Go_Backend/utils"
)

var jwtKey = []byte(config.LoadConfig().JwtKey)
var userCollection *mongo.Collection

// SetupPublicRoutes registers all public endpoints.
func SetupPublicRoutes(group *gin.RouterGroup) {
	userCollection = utils.GetCollection("users")

	group.POST("/usersignup", UserSignup)
	group.GET("/emailverify/:token", EmailVerify)
	group.POST("/usersignin", UserSignin)
	group.POST("/forgotpassword", ForgotPassword)
}
func UserSignup(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input"})
		return
	}

	// Check if user already exists.
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": newUser.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"msg": "User already exists, please sign in"})
		return
	}

	// Hash password.
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error hashing password"})
		return
	}
	newUser.Password = string(hashedPass)

	// Generate verification token.
	tokenStr := utils.GenerateRandomToken()
	newUser.UserVerifyToken.Email = &tokenStr

	// Insert user.
	_, err = userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
		return
	}

	// ✅ Send verification email in background (non-blocking)
	go func(user models.User, token string) {
		emailData := utils.EmailData{
			From:    config.LoadConfig().Email,
			To:      user.Email,
			Subject: "Verification Link",
			Text:    config.LoadConfig().URL + "/api/public/emailverify/" + token,
			HTML:    "<p>Click the link to verify your email: <a href='" +
				config.LoadConfig().URL + "/api/public/emailverify/" + token + "'>Verify Email</a></p>",
		}
		_ = utils.SendEmail(emailData) // Ignore error for now or log it if needed
	}(newUser, tokenStr)

	c.JSON(http.StatusOK, gin.H{"msg": "You'll be registered once you verify your email! 🙌"})
}


func EmailVerify(c *gin.Context) {
	token := c.Param("token")
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"userVerifyToken.email": token}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Invalid token ❌"})
		return
	}

	if user.UserVerified.Email {
		c.JSON(http.StatusOK, gin.H{"msg": "User email already verified"})
		return
	}

	update := bson.M{
		"$set":   bson.M{"userVerified.email": true},
		"$unset": bson.M{"userVerifyToken.email": ""},
	}
	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"email": user.Email}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error updating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "User email verified successfully! ✅"})
}

func UserSignin(c *gin.Context) {
	var loginUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input"})
		return
	}

	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": loginUser.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Email doesn't exist"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid password"})
		return
	}

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

func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input"})
		return
	}

	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Email not found"})
		return
	}

	newPass := utils.GenerateRandomToken()
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error generating password"})
		return
	}

	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"email": user.Email}, bson.M{"$set": bson.M{"password": string(hashedPass)}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error updating password"})
		return
	}

	emailData := utils.EmailData{
		From:    config.LoadConfig().Email,
		To:      req.Email,
		Subject: "New Password",
		Text:    "Your new password is: " + newPass,
		HTML:    "<h3>Your new password is:</h3><p><b>" + newPass + "</b></p>",
	}
	if err := utils.SendEmail(emailData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Password updated, but email failed to send"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "New password sent to your email successfully!"})
}