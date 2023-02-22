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
				" Specifically, output a table with three columns." +
				" The first column the will be the name of the product / service," +
				" the second column a link to that product / service," +
				" and the third column character ranges where the product / service appears in my query.",
		},
	}

	for _, c := range cases {
		pq := ProductRecommendationsQuery(c[0], c[1])
		if c[2] != pq {
			t.Fatalf(`ProductRecommendationsQuery(%s, %s)=%s. Want: %s`, c[0], c[1], pq, c[2])
		}
	}
}
