package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// TODO parsing of first response and then using that to get recipe carbs, proteins and calories from summed ingredients
type Recipe struct {
	Id                int          `json:"id"`
	Title             string       `json:"title"`
	UsedIngredients   []Ingredient `json:"usedIngredients"`
	MissedIngredients []Ingredient `json:"missedIngredients"`
	Carbs             float64      `json:"carbs"`
	Proteins          float64      `json:"protein"`
	Calories          float64      `json:"calories"`
}

type Ingredient struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
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
		log.Printf("Could not make request: %v", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Could not read response body. %v", err)
	}

	var recipes []Recipe
	err = json.Unmarshal(body, &recipes)
	if err != nil {
		fmt.Printf("Could not unmarshal response body. %v", err)
	}

	for _, recipe := range recipes {
		url = fmt.Sprintf("https://api.spoonacular.com/recipes/%d/information?%s&includeNutrition=true", recipe.Id, apiKeyUrl)
		response, err = http.Get(url)
		if err != nil {
			log.Printf("Could not make request: %v", err)
			return
		}
		defer response.Body.Close()

		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Could not read response body. %v", err)
		}

		var data struct {
			Nutrition struct {
				Nutrients []struct {
					Name   string
					Amount float64
					Unit   string
				}
			}
		}
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Printf("Could not Unmarshal response body. %v", err)
		}

		for _, nutrient := range data.Nutrition.Nutrients { // TODO przerobić aby tylko wyłuskać te 3 informacje
			switch nutrient.Name {
			case "Calories":
				recipe.Calories = nutrient.Amount
			case "Carbohydrates":
				recipe.Carbs = nutrient.Amount
			case "Protein":
				recipe.Proteins = nutrient.Amount
			}
		}

		// err = ioutil.WriteFile(jsonFileRecipe, body, 0644)
		// if err != nil {
		// 	fmt.Printf("Could not write response to file. %v", err)
		// 	return
		// }
		// fmt.Printf("Response written to %s\n", jsonFileRecipe)
	}

	printRecipes(recipes)

	// err = ioutil.WriteFile(jsonFileRecipes, body, 0644)
	// if err != nil {
	// 	fmt.Printf("Could not write response to file. %v", err)
	// 	return
	// }
	// fmt.Printf("Response written to %s\n", jsonFileRecipes)
}

func printRecipes(recipes []Recipe) {
	for _, recipe := range recipes {
		fmt.Printf("Recipe: %s\n", recipe.Title)
		fmt.Printf("Carbs=  Proteins= Calories=")
		fmt.Println("Used Ingredients:")
		for _, usedIngredient := range recipe.UsedIngredients {
			fmt.Printf("- %s\n", usedIngredient.Name)
		}
		fmt.Println("Missing Ingredients:")
		for _, missedIngredient := range recipe.MissedIngredients {
			fmt.Printf("- %s\n", missedIngredient.Name)
		}
		fmt.Println("-----------------------------")
	}
}
