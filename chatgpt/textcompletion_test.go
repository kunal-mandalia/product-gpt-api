package chatgpt

import "testing"

func Test_ProductRecommendationsQuery(t *testing.T) {
	cases := [][]string{
		{
			"What is the square root of 9?",
			"The square root of 9 is 3",
			"\ncontext:\n" +
				"What is the square root of 9?\n" +
				"The square root of 9 is 3\n\n" +
				"prompt:\n" +
				"Identify and provide information on products and services found within the context above." +
				" Specifically, output an array of JSON objects which include the following:" +
				" key: name, value: name of product or service" +
				" key: link, value: a link to that product or service," +
				" key: isProduct, value: true if the link refers to a product, false if it refers to a service" +
				" key: characterRange, value: an array of JSON objects. The JSON objects should include a key called startChar and endChar for where the product or service appears in my query",
		},
	}

	for _, c := range cases {
		pq := ProductRecommendationsQuery(c[0], c[1])
		if c[2] != pq {
			t.Fatalf(`ProductRecommendationsQuery(%s, %s)=%s. Want: %s`, c[0], c[1], pq, c[2])
		}
	}
}
