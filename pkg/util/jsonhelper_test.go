package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinifyJSON(t *testing.T) {
	t.Parallel()

	inputAndExpected := [][]string{
		// {foo:"bar"}
		{
			"{\"foo\":\"bar\"}", // input
			"{\"foo\":\"bar\"}", // expected
		},
		{
			"{  \"foo\"  :  \"bar\"  }",
			"{\"foo\":\"bar\"}",
		},
		{
			"{\n  \"foo\"  :  \"bar\"  \n}",
			"{\"foo\":\"bar\"}",
		},

		// {foo:1, foo2:1.2, foo3:true, foo4:null}
		{
			"{\"foo\":1,\"foo\":1.2,\"foo\":true,\"foo\":null}",
			"{\"foo\":1,\"foo\":1.2,\"foo\":true,\"foo\":null}",
		},
		{
			"{\n  \"foo\":1,\n  \"foo\":1.2,\n  \"foo\":true,\n  \"foo\":null\n}",
			"{\"foo\":1,\"foo\":1.2,\"foo\":true,\"foo\":null}",
		},

		// {foo:"bar"  <= uncompleted json
		{
			"{\"foo\":\"bar\"",
			"{\"foo\":\"bar\"",
		},
		{
			"{\n\"foo\":\"bar\"",
			"{\"foo\":\"bar\"", // => the new line should be erased
		},

		// {foo:"ba\"r ba\"r    ba\"r"}
		{
			"{\n  \"foo\"  :  \"ba\\\"r ba\\\"r    ba\\\"r\"  \n}",
			"{\"foo\":\"ba\\\"r ba\\\"r    ba\\\"r\"}",
		},

		// {foo:"bar\nbar"}
		{
			"{\"foo\":\"bar\\nbar\"}",
			"{\"foo\":\"bar\\nbar\"}",
		},

		// {foo:"bar",foo2:{foo3:"bar"}}
		{
			"{\"foo\":\"bar\",\"foo2\":{\"foo3\":\"bar\"}}",
			"{\"foo\":\"bar\",\"foo2\":{\"foo3\":\"bar\"}}",
		},
		{
			"{\n  \"foo\"  :  \"bar\",\n  \"foo2\"  :  {\n  \"foo3\"  :  \"bar\"\n  }  \n}",
			"{\"foo\":\"bar\",\"foo2\":{\"foo3\":\"bar\"}}",
		},

		// {foo:"bar",foo2:{foo3:"bar"}  <= uncompleted json
		{
			"{\"foo\":\"bar\",\"foo2\":{\"foo3\":\"bar\"}",
			"{\"foo\":\"bar\",\"foo2\":{\"foo3\":\"bar\"}",
		},

		// {foo:"bar",foo2:{foo3:["bar","bar"]}}
		{
			"{  \"foo\"  :  \"bar\",  \"foo2\"  :  {  \"foo3\"  :  [  \"bar\",  \"bar\"  ]  }  }",
			"{\"foo\":\"bar\",\"foo2\":{\"foo3\":[\"bar\",\"bar\"]}}",
		},
	}

	for i := 0; i < len(inputAndExpected); i++ {
		input := inputAndExpected[i][0]
		expected := inputAndExpected[i][1]

		bytes := []byte(input)
		result := MinifyJSON(bytes)
		fmt.Println("-------------------------------")
		fmt.Printf("input: %v\n", input)
		fmt.Printf("expected: %v\n", expected)
		fmt.Printf("result: %v\n", result)

		assert.Equal(t, expected, result)
	}
}
