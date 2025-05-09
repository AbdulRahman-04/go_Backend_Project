package private

import (
	"context"
	"log"
	"net/http"
	"os"
	

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"Go_Backend/models"
	"Go_Backend/utils"
)

// ✅ Add Todo (POST)
func AddTodo(c *gin.Context) {
	collection := utils.GetCollection("todos")

	var newTodo models.Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input ❌"})
		return
	}

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	if file, err := c.FormFile("file"); err == nil {
		filePath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			log.Println("❌ File upload failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "File upload failed ❌"})
			return
		}
		newTodo.Image = filePath
	}

	newTodo.ID = primitive.NewObjectID()
	if _, err := collection.InsertOne(context.TODO(), newTodo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error ❌"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Todo added successfully ✅", "todo": newTodo})
}

// ✅ Get One Todo
func GetOneTodo(c *gin.Context) {
	collection := utils.GetCollection("todos")
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌"})
		return
	}

	var todo models.Todo
	if err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&todo); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Todo not found ❌"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

// ✅ Get All Todos
func GetAllTodos(c *gin.Context) {
	collection := utils.GetCollection("todos")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error ❌"})
		return
	}
	defer cursor.Close(context.TODO())

	var todos []models.Todo
	if err := cursor.All(context.TODO(), &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error fetching todos ❌"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

// ✅ Edit Todo
func EditTodo(c *gin.Context) {
	collection := utils.GetCollection("todos")
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌"})
		return
	}

	var updatedTodo models.Todo
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input ❌"})
		return
	}

	if file, err := c.FormFile("fileUpload"); err == nil {
		filePath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			log.Println("❌ File upload failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "File upload failed ❌"})
			return
		}
		updatedTodo.Image = filePath
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"date":            updatedTodo.Date,
			"taskTitle":       updatedTodo.TaskTitle,
			"taskDescription": updatedTodo.TaskDescription,
			"image":           updatedTodo.Image,
		}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error updating todo ❌"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Todo updated successfully ✅", "updatedTodo": updatedTodo})
}

// ✅ Delete One Todo
func DeleteTodo(c *gin.Context) {
	collection := utils.GetCollection("todos")
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌"})
		return
	}

	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error deleting todo ❌"})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Todo not found ❌"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Todo deleted successfully ✅"})
}

// ✅ Delete All Todos
func DeleteAllTodos(c *gin.Context) {
	collection := utils.GetCollection("todos")
	if _, err := collection.DeleteMany(context.TODO(), bson.M{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error deleting todos ❌"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "All todos deleted successfully ✅"})
}
