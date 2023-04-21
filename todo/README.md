# Todo CLI

This is a command-line interface for managing a todo list. It allows you to add, complete, and delete tasks, as well as list all current tasks.

It was created as an introduction to the Go and CLI language included in [this youtube video](https://youtu.be/j1CXoOQXbco).

## Prerequisites
+ Go 1.16 or higher

## Installation
After cloning or downloading the repository folder - build the executable by running makefile command:
```
make build
```

To build the executable and list tasks, run:
```
make all
```

## Usage
+ To add a task, run:
    ```
    make add
    ```
    You will be prompted to enter the task you want to add.

+ To complete a task, run:
    ```
    make complete
    ```
    You will be prompted to enter the number of the task that you want to complete.

+ To delete a task, run:
    ```
    make delete
    ```
    You will be prompted to enter the number of the task that you want to delete.

+ To list all tasks, run:
    ```
    make list
    ```