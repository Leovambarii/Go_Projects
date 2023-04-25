package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Recipe struct {
	Id                int          `json:"id"`
	Title             string       `json:"title"`
	UsedIngredients   []Ingredient `json:"usedIngredients"`
	MissedIngredients []Ingredient `json:"missedIngredients"`
	Carbs             Nutrient     `json:"carbs"`
	Proteins          Nutrient     `json:"protein"`
	Calories          Nutrient     `json:"calories"`
}

type Ingredient struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Nutrient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Nutrition struct {
	Nutrients []Nutrient `json:"nutrients"`
}

type NutritionInfo struct {
	Nutrition `json:"nutrition"`
}

const apiEndpoint string = "https://api.spoonacular.com/recipes/findByIngredients"
const apiKey string = "1378ac4213764197b9ed83d4b968af53"

func main() {
	// Define the flags
	flagIngredients := flag.String("ingredients", "", "A comma-separated list of ingredients")
	flagRecipesNumber := flag.Int("numberOfRecipes", 1, "The number of recipes to display")

	// Parse the command line arguments
	flag.Parse()

	// Check for errors in the flag values
	if *flagIngredients == "" {
		fmt.Println("Error: You must provide a list of ingredients")
		return
	}

	if *flagRecipesNumber <= 0 {
		fmt.Println("Error: The number of recipes must be greater than zero")
		return
	}

	// Create string parameters for url request
	apiKeyUrl := fmt.Sprintf("apiKey=%s", apiKey)
	ingredientsParameter := fmt.Sprintf("ingredients=%s", *flagIngredients)
	numberRecipesParameter := fmt.Sprintf("number=%d", *flagRecipesNumber)
	rankingParameter := "ranking=2"

	url := fmt.Sprintf("%s?%s&%s&%v&%v", apiEndpoint, apiKeyUrl, ingredientsParameter, numberRecipesParameter, rankingParameter)

	// Get recipes from API
	err, recipes := getRecipes(url)
	if err != nil {
		fmt.Printf("Error getting recipes: %v", err)
		return
	}

	// Get nutrition info for recipes
	err = getNutritionInfo(apiKeyUrl, recipes)
	if err != nil {
		fmt.Printf("Error getting nutrition information: %v", err)
		return
	}

	// Print recipes with their info in terminal
	printRecipes(recipes)
}

func getRecipes(url string) (error, []Recipe) {
	var recipes []Recipe

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Could not make request: %v\n", err), recipes
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Could not read response body. %v\n", err), recipes
	}

	err = json.Unmarshal(body, &recipes)
	if err != nil {
		return fmt.Errorf("Could not unmarshal response body. %v\n", err), recipes
	}

	return nil, recipes
}

func getNutritionInfo(apiKeyUrl string, recipes []Recipe) error {
	for i := range recipes {
		url := fmt.Sprintf("https://api.spoonacular.com/recipes/%d/information?%s&includeNutrition=true", recipes[i].Id, apiKeyUrl)
		response, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to make request: %v", err)
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}

		var nutritions NutritionInfo
		err = json.Unmarshal(body, &nutritions)
		if err != nil {
			return fmt.Errorf("failed to unmarshal nutrition information: %v", err)
		}

		var nutrientMap = map[string]*Nutrient{
			"Calories":      &recipes[i].Calories,
			"Carbohydrates": &recipes[i].Carbs,
			"Protein":       &recipes[i].Proteins,
		}

		for _, nutrient := range nutritions.Nutrition.Nutrients {
			if nutrientPtr, ok := nutrientMap[nutrient.Name]; ok {
				*nutrientPtr = nutrient
			}
		}
	}

	return nil
}

func printRecipes(recipes []Recipe) {
	for _, recipe := range recipes {
		fmt.Println("+=======================================+")
		fmt.Printf("+ %s\n", recipe.Title)
		fmt.Println("+---------------------------------------+")

		fmt.Printf("+ %s = %.2f%s\n", recipe.Carbs.Name, recipe.Carbs.Amount, recipe.Carbs.Unit)
		fmt.Printf("+ %s = %.2f%s\n", recipe.Proteins.Name, recipe.Proteins.Amount, recipe.Proteins.Unit)
		fmt.Printf("+ %s = %.2f%s\n", recipe.Calories.Name, recipe.Calories.Amount, recipe.Calories.Unit)

		fmt.Println("+---------------------------------------+")

		fmt.Println("+ Used Ingredients:")
		for _, usedIngredient := range recipe.UsedIngredients {
			fmt.Printf("+ %.2f %s %s\n", usedIngredient.Amount, usedIngredient.Unit, usedIngredient.Name)
		}

		fmt.Println("+---------------------------------------+")

		fmt.Println("+ Missing Ingredients:")
		for _, missedIngredient := range recipe.MissedIngredients {
			fmt.Printf("+ %.2f %s %s\n", missedIngredient.Amount, missedIngredient.Unit, missedIngredient.Name)
		}

		fmt.Println("+=======================================+\n ")

	}
}
