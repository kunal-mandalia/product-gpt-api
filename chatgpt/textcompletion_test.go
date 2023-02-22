package chatgpt

import "testing"

func Test_ProductRecommendationsQuery(t *testing.T) {
	cases := [][]string{
		{
			"What is the square root of 9?",
			"The square root of 9 is 3",
		},
		{
			"I'm a middle aged man. What should I do to get back into shape?",
			"n\n1. Start with a physical activity that you enjoy. This could be walking, running, swimming, biking, or any other activity that you find enjoyable.\n\n2. Set realistic goals for yourself. Start with small goals and gradually increase the intensity and duration of your workouts.\n\n3. Make sure to get enough rest and recovery time. This will help your body to recover and rebuild after each workout.\n\n4. Eat a balanced diet that is rich in fruits, vegetables, lean proteins, and healthy fats.\n\n5. Stay hydrated by drinking plenty of water throughout the day.\n\n6. Incorporate strength training into your routine. This will help to build muscle and burn fat.\n\n7. Make sure to stretch before and after each workout. This will help to prevent injury and improve flexibility.\n\n8. Track your progress and celebrate your successes. This will help to keep you motivated and on track.",
		},
	}

	for _, c := range cases {
		pq := ProductRecommendationsQuery(c[0], c[1])
		expected := "\ncontext:\n" +
			c[0] + "\n" +
			c[1] + "\n\n" +
			"prompt:\n" +
			"Identify and provide information on products and services found within the context above." +
			" Specifically, output an array of JSON objects which include the following:" +
			" key: name, value: name of product or service" +
			" key: link, value: a link to that product or service," +
			" key: isProduct, value: true if the link refers to a product, false if it refers to a service" +
			" key: characterRange, value: an array of JSON objects. The JSON objects should include a key called startChar and endChar for where the product or service appears in my query"

		if pq != expected {
			t.Fatalf(`ProductRecommendationsQuery(%s, %s)=%s. Want: %s`, c[0], c[1], pq, expected)
		}
	}
}
