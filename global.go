package main


type Handler struct {
	IcsFields []string
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

var (
	offset int64 = 0
	counts int64 = 0
	htmlIndex = `
		<html>
		<head>
		</head>
		<body>
			<div align="center">
			<a href="/login">Download *.ics File</a>
			</div>
		</body>
		</html>`
	host = "http://127.0.0.1:8080"
	url  = "https://oauth.vk.com/authorize?client_id=" + APP_ID + "&display=page&response_type=code&redirect_uri=" + host + "/result&scope=friends,offline&v=5.52"
	
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
