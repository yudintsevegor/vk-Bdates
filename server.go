package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	vkapi "github.com/Dimonchik0036/vk-api"
	"github.com/pkg/errors"
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

func (h *Handler) handleDownLoad(w http.ResponseWriter, r *http.Request) {
	content := h.Content
	w.Header().Set("Content-type", "text/calendar")
	w.Header().Add("Content-Disposition", "Attachment; filename=vkBdates.ics")
	//http.ServeContent(w, r, "", time.Now(), bytes.NewReader([]byte(content)))
	fmt.Fprintf(w, content)
}

func (h *Handler) handleResult(w http.ResponseWriter, r *http.Request) {
	tokenUser, status, err := getToken(r.FormValue("code"))
	if err != nil {
		SendError(w, status, err)
	}

	client, err := vkapi.NewClientFromToken(tokenUser.Token)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	content, err := h.getContent(client, tokenUser.UserId)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	h.Content = content
	fmt.Fprintf(w, htmlMain)
}

func getToken(code string) (TokenUser, int, error) {
	urlToken := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=%v&client_secret=%v&redirect_uri=%v/result&code=%v",
		APP_ID, CLIENT_SECRET, host, code)

	resp, err := http.Get(urlToken)
	if err != nil {
		return TokenUser{}, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TokenUser{}, http.StatusInternalServerError, err
	}

	tokenUser := new(TokenUser)
	if err = json.Unmarshal(body, tokenUser); err != nil {
		return TokenUser{}, http.StatusInternalServerError, err
	}

	if tokenUser.Token == "" {
		return TokenUser{}, http.StatusUnauthorized, errors.Errorf("access token not found.")
	}

	return *tokenUser, http.StatusOK, nil
}
