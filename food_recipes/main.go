package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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
const dataSourceName string = "root:password@tcp(localhost)/food_recipes"

// To build and run in Windows:
// go build -o recipeFinder.exe main.go
// recipeFinder.exe --ingredients=tomatoes,eggs,pasta --numberOfRecipes=1

func main() {
	// Define the flags
	flagIngredients := flag.String("ingredients", "", "A comma-separated list of ingredients")
	flagRecipesNumber := flag.Int("numberOfRecipes", 1, "The number of recipes to display")

	// Parse the command line arguments
	flag.Parse()

	recipesNumber := *flagRecipesNumber
	sortedIngredients := sortIngredients(*flagIngredients)
	fmt.Printf("Ingredients: %s\n", sortedIngredients)
	fmt.Printf("Recipes number: %d\n\n", recipesNumber)

	// Check for errors in the flag values
	if sortedIngredients == "" {
		fmt.Println("Error: You must provide a list of ingredients")
		return
	}

	if recipesNumber <= 0 {
		fmt.Println("Error: The number of recipes must be greater than zero")
		return
	}

	// Check for similar input in database
	err, found := checkForSameInput(sortedIngredients, recipesNumber)
	if err != nil {
		fmt.Printf("Error checking input in database: %v", err)
		return
	}

	if found {
		fmt.Println("Input already exists in the database")
	} else {
		fmt.Println("Getting information from API")

		// Create string parameters for url request
		apiKeyUrl := fmt.Sprintf("apiKey=%s", apiKey)
		ingredientsParameter := fmt.Sprintf("ingredients=%s", sortedIngredients)
		numberRecipesParameter := fmt.Sprintf("number=%d", recipesNumber)
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

		// Cache information in database
		err = storeRecipesInDatabase(sortedIngredients, recipesNumber, recipes)
		if err != nil {
			fmt.Printf("Error caching data: %v", err)
			return
		}
	}

}

func sortIngredients(ingredients string) string {
	// Split the comma-separated string of ingredients into a slice of strings
	ingredientsSlice := strings.Split(ingredients, ",")

	// Sort the slice of strings so every input with same ingredients has the same string
	sort.Strings(ingredientsSlice)

	// Create a comma-separated string from the sorted slice of strings
	sortedIngredientsString := strings.Join(ingredientsSlice, ",")

	return sortedIngredientsString
}

// Database related functions

func checkForSameInput(ingredients string, recipesNumber int) (error, bool) {
	// Open the database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("Could not open database: %v\n", err), false
	}
	defer db.Close()

	// Prepare the SQL statement to select rows with matching IngredientText and RecipesNumber
	stmt, err := db.Prepare("SELECT IngredientText, RecipesNumber FROM argumentinput WHERE IngredientText = ? AND RecipesNumber = ?")
	if err != nil {
		return fmt.Errorf("Could not prepare statement: %v\n", err), false
	}
	defer stmt.Close()

	// Execute the statement with the values passed as arguments
	rows, err := stmt.Query(ingredients, recipesNumber)
	if err != nil {
		return fmt.Errorf("Could not execute statement: %v\n", err), false
	}
	defer rows.Close()

	// Check if any rows were returned
	if rows.Next() {
		return nil, true
	} else {
		return nil, false
	}
}

func storeRecipesInDatabase(ingredients string, recipesNumber int, recipes []Recipe) error {
	fmt.Println("Caching the data from API into local database")

	// Open the database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("Could not open database: %v\n", err)
	}
	defer db.Close()

	// Prepare the SQL statement insert input data
	stmt, err := db.Prepare("INSERT INTO argumentinput (IngredientText, RecipesNumber) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("Could not prepare statement: %v\n", err)
	}
	defer stmt.Close()

	// Insert input data
	result, err := stmt.Exec(ingredients, recipesNumber)
	if err != nil {
		return fmt.Errorf("Could not insert input data: %v\n", err)
	}

	// Get the id of the inserted row
	argumentInputID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Could not get id of inserted row: %v\n", err)
	}

	// Loop through the recipe slice and insert each recipe into the database
	for _, recipe := range recipes {
		// Prepare the SQL statement insert recipe data
		stmt, err = db.Prepare("INSERT INTO recipe (RecipeId, ArgumentId, Title) VALUES (?, ?, ?)")
		if err != nil {
			return fmt.Errorf("Could not prepare statement: %v\n", err)
		}
		defer stmt.Close()

		// Insert recipe data
		result, err := stmt.Exec(recipe.Id, argumentInputID, recipe.Title)
		if err != nil {
			return fmt.Errorf("Could not insert recipe data: %v\n", err)
		}

		// Get the id of the inserted row
		recipeID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("Could not get id of inserted row: %v\n", err)
		}

		// Prepare the SQL statement insert nutrients data
		stmt, err = db.Prepare("INSERT INTO nutrients (RecipeId, Name, Amount, Unit) VALUES (?, ?, ?, ?)")
		if err != nil {
			return fmt.Errorf("Could not prepare nutrients statement: %v\n", err)
		}
		defer stmt.Close()
		// Insert carbs nutrient data
		_, err = stmt.Exec(recipeID, recipe.Carbs.Name, recipe.Carbs.Amount, recipe.Carbs.Unit)
		if err != nil {
			return fmt.Errorf("Could not insert carbs nutrient data: %v\n", err)
		}
		// Insert proteins nutrient data
		_, err = stmt.Exec(recipeID, recipe.Proteins.Name, recipe.Proteins.Amount, recipe.Proteins.Unit)
		if err != nil {
			return fmt.Errorf("Could not insert proteins nutrient data: %v\n", err)
		}
		// Insert calories nutrient data
		_, err = stmt.Exec(recipeID, recipe.Calories.Name, recipe.Calories.Amount, recipe.Calories.Unit)
		if err != nil {
			return fmt.Errorf("Could not insert calories nutrient data: %v\n", err)
		}

		// Prepare the SQL statement insert usedingredient data
		stmt, err = db.Prepare("INSERT INTO usedingredient (IngredientId, RecipeId, Name, Amount, Unit) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			return fmt.Errorf("Could not prepare usedingredient statement: %v\n", err)
		}
		defer stmt.Close()

		for _, usedIngredient := range recipe.UsedIngredients {
			_, err := stmt.Exec(usedIngredient.Id, recipeID, usedIngredient.Name, usedIngredient.Amount, usedIngredient.Unit)
			if err != nil {
				return fmt.Errorf("Could not insert usedingredient data: %v\n", err)
			}
		}

		// Prepare the SQL statement insert usedingredient data
		stmt, err = db.Prepare("INSERT INTO missingingredient (IngredientId, RecipeId, Name, Amount, Unit) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			return fmt.Errorf("Could not prepare missingingredient statement: %v\n", err)
		}
		defer stmt.Close()

		for _, missedIngredient := range recipe.MissedIngredients {
			_, err := stmt.Exec(missedIngredient.Id, recipeID, missedIngredient.Name, missedIngredient.Amount, missedIngredient.Unit)
			if err != nil {
				return fmt.Errorf("Could not insert missingingredient data: %v\n", err)
			}
		}

	}

	return nil
}

// API related functions

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

// Print information stored in given Recipe slice
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
