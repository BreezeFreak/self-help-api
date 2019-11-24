package app

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"strings"
)

type Record struct {
	ID         bson.ObjectId `json:"_id" bson:"_id"`
	IP         string        `json:"ip" bson:"ip"`
	Collection string        `json:"collection" bson:"collection"`
}

// /[collection]/[_id]
// /[collection]:[action]

type UrlDivide struct {
	Collection string
	Id         string
	Action     string
}

func (u *UrlDivide) ParseUrl(url string) error {
	urlParts := strings.Split(url, "/")

	for i, j := range urlParts {
		switch i {
		case 0: //  / -> ["", ""], never match 0, but needs to stay
		case 1:
			if strings.Contains(j, ":") {
				split := strings.Split(j, ":")
				if split[1] == "" {
					return errors.New(":empty action")
				}
				u.Action = split[1]
				u.Collection = split[0]
			} else {
				u.Collection = j
			}
		case 2:
			if u.Action != "" {
				return errors.New(":action doesn't support '/{_id}'") // TODO maybe can?
			}
			if !bson.IsObjectIdHex(j) {
				return errors.New("invalid id")
			}
			u.Id = j
		default:
			return errors.New("invalid url of api")
		}
	}

	return nil
}
