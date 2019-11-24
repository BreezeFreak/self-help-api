package mongodb

import (
	"api/config"
	"api/utils/model"
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

func FindAll(collection string, query bson.M, result interface{}, extParams model.ExtParams) (int, error) {
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
func Patch(collection, _id string, patchData *[]model.PatchParams) error {
	if len(*patchData) == 0 {
		return errors.New("empty patch content")
	}

	id := bson.ObjectIdHex(_id)

	return PatchByQuery(collection, bson.M{"_id": id}, patchData)
}

func PatchByQuery(collection string, query bson.M, patchData *[]model.PatchParams) error {
	for _, patch := range *patchData {
		log.Println("[PATCH]: " + fmt.Sprint(patch))
		var set bson.M

		k := strings.ReplaceAll(patch.Path[1:], "/", ".")
		if e := patch.ParseQuery(&query, &set, k); e != nil {
			return e
		}

		if e := Update(collection, query, set); e!=nil{
			if patch.Ignore() {
				log.Printf("[PATCH ignored]: %v -> %v\n", patch.Op, fmt.Sprint(patch))
			} else {
				return errors.New(fmt.Sprintf("faild patching _id (%v), field (%v): %v", query["_id"], patch.Path, e.(error)))
			}
		}
		delete(query, k)
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

func Update(collection string, query, set bson.M, option ...string) error {
	if e := parseObjectId(&query); e != nil {
		return e
	}

	s := Session.Clone()
	defer s.Close()

	if len(option) != 0 {
		set[option[0]] = set
	}
	return s.DB(DBName).C(collection).Update(query, set)
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
