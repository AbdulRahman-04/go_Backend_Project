package private

import (
	"context"
	"log"
	"net/http"
	"os"

	"Go_Backend/models"
	"Go_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var todoCollection *mongo.Collection

// SetupPrivateRoutes initializes the collection and registers all private routes.
// This function is exported (capitalized) so that main.go can call it.
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

func GetAllTodos(c *gin.Context) {
	limit, skip := utils.GetPaginationParams(c)
	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := todoCollection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error ❌", "error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var todos []models.Todo
	if err := cursor.All(context.Background(), &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error fetching todos ❌", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
		"pagination": gin.H{
			"limit": limit,
			"skip":  skip,
			"count": len(todos),
		},
	})
}

func EditTodo(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌", "error": err.Error()})
		return
	}

	var updated models.Todo
	// Bind JSON or multipart/form-data into the updated struct.
	if err := c.ShouldBind(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input ❌", "error": err.Error()})
		return
	}

	// Handle file upload if present under key "fileUpload".
	if file, err := c.FormFile("fileUpload"); err == nil {
		// Ensure the "uploads" directory exists.
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

	// Log the update result for debugging.
	log.Printf("Update result: Matched=%d, Modified=%d", result.MatchedCount, result.ModifiedCount)

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Todo not found ❌"})
		return
	}

	// If the update was successful (even if ModifiedCount is 0 due to same values),
	// return a success response.
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