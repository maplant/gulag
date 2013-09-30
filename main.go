package main

import (
	"github.com/hoisie/web"
	"time"
	"crypto/md5"
	"io"
	"io/ioutil"
	"html"
	"strconv"
	"image"
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"strings"
	"fmt"
)

func Md5(r io.Reader) string {
    hash := md5.New()
    io.Copy(hash, r)
    return fmt.Sprintf("%x", hash.Sum(nil))
}

func dispThreadHandler(v string) string {
	tid, err := strconv.ParseUint(v, 32, 64)
	if err != nil {
		return ERR404
	}
	t := getThread(tid)
	if t == nil {
		return ERR404
	}
	return dispThread(t, v)
}

func postNewThread(ctx *web.Context) string {
	ctx.Request.ParseMultipartForm(10 * 1024 * 1024)
	form := ctx.Request.MultipartForm
	titlea, ok := form.Value["title"]
	if !ok {
		fmt.Printf("err 1\n")
		return ERRNOPOST
	}
	title := strings.TrimSpace(html.EscapeString(titlea[0]))
	if len(title) > MAX_TITLE_LEN || len(title) == 0 {
		fmt.Printf("err 2\n")
		return ERRNOPOST
	}
	namea, ok := form.Value["name"]
	if !ok {
		fmt.Printf("err 4\n")
		return ERRNOPOST
	}
	name := strings.TrimSpace(html.EscapeString(namea[0]))
	if len(name) > 60 || len(name) == 0 {
		name = "anonymous"
	}
	file, ok := form.File["image"]
	var imgd *Img = nil
	if ok {
		fileHeader := file[0]
		fd, err := fileHeader.Open()
		fdd, _ := ioutil.ReadAll(fd)
		fd.Seek(0, 0)
		_, fe, err := image.Decode(fd)
		if err != nil {
			print(err.Error())
			return ERRNOUPLOAD
		}
		//bounds := img.Bounds()
		fd.Seek(0, 0)
		name := Md5(fd) + "." + fe
		err = ioutil.WriteFile("./img/" + name, fdd, 0644)
		if err != nil {
			print("ERROR: " + err.Error() + "\n")
		}
		imgd = &Img{ 0, 0, name }
	}
	bodya, ok := form.Value["body"]
	if !ok {
		fmt.Printf("err 3\n")
		return ERRNOPOST
	}
	body := strings.TrimSpace(html.EscapeString(bodya[0]))
	if len(body) > MAX_POST_LEN || (len(body) == 0 && imgd == nil) {
		return ERRNOPOST
	}
	t := &Thread{title, 0, []*Post{&Post{time.Now().Format("Mon Jan 2 15:04:05"), name, parse(body), imgd}}, nil, nil }
	t.post()
	return dispThread(t, strconv.FormatUint(t.id, 32))
}

func imgHandler(ctx *web.Context, v string) string {
	b, err := ioutil.ReadFile("./img/" + v)
	if err != nil {
		return ERR404
	}
	return string(b)
}

func postToThread(ctx *web.Context, v string) string {
	ctx.Request.ParseMultipartForm(10 * 1024 * 1024)
	form := ctx.Request.MultipartForm
	tid, err := strconv.ParseUint(v, 32, 64)
	if err != nil {
		return ERR404
	}
	t := getThread(tid)
	if t == nil {
		return ERR404
	}
	file, ok := form.File["image"]
	var imgd *Img = nil
	if ok {
		fileHeader := file[0]
		fd, err := fileHeader.Open()
		fdd, _ := ioutil.ReadAll(fd)
		fd.Seek(0, 0)
		_, fe, err := image.Decode(fd)
		if err != nil {
			print(err.Error())
			return ERRNOUPLOAD
		}
		//bounds := img.Bounds()
		fd.Seek(0, 0)
		name := Md5(fd) + "." + fe
		err = ioutil.WriteFile("./img/" + name, fdd, 0644)
		if err != nil {
			print("ERROR: " + err.Error() + "\n")
		}
		imgd = &Img{ 0, 0, name }
	}
	pva, ok := form.Value["comment"]
	if !ok {
		return ERRNOPOST	/* This is clearly not what should happen. */
	}
	pv := strings.TrimSpace(html.EscapeString(pva[0]))
	if len(pv) > MAX_POST_LEN || (len(pv) == 0 && imgd == nil) {
		return ERRNOPOST	/* Ditto */
	}
	namea, ok := form.Value["name"]
	if !ok {
		return ERRNOPOST
	}
	name := strings.TrimSpace(html.EscapeString(namea[0]))
	if len(name) > 60 || len(name) == 0 {
		name = "anonymous"
	}
	t.bump(&Post{time.Now().Format("Mon Jan 2 15:04:05"), name, parse(pv), imgd})
	t = getThread(tid)	/* Reload the thread. */
	return dispThread(t, v)
}

func main() {
	web.Get("/", dispAllThreads)
	web.Get("/t/(.*)", dispThreadHandler)
	web.Get("/img/(.*)", imgHandler)
	web.Post("/p/(.*)", postToThread)
	web.Post("/new", postNewThread)
	web.Run("0.0.0.0:9999")
}