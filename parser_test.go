package envparser

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse_NonPointer_Error(t *testing.T) {
	t.Setenv("STRING_VAL", "hello")
	type Env struct {
		StringVal string `env:"STRING_VAL"`
	}
	var env Env
	err := Parse(env) // non pointer parsing
	assert.Error(t, err)
}

func TestParse_MissingENV_Error(t *testing.T) {
	type Env struct {
		StringVal string `env:"STRING_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Ignored(t *testing.T) {
	type Env struct {
		StringVal  string `env:"-"`
		IgnoredVal string
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.StringVal, "")
	assert.Equal(t, env.IgnoredVal, "")
}

func TestParse_String(t *testing.T) {
	t.Setenv("STRING_VAL", "hello")
	type Env struct {
		StringVal string `env:"STRING_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.StringVal, "hello")
}

func TestParse_Bool(t *testing.T) {
	t.Setenv("BOOL_VAL", "true")
	type Env struct {
		BoolVal bool `env:"BOOL_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.BoolVal, true)
}

func TestParse_Bool_Error(t *testing.T) {
	t.Setenv("BOOL_VAL", "not true")
	type Env struct {
		BoolVal bool `env:"BOOL_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Int(t *testing.T) {
	t.Setenv("INT_VAL", "2")
	type Env struct {
		IntVal int `env:"INT_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.IntVal, 2)
}

func TestParse_Int_Error(t *testing.T) {
	t.Setenv("INT_VAL", "2e")
	type Env struct {
		IntVal int `env:"INT_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Float(t *testing.T) {
	t.Setenv("FLOAT_VAL", "3.14")
	type Env struct {
		FloatVal float64 `env:"FLOAT_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.FloatVal, 3.14)
}

func TestParse_Float_Error(t *testing.T) {
	t.Setenv("FLOAT_VAL", "3.14e")
	type Env struct {
		FloatVal float64 `env:"FLOAT_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Uint(t *testing.T) {
	t.Setenv("UINT_VAL", "3")
	type Env struct {
		UintVal uint `env:"UINT_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.UintVal, uint(3))
}

func TestParse_Uint_Error(t *testing.T) {
	t.Setenv("UINT_VAL", "-3")
	type Env struct {
		UintVal uint `env:"UINT_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Duration(t *testing.T) {
	t.Setenv("DURATION_VAL", "2h30m")
	type Env struct {
		DurationVal time.Duration `env:"DURATION_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.DurationVal, 2*time.Hour+30*time.Minute)
}

func TestParse_Duration_Error(t *testing.T) {
	t.Setenv("DURATION_VAL", "2hh30mm")
	type Env struct {
		DurationVal time.Duration `env:"DURATION_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Datetime(t *testing.T) {
	t.Setenv("DATETIME_VAL", "2023-10-01T15:04:05Z")
	type Env struct {
		DateTimeVal time.Time `env:"DATETIME_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)

	expectedTime, _ := time.Parse(time.RFC3339, "2023-10-01T15:04:05Z")
	assert.Equal(t, env.DateTimeVal, expectedTime)
}

func TestParse_Datetime_Error(t *testing.T) {
	t.Setenv("DATETIME_VAL", "29-01-2024 15:00:00")
	type Env struct {
		DateTimeVal time.Time `env:"DATETIME_VAL"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_StringSlice(t *testing.T) {
	t.Setenv("STRING_SLICE", "a,b,c")
	type Env struct {
		StringSliceVal []string `env:"STRING_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.StringSliceVal, []string{"a", "b", "c"})
}

func TestParse_IntSlice(t *testing.T) {
	t.Setenv("INT_SLICE", "1,2,3")
	type Env struct {
		IntSlice []int `env:"INT_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.IntSlice, []int{1, 2, 3})
}

func TestParse_IntSlice_Error(t *testing.T) {
	t.Setenv("INT_SLICE", "1,2,3e")
	type Env struct {
		IntSlice []int `env:"INT_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Int32Slice(t *testing.T) {
	t.Setenv("INT32_SLICE", "1,2,3")
	type Env struct {
		IntSlice []int32 `env:"INT32_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.IntSlice, []int32{1, 2, 3})
}

func TestParse_Int32Slice_Error(t *testing.T) {
	t.Setenv("INT32_SLICE", "1,2,3e")
	type Env struct {
		IntSlice []int32 `env:"INT32_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Int64Slice(t *testing.T) {
	t.Setenv("INT64_SLICE", "1,2,3")
	type Env struct {
		IntSlice []int64 `env:"INT64_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.IntSlice, []int64{1, 2, 3})
}

func TestParse_Int64Slice_Error(t *testing.T) {
	t.Setenv("INT64_SLICE", "1,2,3e")
	type Env struct {
		IntSlice []int64 `env:"INT64_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Float32Slice(t *testing.T) {
	t.Setenv("FLOAT32_SLICE", "1.1,2,3.2")
	type Env struct {
		IntSlice []float32 `env:"FLOAT32_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.IntSlice, []float32{1.1, 2, 3.2})
}

func TestParse_Float32Slice_Error(t *testing.T) {
	t.Setenv("FLOAT32_SLICE", "1.1,2,3.2e")
	type Env struct {
		IntSlice []float32 `env:"FLOAT32_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Float64Slice(t *testing.T) {
	t.Setenv("FLOAT64_SLICE", "1.1,2,3.2")
	type Env struct {
		IntSlice []float64 `env:"FLOAT64_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.IntSlice, []float64{1.1, 2, 3.2})
}

func TestParse_Float64Slice_Error(t *testing.T) {
	t.Setenv("FLOAT64_SLICE", "1.1,2,3.2e")
	type Env struct {
		IntSlice []float64 `env:"FLOAT64_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_UintSlice(t *testing.T) {
	t.Setenv("UINT_SLICE", "1,2,3")
	type Env struct {
		IntSlice []uint `env:"UINT_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.IntSlice, []uint{1, 2, 3})
}

func TestParse_UintSlice_Error(t *testing.T) {
	t.Setenv("UINT_SLICE", "1,2,-3")
	type Env struct {
		IntSlice []uint `env:"UINT_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_UintSlice32(t *testing.T) {
	t.Setenv("UINT_SLICE", "1,2,3")
	type Env struct {
		IntSlice []uint32 `env:"UINT_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.IntSlice, []uint32{1, 2, 3})
}

func TestParse_UintSlice32_Error(t *testing.T) {
	t.Setenv("UINT_SLICE", "1,-2,3")
	type Env struct {
		IntSlice []uint32 `env:"UINT_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_UintSlice64(t *testing.T) {
	t.Setenv("UINT_SLICE", "1,2,3")
	type Env struct {
		IntSlice []uint64 `env:"UINT_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.IntSlice, []uint64{1, 2, 3})
}

func TestParse_UintSlice64_Error(t *testing.T) {
	t.Setenv("UINT_SLICE", "1,2,3e")
	type Env struct {
		IntSlice []uint64 `env:"UINT_SLICE"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Nested(t *testing.T) {
	t.Setenv("NESTED", "nested value")
	type Nested struct {
		Value string `env:"NESTED"`
	}
	type Env struct {
		Nested
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.Value, "nested value")
}

func TestParse_Embeded(t *testing.T) {
	t.Setenv("EMBEDED", "embeded value")
	type Embeded struct {
		Value string `env:"EMBEDED"`
	}
	type Env struct {
		Embeded Embeded
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.Embeded.Value, "embeded value")
}

func TestParse_Embeded_Error(t *testing.T) {
	t.Setenv("EMBEDED", "embeded value")
	type Embeded struct {
		Value int `env:"EMBEDED"`
	}
	type Env struct {
		Embeded Embeded
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Encoding_JSON(t *testing.T) {
	t.Setenv("JSON_VAL", `{"field":"jsonvalue"}`)
	type JSONStruct struct {
		Field string `json:"field"`
	}
	type Env struct {
		JSONStruct JSONStruct `env:"JSON_VAL" encoding:"json"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.JSONStruct.Field, "jsonvalue")
}

func TestParse_Encoding_JSON_Error(t *testing.T) {
	t.Setenv("JSON_VAL", `"field"="jsonvalue"`)
	type JSONStruct struct {
		Field string `json:"field"`
	}
	type Env struct {
		JSONStruct JSONStruct `env:"JSON_VAL" encoding:"json"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Encoding_XML(t *testing.T) {
	t.Setenv("XML_VAL", `<XMLStruct><field>xmlvalue</field></XMLStruct>`)
	type XMLStruct struct {
		Field string `xml:"field"`
	}
	type Env struct {
		XMLStruct XMLStruct `env:"XML_VAL" encoding:"xml"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.XMLStruct.Field, "xmlvalue")
}

func TestParse_Encoding_XML_Error(t *testing.T) {
	t.Setenv("XML_VAL", `{"field":"jsonvalue"}`)
	type XMLStruct struct {
		Field string `xml:"field"`
	}
	type Env struct {
		XMLStruct XMLStruct `env:"XML_VAL" encoding:"xml"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Encoding_Form(t *testing.T) {
	t.Setenv("FORM_DATA", `field1=value&field2=val`)
	type Env struct {
		FormVal url.Values `env:"FORM_DATA" encoding:"form"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, env.FormVal.Get("field1"), "value")
	assert.Equal(t, env.FormVal.Get("field2"), "val")
}

func TestParse_Encoding_Form_Error(t *testing.T) {
	t.Setenv("FORM_DATA", "field1=value;field2=val")
	type Env struct {
		FormVal url.Values `env:"FORM_DATA" encoding:"form"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Encoding_Base64(t *testing.T) {
	t.Setenv("BASE64_DATA", `aGVsbG8gd29ybGQ=`)
	type Env struct {
		Data []byte `env:"BASE64_DATA" encoding:"base64"`
	}
	var env Env
	err := Parse(&env)
	assert.NoError(t, err)
	assert.Equal(t, string(env.Data), "hello world")
}

func TestParse_Encoding_Base64_Error(t *testing.T) {
	t.Setenv("BASE64_DATA", `1sx 123 `)
	type Env struct {
		Data []byte `env:"BASE64_DATA" encoding:"base64"`
	}
	var env Env
	err := Parse(&env)
	assert.Error(t, err)
}

func TestParse_Unexported(t *testing.T) {
	t.Setenv("DATA_UNEXPORTED", "data unexported")
	type unexported struct {
		Data string `env:"DATA_UNEXPORTED"`
	}

	var u unexported
	err := Parse(&u)
	assert.NoError(t, err)
	assert.Equal(t, u.Data, "data unexported")
}

func TestParse_Unexported_Member(t *testing.T) {
	t.Setenv("DATA_UNEXPORTED", "data unexported")
	type unexported struct {
		// unexported ignored
		dataUnexported string `env:"DATA_UNEXPORTED"`
	}

	var u unexported
	err := Parse(&u)
	assert.NoError(t, err)
	assert.Equal(t, u.dataUnexported, "")
}
