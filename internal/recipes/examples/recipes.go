package examples

var (
	Question = "Give me a simple recipe that includes the following ingredients: tomato, mushroom. Also, the recipe cannot have the following ingredients: onion, garlic. I would like to have only 4 sections separated by the pipe character |. Something like the following: name: x | ingredients: y | instructions: w | calories per serving: z. Also, split the list of ingredients by semicolon character ;"
	Answer   = "name: %s | \ningredients: 2 tomatoes; 1 cup mushrooms; 1 tablespoon olive oil; 1 tablespoon balsamic vinegar; salt and pepper to taste | \ninstructions: \n1. Slice the tomatoes and mushrooms into bite-sized pieces.\n2. In a mixing bowl, combine the tomatoes and mushrooms.\n3. Drizzle olive oil and balsamic vinegar over the mixture. \n4. Season with salt and pepper to taste.\n5. Gently toss the ingredients in the bowl until well combined.\n6. Serve the tomato mushroom salad immediately.\n| calories per serving: 75."
)
