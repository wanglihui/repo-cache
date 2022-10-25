package repocache_test

import (
	"encoding/json"
	"fmt"

	repocache "github.com/wanglihui/repo-cache"
)

type EntityTest struct {
	ID   repocache.ID `json:"id,omitempty"`
	Name string       `json:"name,omitempty" :"name" json:"name,omitempty"`
	Age  int32        `json:"age" json:"age,omitempty"`
}

func (e *EntityTest) GetID() repocache.ID {
	return e.ID
}

func (e *EntityTest) Serialize() []byte {
	bs, _ := json.Marshal(e)
	return bs
}

func (e *EntityTest) Deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, e); err != nil {
		return err
	}
	fmt.Println(e)
	return nil
}
