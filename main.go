package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"html/template"

	"github.com/gorilla/mux"
)

var templates *template.Template

func main() {

	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", gethandler).Methods("GET")
	r.HandleFunc("/", posthandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
func gethandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "ERROR-404!", http.StatusNotFound)
		return
	}
	templates.ExecuteTemplate(w, "index.html", nil)

}
func posthandler(w http.ResponseWriter, r *http.Request) {

	text := r.FormValue("text")
	font := r.FormValue("fonts")

	for _, v := range text {
		if !(v >= 32 && v <= 126) {

			http.Error(w, "ERROR-400\nBad request!", http.StatusBadRequest)
			return
		}
	}

	file, err := os.Open(FormatType(font))

	if err != nil {
		http.Error(w, "Internal Server Error!!!\nERROR-500", http.StatusInternalServerError)
		return
	}

	defer file.Close()
	banners := [][]string{}
	banner := []string{}

	scanner := bufio.NewScanner(file)
	i := 0

	for scanner.Scan() {
		i++
		banner = append(banner, scanner.Text())

		if i == 9 {
			banners = append(banners, banner)
			banner = []string{}
			i = 0
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error")
		return
	}
	c := ""
	str := []string{}
	array := strings.Split(text, "\\n")
	for i := 0; i < len(array); i++ {

		for j := 1; j <= 8; j++ {
			for _, value := range array[i] {

				str = banners[int(value)-32]

				c = c + str[j]
			}

			if len(array[i]) != 0 {
				c = c + "\n"
			}
		}
	}

	templates.ExecuteTemplate(w, "index.html", c)
}

func FormatType(fs string) string {
	if fs == "shadow" {
		return "shadow.txt"
	}
	if fs == "thinkertoy" {
		return "thinkertoy.txt"
	}
	if fs == "standard" {
		return "standard.txt"
	}
	return "Error 500"
}
