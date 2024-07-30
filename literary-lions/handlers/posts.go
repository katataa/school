package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Post struct {
	ID        int
	Title     string
	Content   string
	Category  string
	PostType  sql.NullString
	ImagePath sql.NullString
	CreatedAt string
	Username  string
	Comments  []Comment
	Likes     int
	Dislikes  int
}

type Comment struct {
	ID        int
	Content   string
	CreatedAt string
	Username  string
	Likes     int
	Dislikes  int
}

func DisplayPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Retrieving posts from the database")
		rows, err := db.Query(`
            SELECT posts.id, posts.title, posts.content, posts.category, posts.post_type, posts.image_path, posts.created_at, users.username,
            (SELECT COUNT(*) FROM likes_dislikes WHERE post_id = posts.id AND is_like = 1) as likes,
            (SELECT COUNT(*) FROM likes_dislikes WHERE post_id = posts.id AND is_like = 0) as dislikes
            FROM posts
            INNER JOIN users ON posts.user_id = users.id
            ORDER BY posts.created_at DESC`)
		if err != nil {
			log.Println("Error retrieving posts:", err)
			http.Error(w, "Server error, unable to retrieve posts.", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var post Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.PostType, &post.ImagePath, &post.CreatedAt, &post.Username, &post.Likes, &post.Dislikes); err != nil {
				log.Println("Error scanning post:", err)
				http.Error(w, "Server error, unable to retrieve posts.", http.StatusInternalServerError)
				return
			}

			log.Printf("Retrieving comments for post ID %d", post.ID)
			commentsRows, err := db.Query(`
                SELECT comments.id, comments.content, comments.created_at, users.username,
                (SELECT COUNT(*) FROM likes_dislikes WHERE comment_id = comments.id AND is_like = 1) as likes,
                (SELECT COUNT(*) FROM likes_dislikes WHERE comment_id = comments.id AND is_like = 0) as dislikes
                FROM comments
                INNER JOIN users ON comments.user_id = users.id
                WHERE comments.post_id = ?
                ORDER BY comments.created_at ASC`, post.ID)
			if err != nil {
				log.Println("Error retrieving comments:", err)
				http.Error(w, "Server error, unable to retrieve comments.", http.StatusInternalServerError)
				return
			}
			defer commentsRows.Close()

			var comments []Comment
			for commentsRows.Next() {
				var comment Comment
				if err := commentsRows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.Username, &comment.Likes, &comment.Dislikes); err != nil {
					log.Println("Error scanning comment:", err)
					http.Error(w, "Server error, unable to retrieve comments.", http.StatusInternalServerError)
					return
				}
				comments = append(comments, comment)
			}
			post.Comments = comments

			posts = append(posts, post)
		}

		if err = rows.Err(); err != nil {
			log.Println("Error iterating over rows:", err)
			http.Error(w, "Server error, unable to retrieve posts.", http.StatusInternalServerError)
			return
		}

		log.Println("Passing Posts to Template:", posts)
		templateData := TemplateData{
			Data: posts,
		}
		RenderTemplate(w, r, "discussion.html", templateData)
	}
}

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromSession(db, r)
		if err != nil {
			log.Println("No session token, redirecting to login prompt.")
			http.Redirect(w, r, "/login_prompt", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			title := r.FormValue("title")
			content := r.FormValue("content")
			category := r.FormValue("category")
			postType := r.FormValue("post_type")

			trimmedTitle := strings.TrimSpace(title)
			trimmedContent := strings.TrimSpace(content)

			if trimmedTitle == "" || trimmedContent == "" {
				log.Println("Title or content is empty or whitespace-only.")
				http.Redirect(w, r, "/discussion", http.StatusSeeOther)
				return
			}

			var imagePath sql.NullString
			file, handler, err := r.FormFile("image")
			if err == nil {
				defer file.Close()
				if handler.Filename != "" {
					imagePath.String = fmt.Sprintf("static/uploads/%d_%s", userID, handler.Filename)
					imagePath.Valid = true

					f, err := os.OpenFile(imagePath.String, os.O_WRONLY|os.O_CREATE, 0666)
					if err != nil {
						log.Println("Error saving image:", err)
						http.Error(w, "Server error, unable to save image.", http.StatusInternalServerError)
						return
					}
					defer f.Close()
					_, err = io.Copy(f, file)
					if err != nil {
						log.Println("Error copying image:", err)
						http.Error(w, "Server error, unable to save image.", http.StatusInternalServerError)
						return
					}
				}
			}

			query := "INSERT INTO posts(user_id, title, content, category, post_type, created_at"
			values := "VALUES(?, ?, ?, ?, ?, ?"
			args := []interface{}{userID, trimmedTitle, trimmedContent, category, postType, time.Now().Format(time.RFC3339)}

			if imagePath.Valid {
				query += ", image_path"
				values += ", ?"
				args = append(args, imagePath.String)
			}

			query += ") " + values + ")"

			stmt, err := db.Prepare(query)
			if err != nil {
				log.Println("Error preparing statement for creating post:", err)
				http.Error(w, "Server error, unable to create your post.", http.StatusInternalServerError)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(args...)
			if err != nil {
				log.Println("Error executing statement for creating post:", err)
				http.Error(w, "Server error, unable to create your post.", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/discussion", http.StatusSeeOther)
			return
		}

		RenderTemplate(w, r, "create_post.html", nil)
	}
}

func CreateCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			userID, err := GetUserIDFromSession(db, r)
			if err != nil {
				log.Println("No session token, redirecting to login prompt.")
				http.Redirect(w, r, "/login_prompt", http.StatusSeeOther)
				return
			}

			postID := r.FormValue("post_id")
			content := r.FormValue("content")

			trimmedContent := strings.TrimSpace(content)
			if trimmedContent == "" {
				log.Println("Comment content is empty or whitespace-only.")
				http.Error(w, "Comment cannot be empty or just whitespace.", http.StatusBadRequest)
				return
			}

			stmt, err := db.Prepare("INSERT INTO comments(post_id, user_id, content, created_at) VALUES(?, ?, ?, ?)")
			if err != nil {
				log.Println("Error preparing statement for creating comment:", err)
				http.Redirect(w, r, "/login_prompt", http.StatusSeeOther)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(postID, userID, trimmedContent, time.Now().Format(time.RFC3339))
			if err != nil {
				log.Println("Error executing statement for creating comment:", err)
				http.Redirect(w, r, "/login_prompt", http.StatusSeeOther)
				return
			}

			http.Redirect(w, r, "/discussion", http.StatusSeeOther)
			return
		}

		RenderTemplate(w, r, "discussion.html", nil)
	}
}

func MyPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromSession(db, r)
		if err != nil {
			log.Println("No session token, redirecting to login prompt.")
			http.Redirect(w, r, "/login_prompt", http.StatusSeeOther)
			return
		}

		log.Println("Retrieving posts for user ID:", userID)

		rows, err := db.Query(`
            SELECT posts.id, posts.title, posts.content, posts.category, posts.post_type, posts.image_path, posts.created_at, users.username,
            (SELECT COUNT(*) FROM likes_dislikes WHERE post_id = posts.id AND is_like = 1) as likes,
            (SELECT COUNT(*) FROM likes_dislikes WHERE post_id = posts.id AND is_like = 0) as dislikes
            FROM posts
            INNER JOIN users ON posts.user_id = users.id
            WHERE posts.user_id = ?
            ORDER BY posts.created_at DESC`, userID)
		if err != nil {
			log.Println("Error retrieving user's posts:", err)
			http.Error(w, "Server error, unable to retrieve posts.", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var post Post
			var createdAt string
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.PostType, &post.ImagePath, &createdAt, &post.Username, &post.Likes, &post.Dislikes); err != nil {
				log.Println("Error scanning post:", err)
				http.Error(w, "Server error, unable to retrieve posts.", http.StatusInternalServerError)
				return
			}
			post.CreatedAt = createdAt
			posts = append(posts, post)
		}

		data := TemplateData{
			Title: "My Posts",
			Data:  posts,
		}

		if len(posts) == 0 {
			data.Data = []Post{}
		}

		RenderTemplate(w, r, "my_posts.html", data)
	}
}
