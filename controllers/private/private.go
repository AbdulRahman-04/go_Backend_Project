package private

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"
	"github.com/go-redis/redis/v8"
	"errors"
	"Go_Backend/models"
	"Go_Backend/utils"
    "strconv"
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

    // ✅ Bind JSON or multipart/form-data
    if err := c.ShouldBindJSON(&newTodo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input ❌", "error": err.Error()})
        return
    }

    // ✅ Retrieve file path from middleware
    filePath := c.GetString("filePath")
    if filePath != "" {
        newTodo.Image = filePath
    }

    // ✅ Generate MongoDB ObjectID
    newTodo.ID = primitive.NewObjectID()

    // ✅ Insert into MongoDB
    if _, err := todoCollection.InsertOne(context.Background(), newTodo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error ❌", "error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"msg": "Todo added successfully ✅", "todo": newTodo})
}

func GetAllTodos(c *gin.Context) {
	limit, skip := utils.GetPaginationParams(c)
	cacheKey := fmt.Sprintf("todos_limit_%d_skip_%d", limit, skip)
	ctx := context.Background()

	cached, err := utils.RedisClient.Get(ctx, cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("Redis GET error: %v", err)
	}

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
			log.Printf("Unmarshal error: %v", err)
		}
	}

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

	todosJSON, err := json.Marshal(todos)
	if err == nil {
		if err := utils.RedisClient.Set(ctx, cacheKey, todosJSON, 5*time.Minute).Err(); err != nil {
			log.Printf("Redis SET error: %v", err)
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


func GetPaginationParams(c *gin.Context) (int, int) {
	limitStr := c.DefaultQuery("limit", "10")
	skipStr := c.DefaultQuery("skip", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	skip, err := strconv.Atoi(skipStr)
	if err != nil || skip < 0 {
		skip = 0
	}
	return limit, skip
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
    // Convert URL parameter to ObjectID
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌", "error": err.Error()})
        return
    }

    // Use a dedicated input struct for binding update fields
    var input struct {
        Date            string `json:"date" binding:"required"`
        TodoNo          int    `json:"todoNo" binding:"required"`
        TaskTitle       string `json:"taskTitle" binding:"required"`
        TaskDescription string `json:"taskDescription" binding:"required"`
    }
    
    // Bind JSON data – ensure keys are in lower-case (e.g., "date", "todoNo")
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input ❌", "error": err.Error()})
        return
    }
    
    // Handle file upload if provided under "fileUpload"
    filePath := "" // Will store new image path if uploaded
    if file, err := c.FormFile("fileUpload"); err == nil {
        // Ensure "uploads" folder exists
        if _, err := os.Stat("uploads"); os.IsNotExist(err) {
            os.Mkdir("uploads", os.ModePerm)
        }
        path := "uploads/" + file.Filename
        if err := c.SaveUploadedFile(file, path); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"msg": "File upload failed ❌", "error": err.Error()})
            return
        }
        filePath = path
    }
    
    // Build update document – include file path if new file is uploaded
    updateData := bson.M{
        "date":            input.Date,
        "todoNo":          input.TodoNo,
        "taskTitle":       input.TaskTitle,
        "taskDescription": input.TaskDescription,
    }
    if filePath != "" {
        updateData["image"] = filePath
    }
    
    updateDoc := bson.M{"$set": updateData}
    
    // Perform the update operation in MongoDB
    result, err := todoCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateDoc)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error updating todo ❌", "error": err.Error()})
        return
    }
    
    if result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"msg": "Todo not found ❌"})
        return
    }
    
    // Build the response todo with the correct ID
    // Note: If file wasn't updated, we can fetch the existing file path if needed.
    var existing models.Todo
    if err := todoCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&existing); err != nil {
        // Not fatal – we already updated the document. Use empty string if not found.
        existing.Image = ""
    }
    
    respTodo := models.Todo{
        ID:              objID,
        Date:            input.Date,
        TodoNo:          input.TodoNo,
        TaskTitle:       input.TaskTitle,
        TaskDescription: input.TaskDescription,
        Image:           filePath,
    }
    
    // If no new file was uploaded, ensure we return the existing image path
    if respTodo.Image == "" {
        respTodo.Image = existing.Image
    }
    
    c.JSON(http.StatusOK, gin.H{
        "msg":         "Todo updated successfully ✅", 
        "updatedTodo": respTodo,
    })
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