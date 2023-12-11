package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"postgresql-blog/database"
	"postgresql-blog/models"
	"postgresql-blog/repository"
	"postgresql-blog/service"

	"gorm.io/gorm"
)

func main() {
	displayMenu()
}

func displayMenu() {
	// Initialize the database and repositories
	db, errDb := database.RunDatabase()
	if errDb != nil {
		fmt.Println("Error setting up the database:", errDb)
		os.Exit(1)
	}
	// close db connection
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Println("Error getting underlying *sql.DB instance:", err)
			os.Exit(1)
		}
		sqlDB.Close()
	}()

	// Migrate the database
	// MigrateDatabase(db)

	// show menus
	fmt.Printf("\n\n")
	fmt.Println("===========================================================================")
	fmt.Println("Hello! Please choose one of the options below:")
	fmt.Println("A. Users")
	fmt.Println("B. Posts")
	fmt.Println("C. Comments")
	fmt.Println("D. Exit")
	fmt.Println("===========================================================================")
	fmt.Println("Please choose one of the options above by typing the letter (A/B/C):")

	// read input
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("===========================================================================")
		os.Exit(1)
	}

	// error handling input
	if input == "" {
		fmt.Println("Empty input. Please try once more!")
		fmt.Println("===========================================================================")
		displayMenu() // Recursively call the menu to try again
		return
	}

	switch input {
	case "A", "a":
		displayUser(db)
	case "B", "b":
		displayPost(db)
	case "C", "c":
		displayComment(db)
	case "D", "d":
		fmt.Println("Exited!")
		os.Exit(0)
	default:
		fmt.Println("Invalid input. Please try once more!")
		fmt.Println("===========================================================================")
		displayMenu() // Recursively call the menu to try again
	}
}

func displayUser(db *gorm.DB) {
	// Create a repository instance and provide it to the service
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, db)

	fmt.Printf("\n")
	fmt.Println("Showing Options for Users")
	fmt.Println("===========================================================================")
	fmt.Println("A. See all users")
	fmt.Println("B. See a user")
	fmt.Println("C. Add a new user")
	fmt.Println("D. Update a user")
	fmt.Println("E. Delete a user")
	fmt.Println("F. Exit")
	fmt.Println("===========================================================================")
	fmt.Println("Please choose one of the options above by typing the letter (A/B/C/D/E/F):")

	// read input
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("===========================================================================")
		os.Exit(1)
	}

	switch input {
	case "A", "a":
		fmt.Println("\n---------------------------------------------------------------------------")
		fmt.Println("Get all users")
		fmt.Println("---------------------------------------------------------------------------")
		GetAllUsers(*userService)
		fmt.Println("---------------------------------------------------------------------------")
		displayUser(db)
	case "B", "b":
		fmt.Println("\nGet a user by email")
		fmt.Println("===========================================================================")
		GetUserByEmail(*userService)
		displayUser(db)
	case "C", "c":
		fmt.Println("\nCreate a new user")
		fmt.Println("===========================================================================")
		CreateUser(*userService)
		displayUser(db)
	case "D", "d":
		fmt.Println("\nUpdate a user")
		fmt.Println("===========================================================================")
		UpdateUser(*userService)
		displayUser(db)
	case "E", "e":
		fmt.Println("\nDelete a user")
		fmt.Println("===========================================================================")
		DeleteUser(*userService)
		displayUser(db)
	case "F", "f":
		fmt.Println("Exited!")
		displayMenu()
	default:
		fmt.Println("Invalid input. Please try once more!")
		fmt.Println("===========================================================================")
		displayUser(db) // Recursively call the menu to try again
	}
}

func displayPost(db *gorm.DB) {
	// Create a repository instance and provide it to the service
	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository, db)

	username, password := readUsernameAndPassword()

	fmt.Println("---------------------------------------------------------------------------")
	userid := GetUserByUsernameAndPassword(db, username, password)
	if userid == 0 {
		fmt.Println("Username or Password not found. Please try again!")
		fmt.Println("---------------------------------------------------------------------------")
		displayPost(db)
		return
	}
	fmt.Println("---------------------------------------------------------------------------")
	displayPostSubmenu(db, *postService, userid)
}

