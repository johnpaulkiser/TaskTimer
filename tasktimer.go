package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Task represents a task by name and length (minutes)
type Task struct {
	name    string
	minutes int
}

// TaskList represents a list of Task structs
type TaskList []Task

// toString for Task
func (t Task) String() string {
	minOrMins := "minutes"
	if t.minutes == 1 {
		minOrMins = "minute"
	}
	return fmt.Sprintf("%s for %d %s", t.name, t.minutes, minOrMins)
}

// ToCSV returns the csv line representation of a Task -> "chem,45\n"
func (t Task) ToCSV() []string {
	return []string{t.name, strconv.Itoa(t.minutes)}
}

func (tL *TaskList) createTask() {

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
	*tL = append(*tL, Task{name, mins})
}

func (tL *TaskList) browseTasks() {

	var selection int

	if len(*tL) == 0 {
		printHomeScreen("There are no tasks queued create one first.")
		return
	}
	clearTerm()

	for { // task selection menu

		for i, err := range *tL {
			fmt.Printf("%d. %s\n", i+1, err) //print tasks
		}
		fmt.Printf("%d. Cancel\n", len(*tL)+1) //print cancel option

		_, err := fmt.Scan(&selection) //get selection

		if err == nil && selection == len(*tL)+1 { //check if user cancels selection
			clearTerm()
			printHomeScreen()
			return
		}
		if err == nil && selection <= len(*tL) && selection > 0 {
			tL.startTask(selection - 1)
			return
		}
		clearTerm()
		fmt.Println("User Input Error: Enter number corresponding to task")

	}

}

func (tL *TaskList) startTask(selected int) {
	task := (*tL)[selected]
	timer := time.NewTimer(time.Duration(task.minutes) * time.Minute)

	clearTerm()
	fmt.Println("Started task", task)
	<-timer.C
	fmt.Println("Finished Task:", task)
	*tL = tL.removeTask(selected)
	cont()
	clearTerm()
	printHomeScreen()
}

func (tL TaskList) removeTask(toRemove int) TaskList {
	copy(tL[toRemove:], tL[toRemove+1:])
	return tL[:len(tL)-1]
}

func main() {

	var selection int
	filename := "tasks.csv"
	taskList := readTaskFromFile(filename)
	printHomeScreen() //prompt inital selection screen

	for {
		_, err := fmt.Scan(&selection)
		if err == nil && selection < 5 && selection > 0 {

			switch selection {
			case 1:
				taskList.createTask()
				printHomeScreen("Created Task: " + taskList[len(taskList)-1].String()) // -> calls String()
			case 2:
				taskList.browseTasks()
			case 3:
				taskList.startTask(rand.Intn(len(taskList))) //do random task
			case 4:
				clearTerm()
				fmt.Println("Bye!")
				dumpTasksToFile(filename, taskList)
				return
			}

		} else {
			printHomeScreen("User Input Error: Please select an appropriate number")
		}

	}
}

func readTaskFromFile(filename string) TaskList {
	// file format: taskName,numMinutes\n

	list := make(TaskList, 0)
	csvfile, err := os.Open(filename)

	if err != nil {
		fmt.Println("No tasks found. Initializing empty task list.")
		return list
	}

	data := csv.NewReader(csvfile)

	for {
		// Read each record from csv
		record, err := data.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Your tasks.csv file has been corrupted", err)
			break
		}
		mins, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Println("Your tasks.csv file has been corrupted", err)
			break
		}
		list = append(list, Task{record[0], mins})
	}
	return list
}

func dumpTasksToFile(filename string, list TaskList) {
	csvfile, err := os.Create(filename)

	if err != nil {
		fmt.Printf("Could not create file: %s\n", err)
	}

	writer := csv.NewWriter(csvfile)

	for _, task := range list {
		fmt.Println(task.ToCSV())
		err = writer.Write(task.ToCSV())
		if err != nil {
			fmt.Println(err)
		}

	}
	writer.Flush()
	csvfile.Close()

}

func printHomeScreen(errs ...string) {

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

func cont() {
	fmt.Print("\nPress 'Enter' to go to main menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
