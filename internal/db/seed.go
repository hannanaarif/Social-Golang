package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/hannanaarif/Social/internal/store"
)

var usernames = []string{
	"alice", "bob", "dave", "charlie", "eve",
	"frank", "grace", "heidi", "ivan", "judy",
	"mallory", "niaj", "oscar", "peggy", "quinn",
	"rupert", "sybil", "trent", "ursula", "victor",
	"wendy", "xavier", "yvonne", "zach", "adam",
	"bella", "caleb", "diana", "ethan", "fiona",
	"george", "hannah", "isaac", "jasmine", "kevin",
	"luna", "mason", "nora", "owen", "paula",
	"quentin", "rachel", "sam", "tina", "ulrich",
	"violet", "will", "xena", "yusuf", "zoe",
}

var titles = []string{
	"Getting Started with Go",
	"Understanding Pointers",
	"REST APIs Made Simple",
	"Intro to Concurrency",
	"Clean Code Basics",
	"Mastering Slices",
	"Go vs Other Languages",
	"Structs Explained",
	"Error Handling Tips",
	"Working with JSON",
	"Build CLI Tools",
	"Testing in Go",
	"Maps Deep Dive",
	"Interfaces Simplified",
	"Go Routines Guide",
	"Dependency Injection",
	"Logging Best Practices",
	"File Handling in Go",
	"Writing Middleware",
	"Optimizing Performance",
}

var contents = []string{
	"Go makes backend development simple, fast, and highly efficient for scalable systems.",
	"Concurrency in Go allows programs to handle thousands of tasks simultaneously.",
	"Understanding slices is essential for writing idiomatic and memory-efficient Go code.",
	"Error handling in Go encourages explicit and predictable program behavior.",
	"Structs help organize data logically and make code easier to maintain.",
	"Interfaces enable flexible and decoupled software design patterns.",
	"Testing ensures your application remains stable during future changes.",
	"Middleware is useful for logging, authentication, and request validation.",
	"Go routines are lightweight threads managed by the Go runtime.",
	"Channels help communicate safely between concurrent processes.",
	"Clean architecture improves readability and long-term maintainability.",
	"JSON handling is straightforward with Go’s encoding libraries.",
	"Dependency injection improves testability and modular design.",
	"Logging is crucial for debugging and monitoring production systems.",
	"Proper folder structure keeps projects organized and scalable.",
	"Using context helps manage timeouts and cancellations effectively.",
	"Database migrations keep schema changes consistent across environments.",
	"Writing reusable packages reduces duplicated logic.",
	"Profiling helps identify performance bottlenecks in applications.",
	"Consistent naming conventions make code easier to understand.",
}

var tags = []string{
	"go", "golang", "backend", "api", "rest",
	"programming", "coding", "webdev", "database", "sql",
	"concurrency", "performance", "testing", "debugging", "architecture",
	"microservices", "json", "security", "deployment", "scalability",
}

var comments = []string{
	"Great article, really helped me understand the concept!",
	"This was super clear and easy to follow.",
	"I tried this and it worked perfectly.",
	"Thanks for sharing such useful information.",
	"Can you explain this part in more detail?",
	"This saved me hours of debugging.",
	"Very well written and beginner friendly.",
	"I like how you broke everything down step by step.",
	"This is exactly what I was looking for.",
	"Awesome explanation, learned something new today!",
	"I appreciate the practical examples.",
	"This topic always confused me until now.",
	"Looking forward to more posts like this.",
	"Clean and simple explanation.",
	"I tested this approach and it’s really efficient.",
	"Nice write-up, keep it up!",
	"This clarified many doubts I had.",
	"Short, simple, and very helpful.",
	"Your examples make it easy to understand.",
	"Fantastic content, thanks for posting!",
}

func Seed(store *store.Storage, db *sql.DB) {
	ctx := context.Background()
	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("error while creating user", err)
			return
		}
	}
	posts := generatePosts(100, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("error while creating post", err)
			return
		}
	}

	comments := generateComments(100, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("error while creating comment", err)
			return
		}
	}

	for _, user := range users {
		followerID := users[rand.Intn(len(users))].ID
		if followerID == user.ID {
			continue // avoid self follow
		}
		if err := store.Followers.Follow(ctx, followerID, user.ID); err != nil {
			if err.Error() != "resource already exists" {
				log.Println("error while following", err)
			}
		}
	}

	log.Println("Seeding completed successfully")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@okmail.com",
			Password: "12341234",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	// TODO: Implement post generation logic
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cmt := make([]*store.Comment, num)
	// TODO: Implement comment generation logic
	for i := 0; i < num; i++ {
		cmt[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cmt
}
