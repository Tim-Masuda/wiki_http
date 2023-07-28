package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Page struct{
	Title string
	Body []byte
}

// создание
func (p *Page) savePage() error{
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600 ) // 1 - куда, 2 - что записываем, 3 - как(POSIX)
}

// чек и загрузка
func loadPage(title string) (*Page, error){
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename) // поиск
	if err != nil{
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// отображение
func viewHundler(resp http.ResponseWriter, request *http.Request){
	title := request.URL.Path[len("/view/"):] // записываем путь, и обрезаем что бы load смог найти
	p, err := loadPage(title)
	if err != nil{
		http.Redirect(resp, request, "/edit/"+title, http.StatusFound)
		return
	};
	renderTemplate(resp, "view", p)
}

func editHandler(resp http.ResponseWriter, request *http.Request){
	title := request.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil{
		p = &Page{Title: title}
	}
	renderTemplate(resp, "edit", p) 
}

func saveHandler(resp http.ResponseWriter, request *http.Request){
	title := request.URL.Path[len("/save/"):]
	body := request.FormValue("body")
	p := Page{Title: title, Body: []byte(body)}
	p.savePage()
	http.Redirect(resp, request, "/view/"+title, http.StatusFound)
} 

func renderTemplate(resp http.ResponseWriter, tmpl string, p *Page){
	t, _ := template.ParseFiles(tmpl + ".html") // расспарсивываем(читаем) 
	t.Execute(resp, p) // общращаемся к чтмлки, передаем pagef
}

func main(){
	http.HandleFunc("/view/", viewHundler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	fmt.Println("запуск")
	log.Fatal(http.ListenAndServe(":8080", nil))
}