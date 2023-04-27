# **Recipe Finder**

Recipe finder is a CLI application that accepts two arguments: ingredients and numberOfRecipes.

This project allows the user to find recipes based on a list of ingredients using the Spoonacular API. If the same list of ingredients and number of recipes have been previously searched for, the package retrieves the cached information from a MySQL database instead of using the API. The package returns recipes in the terminal, sorted by number of missing ingredients.

# **Setup**
To get the code or app you can download the repository food_recipes folder or clone the repository. You can also use 'go get' command.

To build the app, run the following commands:
+ **Linux**
    ```
    go build -o recipeFinder main.go
    ```
+ **Windows**
    ```
    go build -o recipeFinder.exe main.go
    ```

# **Usage**

To use the application after building it, execute the command with the subsequent flags:
+ --ingredients: A comma-separated list of ingredients
+ --numberOfRecipes: The number of recipes to display

For example:
+ **Linux**:
    ```
    ./recipeFinder --ingredients="tomatoes,eggs,pasta" --numberOfRecipes=1
    ```
+ **Windows**:
    ```
    recipeFinder.exe --ingredients="tomatoes,eggs,pasta" --numberOfRecipes=1
    ```

Please, include " " at the start and end of ingredients flag input to make sure the ingredients containining space do not cut the list of available ingredients.

# **Functions**

The application includes the following functions in order to work:

## **sortIngredients**

Sorts a comma-separated string of ingredients and returns a new comma-separated sorted string. This assures that same inputs of ingredients with different order be recognized as the same one.

## **checkForInputInDatabase**

Searches the input table in the MySQL database for a row that has matching IngredientText and RecipesNumber - values given by the user in the argument flags.

It returns an error if any occurred and the Id value of the matching row, or -1 if no row was found/error occurred.

## **getDataFromDatabase**

Retrieves recipe data from the MySQL database based on the input id. If successful, it returns a nil error and a slice of Recipe structs.

## **storeRecipesInDatabase**

Stores recipe slice data in the database. If successful, it returns a nil error.

## **getRecipies**

Retrieves recipes from a given URL.
It makes a GET request to the specified URL and unmarshals the response body
into a slice of Recipe structs. If any errors occur during the process, an
error is returned along with a nil slice of Recipe structs.

## **getNutritionInfo**

Retrieves nutrition information for a slice of recipes using the Spoonacular API.

## **printRecipes**

Prints the information stored in a given slice of Recipe structs.