func displayPostSubmenu(db *gorm.DB, postService service.PostService, userid int64) {
	fmt.Printf("\n")
	fmt.Println("Showing Options for Posts")
	fmt.Println("===========================================================================")
	fmt.Println("A. See all posts")
	fmt.Println("B. See your posts")
	fmt.Println("C. Search a post by title")
	fmt.Println("D. Add a new post")
	fmt.Println("E. Update a post")
	fmt.Println("F. Delete a post")
	fmt.Println("G. Exit")
	fmt.Println("===========================================================================")
	fmt.Println("Please choose one of the options above by typing the letter (A/B/C/D/E/F/G):")

	// read input
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("===========================================================================")
		os.Exit(1)
	}

	switch input {
	case "A", "a":
		fmt.Println("\n---------------------------------------------------------------------------")
		fmt.Println("Get all posts")
		fmt.Println("---------------------------------------------------------------------------")
		GetAllPosts(postService)
		fmt.Println("---------------------------------------------------------------------------")
		displayPostSubmenu(db, postService, userid)
	case "B", "b":
		fmt.Println("\n---------------------------------------------------------------------------")
		fmt.Println("Get your posts")
		fmt.Println("---------------------------------------------------------------------------")
		GetUserPosts(postService, userid)
		fmt.Println("---------------------------------------------------------------------------")
		displayPostSubmenu(db, postService, userid)
	case "C", "c":
		fmt.Println("\nGet a post by title")
		fmt.Println("===========================================================================")
		GetPostByTitle(postService)
		displayPostSubmenu(db, postService, userid)
	case "D", "d":
		fmt.Println("\nCreate a post")
		fmt.Println("===========================================================================")
		CreatePost(postService, userid)
		displayPostSubmenu(db, postService, userid)
	case "E", "e":
		fmt.Println("\nUpdate a post")
		fmt.Println("===========================================================================")
		UpdatePost(postService, userid)
		displayPostSubmenu(db, postService, userid)
	case "F", "f":
		fmt.Println("\nDelete a post")
		fmt.Println("===========================================================================")
		DeletePost(postService, userid)
		displayPostSubmenu(db, postService, userid)
	case "G", "g":
		fmt.Println("Exited!")
		displayMenu()
	default:
		fmt.Println("Invalid input. Please try once more!")
		fmt.Println("===========================================================================")
		displayPostSubmenu(db, postService, userid) // Recursively call the menu to try again
	}
}

func displayComment(db *gorm.DB) {
	commentRepository := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepository, db)

	username, password := readUsernameAndPassword()

	fmt.Println("---------------------------------------------------------------------------")
	userid := GetUserByUsernameAndPassword(db, username, password)
	if userid == 0 {
		fmt.Println("Username or Password not found. Please try again!")
		fmt.Println("---------------------------------------------------------------------------")
		displayComment(db)
		return
	}
	fmt.Println("---------------------------------------------------------------------------")
	displayCommentSubmenu(db, *commentService, userid)
}

func displayCommentSubmenu(db *gorm.DB, commentService service.CommentService, userid int64) {
	fmt.Printf("\n")
	fmt.Println("Showing Options for Comments")
	fmt.Println("===========================================================================")
	fmt.Println("A. See all comments")
	fmt.Println("B. See your comments")
	fmt.Println("C. See comments by post")
	fmt.Println("D. Add a new comment")
	fmt.Println("E. Update a comment")
	fmt.Println("F. Delete a comment")
	fmt.Println("G. Exit")
	fmt.Println("===========================================================================")
	fmt.Println("Please choose one of the options above by typing the letter (A/B/C/D/E/F/G):")

	// read input
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	switch input {
	case "A", "a":
		fmt.Println("\n---------------------------------------------------------------------------")
		fmt.Println("Get all comments")
		fmt.Println("---------------------------------------------------------------------------")
		GetAllComments(commentService)
		fmt.Println("---------------------------------------------------------------------------")
		displayCommentSubmenu(db, commentService, userid)
	case "B", "b":
		fmt.Println("\n---------------------------------------------------------------------------")
		fmt.Println("Get your comments")
		fmt.Println("---------------------------------------------------------------------------")
		GetUserComments(commentService, userid)
		fmt.Println("---------------------------------------------------------------------------")
		displayCommentSubmenu(db, commentService, userid)
	case "C", "c":
		fmt.Println("\nGet comments by post")
		fmt.Println("===========================================================================")
		postid := GetPostForComment(db)
		GetPostComments(commentService, postid)
		displayCommentSubmenu(db, commentService, userid)
	case "D", "d":
		fmt.Println("\nCreate a comment")
		fmt.Println("===========================================================================")
		postid := GetPostForComment(db)
		CreateComment(commentService, userid, postid)
		displayCommentSubmenu(db, commentService, userid)
	case "E", "e":
		fmt.Println("\nUpdate a comment")
		fmt.Println("===========================================================================")
		UpdateComment(commentService, userid)
		displayCommentSubmenu(db, commentService, userid)
	case "F", "f":
		fmt.Println("\nDelete a comment")
		fmt.Println("===========================================================================")
		DeleteComment(commentService, userid)
		displayCommentSubmenu(db, commentService, userid)
	case "G", "g":
		fmt.Println("Exited!")
		displayMenu()
	default:
		fmt.Println("Invalid input. Please try once more!")
		displayCommentSubmenu(db, commentService, userid)
	}
}

