package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	vkapi "github.com/Dimonchik0036/vk-api"
)

func getUsers(friends []vkapi.Users) []NDY {
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

var (
	offset int64
	counts int64
)

func (h *Handler) getContent(client *vkapi.Client, id int64) (string, error) {
	// name - sorting by name
	// nom - Nominative
	friends, errVk := client.GetFriends(id, "name", counts, offset, "nom", "bdate")
	if errVk != nil {
		return "", errVk
	}

	users := getUsers(friends)
	content, err := h.makeContent(users)
	if err != nil {
		return "", fmt.Errorf("InternalServerError")
	}

	return content, nil
}

func (h *Handler) makeContent(users []NDY) (string, error) {
	endYear := strconv.Itoa(time.Now().AddDate(100, 0, 0).Year())
	var content string
	content += "BEGIN:VCALENDAR\nVERSION:2.0\n"

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

	content += "END:VCALENDAR"

	return content, nil
}
