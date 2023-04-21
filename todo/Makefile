# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build

# Main package
MAIN_PATH=./cmd/todo
MAIN_EXEC=./todo

# Flags for the program
ADD_FLAG=-add
COMPLETE_FLAG=-complete
DELETE_FLAG=-delete
LIST_FLAG=-list

# Default Make target
all: build list

# Build the executable
build:
	$(GOBUILD) -o $(MAIN_EXEC) $(MAIN_PATH)

# Run the program with the "add" flag
add:
	@echo "Enter the task you want to add:" && read -r task && ./$(MAIN_EXEC) $(ADD_FLAG) "$$task"

# Run the program with the "complete" flag
complete:
	@echo "Enter the number of the task that you want to complete:" && read -r number && ./$(MAIN_EXEC) $(COMPLETE_FLAG) "$$number"

# Run the program with the "delete" flag
delete:
	@echo "Enter the number of the task that you want to delete:" && read -r number && ./$(MAIN_EXEC) $(DELETE_FLAG) "$$number"

# Run the program with the "list" flag
list:
	./$(MAIN_EXEC) $(LIST_FLAG)

# Clean up
clean:
	rm -f $(MAIN_EXEC)