package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	
	"github.com/Dimonchik0036/vk-api"
)
type Handler struct{
	
}
type FromResp struct{
	Token string `json:"access_token"`
	Expires string
	UserId int `json:"user_id"`
}

var url = "https://oauth.vk.com/authorize?client_id="+key+"&display=page&response_type=code&redirect_uri=http://127.0.0.1:8080/result&scope=friends,offline&v=5.52"
//var url = "https://oauth.vk.com/authorize?client_id=6920031&display=page&redirect_uri=https://oauth.vk.com/blank.html&scope=friends,offline&response_type=token&v=5.52"
func (h *Handler) handleMain(w http.ResponseWriter, r *http.Request){
	
}

func (h *Handler) handleResult(w http.ResponseWriter, r *http.Request){
	CODE := r.FormValue("code")
	fmt.Println(CODE)
	var urlToken = "https://oauth.vk.com/access_token?client_id="+key+"&client_secret="+CLIENT_SECRET+"&redirect_uri=http://127.0.0.1:8080/result&code="+CODE


//	http.Redirect(w, r, urlToken, http.StatusTemporaryRedirect)
//	return
	resp, err :=  http.Get(urlToken)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("BODY:", string(body))
	var st = FromResp{}
	err = json.Unmarshal(body,&st)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(st)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	switch r.URL.Path{
	case "/":
		h.handleMain(w,r)
	case "/login":
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	case "/result":
		h.handleResult(w,r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main(){
	client, err := vkapi.NewClientFromToken(TOKEN)
	if err != nil {
	    fmt.Println(err)
	}
	var userId int64 = 31567954
	/*
	hints - sort by rating
	nom - Nominative
	*/
	var offset int64 = 0
	var counts int64 = 5000
	friends, err := client.GetFriends(userId, "hints", counts, offset, "nom", "bdate")
	for _, friend := range friends{
		fmt.Println(friend.FirstName)
		fmt.Println(friend.LastName)
		fmt.Println(friend.Bdate)
	}
	
	handler := &Handler{}
	port := "8080"
	fmt.Println("Start listening at port: ", port)
	http.ListenAndServe(":"+port, handler)
//	values := url.Values{}
//	values.Set("user_id", "31567954")
//	values.Set("count", "10")
//
//	res, err := client.Do(vkapi.NewRequest("friends.getOnline", "", values))
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Println(res.Response.String())
}
