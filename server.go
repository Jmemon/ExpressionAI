
package main

import (
    "fmt"
    "net/http"
)


// -------------- Go server to receive audio

func main() {
    http.HandleFunc("/upload_audio", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
            return
        }

        // Assuming a form field named 'audio' is used for the upload
        file, _, err := r.FormFile("audio")
        if err != nil {
            http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
            return
        }
        defer file.Close()

        // Process the file here. For now, just print a message.
        fmt.Fprintf(w, "File uploaded successfully")
    })

    fmt.Printf("Server starting on port 8080\n")
    http.ListenAndServe(":8080", nil)
}
