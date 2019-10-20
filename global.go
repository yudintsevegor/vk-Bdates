package main

import (
	"fmt"
)

type Handler struct {
	IcsFields []string
	Content   string
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

const (
	alarm = `BEGIN:VALARM
	TRIGGER:-P0D
	DESCRIPTION:reminder
	ACTION:DISPLAY
	END:VALARM
`
)

var (
	url = fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%v&display=popup&response_type=code&redirect_uri=%v/result&scope=friends,offline&v=5.52",
		APP_ID, host)

	htmlMain = `
		<html>
		<head>
		</head>
		<body>
			<div align="center">
			<a href="/download">Download *.ics File</a>
			</div>
			<div align="center">
			<a href="https://vk.com">VK</a>
			</div>
		</body>
		</html>`

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
)

var (
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
