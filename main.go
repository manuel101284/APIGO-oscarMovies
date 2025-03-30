package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	println("...")
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Bienvenido a esta API de películas ganadoras de Oscar",
		})
	})

	router.GET("/movies", getMovies)

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
