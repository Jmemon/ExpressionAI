package main

import (
    "context"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const (
    mongoURI = "mongodb://localhost:27017"
	maxChunkSize = 15 * 1024 * 1024  // 15 Mb
)

var (
    upgrader      = websocket.Upgrader{}
    mongoClient   *mongo.Client
	audioDataCollection *mongo.Collection
    buffer1       []byte
    buffer2       []byte
    currentBuffer *[]byte
    bufferMutex   sync.Mutex
)

func init() {
    var err error
    mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
	audioDataCollection = mongoClient.Database("audioDB").Collection("audioChunks")

    // Initialize buffers
    buffer1 = make([]byte, 0, maxChunkSize)
    buffer2 = make([]byte, 0, maxChunkSize)
    currentBuffer = &buffer1
}

func handleAudioStream(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer conn.Close()

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            break
        }

        bufferMutex.Lock()
        if len(*currentBuffer) + len(message) > maxChunkSize {
            go storeChunkAndSwapBuffers() // Store current buffer and swap buffers in a goroutine
        }
        *currentBuffer = append(*currentBuffer, message...)
        bufferMutex.Unlock()
    }
}

func storeChunkAndSwapBuffers() {
    bufferMutex.Lock()

    // Determine which buffer is currently passive
    passiveBuffer := currentBuffer
    if currentBuffer == &buffer1 {
        currentBuffer = &buffer2
    } else {
        currentBuffer = &buffer1
    }

    // Copy data to avoid holding the lock while accessing the database
    dataToStore := make([]byte, len(*passiveBuffer))
    copy(dataToStore, *passiveBuffer)

    // Reset the passive buffer
    *passiveBuffer = (*passiveBuffer)[:0]
    bufferMutex.Unlock()

    // Create a document with the data
    doc := bson.M{
        "data": dataToStore,
    }

    // Insert the chunk into the MongoDB
    _, err := audioDataCollection.InsertOne(context.Background(), doc)
    if err != nil {
    	log.Printf("Failed to store audio chunk: %v", err)
    }
}

func main() {
    defer func() {
        if err := mongoClient.Disconnect(context.TODO()); err != nil {
            log.Fatalf("Error disconnecting from MongoDB: %v", err)
        }
    }()

    http.HandleFunc("/audio", handleAudioStream)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
