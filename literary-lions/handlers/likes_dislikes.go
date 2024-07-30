package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LikeDislikeRequest struct {
	PostID    int  `json:"post_id"`
	CommentID int  `json:"comment_id"`
	IsLike    bool `json:"is_like"`
}

func LikeDislikeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			userID, err := GetUserIDFromSession(db, r)
			if err != nil {
				log.Println("No session token, redirecting to login prompt.")
				http.Redirect(w, r, "/login_prompt", http.StatusSeeOther)
				return
			}

			var req LikeDislikeRequest
			if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Println("Error decoding request body:", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			if err = processLikeDislike(db, req, userID); err != nil {
				log.Println("Error processing like/dislike:", err)
				http.Error(w, "Server error, unable to process your request.", http.StatusInternalServerError)
				return
			}

			likes, dislikes := getUpdatedCounts(db, req)

			response := map[string]interface{}{
				"post_id":    req.PostID,
				"comment_id": req.CommentID,
				"likes":      likes,
				"dislikes":   dislikes,
			}
			if err = json.NewEncoder(w).Encode(response); err != nil {
				log.Println("Error encoding response:", err)
				http.Error(w, "Server error, unable to process your request.", http.StatusInternalServerError)
				return
			}
		}
	}
}

func processLikeDislike(db *sql.DB, req LikeDislikeRequest, userID int) error {
	if req.PostID != 0 {
		return toggleLikeDislike(db, "post_id", req.PostID, userID, req.IsLike)
	} else if req.CommentID != 0 {
		return toggleLikeDislike(db, "comment_id", req.CommentID, userID, req.IsLike)
	}
	return nil
}

func getUpdatedCounts(db *sql.DB, req LikeDislikeRequest) (int, int) {
	if req.PostID != 0 {
		return getCount(db, "post_id", req.PostID, true), getCount(db, "post_id", req.PostID, false)
	}
	return getCount(db, "comment_id", req.CommentID, true), getCount(db, "comment_id", req.CommentID, false)
}

func LikedPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromSession(db, r)
		if err != nil {
			log.Println("No session token, redirecting to login prompt.")
			http.Redirect(w, r, "/login_prompt", http.StatusSeeOther)
			return
		}

		log.Println("Retrieving liked posts for user ID:", userID)

		rows, err := db.Query(`
            SELECT posts.id, posts.title, posts.content, posts.category, posts.post_type, posts.image_path, posts.created_at, users.username,
            (SELECT COUNT(*) FROM likes_dislikes WHERE post_id = posts.id AND is_like = 1) as likes,
            (SELECT COUNT(*) FROM likes_dislikes WHERE post_id = posts.id AND is_like = 0) as dislikes
            FROM likes_dislikes
            INNER JOIN posts ON likes_dislikes.post_id = posts.id
            INNER JOIN users ON posts.user_id = users.id
            WHERE likes_dislikes.user_id = ? AND likes_dislikes.is_like = 1
            ORDER BY posts.created_at DESC`, userID)
		if err != nil {
			log.Println("Error retrieving liked posts:", err)
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
			Title: "Liked Posts",
			Data:  posts,
		}

		if len(posts) == 0 {
			data.Data = []Post{}
		}

		RenderTemplate(w, r, "liked_posts.html", data)
	}
}

func toggleLikeDislike(db *sql.DB, column string, id int, userID int, isLike bool) error {
	liked, disliked := checkLikeDislikeStatus(db, column, id, userID)
	if isLike && !liked || !isLike && !disliked {
		if err := insertLikeDislike(db, column, id, userID, isLike); err != nil {
			return err
		}
		if isLike && disliked || !isLike && liked {
			return removeLikeDislike(db, column, id, userID, !isLike)
		}
	} else {
		return removeLikeDislike(db, column, id, userID, isLike)
	}
	return nil
}

func insertLikeDislike(db *sql.DB, column string, id int, userID int, isLike bool) error {
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO likes_dislikes (%s, user_id, is_like) VALUES (?, ?, ?)", column))
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, userID, isLike)
	return err
}

func removeLikeDislike(db *sql.DB, column string, id int, userID int, isLike bool) error {
	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM likes_dislikes WHERE %s = ? AND user_id = ? AND is_like = ?", column))
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, userID, isLike)
	return err
}

func checkLikeDislikeStatus(db *sql.DB, column string, id int, userID int) (liked bool, disliked bool) {
	var likes, dislikes int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM likes_dislikes WHERE %s = ? AND user_id = ? AND is_like = 1", column), id, userID).Scan(&likes)
	if err != nil {
		fmt.Println("Error scanning row:", err)
	}
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM likes_dislikes WHERE %s = ? AND user_id = ? AND is_like = 0", column), id, userID).Scan(&dislikes)
	if err != nil {
		fmt.Println("Error scanning row:", err)
	}
	return likes > 0, dislikes > 0
}

func getCount(db *sql.DB, column string, id int, isLike bool) int {
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(%s) FROM likes_dislikes WHERE %s = ? AND is_like = ?", column, column), id, isLike).Scan(&count)
	if err != nil {
		fmt.Println("Error scanning row:", err)
	}
	return count
}

func GetUserLikedPostList(db *sql.DB, userID int) (list []int) {
	rows, err := db.Query("SELECT post_id FROM likes_dislikes WHERE user_id = ? AND is_like = 1", userID)
	if err != nil {
		return list
	}
	defer rows.Close()

	for rows.Next() {
		var postID int
		err = rows.Scan(&postID)
		if err != nil {
			fmt.Println("Error scanning row:", err)
		}
		list = append(list, postID)
	}
	return list
}
