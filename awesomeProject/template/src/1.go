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

type pageContent struct {
	Header []string
	Body   []interface{}
	Footer string
}

func ExecuteTemplate() {

	log.Println("********************* Template Section ***************************")
	/* Write html to file */
	html := "<html>" +
		"<h1>HelloWorld T1</h1>" +
		"</html>"

	nf, err := os.Create("template/resources/1.gohtml")

	if err != nil {
		log.Fatalln("Failed to create file")
	}
	_, err = io.Copy(nf, strings.NewReader(html))

	if err != nil {
		log.Fatalln("Failed to write html to file")
	}

	/* Read html file & write to console */

	fMap := template.FuncMap{
		"trim": func(s string) string {
			return strings.Trim(s, " ")
		},
		"caps": func(s string) string {
			return strings.ToUpper(s)
		},
	}

	/* Passing helper funcs to template */
	t = template.Must(template.New("").Funcs(fMap).ParseGlob("template/resources/*.gohtml"))
	err = t.ExecuteTemplate(os.Stdout, "2.gohtml", nil)
	if err != nil {
		panic("Boom!!")
	}
	fmt.Println()

	/* Placeholders */
	/* Note make the fields caps to make it available in template file */
	page := pageContent{
		Header: []string{"Hone", "Htwo", "Hthree"},
		Body:   []interface{}{"abc", 123, 10.5},
		Footer: "footer",
	}

	nf, err = os.Create("template/resources/templatized.html")
	if err != nil {
		log.Fatalln("Failed to create templatised file")
	}
	err = t.ExecuteTemplate(nf, "3.gohtml", page)
	if err != nil {
		panic("Boom!!")
	}
	log.Println("********************* Template Section ends ***************************")

}
