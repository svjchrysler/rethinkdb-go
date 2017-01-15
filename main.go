package main

import (
	r "gopkg.in/gorethink/gorethink.v2"
	"log"
	"github.com/kataras/iris"
	"fmt"
)

type Author struct {
	ID string `gorethink:"id,omitempty"`
	Name string `gorethink:"name"`
}


func session() *r.Session {
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})

	if err != nil {
		log.Fatalln("Servidor: " + err.Error())
	}

	return session
}

func GETAuthor(ctx *iris.Context) {

	session := session()

	data, err := r.DB("test").Table("authors").Run(session)

	if err != nil {
		log.Fatal("Error query: " + err.Error())
	}

	var rows []Author

	data.All(&rows)

	ctx.Render("author.html", rows)

}


func GETCreate(ctx *iris.Context) {
	ctx.Render("create.html", nil)
}

func POSTStore(ctx *iris.Context)  {

	session := session()
	name := ctx.PostValues("name")

	var author Author
	author.Name = name[0]

	e := r.DB("test").Table("authors").Insert(author).Exec(session)

	if e != nil {
		log.Fatalln("Insertar: " + e.Error())
	}

	ctx.Redirect("/author/")

}

func DELETEAuthor(ctx *iris.Context) {
	session := session()
	id := ctx.Param("id")
	fmt.Println("id : " +  id)
	r.Table("authors").Get(id).Delete().Run(session)
	ctx.Redirect("/author/")
}

func GETUpdate(ctx *iris.Context)  {
	session := session()
	id := ctx.Param("id")
	row, err := r.Table("authors").Get(id).Run(session)
	if err != nil {
		log.Fatal("query: " + err.Error())
	}

	var author Author

	row.One(&author)

	ctx.Render("edit.html", author)

}

func PATCHAuthor(ctx *iris.Context) {
	session := session()
	name := ctx.PostValue("name")
	id := ctx.PostValue("id")

	author := Author{ID:id, Name:name}

	r.Table("authors").Update(author).Run(session)

	ctx.Redirect("/author/")

}

func main()  {

	iris.Get("author/", GETAuthor)
	iris.Get("create/", GETCreate)
	iris.Post("create/store/", POSTStore)
	iris.Get("delete/:id", DELETEAuthor)
	iris.Get("edit/:id", GETUpdate)
	iris.Post("edit/update", PATCHAuthor)

	iris.Listen(":9000")

}
