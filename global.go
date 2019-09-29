package main

import (
	"fmt"
)
type Handler struct {
	IcsFields []string
	Content string
}

type TokenUser struct {
	Token  string `json:"access_token"`
	UserId int64  `json:"user_id"`
}

type NDY struct {
	Name string
	Date string
	Year string
}

type FromServer struct {
	Error string
}

var alarm = `BEGIN:VALARM
TRIGGER:-P0D
DESCRIPTION:reminder
ACTION:DISPLAY
END:VALARM
`

var html = `
		<html>
		<head>
		</head>
		<body>
			<div align="center">
			<a href="/download">Download *.ics File</a>
			</div>
		</body>
		</html>`
		

var (
	offset int64 = 0
	counts int64 = 0
	htmlIndex = `
		<html>
		<head>
		</head>
		<body>
			<div align="center">
			<a href="/login">Start</a>
			</div>
		</body>
		</html>`
		
	host = "http://127.0.0.1:8080"
	url  = fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%v&display=popup&response_type=code&redirect_uri=%v/result&scope=friends,offline&v=5.52", APP_ID, host)
	zeroNum = map[string]string{
		"1": "01",
		"2": "02",
		"3": "03",
		"4": "04",
		"5": "05",
		"6": "06",
		"7": "07",
		"8": "08",
		"9": "09",
	}
)
