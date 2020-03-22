package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//Task represents a task by name and length (minutes)
type Task struct {
	name    string
	minutes int
}

//toString for Task
func (t Task) String() string {
	return fmt.Sprintf("%s for %d minutes", t.name, t.minutes)
}

func main() {

	var selection int

	taskList := make([]Task, 0)

	printHomeScreen() //prompt inital selection screen

	for {
		_, err := fmt.Scan(&selection)
		if err == nil && selection < 5 && selection > 0 {

			switch selection {
			case 1:
				taskList = createTask(taskList)
				printHomeScreen("Created Task: " + taskList[len(taskList)-1].String()) // -> calls String()
			case 2:
				browseTasks(taskList)

			}

		} else {
			printHomeScreen("User Input Error: Please select an appropriate number")
		}

	}

}

func printHomeScreen(errs ...string) {
	//TODO

	clearTerm()

	// print out each error
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err)
			fmt.Println()
		}
	}

	fmt.Println("Make Selection:")
	fmt.Println("1. Create new task")
	fmt.Println("2. Browse tasks")
	fmt.Println("3. Start random task")
	fmt.Println("4. Quit")

}

//clears terminal linux -> needs support for all OS
func clearTerm() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func createTask(taskList []Task) []Task {

	var mins int
	clearTerm()
	fmt.Println("Enter task name:")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSuffix(name, "\n")

	//loop until user input is correct
	for {
		fmt.Println("Enter time to completion (minutes):")
		_, err := fmt.Scan(&mins)

		if err == nil && mins > 0 {
			break
		}

		clearTerm()
		fmt.Println("User Input Error:")
		fmt.Println(mins)
	}
	return append(taskList, Task{name, mins})
}

func browseTasks(taskList []Task) {

	var selection int

	if len(taskList) == 0 {
		printHomeScreen("There are no tasks queued create one first.")
		return
	}
	clearTerm()

	// task selection menu
	for {
		//print tasks
		for i, err := range taskList {
			fmt.Println(i+1, ".", err)
		}
		fmt.Println(len(taskList)+1, ". Cancel")

		//get selection
		_, err := fmt.Scan(&selection)

		//check if user cancels selection
		if err == nil && selection == len(taskList)+1 {
			return
		}
		if err == nil && selection < len(taskList) && selection > 0 {
			// startTask(taskList, selection-1)
			return
		} else {
			clearTerm()
			fmt.Println("User Input Error: Enter number corresponding to task")
		}

	}

}
