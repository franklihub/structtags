package gtags

import (
	"encoding/json"
	"fmt"
	"testing"

	"gotest.tools/assert"
)

type HeartCfg int
type NestedParams struct {
	Host    string `json:"nested_host" `
	TimeOut int    `json:"time_out" d:"10" max:"1000" min:"10"`
}
type CfgParams struct {
	Host    string `json:"host" v:"required"`
	TimeOut int    `json:"time_out" d:"10" max:"1000" min:"10"`
	NestedParams
	Conn     NestedParams `json:"conn"`
	Port     int          `json:"port,omitempty"`
	Username string       `json:"username,omitempty" v:"required" d:"user"`
	Password string       `json:"password" v:"required"`
	Heart    HeartCfg     `json:"heart" d:"allways"`
}

func Test_ParseTag(t *testing.T) {
	//
	cfg := &CfgParams{
		NestedParams: NestedParams{
			Host: "nested_host",
		},
		Host: "host",
	}
	d, _ := json.Marshal(cfg)
	fmt.Println(string(d))
	stags := ParseStructTags(cfg)
	tag := stags.FieldByName("Host").Tags().Get("json")
	assert.Equal(t, "json", tag.Key())
	assert.Equal(t, "host", tag.Val())
	///
	assert.Equal(t, stags.FieldByName("NestedParams") == nil, true)
	tag = stags.FieldByName("nested_host").Tags().Get("json")
	assert.Equal(t, "nested_host", tag.Val())
	//todo: nested field
	// tag = stags.Field("nested_host").Tags().Get("json")
	// assert.Equal(t, "nested_host", tag.Val())

	conn := stags.NestedByName("Conn")
	field := conn.FieldByName("TimeOut")
	tags := field.Tags()
	assert.Equal(t, "json", tags.Get("json").Key())
	assert.Equal(t, "time_out", tags.Get("json").Val())
	assert.Equal(t, "1000", tags.Get("max").Val())
	assert.Equal(t, "10", tags.Get("d").Val())
	assert.Equal(t, "10", tags.Get("min").Val())

	field = conn.FieldByName("Absent")
	assert.Equal(t, field == nil, true)

}

func Test_ParseStructTag(t *testing.T) {
	//
	cfg := &CfgParams{
		NestedParams: NestedParams{
			Host: "nested_host",
		},
		Host: "host",
	}
	d, _ := json.Marshal(cfg)
	fmt.Println(string(d))
	stags := ParseStructTags(cfg)
	tag := stags.FieldByName("Host").Tags().Get("json")
	assert.Equal(t, "json", tag.Key())
	assert.Equal(t, "host", tag.Val())
	///
	//todo: nested field
	// tag = stags.Field("nested_host").Tags().Get("json")
	// assert.Equal(t, "nested_host", tag.Val())

	conn := stags.NestedByName("Conn")
	field := conn.FieldByName("TimeOut")
	tags := field.Tags()
	assert.Equal(t, "json", tags.Get("json").Key())
	assert.Equal(t, "time_out", tags.Get("json").Val())
	assert.Equal(t, "1000", tags.Get("max").Val())
	assert.Equal(t, "10", tags.Get("d").Val())
	assert.Equal(t, "10", tags.Get("min").Val())

	field = conn.FieldByName("Absent")
	tags = field.Tags()
	assert.Equal(t, "", tags.Get("json").Key())

}

///
func Test_Field_index(t *testing.T) {
	stags := ParseStructTags(&CfgParams{})
	///
	for _, f := range stags.Fields() {
		i := len(f.Index())
		assert.Equal(t, i, 1)
	}
	///
	for _, nested := range stags.Nesteds() {
		i := len(nested.Index())
		assert.Equal(t, i, 1)
		for _, f := range nested.Fields() {
			i := len(f.Index())
			assert.Equal(t, i, 2)
		}
	}
}