package model

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"strings"
)

type patchOp struct {
	Add     string
	Move    string
	Replace string
	Remove  string
	Copy    string
	Test    string
	AddUp   string
}

var PatchOp patchOp

// add move replace remove copy test
type PatchParams struct {
	Op    string      `json:"op" bson:"op"`
	Path  string      `json:"path" bson:"path"` // FIXME what if front parsing a value didn't startswith '/'
	Value interface{} `json:"value" bson:"value"`
}

type ExtParams struct {
	PageNum  int
	PageSize int
	Sort     string
}

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

func init() {
	PatchOp.Add = "add"
	PatchOp.Move = "move"
	PatchOp.Replace = "replace"
	PatchOp.Remove = "remove"
	PatchOp.Copy = "copy"
	PatchOp.Test = "test"
	PatchOp.AddUp = "add_up"
}

func (u *UrlDivide) ParseUrl(url string) error {
	urlParts := strings.Split(url, "/")

	for i, j := range urlParts {
		switch i {
		case 0: //  / -> ["", ""], never match 0, but needs to stay
		case 1:
			if e := u.getColAndAction(j); e != nil {
				return e
			}
		case 2:
			if e := u.getId(j); e != nil {
				return e
			}
		default:
			return errors.New("invalid url of api")
		}
	}

	return nil
}

func (u *UrlDivide) Validate(method string) bool {
	switch method {
	case "GET":
		return true
		//fallthrough
	case "POST":
		if u.Id != "" {
			return false
		}
		fallthrough
	case "PUT":
		if u.Id == "" {
			return false
		}
		fallthrough
	case "PATCH":
		if u.Id == "" {
			return false
		}
		fallthrough
	case "DELETE":
		if u.Id == "" {
			return false
		}
		fallthrough
	default:
		return true
	}

}

func (u *UrlDivide) getColAndAction(path string) error {
	if strings.Contains(path, ":") {
		split := strings.Split(path, ":")
		if split[1] == "" {
			return errors.New(":empty action")
		}
		u.Action = split[1]
		u.Collection = split[0]
	} else {
		u.Collection = path
	}
	return nil
}

func (u *UrlDivide) getId(path string) error {
	if u.Action != "" {
		return errors.New(":action doesn't support '/{_id}'") // TODO maybe can?
	}
	if !bson.IsObjectIdHex(path) {
		return errors.New("invalid id")
	}
	u.Id = path
	return nil
}

func (p *PatchParams) ParseQuery(query *bson.M, set *bson.M, k string) error {
	if p.Path == "/_id"{
		return errors.New("can do nothing on '_id'")
	}

	var e error
	switch p.Op {
	case PatchOp.Add:
		(*query)[k] = bson.M{"$exists": false}
	case PatchOp.Replace:
		fallthrough
	case PatchOp.Remove:
		fallthrough
	case PatchOp.AddUp:
		(*query)[k] = bson.M{"$exists": true}
	case PatchOp.Move:
		e = errors.New("TODO")
	case PatchOp.Copy:
		e = errors.New("TODO")
	case PatchOp.Test:
		e = errors.New("TODO")
	default:
		e = errors.New("patch op invalid: " + p.Op)
	}

	if e != nil { // TODO remove after all done
		if e.Error() != "TODO" {
			return e
		}
		e = nil
	}

	var updateOp string
	switch p.Op {
	case PatchOp.AddUp:
		updateOp = "inc"
	case PatchOp.Remove:
		updateOp = "unset"
	default:
		updateOp = "set"
	}
	*set = bson.M{"$" + updateOp: bson.M{k: p.Value}}

	return e
}

func (p *PatchParams) Ignore() bool {
	switch p.Op {
	case PatchOp.Add:
		fallthrough
	case PatchOp.Remove:
		return true
	default:
		return false
	}
}
