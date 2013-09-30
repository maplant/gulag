package main

import (
	"strings"
	"strconv"
)

const HEADER = `<!DOCTYPE html>
	<html>
	<head>
	<title>Gulag.co - Speak freely</title>
	</head>
	<body>
	<div style="width:800px; margin:0 auto; font-size:small; font-family:verdana,arial,helvetica,sans-serif;">
	<div style="text-align:center"><h1><a href="/">GULAG</a> - Speak boy, speak</h1></div>`

const FOOTER = `<a href="/">Return</a></div></div></body></html>`

const ERRNOPOST = (HEADER +
	"<h1>Error!</h1><br /><b>Could not post!</b>" +
	FOOTER)

const ERRNOUPLOAD = (HEADER +
	"<h1>Error!</h1><br /><b>Could not upload!</b>" +
	FOOTER)

const ERR404 = (HEADER +
	"<h1>404!</h1><br /><b>This shit doesn't exist</b>" +
	FOOTER)


func dispPost(p *Post, id int) string {
	n := strconv.Itoa(id)
	s := []string{`<div style="background:#C2C2C2;padding:5px;overflow:hidden" id="` + n + `">
 		<span style="float:left">` + p.name + `<div style="color:white"> ` + p.pd +` </div></span>
		<span style="float:right">` + n + `</span>
		</div>
			<div style="background:#D6D6D6;height:auto;padding:5px;overflow:hidden">`}
	if p.img != nil {
		s = append(s, "<div style=\"float:left; margin-right:5px\"><img src=\"/img/" + p.img.fn + "\"></img></div>")
	}
	s = append(s, p.comm + `</div>`)
	return strings.Join(s, "")
}

func dispThread(t *Thread, v string) string {
	s := []string{HEADER + `<b>Topic: ` + t.title + `</b>`}
	for i, p := range t.posts {
		s = append(s, dispPost(p, i + 1))
	}
	//  enctype="multipart/form-data">
	return strings.Join(append(s,
		`<form action="/p/` + v + `" method="post" enctype="multipart/form-data">`,
		`<div style="background:#C2C2C2;padding:5px">Reply: </div>`,
	        `<div style="background:#D6D6D6;padding:5px">`,
		`Name: <input type="text" name="name" value="anonymous"><br />`,
		`<textarea wrap="soft" rows="4" cols="48" name="comment"></textarea></div>`,
                `<div style="background:#C2C2C2;padding:5px>`,
 		`<span style="float:left"><input type="file" name="image"><span>`,
		`<span style="float:right"><input type="submit" value="Post"><span>`,
		`</div></form>` + FOOTER), "")
}

func dispAllThreads() string {
	s := []string{HEADER}
	ps := getReleventThreads()
	i := 1
	for _, p := range ps {
		if p == nil || i > MAX_POSTS {
			break;
		}
		var color string
		if i & 1 == 1 {
			color = "#C2C2C2"
		} else {
			color = "#D6D6D6"
		}
		s = append(s, `<div style="background:` + color + `;padding:5px"><b>` +
			strconv.Itoa(i) + `: </b><a href="/t/` +
			strconv.FormatUint(p.id, 32) + `">` + p.title + `</a></div>`)
		i++
	}
	return strings.Join(append(s,
		`<form action="/new" method="post" enctype="multipart/form-data">`,
		`<div style="background:#C2C2C2;padding:5px">New Post: </div>`,
	        `<div style="background:#D6D6D6;padding:5px">`,
		`Title: <input type="text" name="title" value=""><br />`,
		`Name: <input type="text" name="name" value="anonymous"><br />`,
		`<textarea wrap="soft" rows="4" cols="48" name="body"></textarea></div>`,
                `<div style="background:#C2C2C2;padding:5px>`,
 		`<span style="float:left"><input type="file" name="image"><span>`,
		`<span style="float:right"><input type="submit" value="Post"><span>`,
		`</div></form>` + FOOTER), "")
}
			
