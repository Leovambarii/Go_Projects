DROP TABLE IF EXISTS `MissingIngredient`;
DROP TABLE IF EXISTS `UsedIngredient`;
DROP TABLE IF EXISTS `Nutrients`;
DROP TABLE IF EXISTS `Recipe`;
DROP TABLE IF EXISTS `ArgumentInput`;

CREATE TABLE `ArgumentInput`(
    `Index` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `IngredientText` text NOT NULL,
    `RecipesNumber` int NOT NULL,
    PRIMARY KEY (`Index`)
);

CREATE TABLE `Recipe`(
    `Index` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `RecipeId` int UNSIGNED NOT NULL,
    `ArgumentId` int UNSIGNED NOT NULL,
    `Title` text NOT NULL,
    `Servings` int UNSIGNED NOT NULL,
    PRIMARY KEY (`Index`),
    KEY `recipe_argumentid_foreign` (`ArgumentId`),
    CONSTRAINT `recipe_argumentid_foreign` FOREIGN KEY (`ArgumentId`) REFERENCES `ArgumentInput` (`Index`)
);

CREATE TABLE `MissingIngredient`(
    `Index` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `IngredientId` int UNSIGNED NOT NULL,
    `RecipeIndex` int UNSIGNED NOT NULL,
    `Name` text NOT NULL,
    `Amount` double(8, 2) NOT NULL,
    `Unit` text,
    PRIMARY KEY (`Index`),
    KEY `missingingredient_recipeid_foreign` (`RecipeIndex`),
    CONSTRAINT `missingingredient_recipeid_frg` FOREIGN KEY (`RecipeIndex`) REFERENCES `Recipe` (`Index`)
);

CREATE TABLE `UsedIngredient`(
    `Index` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `IngredientId` int UNSIGNED NOT NULL,
    `RecipeIndex` int UNSIGNED NOT NULL,
    `Name` text NOT NULL,
    `Amount` double(8, 2) NOT NULL,
    `Unit` text,
    PRIMARY KEY (`Index`),
    KEY `usedingredient_recipeid_foreign` (`RecipeIndex`),
    CONSTRAINT `usedingredient_recipeid_frg` FOREIGN KEY (`RecipeIndex`) REFERENCES `Recipe` (`Index`)
);

CREATE TABLE `Nutrients`(
    `Index` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `RecipeIndex` int UNSIGNED NOT NULL,
    `Name` text NOT NULL,
    `Amount` double(8, 2) NOT NULL,
    `Unit` text NOT NULL,
    PRIMARY KEY (`Index`),
    KEY `nutrients_recipeid_foreign` (`RecipeIndex`),
    CONSTRAINT `nutrients_recipeid_frg` FOREIGN KEY (`RecipeIndex`) REFERENCES `Recipe` (`Index`)
);