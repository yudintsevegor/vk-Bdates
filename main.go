package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

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
	fields := h.IcsFields
	users := make([]NDY, 0, 1)
	
	now := strconv.Itoa(time.Now().Year())
	endYear := strconv.Itoa(time.Now().AddDate(100,0,0).Year())
	//name - sorting by name
	//nom - Nominative
	friends, err := client.GetFriends(id, "name", counts, offset, "nom", "bdate")
	if err != nil{
		return "", err
	}
	
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

	content := ""
	begin := "BEGIN:VCALENDAR\nVERSION:2.0\n"
	content += begin
	end := "END:VCALENDAR"
	for _, user := range users {
		for _, field := range fields {
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
				return "", fmt.Errorf("InternalServerError")
			}
		}
	}
	content += end
	return content, nil
}

func (h *Handler) handleResult(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	urlToken := "https://oauth.vk.com/access_token?client_id=" + APP_ID + "&client_secret=" + CLIENT_SECRET + "&redirect_uri=" + host + "/result&code=" + code

	resp, err := http.Get(urlToken)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
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
	content, err := h.getContent(client, st.UserId)
	if err != nil{
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-type", "text/calendar")
	w.Header().Add("Content-Disposition", "Attachment; filename=vkBdates.ics")
	http.ServeContent(w, r, "", time.Now(), bytes.NewReader([]byte(content)))
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
