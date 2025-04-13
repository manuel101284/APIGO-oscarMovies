package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://manuel101284:manuel101284@cluster0.elm9zh7.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

var mongoClient *mongo.Client

func init() {
	if err := connectToMongoDB(); err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	println("----->>>    CONECTADO A LA BASE DE DATOS    <<<-----")
	println(".............................................................................................")
	println(".............................................................................................")
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Bienvenido a esta API de películas ganadoras de Oscar",
		})
	})

	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovieByID)
	router.POST("/movies/aggregate", aggregateMovies)

	router.Run()
}

func connectToMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	ops := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), ops)

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)

	mongoClient = client
	return err
}

// Funciones para hacer las petiones a la base de datos
func getMovies(c *gin.Context) {
	cursor, err := mongoClient.Database("oscarMovies").Collection("movies").Find(context.TODO(), bson.D{{}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las películas " + err.Error()})
		return
	}

	var movies []bson.M

	if err = cursor.All(context.TODO(), &movies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer las películas " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func aggregateMovies(c *gin.Context) {
	var pipeline interface{}

	if err := c.ShouldBindJSON(&pipeline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON" + err.Error()})
		return
	}

	cursor, err := mongoClient.Database("oscarMovies").Collection("movies").Aggregate(context.TODO(), pipeline)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las películas " + err.Error()})
		return
	}

	var results []bson.M

	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer las películas " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func getMovieByID(c *gin.Context) {
	idMovieStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idMovieStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format" + err.Error()})
		return
	}

	var movie bson.M

	err = mongoClient.Database("oscarMovies").Collection("movies").FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&movie)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func addBook(c *gin.Context) {
	var pipeline interface{}

	if err := c.ShouldBindJSON(&pipeline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON" + err.Error()})
		return
	}

	result, err := mongoClient.Database("oscarMovies").Collection("movies").InsertOne(context.TODO(), pipeline)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al agregar la película " + err.Error()})
		return
	}

	insertedID := result.InsertedID

	c.JSON(http.StatusOK, gin.H{"message": "Movie added successfully", "inserted_id": insertedID})
}
