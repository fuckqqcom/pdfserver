package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

//	eng.C.Redirect(")
func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("./templates/*")

	//router.StaticFS("/static", http.Dir("static"))

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	r := router.Group("down")
	{
		r.GET("/ping", server)
		r.GET("/pong", pong)

	}
	router.Run(":80")
}

func pong(c *gin.Context) {
	filename := c.Query("filename")
	name := strings.Split(filename, "/")
	name = name[len(name)-1:]
	a := strings.Split(name[0], ".")
	suffix := a[1]

	switch strings.ToUpper(suffix) {
	case "PDF":
		pdf(c, filename, name[0])
	case "DOC", "DOCX":
		office(c, filename)
	}
}

func office(c *gin.Context, filename string) {

	src := "https://docview.mingdao.com/op/view.aspx?src=" + filename
	m := " <iframe id='iframe' frameborder='0' src='" + src + "' style='width:100%;'></iframe>"

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data": template.HTML(m),
	})

}

func pdf(c *gin.Context, filename, name string) {
	req, _ := http.NewRequest("GET", filename, nil)

	resp, _ := http.DefaultClient.Do(req)
	f, err := os.OpenFile("./static/pdf/"+name, os.O_RDWR|os.O_CREATE, 0755)
	if err == nil {
		io.Copy(f, resp.Body)
		f.Close()
	}
	defer resp.Body.Close()

	m := " <iframe id='iframe' frameborder='0' src='/static/web/viewer.html?file=/static/pdf/" + name + "' style='width:100%;'></iframe>"
	fmt.Println(m)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data": template.HTML(m),
	})
}

func server(c *gin.Context) {
	filename := c.Query("filename") // shortcut for c.Request.URL.Query().Get("lastname")
	fmt.Println(filename)
	name := strings.Split(filename, "/")
	name = name[len(name)-1:]
	resp, err := http.Get(filename)
	defer resp.Body.Close()
	if err == nil {
		f, err := os.Create("/usr/local/openresty/nginx/html/pdf/pdf/" + name[0])
		fmt.Println("err->", err)
		io.Copy(f, resp.Body)
		f.Close()
	}

	pdf := "http://58.87.64.219/web/viewer.html?file=http://58.87.64.219/pdf/pdf/" + name[0]
	//c.SaveUploadedFile()
	//c.Redirect(http.StatusMovedPermanently, ")
	c.Redirect(http.StatusMovedPermanently, pdf)
}
