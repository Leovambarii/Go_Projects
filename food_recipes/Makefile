# Makefile for building and testing the recipeFinder program

# build the program for Linux/macOS
build:
	go build -o recipeFinder main.go

# build the program for Windows
build_win:
	go build -o recipeFinder.exe main.go

# run simple test on Linux/macOS
run:
	./recipeFinder --ingredients="tomatoes,eggs,pasta" --numberOfRecipes=1

# run simple test on Windows
run_win:
	recipeFinder.exe --ingredients="tomatoes,eggs,pasta" --numberOfRecipes=1