package private

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"Go_Backend/models"
	"Go_Backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// todoCollection holds the MongoDB collection instance for todos.
var todoCollection *mongo.Collection

// SetupPrivateRoutes initializes the collection and registers all private routes.
func SetupPrivateRoutes(rg *gin.RouterGroup) {
	// Initialize the "todos" collection.
	todoCollection = utils.GetCollection("todos")

	// Register endpoints for todos.
	rg.POST("/addtodo", AddTodo)
	rg.GET("/alltodos", GetAllTodos)
	rg.GET("/getone/:id", GetOneTodo)
	rg.PUT("/editone/:id", EditTodo)
	rg.DELETE("/deleteone/:id", DeleteTodo)
	rg.DELETE("/deleteall", DeleteAllTodos)
}

func AddTodo(c *gin.Context) {
	var newTodo models.Todo

	// Bind JSON or multipart/form-data.
	if err := c.ShouldBind(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input ❌", "error": err.Error()})
		return
	}

	// Ensure the "uploads" folder exists.
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	// Handle file upload if provided under the key "file".
	if file, err := c.FormFile("file"); err == nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			log.Println("❌ File upload failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "File upload failed ❌", "error": err.Error()})
			return
		}
		newTodo.Image = path
	}

	newTodo.ID = primitive.NewObjectID()
	if _, err := todoCollection.InsertOne(context.Background(), newTodo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error ❌", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Todo added successfully ✅", "todo": newTodo})
}

func GetAllTodos(c *gin.Context) {
	// Get pagination parameters (limit & skip) from query parameters.
	limit, skip := utils.GetPaginationParams(c)
	cacheKey := fmt.Sprintf("todos_limit_%d_skip_%d", limit, skip)
	ctx := context.Background()

	// Attempt to retrieve cached data.
	cached, err := utils.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var todos []models.Todo
		if err := json.Unmarshal([]byte(cached), &todos); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"todos": todos,
				"pagination": gin.H{
					"limit": limit,
					"skip":  skip,
					"count": len(todos),
				},
				"cache": true,
			})
			return
		} else {
			log.Printf("Error unmarshalling cached data: %v", err)
		}
	}

	// Cache miss: fetch data from MongoDB.
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	cursor, err := todoCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error ❌", "error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var todos []models.Todo
	if err := cursor.All(ctx, &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error fetching todos ❌", "error": err.Error()})
		return
	}

	// Marshal the result to JSON and store in Redis for 5 minutes.
	todosJSON, err := json.Marshal(todos)
	if err != nil {
		log.Printf("Error marshalling todos: %v", err)
	} else {
		if err := utils.RedisClient.Set(ctx, cacheKey, todosJSON, 5*time.Minute).Err(); err != nil {
			log.Printf("Error setting todos to cache: %v", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
		"pagination": gin.H{
			"limit": limit,
			"skip":  skip,
			"count": len(todos),
		},
		"cache": false,
	})
}

func GetOneTodo(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌", "error": err.Error()})
		return
	}
	var todo models.Todo
	if err := todoCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&todo); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Todo not found ❌", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func EditTodo(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌", "error": err.Error()})
		return
	}

	var updated models.Todo
	// Bind JSON or multipart/form-data.
	if err := c.ShouldBind(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input ❌", "error": err.Error()})
		return
	}

	// Handle file upload if present under the key "fileUpload".
	if file, err := c.FormFile("fileUpload"); err == nil {
		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			os.Mkdir("uploads", os.ModePerm)
		}
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "File upload failed ❌", "error": err.Error()})
			return
		}
		updated.Image = path
	}

	// Create the update document.
	updateDoc := bson.M{
		"$set": bson.M{
			"date":            updated.Date,
			"taskTitle":       updated.TaskTitle,
			"taskDescription": updated.TaskDescription,
			"image":           updated.Image,
		},
	}

	// Perform the update.
	result, err := todoCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error updating todo ❌", "error": err.Error()})
		return
	}
	log.Printf("Update result: Matched=%d, Modified=%d", result.MatchedCount, result.ModifiedCount)
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Todo not found ❌"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Todo updated successfully ✅", "updatedTodo": updated})
}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌", "error": err.Error()})
		return
	}
	res, err := todoCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error deleting todo ❌", "error": err.Error()})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Todo not found ❌"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Todo deleted successfully ✅"})
}

func DeleteAllTodos(c *gin.Context) {
	res, err := todoCollection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error deleting todos ❌", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "All todos deleted successfully ✅", "deletedCount": res.DeletedCount})
}