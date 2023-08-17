package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

// Struct für die Datenstruktur
type UserData struct {
	Name  string
	Email string
}

// Methode zum Laden von Benutzerdaten
func (u *UserData) Load() {
	// Hier könnten Sie den Ladevorgang aus einer Datenbank oder einer Datei implementieren
	u.Name = "Herbert Robisch"
	u.Email = "herbertrobisch@example.com"
}

// Methode zum Speichern von Benutzerdaten
func (u *UserData) Save() {
	// Hier könnten Sie den Speichervorgang in einer Datenbank oder einer Datei implementieren
	fmt.Printf("Saving user data: Name: %s, Email: %s\n", u.Name, u.Email)
}

// Hauptfunktion
func main() {
	// Datenstruktur erstellen
	userData := &UserData{}

	// Webserver-Endpunkt zum Laden und Anzeigen von Benutzerdaten
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userData.Load()

		tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>User Data</title>
		</head>
		<body>
			<h1>User Data</h1>
			<p>Name: {{.Name}}</p>
			<p>Email: {{.Email}}</p>
		</body>
		</html>
		`

		t := template.Must(template.New("userData").Parse(tmpl))
		t.Execute(w, userData)
	})

	// Webserver-Endpunkt zum Bearbeiten von Benutzerdaten
	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			newName := r.FormValue("name")
			newEmail := r.FormValue("email")

			// Validierung der Benutzereingabe mit Regex
			validEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`).MatchString(newEmail)
			if !validEmail {
				http.Error(w, "Invalid email address", http.StatusBadRequest)
				return
			}

			userData.Name = newName
			userData.Email = newEmail
			userData.Save()

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			tmpl := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Edit User Data</title>
			</head>
			<body>
				<h1>Edit User Data</h1>
				<form method="post" action="/edit">
					<label for="name">Name:</label>
					<input type="text" id="name" name="name" value="{{.Name}}" required><br>
					<label for="email">Email:</label>
					<input type="email" id="email" name="email" value="{{.Email}}" required><br>
					<button type="submit">Save</button>
				</form>
			</body>
			</html>
			`

			t := template.Must(template.New("editUserData").Parse(tmpl))
			t.Execute(w, userData)
		}
	})

	// Den Webserver starten
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
