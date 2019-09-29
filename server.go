package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"errors"
	//"bytes"
	//"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"log"

	"github.com/Dimonchik0036/vk-api"
)

func SendError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	jsonError, _ := json.Marshal(FromServer{Error: err.Error()})
	w.Write(jsonError)
}

func (h *Handler) handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func (h *Handler) getContent(client *vkapi.Client, id int64) (string, error){
	// name - sorting by name
	// nom - Nominative
	friends, errVk := client.GetFriends(id, "name", counts, offset, "nom", "bdate")
	if errVk != nil{
		return "", errors.New("bad request")
	}
	
	users := getUsers(friends)
	content, err := h.makeContent(users)
	if err != nil {
		return "", fmt.Errorf("InternalServerError")
	}
	
	return content, nil
}

func getUsers(friends []vkapi.Users) []NDY{
	now := strconv.Itoa(time.Now().Year())
	users := make([]NDY, 0, 1)
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
		user := NDY{
			Name: fmt.Sprintf("%v %v", friend.FirstName, friend.LastName),
			Date: now + month + day,
			Year: year,
		}
		users = append(users, user)
	}
	
	return users
}

func (h *Handler) makeContent(users []NDY) (string, error){
	endYear := strconv.Itoa(time.Now().AddDate(100,0,0).Year())
	content := ""
	begin := "BEGIN:VCALENDAR\nVERSION:2.0\n"
	content += begin
	end := "END:VCALENDAR"
	for _, user := range users {
		for _, field := range h.IcsFields {
			switch field {
			case "BEGIN:":
				content += field + "VEVENT\n"
				continue
			case "SUMMARY:":
				content += field + user.Name + "'s B-Day\n"
				continue
			case "DTSTART;VALUE=DATE:":
				content += field + user.Date + "\n"
				continue
			case "DTEND;VALUE=DATE:":
				content += field + user.Date + "\n"
				continue
			case "RRULE:FREQ=YEARLY;UNTIL=":
				content += field + endYear + "0101\n"
				continue
			case "DESCRIPTION:":
				content += field + "Year of Birth: " + user.Year + "\n"
				continue
			case "END:":
				content += field + "VEVENT\n"
				continue
			default:
				return "", errors.New("unknown field")
			}
		}
		content += alarm
	}
	content += end

	return content, nil
}


func (h *Handler) handleDownLoad(w http.ResponseWriter, r *http.Request){
	log.Println("HERE")
	content := h.Content
	w.Header().Set("Content-type", "text/calendar")
	w.Header().Add("Content-Disposition", "Attachment; filename=vkBdates.ics")
	//http.ServeContent(w, r, "", time.Now(), bytes.NewReader([]byte(content)))
	fmt.Fprintf(w, content)
}

func (h *Handler) handleResult(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	urlToken := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%v&client_secret=%v&redirect_uri=%v/result&code=%v", APP_ID, CLIENT_SECRET, host, code)
	
	resp, err := http.Get(urlToken)
	if err != nil{
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	st := TokenUser{}
	err = json.Unmarshal(body, &st)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	
	if st.Token == "" {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("Access token not found."))
		return
	}
	
	client, err := vkapi.NewClientFromToken(st.Token)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Println(st)
	content, err := h.getContent(client, st.UserId)
	if err != nil{
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	h.Content = content
	fmt.Fprintf(w, html)
}

