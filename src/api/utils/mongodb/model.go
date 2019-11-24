package mongodb

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

func init() {
	PatchOp.Add = "add"
	PatchOp.Move = "move"
	PatchOp.Replace = "replace"
	PatchOp.Remove = "remove"
	PatchOp.Copy = "copy"
	PatchOp.Test = "test"
	PatchOp.AddUp = "add_up"
}