func readUsernameAndPassword() (string, string) {
	fmt.Printf("\nPlease type the username and password to continue\n")
	fmt.Println("===========================================================================")
	reader := bufio.NewReader(os.Stdin)

	var username, password string

	fmt.Print("Enter Username: ")
	_, err := fmt.Scan(&username)
	username = strings.TrimSpace(username)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return "", ""
	}
	reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return "", ""
	}
	reader.ReadString('\n')

	return username, password
}

func MigrateDatabase(db *gorm.DB) {
	ctx := context.Background()
	dbIn := database.NewPostgreSQLGORMRepository(db)
	if err := dbIn.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("MIGRATE REPOSITORY")
}

func CreateUser(userService service.UserService) {
	var newUser models.User
	reader := bufio.NewReader(os.Stdin)

	// ERROR TERUS COK
	fmt.Print("Enter Name: ")
	_, err := fmt.Scan(&newUser.Name)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	fmt.Print("Enter Email: ")
	_, err = fmt.Scan(&newUser.Email)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	_, err = fmt.Scan(&newUser.Password)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	fmt.Print("Enter Username: ")
	_, err = fmt.Scan(&newUser.Username)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	fmt.Println("---------------------------------------------------------------------------")
	// Call the user creation service method
	createdUser, err := userService.CreateUser(context.Background(), newUser)
	if err != nil {
		fmt.Println("Error creating user:", err)
	} else {
		fmt.Println("User created successfully!")
	}
	fmt.Println("Created record:")
	fmt.Printf("ID:       %d\n", createdUser.ID)
	fmt.Printf("Name:     %s\n", createdUser.Name)
	fmt.Printf("Email:    %s\n", createdUser.Email)
	fmt.Printf("Username: %s\n", createdUser.Username)
	fmt.Println("---------------------------------------------------------------------------")
}

func GetAllUsers(userService service.UserService) {
	all, err := userService.GetAllUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range all {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Username: %s\n", user.ID, user.Name, user.Email, user.Username)
	}
}

func GetUserByEmail(userService service.UserService) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Email to find user: ")

	var email string
	var err error
	_, err = fmt.Scan(&email)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	email = strings.TrimSpace(email)
	reader.ReadString('\n')

	fmt.Println("---------------------------------------------------------------------------")
	user, err := userService.GetUserByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			fmt.Printf("User with email '%s' does not exist in the repository\n", email)
		} else {
			fmt.Printf("Error finding user by email '%s': %v\n", email, err)
		}
	} else {
		fmt.Printf("User found:\nID: %d, Name: %s, Email: %s, Username: %s\n", user.ID, user.Name, user.Email, user.Username)
	}
	fmt.Println("---------------------------------------------------------------------------")
}

func GetUserByUsernameAndPassword(db *gorm.DB, username, password string) int64 {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, db)

	user, err := userService.GetUserByUsernameAndPassword(context.Background(), username, password)
	if err != nil {
		fmt.Printf("Error getting user by username and password: %v\n", err)
		return 0
	} else {
		fmt.Printf("User found:\nID: %d, Name: %s, Email: %s, Username: %s\n", user.ID, user.Name, user.Email, user.Username)
		return user.ID
	}
}

