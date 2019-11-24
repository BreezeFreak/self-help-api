package mongodb

import (
	"api/config"
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
	"strings"
)

var (
	Session *mgo.Session
	Mongo   *mgo.DialInfo
	DBName  string
)

var InvalidObjectId = errors.New("invalid ObjectId")

func Init() {
	uri := config.MongoDB.URI
	db := config.MongoDB.DBName

	mongo, err := mgo.ParseURL(uri)
	s, err := mgo.Dial(uri)
	if err != nil {
		log.Fatalf("[MongoDB] Connecion failed to %v: %v\n", uri, err)
	}
	s.SetSafe(&mgo.Safe{})
	fmt.Println("[MongoDB] Connected: ", uri)
	DBName = db
	Session = s
	Mongo = mongo
}

func Insert(collection string, doc bson.M) error {
	s := Session.Clone()
	defer s.Close()

	return s.DB(DBName).C(collection).Insert(doc)
}

func FindGreater(collection string, query bson.M, result interface{}, sort string, limit int) error {
	s := Session.Clone()
	defer s.Close()

	return s.DB(DBName).C(collection).Find(query).Sort(sort).Limit(limit).All(result)
}

func FindAll(collection string, query bson.M, result interface{}, extParams ExtParams) (int, error) {
	if e := parseObjectId(&query); e != nil {
		return 0, e
	}

	s := Session.Clone()
	defer s.Close()

	cursor := s.DB(DBName).C(collection).Find(query)

	if extParams.Sort != "" {
		cursor.Sort(extParams.Sort)
	}

	count, _ := cursor.Count()

	e := cursor.
		Skip(extParams.PageNum * extParams.PageSize).
		Limit(extParams.PageSize).
		All(result)
	return count, e
}

func FindOne(collection string, query bson.M, result interface{}) error {
	if e := parseObjectId(&query); e != nil {
		return e
	}
	s := Session.Clone()
	defer s.Close()

	return s.DB(DBName).C(collection).Find(query).One(result)
}

// add move replace remove copy test
func Patch(collection, _id string, patchData *[]PatchParams, fieldsAvailable []string) error {
	if len(*patchData) == 0 {
		return errors.New("empty patch content")
	}

	id := bson.ObjectIdHex(_id)

	return PatchByQuery(collection, bson.M{"_id": id}, patchData, fieldsAvailable)
}

func PatchByQuery(collection string, query bson.M, patchData *[]PatchParams, fieldsAvailable []string) error {

	for _, patch := range *patchData {
		log.Println("[mongodb PATCH]: " + fmt.Sprint(patch))
		for _, f := range fieldsAvailable {
			var e interface{} = nil

			k := strings.ReplaceAll(patch.Path[1:], "/", ".")
			field := strings.Split(k, ".")[0]

			if field == f {
				e = getPatchQuery(&query, k, &patch)
			}

			if e != nil { // TODO remove after all done
				if e.(error).Error() != "TODO" {
					return e.(error)
				}
			}

			set := bson.M{k: patch.Value}
			var updateOp string
			if patch.Op == PatchOp.AddUp {
				updateOp = "$inc"
			} else {
				updateOp = "$set"
			}
			e = Update(collection, query, set, updateOp)

			if e != nil {
				if patch.Op == PatchOp.Add {
					log.Println("[mongodb PATCH]: adding exists field ignored: " + fmt.Sprint(patch))
				} else {
					return e.(error)
				}
			}

			delete(query, k)
			break
		}

	}
	return nil
}

func Delete(collection, id string) error {
	_id, e := ToObjectId(id)
	if e != nil {
		return e
	}

	s := Session.Clone()
	defer s.Close()

	if !bson.IsObjectIdHex(id) {
		return InvalidObjectId
	}

	return s.DB(DBName).C(collection).RemoveId(_id)
}

func Update(collection string, query, set bson.M, option string) error {
	s := Session.Clone()
	defer s.Close()

	return s.DB(DBName).C(collection).Update(query, bson.M{option: set})
}

func getPatchQuery(q *bson.M, k string, patch *PatchParams) error {
	var e error
	switch patch.Op {
	case "add":
		(*q)[k] = bson.M{"$exists": false}
	case "replace":
		(*q)[k] = bson.M{"$exists": true}
	case "remove":
		// FIXME should may remove the field instead of value only
		patch.Value = ""
		(*q)[k] = bson.M{"$exists": true}
	case "move":
		e = errors.New("TODO")
	case "copy":
		e = errors.New("TODO")
	case "test":
		e = errors.New("TODO")
	case "add_up":
		(*q)[k] = bson.M{"$exists": true}
	default:
		e = errors.New("patch op invalid: " + patch.Op)
	}
	return e
}

func parseObjectId(query *bson.M) error {
	if v, ok := (*query)["_id"]; ok {
		_id, e := ToObjectId(v)
		if e != nil {
			return e
		}
		(*query)["_id"] = _id
	}
	return nil
}

func ToObjectId(v interface{}) (interface{}, error) {
	if _, ok := v.(string); ok {
		if bson.IsObjectIdHex(fmt.Sprint(v)) {
			return bson.ObjectIdHex(fmt.Sprint(v)), nil
		}
	}
	if _, ok := v.(bson.ObjectId); ok {
		return v, nil
	}
	return nil, InvalidObjectId
}

func NewCollection(collection string)  {
	s := Session.Clone()
	defer s.Close()

	fmt.Println(s.DB(DBName).C(collection).Count())
}
