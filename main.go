package main

import (
	"GoShare/storage"
	"html/template"
	"log"
	"net/http"
)

type UserData struct {
	UserName string
	Error    string
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, "Temlate not found", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка рендеренга шаблона", http.StatusInternalServerError)
		return
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")

		log.Println("🔹 Попытка входа:", username)

		exists, err := storage.ValidateUser(username, password)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if !exists {
			log.Println("❌ Неверный логин или пароль:", username)
			data := UserData{
				UserName: username,
				Error:    "Invalid username or password",
			}
			renderTemplate(w, "login", data)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
			Path:  "/",
		})

		log.Println("✅ Успешный вход:", username)

		http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "login", nil)
}

func registerHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		email := req.FormValue("email")

		if username == "" || password == "" || email == "" {
			data := UserData{
				UserName: username,
				Error:    "All fields are required",
			}
			renderTemplate(w, "register", data)
			return
		}

		exists, err := storage.UserExists(username)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if exists {
			data := UserData{
				UserName: username,
				Error:    "Username already taken",
			}
			renderTemplate(w, "register", data)
			return
		}

		err = storage.CreateUser(username, password, email)
		if err != nil {
			log.Println("Error with adding user:", err)
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}

		log.Println("Пользователь успешно зарегистрирован:", username, email)

		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
			Path:  "/",
		})

		http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "register", nil)
}

func dashboardHandler(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("username")
	if err != nil {
		log.Println("❌ Пользователь не авторизован, редирект на /login")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	log.Println("✅ Вход на Dashboard, пользователь:", cookie.Value)

	data := UserData{
		UserName: cookie.Value,
	}
	renderTemplate(w, "dashboard", data)
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func main() {
	storage.InitDB()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		renderTemplate(w, "index", nil)
	})

	http.HandleFunc("/login", loginHandler)

	http.HandleFunc("/register", registerHandler)

	http.HandleFunc("/dashboard", dashboardHandler)

	http.HandleFunc("/logout", logoutHandler)

	log.Print("Server is listening on port 8080...")

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
