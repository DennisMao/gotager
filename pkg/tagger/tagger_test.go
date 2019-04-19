package tagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag(t *testing.T) {
}

func TestTagConvert(t *testing.T) {
	fieldName := "HelloWorld"
	prefix := "json"

	//	t.Log(tagConvert(fieldName, prefix, STYLE_RAW))
	//	t.Log(tagConvert(fieldName, prefix, STYLE_CAMEL))
	//	t.Log(tagConvert(fieldName, prefix, STYLE_CAMEL_LOWER))
	//	t.Log(tagConvert(fieldName, prefix, STYLE_SNAKE))
	//	t.Log(tagConvert(fieldName, prefix, STYLE_LOWER))

	assert.Equal(t, "json:\"HelloWorld\"", tagConvert(fieldName, prefix, STYLE_RAW))
	assert.Equal(t, "json:\"HelloWorld\"", tagConvert(fieldName, prefix, STYLE_CAMEL))
	assert.Equal(t, "json:\"helloWorld\"", tagConvert(fieldName, prefix, STYLE_CAMEL_LOWER))
	assert.Equal(t, "json:\"hello_world\"", tagConvert(fieldName, prefix, STYLE_SNAKE))
	assert.Equal(t, "json:\"helloworld\"", tagConvert(fieldName, prefix, STYLE_LOWER))
}

func TestTagGenerage(t *testing.T) {
	fieldName := "IdWORLD"
	rawTag := "`json:\"Id\"`"
	prefix := "json"
	style := STYLE_SNAKE

	assert.Equal(t, "`json:\"id_world\"`", tagGenerate(rawTag, fieldName, prefix, style, true))
	assert.Equal(t, rawTag, tagGenerate(rawTag, fieldName, prefix, style, false))

}
