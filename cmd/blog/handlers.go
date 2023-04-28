package main

import (
	"database/sql"
	_ "database/sql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type indexPage struct {
	Title         string
	FeaturedPosts []featuredPostData
	MostRecent    []mostRecentData
}

type postPageData struct {
	Title    string `db:"title"`
	SubTitle string `db:"subtitle"`
	Content  string `db:"content"`
}

type featuredPostData struct {
	ID          string `db:"post_ID"`
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	ImgModifier string `db:"image_url"`
	Author      string `db:"author"`
	AuthorImg   string `db:"author_url"`
	PublishDate string `db:"publish_date"`
	PostURL     string
}

type mostRecentData struct {
	ID          string `db:"post_ID"`
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	ImgModifier string `db:"image_url"`
	Author      string `db:"author"`
	AuthorImg   string `db:"author_url"`
	PublishDate string `db:"publish_date"`
	PostURL     string
}

func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		featuredPosts, err := featuredPosts(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		mostPosts, err := mostRecent(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
			log.Println(err)
			return // Не забываем завершить выполнение ф-ии
		}

		data := indexPage{
			Title:         "Escape",
			FeaturedPosts: featuredPosts,
			MostRecent:    mostPosts,
		}

		err = ts.Execute(w, data) // Заставляем шаблонизатор вывести шаблон в тело ответа
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func post(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := mux.Vars(r)["postID"]
		log.Println(r)
		log.Println(postIDStr)
		postID, err := strconv.Atoi(postIDStr)

		if err != nil {
			http.Error(w, "Invalid order id", 403)
			log.Println(err)
			return
		}

		post, err := postByID(db, postID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Order not found", 404)
				log.Println(err)
				return
			}

			http.Error(w, "Internal Server Error1", 500)
			log.Println(err)
			return
		}
		ts, err := template.ParseFiles("pages/post.html")
		if err != nil {
			http.Error(w, "Internal Server Error2", 500)
			log.Println(err.Error())
			return
		}

		err = ts.Execute(w, post)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func postByID(db *sqlx.DB, postID int) (postPageData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			content
		FROM
			blog.post
		WHERE
			post_ID = ?
	`

	var order postPageData

	// Обязательно нужно передать в параметрах orderID
	err := db.Get(&order, query, postID)
	if err != nil {
		return postPageData{}, err
	}

	return order, nil
}

func featuredPosts(db *sqlx.DB) ([]featuredPostData, error) {
	const query = `
		SELECT
			post_ID,
			title,
			subtitle,
			image_url,
			author,
			author_url,
			publish_date
		FROM
			blog.post
		WHERE featured = 1
		LIMIT 2
	`

	var featurePosts []featuredPostData

	err := db.Select(&featurePosts, query)
	if err != nil {
		return nil, err
	}

	return featurePosts, nil
}

func mostRecent(db *sqlx.DB) ([]mostRecentData, error) {
	const query = `
		SELECT
			post_ID,
			title,
			subtitle,
			image_url,
			author,
			author_url,
			publish_date
		FROM
			blog.post
		WHERE featured = 0
		LIMIT 6
	`
	var mostPosts []mostRecentData

	err := db.Select(&mostPosts, query)
	if err != nil {
		return nil, err
	}

	return mostPosts, nil
}
