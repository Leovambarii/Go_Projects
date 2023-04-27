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

// const apiKey string = "1378ac4213764197b9ed83d4b968af53"
const apiKey string = "5efb2ccd5e6042739782947221c8332b"
const dataSourceName string = "root:password@tcp(localhost)/food_recipes"

// To build and run in Windows:
// go build -o recipeFinder.exe main.go
// recipeFinder.exe --ingredients="tomatoes,eggs,pasta" --numberOfRecipes=1

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
	err, foundId := checkForInputInDatabase(sortedIngredients, recipesNumber)
	if err != nil {
		fmt.Printf("Error checking input in database: %v", err)
		return
	}

	if foundId > 0 {
		fmt.Printf("Input already exists in the database\n\n")

		// Get data from database
		err, recipes := getDataFromDatabase(foundId)
		if err != nil {
			fmt.Printf("Error getting data from database: %v", err)
			return
		}

		// Print recipes with their info in terminal
		printRecipes(recipes)
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

func checkForInputInDatabase(ingredients string, recipesNumber int) (error, int) {
	// Open the database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("Could not open database: %v\n", err), -1
	}
	defer db.Close()

	// Prepare the SQL statement to select rows with matching IngredientText and RecipesNumber
	stmt, err := db.Prepare("SELECT Id FROM argumentinput WHERE IngredientText = ? AND RecipesNumber = ?")
	if err != nil {
		return fmt.Errorf("Could not prepare statement: %v\n", err), -1
	}
	defer stmt.Close()

	// Execute the statement with the values passed as arguments
	rows, err := stmt.Query(ingredients, recipesNumber)
	if err != nil {
		return fmt.Errorf("Could not execute statement: %v\n", err), -1
	}
	defer rows.Close()

	// Check if any rows were returned
	if rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return fmt.Errorf("Could not scan row: %v\n", err), -1
		}
		return nil, id
	} else {
		return nil, -1
	}
}

func getDataFromDatabase(InputId int) (error, []Recipe) {
	var recipes []Recipe

	// Open the database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("Could not open database: %v\n", err), recipes
	}
	defer db.Close()

	// Prepare the SQL statement to select recipe rows with matching InputId
	stmt, err := db.Prepare("SELECT Id, RecipeId, Title FROM recipe WHERE ArgumentId = ?")
	if err != nil {
		return fmt.Errorf("Could not prepare statement select recipe rows: %v\n", err), recipes
	}
	defer stmt.Close()

	// Execute the statement with the values passed as arguments
	rows, err := stmt.Query(InputId)
	if err != nil {
		return fmt.Errorf("Could not execute statement select recipe rows: %v\n", err), recipes
	}
	defer rows.Close()

	for rows.Next() {
		var recipe Recipe
		var id int
		err := rows.Scan(&id, &recipe.Id, &recipe.Title)
		if err != nil {
			return fmt.Errorf("Could not scan row in select recipe rows: %v\n", err), recipes
		}

		// Prepare the SQL statement to select used ingredients data with matching id
		stmt, err := db.Prepare("SELECT IngredientId, Name, Amount, Unit FROM usedingredient WHERE RecipeId = ?")
		if err != nil {
			return fmt.Errorf("Could not prepare statement select used ingredients rows: %v\n", err), recipes
		}
		defer stmt.Close()

		// Execute the statement with the values passed as arguments
		res, err := stmt.Query(id)
		if err != nil {
			return fmt.Errorf("Could not execute statement select used ingredients rows: %v\n", err), recipes
		}
		defer res.Close()

		for res.Next() {
			var usedIngredient Ingredient
			err := res.Scan(&usedIngredient.Id, &usedIngredient.Name, &usedIngredient.Amount, &usedIngredient.Unit)
			if err != nil {
				return fmt.Errorf("Could not scan row in select used ingredients rows: %v\n", err), recipes
			}

			recipe.UsedIngredients = append(recipe.UsedIngredients, usedIngredient)
		}

		// Prepare the SQL statement to select missng ingredients data with matching id
		stmt, err = db.Prepare("SELECT IngredientId, Name, Amount, Unit FROM missingingredient WHERE RecipeId = ?")
		if err != nil {
			return fmt.Errorf("Could not prepare statement select missing ingredients rows: %v\n", err), recipes
		}
		defer stmt.Close()

		// Execute the statement with the values passed as arguments
		res, err = stmt.Query(id)
		if err != nil {
			return fmt.Errorf("Could not execute statement select missing ingredients rows: %v\n", err), recipes
		}
		defer res.Close()

		for res.Next() {
			var missingIngredient Ingredient
			err := res.Scan(&missingIngredient.Id, &missingIngredient.Name, &missingIngredient.Amount, &missingIngredient.Unit)
			if err != nil {
				return fmt.Errorf("Could not scan row: in select missing ingredients rows %v\n", err), recipes
			}

			recipe.MissedIngredients = append(recipe.MissedIngredients, missingIngredient)
		}

		// Prepare the SQL statement to select nutrients data with matching id
		stmt, err = db.Prepare("SELECT Name, Amount, Unit FROM nutrients WHERE RecipeId = ?")
		if err != nil {
			return fmt.Errorf("Could not prepare statement select nutrients rows: %v\n", err), recipes
		}
		defer stmt.Close()

		// Execute the statement with the values passed as arguments
		res, err = stmt.Query(id)
		if err != nil {
			return fmt.Errorf("Could not execute statement select nutrients rows: %v\n", err), recipes
		}
		defer res.Close()

		for res.Next() {
			var nutrient Nutrient
			err := res.Scan(&nutrient.Name, &nutrient.Amount, &nutrient.Unit)
			if err != nil {
				return fmt.Errorf("Could not scan row in select nutrients rows: %v\n", err), recipes
			}

			switch nutrient.Name {
			case "Carbohydrates":
				recipe.Carbs = nutrient
			case "Protein":
				recipe.Proteins = nutrient
			case "Calories":
				recipe.Calories = nutrient
			}
		}

		// Append the recipe to the recipes slice
		recipes = append(recipes, recipe)
	}

	return nil, recipes
}

func storeRecipesInDatabase(ingredients string, recipesNumber int, recipes []Recipe) error {
	fmt.Println("Caching the data from API into local database...")

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

	// Loop through the recipe slice and insert all data from each recipe into the database
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

		// Insert data of all usedingredients
		for _, usedIngredient := range recipe.UsedIngredients {
			// Insert usedigredient data
			_, err := stmt.Exec(usedIngredient.Id, recipeID, usedIngredient.Name, usedIngredient.Amount, usedIngredient.Unit)
			if err != nil {
				return fmt.Errorf("Could not insert usedingredient data: %v\n", err)
			}
		}

		// Prepare the SQL statement insert missingingredient data
		stmt, err = db.Prepare("INSERT INTO missingingredient (IngredientId, RecipeId, Name, Amount, Unit) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			return fmt.Errorf("Could not prepare missingingredient statement: %v\n", err)
		}
		defer stmt.Close()

		// Insert data of all missingingredients
		for _, missedIngredient := range recipe.MissedIngredients {
			// Insert missingingredient data
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
