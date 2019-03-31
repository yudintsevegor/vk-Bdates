package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Dimonchik0036/vk-api"
)

type Handler struct {
	IcsFields []string
}

type FromResp struct {
	Token  string `json:"access_token"`
	UserId int64  `json:"user_id"`
}

var (
	htmlIndex = `
		<html>
		<head>
		</head>
		<body>
			<div align="center">
			<a href="/login">vk Log in</a>
			</div>
		</body>
		</html>`
	host = "http://127.0.0.1:8080"
	url  = "https://oauth.vk.com/authorize?client_id=" + APP_ID + "&display=page&response_type=code&redirect_uri=" + host + "/result&scope=friends,offline&v=5.52"
)

func (h *Handler) handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

type ND struct {
	Name string
	Date string
	Year string
}

func (h *Handler) handleResult(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	urlToken := "https://oauth.vk.com/access_token?client_id=" + APP_ID + "&client_secret=" + CLIENT_SECRET + "&redirect_uri=" + host + "/result&code=" + code

	resp, err := http.Get(urlToken)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	st := FromResp{}
	err = json.Unmarshal(body, &st)
	if err != nil {
		fmt.Println(err)
	}
	//	fmt.Fprintf(w, "ACCESS_TOKEN: %v\nUser_Id: %v\n", st.Token, st.UserId)

	client, err := vkapi.NewClientFromToken(st.Token)
	if err != nil {
		fmt.Println(err)
	}
	/*
		hints - sort by rating
		nom - Nominative
	*/
	var offset int64 = 0
	var counts int64 = 5000
	users := make([]ND, 0, 1)
	now := "2019"

	zeroNum := map[string]string{
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
	friends, err := client.GetFriends(st.UserId, "hints", counts, offset, "nom", "bdate")
	for _, friend := range friends {
		if friend.Bdate == "" {
			continue
		}
		DMY := strings.Split(friend.Bdate, ".")
		year := ""
		if len(DMY) == 2 {
			year = "Unknown"
		} else {
			year = DMY[2]
		}
		month := DMY[1]
		if _, ok := zeroNum[month]; ok {
			month = zeroNum[month]
		}
		day := DMY[0]
		if _, ok := zeroNum[day]; ok {
			day = zeroNum[day]
		}
		user := ND{
			Name: fmt.Sprintf("%v %v", friend.FirstName, friend.LastName),
			Date: now + month + day,
			Year: year,
		}
		//		fmt.Fprintf(w, "%v %v\n", friend.FirstName, friend.LastName)
		//		fmt.Fprintf(w, friend.Bdate+"\n")
		users = append(users, user)
	}
	//	fmt.Println(users)
	file, err := os.Create("kek.ics")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fields := h.IcsFields
	writer := bufio.NewWriter(file)
	begin := "BEGIN:VCALENDAR\nVERSION:2.0\n"
	end := "END:VCALENDAR"
	writer.WriteString(begin)
	for _, user := range users {
		for _, field := range fields {
			switch field {
			case "BEGIN:":
				writer.WriteString(field + "VEVENT")
				writer.WriteString("\n")
				continue
			case "SUMMARY:":
				writer.WriteString(field + user.Name + "'s B-Day")
				writer.WriteString("\n")
				continue
			case "DTSTART;VALUE=DATE:":
				writer.WriteString(field + user.Date)
				writer.WriteString("\n")
				continue
			case "DTEND;VALUE=DATE:":
				writer.WriteString(field + user.Date)
				writer.WriteString("\n")
				continue
			case "RRULE:FREQ=YEARLY;UNTIL=":
				writer.WriteString(field + "21190101")
				writer.WriteString("\n")
				continue
			case "DESCRIPTION:":
				writer.WriteString(field + "Year of Birth: " + user.Year)
				writer.WriteString("\n")
				continue
			case "END:":
				writer.WriteString(field + "VEVENT")
				writer.WriteString("\n")
				continue
			default:
				fmt.Println("BLYA", field)
				return
			}
		}
	}
	writer.WriteString(end)
	writer.Flush()

	w.Header().Add("Content-Disposition", "Attachment")
	content, err := ioutil.ReadFile("kek.ics")
	if err != nil {
		fmt.Println(err)
	}
	http.ServeContent(w, r, "random.ics", time.Now(), bytes.NewReader(content))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.handleMain(w, r)
	case "/login":
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	case "/result":
		h.handleResult(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	fields := []string{"BEGIN:", "SUMMARY:", "DTSTART;VALUE=DATE:", "DTEND;VALUE=DATE:", "RRULE:FREQ=YEARLY;UNTIL=", "DESCRIPTION:", "END:"}
	handler := &Handler{
		IcsFields: fields,
	}
	port := "8080"
	fmt.Println("Start listening at port: ", port)
	http.ListenAndServe(":"+port, handler)
}
