package main

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	router.Run(":80")
}

func scan(c *gin.Context) {
	filename := c.Query("filename")
	name := strings.Split(filename, "/")
	name = name[len(name)-1:]
	a := strings.Split(name[0], ".")
	suffix := a[1]

	switch strings.ToUpper(suffix) {
	case "PDF":
		println(filename, name[0])
		pdf(c, filename, name[0])
	case "DOC", "DOCX":
		office(c, filename)
	}
}

func pdf(c *gin.Context, filename, name string) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Timeout: 20 * time.Second, Transport: tr}

	req, err := http.NewRequest("GET", filename, nil)
	if err != nil {
		log.Printf("http.NewRequest filename(%s) error(%v)", filename, err)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("http.DefaultClient.Do error(%v)", err)
		return
	}
	defer resp.Body.Close()

	streamPDFbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll error(%v)", err)
		return
	}

	f, err := os.OpenFile("./static/pdf/"+name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("os.OpenFile name(%s) error(%v)", "./static/pdf/"+name, err)
		return
	}
	defer f.Close()

	buf := bytes.NewBuffer(streamPDFbytes)
	piper, pipew := io.Pipe()
	go func() {
		defer pipew.Close()
		io.Copy(pipew, buf)
	}()
	io.Copy(f, piper)
	piper.Close()

	m := " <iframe id='iframe' frameborder='0' src='/static/web/viewer.html?file=/static/pdf/" + name + "' style='width:100%;'></iframe>"
	// fmt.Println(m)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data": template.HTML(m),
	})
}

func office(c *gin.Context, filename string) {
	src := "https://docview.mingdao.com/op/view.aspx?src=" + filename
	m := " <iframe id='iframe' frameborder='0' src='" + src + "' style='width:100%;'></iframe>"

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data": template.HTML(m),
	})

}
