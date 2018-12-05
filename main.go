package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//	eng.C.Redirect(")
func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("./templates/*")

	//router.StaticFS("/static", http.Dir("static"))

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	r := router.Group("file")
	{
		r.GET("/scan", scan)

	}
	router.Run(":8080")
}

func scan(c *gin.Context) {
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

func pdf(c *gin.Context, filename, name string) {
	// req, _ := http.NewRequest("GET", filename, nil)

	// resp, _ := http.DefaultClient.Do(req)
	// f, err := os.OpenFile("./static/pdf/"+name, os.O_RDWR|os.O_CREATE, 0755)
	// if err == nil {
	// 	// io.Copy(f, resp.Body)
	// 	f.Close()
	// }
	// defer resp.Body.Close()

	streamPDFbytes, err := ioutil.ReadFile("./static/pdf/" + name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf := bytes.NewBuffer(streamPDFbytes)
	c.Writer.Header().Set("Content-type", "application/pdf")

	//1. buf.WriteTo
	// if _, err := buf.WriteTo(c.Writer); err != nil {
	// 	fmt.Fprintf(c.Writer, "%s", err)
	// }

	//2. 使用io.Pipe()开启读写双通道， io.Copy
	piper, pipew := io.Pipe()
	go func() {
		defer pipew.Close()
		io.Copy(pipew, buf)
	}()
	io.Copy(c.Writer, piper)
	piper.Close()

	// m := " <iframe id='iframe' frameborder='0' src='/static/web/viewer.html?file=/static/pdf/" + name + "' style='width:100%;'></iframe>"
	// fmt.Println(m)
	// c.HTML(http.StatusOK, "index.tmpl", gin.H{
	// 	"data": template.HTML(m),
	// })
}

func office(c *gin.Context, filename string) {
	src := "https://docview.mingdao.com/op/view.aspx?src=" + filename
	m := " <iframe id='iframe' frameborder='0' src='" + src + "' style='width:100%;'></iframe>"

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data": template.HTML(m),
	})

}
