package src

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

var t *template.Template

func ExecuteTemplate() {

	log.Println("********************* Template Section ***************************")
	/* Write html to file */
	html := "<html>" +
		"<h1>HelloWorld T1</h1>" +
		"</html>"

	nf, err := os.Create("template/resources/1.gohtml")

	if err !=nil {
		log.Fatalln("Failed to create file")
	}
	_, err = io.Copy(nf, strings.NewReader(html))

	if err != nil {
		log.Fatalln("Failed to write html to file")
	}

	/* Read html file & write to console */
	t = template.Must(template.ParseGlob("template/resources/*.gohtml"))
	t.ExecuteTemplate(os.Stdout, "2.gohtml", nil)
	fmt.Println()

	log.Println("********************* Template Section ends ***************************")

 }


