package private

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

    "Go_Backend/models"
    "Go_Backend/utils"
)

var todoCollection *mongo.Collection

// SetupPrivateRoutes initializes the collection and routes (call after ConnectDB)
func SetupPrivateRoutes(rg *gin.RouterGroup) {
    todoCollection = utils.GetCollection("todos")

    rg.POST("/addtodo", AddTodo)
    rg.GET("/alltodos", GetAllTodos)
    rg.GET("/getone/:id", GetOneTodo)
    rg.PUT("/editone/:id", EditTodo)
    rg.DELETE("/deleteone/:id", DeleteTodo)
    rg.DELETE("/deleteall", DeleteAllTodos)
}

func AddTodo(c *gin.Context) {
    var newTodo models.Todo
    if err := c.ShouldBindJSON(&newTodo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input ❌"})
        return
    }

    if _, err := os.Stat("uploads"); os.IsNotExist(err) {
        os.Mkdir("uploads", os.ModePerm)
    }

    if file, err := c.FormFile("file"); err == nil {
        path := "uploads/" + file.Filename
        if err := c.SaveUploadedFile(file, path); err != nil {
            log.Println("❌ File upload failed:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"msg": "File upload failed ❌"})
            return
        }
        newTodo.Image = path
    }

    newTodo.ID = primitive.NewObjectID()
    if _, err := todoCollection.InsertOne(context.TODO(), newTodo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error ❌"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"msg": "Todo added successfully ✅", "todo": newTodo})
}

func GetOneTodo(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌"})
        return
    }
    var todo models.Todo
    if err := todoCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&todo); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"msg": "Todo not found ❌"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func GetAllTodos(c *gin.Context) {
    cursor, err := todoCollection.Find(context.TODO(), bson.M{})
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

func EditTodo(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌"})
        return
    }

    var updated models.Todo
    if err := c.ShouldBindJSON(&updated); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input ❌"})
        return
    }

    if file, err := c.FormFile("fileUpload"); err == nil {
        path := "uploads/" + file.Filename
        if err := c.SaveUploadedFile(file, path); err != nil {
            log.Println("❌ File upload failed:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"msg": "File upload failed ❌"})
            return
        }
        updated.Image = path
    }

    _, err = todoCollection.UpdateOne(
        context.TODO(),
        bson.M{"_id": objID},
        bson.M{"$set": bson.M{
            "date":            updated.Date,
            "taskTitle":       updated.TaskTitle,
            "taskDescription": updated.TaskDescription,
            "image":           updated.Image,
        }},
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error updating todo ❌"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"msg": "Todo updated successfully ✅", "updatedTodo": updated})
}

func DeleteTodo(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID format ❌"})
        return
    }
    res, err := todoCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error deleting todo ❌"})
        return
    }
    if res.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"msg": "Todo not found ❌"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"msg": "Todo deleted successfully ✅"})
}

func DeleteAllTodos(c *gin.Context) {
    if _, err := todoCollection.DeleteMany(context.TODO(), bson.M{}); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error deleting todos ❌"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"msg": "All todos deleted successfully ✅"})
}
