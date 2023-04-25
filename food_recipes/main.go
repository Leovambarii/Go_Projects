package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Nutrient struct {
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
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
const jsonFileRecipes string = "recipes.json"
const jsonFileRecipe string = "recipe.json"

func main() {

	flagIngredients := flag.String("ingredients", "", "A comma-separated list of ingredients")
	flagRecipesNumber := flag.Int("numberOfRecipes", 1, "The number of recipes to display")

	flag.Parse()

	apiKeyUrl := fmt.Sprintf("apiKey=%s", apiKey)
	ingredientsParameter := fmt.Sprintf("ingredients=%s", *flagIngredients)
	numberRecipesParameter := fmt.Sprintf("number=%d", *flagRecipesNumber)
	rankingParameter := "ranking=2"

	url := fmt.Sprintf("%s?%s&%s&%v&%v", apiEndpoint, apiKeyUrl, ingredientsParameter, numberRecipesParameter, rankingParameter)

	response, err := http.Get(url)
	if err != nil {
		log.Printf("Could not make request: %v\n", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Could not read response body. %v\n", err)
	}

	var recipes []Recipe
	err = json.Unmarshal(body, &recipes)
	if err != nil {
		fmt.Printf("Could not unmarshal response body. %v\n", err)
	}

	for i := range recipes {
		url = fmt.Sprintf("https://api.spoonacular.com/recipes/%d/information?%s&includeNutrition=true", recipes[i].Id, apiKeyUrl)
		response, err = http.Get(url)
		if err != nil {
			log.Printf("Could not make request: %v\n", err)
			continue
		}
		defer response.Body.Close()

		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Could not read response body. %v\n", err)
			continue
		}

		var nutritions NutritionInfo
		err = json.Unmarshal(body, &nutritions)
		if err != nil {
			fmt.Printf("Could not unmarshal nutrition information. %v\n", err)
			continue
		}

		nutrients := nutritions.Nutrition.Nutrients

		for _, nutrient := range nutrients {
			switch nutrient.Name {
			case "Calories":
				recipes[i].Calories = nutrient
			case "Carbohydrates":
				recipes[i].Carbs = nutrient
			case "Protein":
				recipes[i].Proteins = nutrient
			}
		}
	}

	printRecipes(recipes)
}

func printRecipes(recipes []Recipe) {
	for _, recipe := range recipes {
		fmt.Println("+=======================================+")
		fmt.Printf("+ %s\n", recipe.Title)
		fmt.Println("+---------------------------------------+")

		fmt.Printf("+ %s = %.2f%s\n", recipe.Carbs.Name, recipe.Carbs.Amount, recipe.Carbs.Unit)
		fmt.Printf("+ %s = %.2f%s\n", recipe.Calories.Name, recipe.Calories.Amount, recipe.Calories.Unit)
		fmt.Printf("+ %s = %.2f%s\n", recipe.Proteins.Name, recipe.Proteins.Amount, recipe.Proteins.Unit)

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
