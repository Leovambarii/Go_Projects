DROP TABLE IF EXISTS `MissingIngredient`;
DROP TABLE IF EXISTS `UsedIngredient`;
DROP TABLE IF EXISTS `Nutrients`;
DROP TABLE IF EXISTS `Recipe`;
DROP TABLE IF EXISTS `ArgumentInput`;

CREATE TABLE `ArgumentInput`(
    `Idx` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `IngredientText` text NOT NULL,
    `RecipesNumber` int NOT NULL,
    PRIMARY KEY (`Idx`)
);

CREATE TABLE `Recipe`(
    `Idx` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `RecipeId` int UNSIGNED NOT NULL,
    `ArgumentIdx` int UNSIGNED NOT NULL,
    `Title` text NOT NULL,
    `Servings` int UNSIGNED NOT NULL,
    PRIMARY KEY (`Idx`),
    KEY `recipe_argumentid_foreign` (`ArgumentIdx`),
    CONSTRAINT `recipe_argumentid_foreign` FOREIGN KEY (`ArgumentIdx`) REFERENCES `ArgumentInput` (`Idx`)
);

CREATE TABLE `MissingIngredient`(
    `Idx` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `IngredientId` int UNSIGNED NOT NULL,
    `RecipeIdx` int UNSIGNED NOT NULL,
    `Name` text NOT NULL,
    `Amount` double(8, 2) NOT NULL,
    `Unit` text,
    PRIMARY KEY (`Idx`),
    KEY `missingingredient_recipeid_foreign` (`RecipeIdx`),
    CONSTRAINT `missingingredient_recipeid_frg` FOREIGN KEY (`RecipeIdx`) REFERENCES `Recipe` (`Idx`)
);

CREATE TABLE `UsedIngredient`(
    `Idx` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `IngredientId` int UNSIGNED NOT NULL,
    `RecipeIdx` int UNSIGNED NOT NULL,
    `Name` text NOT NULL,
    `Amount` double(8, 2) NOT NULL,
    `Unit` text,
    PRIMARY KEY (`Idx`),
    KEY `usedingredient_recipeid_foreign` (`RecipeIdx`),
    CONSTRAINT `usedingredient_recipeid_frg` FOREIGN KEY (`RecipeIdx`) REFERENCES `Recipe` (`Idx`)
);

CREATE TABLE `Nutrients`(
    `Idx` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `RecipeIdx` int UNSIGNED NOT NULL,
    `Name` text NOT NULL,
    `Amount` double(8, 2) NOT NULL,
    `Unit` text NOT NULL,
    PRIMARY KEY (`Idx`),
    KEY `nutrients_recipeid_foreign` (`RecipeIdx`),
    CONSTRAINT `nutrients_recipeid_frg` FOREIGN KEY (`RecipeIdx`) REFERENCES `Recipe` (`Idx`)
);