package hack

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFields(t *testing.T) {
	type St struct{
		Str string
		Num int
		LargeNum int64
		Float float64
		Bool bool
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "Str",
				Type: FieldString,
			},
			{
				Name: "Num",
				Type: FieldInteger,
			},
			{
				Name: "LargeNum",
				Type: FieldInteger,
			},
			{
				Name: "Float",
				Type: FieldNumber,
			},
			{
				Name: "Bool",
				Type: FieldBoolean,
			},
		},
	}, schema)
}

func TestRequiredField(t *testing.T) {
	type St struct{
		Str string `binding:"required"`
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "Str",
				Type: FieldString,
				Required: true,
			},
		},
	}, schema)
}

func TestJsonField(t *testing.T) {
	type St struct{
		Str string `json:"my_str"`
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "my_str",
				Type: FieldString,
			},
		},
	}, schema)
}

func TestJsonFieldOmitEmpty(t *testing.T) {
	type St struct{
		Str string `json:"my_str,omitempty"`
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "my_str",
				Type: FieldString,
				OmitEmpty: true,
			},
		},
	}, schema)
}

func TestJsonFieldIgnore(t *testing.T) {
	type St struct{
		Str string `json:"-"`
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: nil,
	}, schema)
}

func TestFormField(t *testing.T) {
	type St struct{
		Str string `form:"my_str"`
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "my_str",
				Type: FieldString,
			},
		},
	}, schema)
}

func TestFormFieldIgnore(t *testing.T) {
	type St struct{
		Str string `form:"-"`
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: nil,
	}, schema)
}

func TestJsonRawMessage(t *testing.T) {
	type St struct{
		RawJson json.RawMessage
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "RawJson",
				Type: FieldString,
			},
		},
	}, schema)
}

func TestArray(t *testing.T) {
	type St struct{
		ManyStrings []string
		ManyInts []int
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "ManyStrings",
				Type: FieldArray,
				ArrayType: &Schema{
					Type: FieldString,
				},
			},
			{
				Name: "ManyInts",
				Type: FieldArray,
				ArrayType: &Schema{
					Type: FieldInteger,
				},
			},
		},
	}, schema)
}

func TestAnonStruct(t *testing.T) {
	st := struct{
		Foo string
	}{
		"",
	}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "Foo",
				Type: FieldString,
			},
		},
	}, schema)
}

func TestNestedStruct(t *testing.T) {
	type Inner struct{
		Bar string
	}
	type Outer struct {
		Foo Inner
	}
	st := Outer{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "Foo",
				Type: FieldObject,
				Children: []Schema{
					{
						Name: "Bar",
						Type: FieldString,
					},
				},
			},
		},
	}, schema)
}

func TestArrayJsonBody(t *testing.T) {
	var st []string
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldArray,
		ArrayType: &Schema{
			Type: FieldString,
		},
	}, schema)
}

func TestPointer(t *testing.T) {
	type St struct{
		Str *string
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "Str",
				Type: FieldString,
			},
		},
	}, schema)
}

func TestRecursiveEmptyStruct(t *testing.T) {
	type St struct {
		Recursive *St
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "Recursive",
				Type: FieldObject,
				Recursive: true,
			},
		},
	}, schema)
}

func TestRecursiveStruct(t *testing.T) {
	type St struct {
		Recursive *St
		Str string
	}
	st := St{}
	schema := StructToSchema(st)
	assert.Equal(t, Schema{
		Name: "",
		Type: FieldObject,
		Children: []Schema{
			{
				Name: "Recursive",
				Type: FieldObject,
				Recursive: true,
			},
			{
				Name: "Str",
				Type: FieldString,
			},
		},
	}, schema)
}
