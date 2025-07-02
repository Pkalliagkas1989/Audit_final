# Nexxus Forum

A modern, full-stack web forum application featuring secure user authentication, post creation, commenting, and reaction systems. Built with Go for the backend and a custom, responsive frontend. The project is fully containerized for streamlined development and deployment.

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Installation](#installation)
- [Running the Project](#running-the-project)
- [API Usage](#api-usage)
- [Frontend Usage](#frontend-usage)
- [Development & Maintenance](#development--maintenance)
- [Troubleshooting](#troubleshooting)
- [Team](#team)

---

## Features

- User registration, login, and logout
- OAuth login via Google and GitHub
- Create, edit, and delete posts
- Comment on posts
- Like/dislike posts and comments
- Category-based organization
- Responsive UI for desktop and mobile
- Dockerized for easy deployment

---

## Tech Stack

- **Backend:** Go (Golang), SQLite
- **Frontend:** HTML, CSS, JavaScript (Vanilla)
- **Containerization:** Docker, Docker Compose

---

## Installation

### Prerequisites

- [Docker](https://www.docker.com/get-started) and Docker Compose installed
- (For development) [Go](https://golang.org/dl/) 1.18+ installed

### Clone the Repository

```sh
git clone https://github.com/yourusername/nexxus-forum.git
cd forum
```

### OAuth Environment Variables
Before running the backend you need to provide OAuth credentials:

```sh
export GOOGLE_CLIENT_ID=<your-google-client-id>
export GOOGLE_CLIENT_SECRET=<your-google-client-secret>
export GITHUB_CLIENT_ID=<your-github-client-id>
export GITHUB_CLIENT_SECRET=<your-github-client-secret>
```

## Running the Project
### Using Docker (Recommended)
### From the root directory:

```sh
docker compose up --build
```
### To build manually:
```sh
docker image build [OPTIONS] PATH | URL | -
```
### After building, list your images:
```sh
docker images
```
Example output:
```sh
student$ docker images
REPOSITORY              TAG       IMAGE ID       CREATED          SIZE
<name of the image>     latest    85a65d66ca39   7 seconds ago    795MB
```
## Run Docker Containers

Start the containers using Docker Compose (recommended):
```sh
docker compose up
```
### To run manually:
```sh
docker container run [OPTIONS] IMAGE [COMMAND] [ARG...]
```

List running and stopped containers:

```sh
docker ps -a
```

Example output:

```sh
student$ docker ps -a
CONTAINER ID   IMAGE                COMMAND        CREATED          STATUS          PORTS                    NAMES
cc8f5dcf760f   <name of the image>  "./server"     6 seconds ago    Up 6 seconds    0.0.0.0:8080->8080/tcp   forum
```

### To stop and remove containers:
```sh
docker container stop [CONTAINER_ID]
```

### This will build and start both the backend and frontend services.

Backend: http://localhost:8080

Frontend: http://localhost:8081

Local Development (Without Docker)

### Backend
```sh
cd API
go run cmd/main.go
```

### Frontend
```sh
cd ui
go run cmd/main.go
```

## API Usage

### Guest View
```sh
curl http://localhost:8080/forum/api/guest
```
### Get Categories
```sh
curl http://localhost:8080/forum/api/categories
```
### Register a New User
```sh
curl -X POST http://localhost:8080/forum/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
  ```
### Login
```sh
curl -X POST http://localhost:8080/forum/api/session/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  -c cookies.txt
```
### OAuth Login
Open one of the following URLs in your browser:

- `http://localhost:8080/forum/api/oauth/google/login`
- `http://localhost:8080/forum/api/oauth/github/login`

After granting access, the server creates a session and redirects back.
### Logout
```sh
curl -X POST http://localhost:8080/forum/api/session/logout \
  -b cookies.txt
```
### Create a Post
```sh
curl -X POST http://localhost:8080/forum/api/posts/create \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"title":"My first TITLE","content":"Hello new forum!","category_ids":[1,2]}'
```
### React to a Post or Comment
```sh
curl -X POST http://localhost:8080/forum/api/react \
  -H "Content-Type: application/json" \
  -d '{"target_id":"<TARGET_ID>","target_type":"post","reaction_type":1}' \
  -b cookies.txt
  ```

  # Frontend Usage

  Visit http://localhost:8081 in your browser.

  Register a new account or log in.

  Create posts, comment, and react to content.

  # Development & Maintenance
  Backend: All Go code is in the API/ directory. Use go mod tidy to manage dependencies.

  Frontend: Static assets and JS are in ui/static/.

  Database: The SQLite database is at API/database/forum.db.

  Configuration: API endpoints are configured in API/config/.

  ## Linting & Formatting

  Use gofmt for Go code.

  Use a code formatter (like Prettier) for JS/CSS/HTML.

  ## Testing

  Add Go tests in the API/ directory as needed.
  
  Manual testing can be done via the API or UI.

  ## Updating Dependencies

  For Go: go get -u and go mod tidy

  For frontend: update static files as needed.

# Troubleshooting

  Port conflicts: Make sure ports 8080 (backend) and 8081 (frontend, if used) are free.

  Database issues: Delete API/database/forum.db to reset the database (will lose all data).

  Docker issues: Run docker compose down -v to remove volumes and reset containers.



## Running the Project

### Using Docker (Recommended)

From the root directory:

```sh
docker compose up --build
```

This will build and start both the backend and frontend services.

- Backend: [http://localhost:8080](http://localhost:8080)
- Frontend: [http://localhost:8081](http://localhost:8081)

### Local Development (Without Docker)

#### Backend

```sh
cd API
go run cmd/main.go
```

#### Frontend

You can serve the `ui/static` directory using any static file server, or run the Go frontend if provided:

```sh
cd ui
go run cmd/main.go
```

---

## API Usage

### Guest View

```sh
curl http://localhost:8080/forum/api/guest
```

### Get Categories

```sh
curl http://localhost:8080/forum/api/categories
```

### Register a New User

```sh
curl -X POST http://localhost:8080/forum/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

### Login

```sh
curl -X POST http://localhost:8080/forum/api/session/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  -c cookies.txt
```

### Logout

```sh
curl -X POST http://localhost:8080/forum/api/session/logout \
  -b cookies.txt
```

### Create a Post

```sh
curl -X POST http://localhost:8080/forum/api/posts/create \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"title":"My first TITLE","content":"Hello new forum!","category_ids":[1,2]}'
```

### Create a Comment

```sh
curl -X POST http://localhost:8080/forum/api/comments \
  -H "Content-Type: application/json" \
  -d '{"post_id":"<POST_ID>","content":"Nice post!"}' \
  -b cookies.txt
```

### React to a Post or Comment

```sh
curl -X POST http://localhost:8080/forum/api/react \
  -H "Content-Type: application/json" \
  -d '{"target_id":"<TARGET_ID>","target_type":"post","reaction_type":1}' \
  -b cookies.txt
```

---

## Frontend Usage

- Visit [http://localhost:8081](http://localhost:8081) in your browser.
- Register a new account or log in.
- Create posts, comment, and react to content.

---
## Verifying SQLite Queries and Data Persistence
  This section guides you through verifying that the forum backend performs the expected SQL operations (CREATE, INSERT, SELECT) and that your actions in the forum UI are correctly reflected in the SQLite database.

### SQL Query Coverage

- **CREATE queries:**  
  The backend codebase contains SQL `CREATE` statements to initialize tables such as `users`, `posts`, and `comments` when the application starts or during migrations.

- **INSERT queries:**  
  When you register, create a post, or add a comment, the backend executes SQL `INSERT` statements to store your data in the respective tables.

- **SELECT queries:**  
  The backend uses SQL `SELECT` statements to fetch users, posts, comments, and other data for API responses.

### Manual Verification Steps

#### Register a User and Verify in SQLite
   **Register via API or UI:**  
   Use the registration form in the forum UI or the following curl command:
   ```sh
   curl -X POST http://localhost:8080/forum/api/register \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
  ```
### Open the SQLite database:
```sh
sqlite3 API/database/forum.db
```
### Query all users:
```sh
SELECT * FROM users;
```
Expected: The user you registered should appear in the results.

### Create a Post and Verify in SQLite
Create a post via UI or API.
Example curl:
```sh
curl -X POST http://localhost:8080/forum/api/posts/create \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"title":"My first post","content":"Hello forum!","category_ids":[1]}'
  ```
### Query all posts:
```sh
SELECT * FROM posts;
```
Expected: The post you created should appear in the results.

### Create a Comment and Verify in SQLite

Create a comment via UI or API.
Example curl:
```sh
curl -X POST http://localhost:8080/forum/api/comments \
  -H "Content-Type: application/json" \
  -d '{"post_id":"<POST_ID>","content":"Nice post!"}' \
  -b cookies.txt
  ```
### Query all comments:
```sh
SELECT * FROM comments;
```
Expected: The comment you created should appear in the results.

## Development & Maintenance

- **Backend:** All Go code is in the `API/` directory. Use `go mod tidy` to manage dependencies.
- **Frontend:** Static assets and JS are in `ui/static/`.
- **Database:** The SQLite database is at `API/database/forum.db`.
- **Configuration:** API endpoints are configured in `API/config/`.

### Linting & Formatting

- Use `gofmt` for Go code.
- Use a code formatter (like Prettier) for JS/CSS/HTML.

### Testing

- Add Go tests in the `API/` directory as needed.
- Manual testing can be done via the API or UI.

### Updating Dependencies

- For Go: `go get -u` and `go mod tidy`
- For frontend: update static files as needed.

---

## Troubleshooting

- **Port conflicts:** Make sure ports 8080 (backend) and 8081 (frontend, if used) are free.
- **Database issues:** Delete `API/database/forum.db` to reset the database (will lose all data).
- **Docker issues:** Run `docker compose down -v` to remove volumes and reset containers.

---

## Team

- [mkouvara](https://platform.zone01.gr/git/mkouvara)

- [pkalliag](https://platform.zone01.gr/git/pkalliag)

- [cemvalot](https://platform.zone01.gr/git/cemvalot)

- [gpatoula](https://platform.zone01.gr/git/gpatoula)