func UpdateUser(userService service.UserService) {
	reader := bufio.NewReader(os.Stdin)

	GetAllUsers(userService)
	fmt.Println("---------------------------------------------------------------------------")

	fmt.Print("Enter the ID of the user you want to update: ")
	var userID int64
	_, err := fmt.Scan(&userID)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	// Check if the user with the entered ID exists
	user, err := userService.GetUserByID(context.Background(), userID)
	if err != nil {
		fmt.Printf("Error finding user by ID %d: %v\n", userID, err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}

	fmt.Printf("User with ID %d found!\n", user.ID)
	fmt.Println("---------------------------------------------------------------------------")
	fmt.Printf("Enter new values for the user with ID %d:\n", user.ID)

	var updatedUser models.User

	fmt.Print("New Name: ")
	_, err = fmt.Scan(&updatedUser.Name)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	fmt.Print("New Email: ")
	_, err = fmt.Scan(&updatedUser.Email)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	fmt.Print("New Password: ")
	_, err = fmt.Scan(&updatedUser.Password)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	fmt.Print("New Username: ")
	_, err = fmt.Scan(&updatedUser.Username)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	updatedUser.ID = userID

	fmt.Println("---------------------------------------------------------------------------")
	_, err = userService.UpdateUserByID(context.Background(), updatedUser)
	if err != nil {
		fmt.Printf("Error updating user with ID %d: %v\n", userID, err)
	} else {
		fmt.Println("User updated successfully!")
	}
	fmt.Println("---------------------------------------------------------------------------")
}

func DeleteUser(userService service.UserService) {
	reader := bufio.NewReader(os.Stdin)

	GetAllUsers(userService)
	fmt.Println("---------------------------------------------------------------------------")

	fmt.Print("Enter the ID of the user you want to delete: ")
	var userID int64
	_, err := fmt.Scan(&userID)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	// Check if the user with the entered ID exists
	user, err := userService.GetUserByID(context.Background(), userID)
	if err != nil {
		fmt.Printf("Error finding user by ID %d: %v\n", userID, err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	fmt.Printf("User with ID %d found!\n", user.ID)
	fmt.Println("---------------------------------------------------------------------------")

	fmt.Printf("Are you sure you want to delete the user with ID %d? (Y/N): ", user.ID)
	var confirmation string
	_, err = fmt.Scan(&confirmation)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	if confirmation == "Y" || confirmation == "y" {
		err = userService.DeleteUserByID(context.Background(), userID)
		if err != nil {
			fmt.Printf("Error deleting user with ID %d: %v\n", userID, err)
		} else {
			fmt.Println("User deleted successfully!")
		}
	} else {
		fmt.Println("User not deleted!")
	}
	fmt.Println("---------------------------------------------------------------------------")
}

func CreatePost(postService service.PostService, userid int64) {
	// Prompt the user to enter post information
	var newPost models.Post
	newPost.UserID = uint64(userid)
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Title: ")
	_, err := fmt.Scan(&newPost.Title)
	newPost.Title = strings.TrimSpace(newPost.Title)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	fmt.Print("Enter Content: ")
	_, err = fmt.Scan(&newPost.Content)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	fmt.Println("---------------------------------------------------------------------------")
	// Call the post creation service method
	createdPost, err := postService.CreatePost(context.Background(), newPost)
	if err != nil {
		fmt.Println("Error creating post:", err)
	} else {
		fmt.Println("Post created successfully!")
	}
	fmt.Println("Created record:")
	fmt.Printf("ID:		%d\n", createdPost.ID)
	fmt.Printf("User ID:	%d\n", createdPost.UserID)
	fmt.Printf("Title:		%s\n", createdPost.Title)
	fmt.Printf("Content:	%s\n", createdPost.Content)
	fmt.Println("---------------------------------------------------------------------------")
}

func GetAllPosts(postService service.PostService) {
	all, err := postService.GetAllPosts(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range all {
		fmt.Printf("ID: %d, User ID: %d, Title: %s, Content: %s\n", post.ID, post.UserID, post.Title, post.Content)
	}
}

func GetUserPosts(postService service.PostService, userid int64) {
	all, err := postService.GetPostByUserID(context.Background(), userid)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range all {
		fmt.Printf("ID: %d, Title: %s, Content: %s\n", post.ID, post.Title, post.Content)
	}
}

func GetPostByTitle(postService service.PostService) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Title to find post: ")

	var title string
	var err error
	_, err = fmt.Scan(&title)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	title = strings.TrimSpace(title)
	fmt.Println("---------------------------------------------------------------------------")
	reader.ReadString('\n')

	post, err := postService.GetPostByTitle(context.Background(), title)
	if err != nil {
		fmt.Printf("Error finding post by title '%s': %v\n", title, err)
	} else {
		fmt.Printf("Post found by title '%s':\n", title)
		fmt.Printf("ID: %d, User ID: %d, Title: %s, Content: %s\n", post.ID, post.UserID, post.Title, post.Content)
	}
	fmt.Println("---------------------------------------------------------------------------")
}

func UpdatePost(postService service.PostService, userid int64) {
	reader := bufio.NewReader(os.Stdin)

	GetUserPosts(postService, userid)
	fmt.Println("---------------------------------------------------------------------------")

	fmt.Print("Enter the ID of the post you want to update: ")
	var postID int64
	_, err := fmt.Scan(&postID)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	// Check if the post with the entered ID exists
	post, err := postService.GetPostByID(context.Background(), postID)
	if err != nil {
		fmt.Printf("Error finding post by ID %d: %v\n", postID, err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}

	fmt.Printf("Post with ID %d found!\n", post.ID)
	fmt.Println("---------------------------------------------------------------------------")
	fmt.Printf("Enter new values for the post with ID %d:\n", post.ID)

	var updatedPost models.Post

	fmt.Print("New Title: ")
	_, err = fmt.Scan(&updatedPost.Title)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	fmt.Print("New Content: ")
	_, err = fmt.Scan(&updatedPost.Content)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	updatedPost.ID = postID
	updatedPost.UserID = uint64(userid)

	fmt.Println("---------------------------------------------------------------------------")
	_, err = postService.UpdatePostByID(context.Background(), updatedPost)
	if err != nil {
		fmt.Printf("Error updating post with ID %d: %v\n", postID, err)
	} else {
		fmt.Println("Post updated successfully!")
	}
	fmt.Println("---------------------------------------------------------------------------")
}

func DeletePost(postService service.PostService, userid int64) {
	reader := bufio.NewReader(os.Stdin)

	GetUserPosts(postService, userid)
	fmt.Println("---------------------------------------------------------------------------")

	fmt.Print("Enter the ID of the post you want to delete: ")
	var postID int64
	_, err := fmt.Scan(&postID)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	// Check if the post with the entered ID exists
	post, err := postService.GetPostByID(context.Background(), postID)
	if err != nil {
		fmt.Printf("Error finding post by ID %d: %v\n", postID, err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}

	fmt.Printf("Post with ID %d found!\n", post.ID)
	fmt.Println("---------------------------------------------------------------------------")

	fmt.Printf("Are you sure you want to delete the post with ID %d? (Y/N): ", post.ID)
	var confirmation string
	_, err = fmt.Scan(&confirmation)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	if confirmation == "Y" || confirmation == "y" {
		err = postService.DeletePostByID(context.Background(), postID)
		if err != nil {
			fmt.Printf("Error deleting post with ID %d: %v\n", postID, err)
		} else {
			fmt.Println("Post deleted successfully!")
		}
	} else {
		fmt.Println("Post not deleted!")
	}
	fmt.Println("---------------------------------------------------------------------------")
}

func GetPostForComment(db *gorm.DB) int64 {
	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository, db)

	GetAllPosts(*postService)
	fmt.Println("---------------------------------------------------------------------------")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the ID of the post: ")
	var postID int64
	_, err := fmt.Scan(&postID)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return 0
	}
	reader.ReadString('\n')

	// Check if the post with the entered ID exists
	post, err := postService.GetPostByID(context.Background(), postID)
	if err != nil {
		fmt.Printf("Error finding post by ID %d: %v\n", postID, err)
		fmt.Println("---------------------------------------------------------------------------")
		return 0
	}
	fmt.Printf("Post with ID %d found!\n", post.ID)
	fmt.Println("---------------------------------------------------------------------------")
	return post.ID
}

func CreateComment(commentService service.CommentService, userid int64, postid int64) {
	// Prompt the user to enter comment information
	var newComment models.Comment
	newComment.UserID = uint64(userid)
	newComment.PostID = uint64(postid)
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Content: ")
	_, err := fmt.Scan(&newComment.Content)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	// Call the comment creation service method
	createdComment, err := commentService.CreateComment(context.Background(), newComment)
	if err != nil {
		fmt.Println("Error creating comment:", err)
		fmt.Println("---------------------------------------------------------------------------")
	} else {
		fmt.Println("Comment created successfully!")
		fmt.Println("---------------------------------------------------------------------------")
	}
	fmt.Printf("created record: %+v\n", createdComment)
}

func GetAllComments(commentService service.CommentService) {
	all, err := commentService.GetAllComments(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, comment := range all {
		fmt.Printf("ID: %d, User ID: %d, Post ID: %d, Content: %s\n", comment.ID, comment.UserID, comment.PostID, comment.Content)
	}
}

func GetUserComments(commentService service.CommentService, userid int64) {
	all, err := commentService.GetCommentByUserID(context.Background(), userid)
	if err != nil {
		log.Fatal(err)
	}

	for _, comment := range all {
		fmt.Printf("ID: %d, Post ID: %d, Content: %s\n", comment.ID, comment.PostID, comment.Content)
	}
}

func GetPostComments(commentService service.CommentService, postid int64) {
	all, err := commentService.GetCommentByPostID(context.Background(), postid)
	if err != nil {
		log.Fatal(err)
	}

	for _, comment := range all {
		fmt.Printf("ID: %d, User ID: %d, Content: %s\n", comment.ID, comment.UserID, comment.Content)
	}
}

func UpdateComment(commentService service.CommentService, userid int64) {
	reader := bufio.NewReader(os.Stdin)

	GetUserComments(commentService, userid)
	fmt.Println("---------------------------------------------------------------------------")

	fmt.Print("Enter the ID of the comment you want to update: ")
	var commentID int64
	_, err := fmt.Scan(&commentID)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	// Check if the comment with the entered ID exists
	comment, err := commentService.GetCommentByID(context.Background(), commentID)
	if err != nil {
		fmt.Printf("Error finding comment by ID %d: %v\n", commentID, err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}

	fmt.Printf("Comment with ID %d found!\n", comment.ID)
	fmt.Println("---------------------------------------------------------------------------")
	fmt.Printf("Enter new values for the comment with ID %d:\n", comment.ID)

	var updatedComment models.Comment

	fmt.Print("New Content: ")
	_, err = fmt.Scan(&updatedComment.Content)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	updatedComment.ID = commentID
	updatedComment.UserID = uint64(userid)

	fmt.Println("---------------------------------------------------------------------------")
	_, err = commentService.UpdateCommentByID(context.Background(), updatedComment)
	if err != nil {
		fmt.Printf("Error updating comment with ID %d: %v\n", commentID, err)
	} else {
		fmt.Println("Comment updated successfully!")
	}
	fmt.Println("---------------------------------------------------------------------------")
}

func DeleteComment(commentService service.CommentService, userid int64) {
	reader := bufio.NewReader(os.Stdin)

	GetUserComments(commentService, userid)
	fmt.Println("---------------------------------------------------------------------------")

	fmt.Print("Enter the ID of the comment you want to delete: ")
	var commentID int64
	_, err := fmt.Scan(&commentID)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	// Check if the comment with the entered ID exists
	comment, err := commentService.GetCommentByID(context.Background(), commentID)
	if err != nil {
		fmt.Printf("Error finding comment by ID %d: %v\n", commentID, err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}

	fmt.Printf("Comment with ID %d found!\n", comment.ID)
	fmt.Println("---------------------------------------------------------------------------")

	fmt.Printf("Are you sure you want to delete the comment with ID %d? (Y/N): ", comment.ID)
	var confirmation string
	_, err = fmt.Scan(&confirmation)
	if err != nil {
		fmt.Println("Error reading input:", err)
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	reader.ReadString('\n')

	if confirmation == "Y" || confirmation == "y" {
		err = commentService.DeleteCommentByID(context.Background(), commentID)
		if err != nil {
			fmt.Printf("Error deleting comment with ID %d: %v\n", commentID, err)
		} else {
			fmt.Println("Comment deleted successfully!")
		}
	} else {
		fmt.Println("Comment not deleted!")
	}
	fmt.Println("---------------------------------------------------------------------------")
}
