package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	APIBaseURL    = "http://localhost:8080/forum/api"
	AuthURI       = APIBaseURL + "/session/verify"
	DataURI       = APIBaseURL + "/allData"
	LoginURI      = APIBaseURL + "/session/login"
	LogoutURI     = APIBaseURL + "/session/logout"
	RegisterURI   = APIBaseURL + "/register"
	CategoriesURI = APIBaseURL + "/categories"
	ReactionsURI  = APIBaseURL + "/react"
	CommentsURI   = APIBaseURL + "/comments"
	CreatePostURI = APIBaseURL + "/posts/create"
	MyPostsURI    = APIBaseURL + "/user/posts"
	LikedPostsURI = APIBaseURL + "/user/liked"
)

// Valid routes that should serve specific pages
var validRoutes = map[string]string{
	"/":        "./static/templates/index.html",
	"/login":   "./static/templates/login.html",
	"/register": "./static/templates/register.html",
	"/guest":   "./static/templates/guest.html",
	"/user":    "./static/templates/user.html",
	"/error":   "./static/templates/error.html",
}

func main() {
	// Serve static assets (css, js, images)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle config endpoint
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		config := map[string]string{
			"APIBaseURL":    APIBaseURL,
			"AuthURI":       AuthURI,
			"DataURI":       DataURI,
			"LoginURI":      LoginURI,
			"LogoutURI":     LogoutURI,
			"RegisterURI":   RegisterURI,
			"CategoriesURI": CategoriesURI,
			"ReactionsURI":  ReactionsURI,
			"CommentsURI":   CommentsURI,
			"CreatePostURI": CreatePostURI,
			"MyPostsURI":    MyPostsURI,
			"LikedPostsURI": LikedPostsURI,
		}
		json.NewEncoder(w).Encode(config)
	})

	// Main handler for all other routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		
		// Check if it's a valid route
		if templatePath, exists := validRoutes[path]; exists {
			http.ServeFile(w, r, templatePath)
			return
		}
		
		// Check if it's a static asset request (should be handled by the static file server)
		if strings.HasPrefix(path, "/static/") {
			http.NotFound(w, r)
			return
		}
		
		// For any other path, serve 404 error page
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		// Read and serve the error.html file
		content, err := os.ReadFile("./static/templates/error.html")
		if err != nil {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		w.Write(content)
	})

	// Start the server
	log.Println("Serving on http://localhost:8081/")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
    }
}

