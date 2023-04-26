DROP TABLE IF EXISTS `MissingIngredient`;
DROP TABLE IF EXISTS `UsedIngredient`;
DROP TABLE IF EXISTS `Nutrients`;
DROP TABLE IF EXISTS `Recipe`;
DROP TABLE IF EXISTS `ArgumentInput`;

CREATE TABLE `ArgumentInput`(
    `Id` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `IngredientText` text NOT NULL,
    `RecipesNumber` int NOT NULL,
    PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `Recipe`(
    `Id` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `RecipeId` int UNSIGNED NOT NULL,
    `ArgumentId` int UNSIGNED NOT NULL,
    `Title` text NOT NULL,
    PRIMARY KEY (`Id`),
    KEY `recipe_argumentid_foreign` (`ArgumentId`),
    CONSTRAINT `recipe_argumentid_foreign` FOREIGN KEY (`ArgumentId`) REFERENCES `ArgumentInput` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `MissingIngredient`(
    `Id` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `IngredientId` int UNSIGNED NOT NULL,
    `RecipeId` int UNSIGNED NOT NULL,
    `Name` text NOT NULL,
    `Amount` double(8, 2) NOT NULL,
    `Unit` text NOT NULL,
    PRIMARY KEY (`Id`),
    KEY `missingingredient_recipeid_foreign` (`RecipeId`),
    CONSTRAINT `missingingredient_recipeid_frg` FOREIGN KEY (`RecipeId`) REFERENCES `Recipe` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `UsedIngredient`(
    `Id` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `IngredientId` int UNSIGNED NOT NULL,
    `RecipeId` int UNSIGNED NOT NULL,
    `Name` text NOT NULL,
    `Amount` double(8, 2) NOT NULL,
    `Unit` text NOT NULL,
    PRIMARY KEY (`Id`),
    KEY `usedingredient_recipeid_foreign` (`RecipeId`),
    CONSTRAINT `usedingredient_recipeid_frg` FOREIGN KEY (`Id`) REFERENCES `Recipe` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `Nutrients`(
    `Id` int UNSIGNED NOT NULL AUTO_INCREMENT,
    `RecipeId` int UNSIGNED NOT NULL,
    `Name` text NOT NULL,
    `Amount` double(8, 2) NOT NULL,
    `Unit` text NOT NULL,
    PRIMARY KEY (`Id`),
    KEY `nutrients_recipeid_foreign` (`RecipeId`),
    CONSTRAINT `nutrients_recipeid_frg` FOREIGN KEY (`RecipeId`) REFERENCES `Recipe` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;