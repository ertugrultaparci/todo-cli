package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gonuts/flag"
	"github.com/olekukonko/tablewriter"
)

const (
	todoFilename = ".todo"
)

type item struct {
	ID   int
	Name string
	Date string
	Done bool
}

func main() {

	// To save input:

	filename := ""
	existCurTodo := false
	curDir, err := os.Getwd()
	if err == nil {
		filename = filepath.Join(curDir, todoFilename)
		_, err = os.Stat(filename)
		if err == nil {
			existCurTodo = true
		}
	}
	if !existCurTodo {
		home := os.Getenv("HOME")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		filename = filepath.Join(home, todoFilename)
	}

	completedfile := "CompletedItemList.todo"
	todolistfile := "TODOList.todo"

	// Command Part:

	todo := flag.NewFlagSet("todo", flag.ExitOnError)
	helpcmd := todo.Bool("h", false, "Help for todo app.")
	versioncmd := todo.Bool("v", false, "Version for todo app.")
	listcmd := todo.Bool("l", false, "List todo items")
	addcmd := todo.Bool("a", false, "Add an item to todo list")
	completedcmd := todo.Bool("c", false, "List of completed item.")
	markedcmd := todo.Bool("m", false, "Mark an item as completed")
	deletecmd := todo.Bool("d", false, "delete an item")

	if len(os.Args) < 2 {
		fmt.Println("expected a todo subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "todo": // if its the 'todo' command
		Handle(todo, todolistfile, completedfile, helpcmd, versioncmd, listcmd, addcmd, completedcmd, markedcmd, deletecmd)
	default: // if we don't understand the input
	}

}

func Handle(todo *flag.FlagSet, filename string, completedfile string, h *bool, v *bool, l *bool, a *bool, c *bool, m *bool, d *bool) {
	todo.Parse(os.Args[2:])

	if *h {
		fmt.Println("help <command> for more information about a command: ")
		fmt.Println("Commands:")
		fmt.Println("     v       : show version of TODO CLI App")
		fmt.Println("     a       : add new todo item")
		fmt.Println("     l       : list all uncompleted items")
		fmt.Println("     m       : mark an item as completed")
		fmt.Println("     c       : list completed items")
		fmt.Println("     d       : delete item")

	}
	if *v {
		fmt.Println("TODO CLI 1.1 Version released at 14.11.2021")
		os.Exit(1)
	}
	if *l {
		List(todo, filename)
		os.Exit(1)
	}
	if *a {
		Add(filename, os.Args[3:])
		os.Exit(1)
	}

	if *c {
		fmt.Println("This is a list of completed item FUNCTION")
		List(todo, completedfile)
		os.Exit(1)
	}

	if *m {
		fmt.Println("This FUNCTION make an item marked!")
		Complete(filename, os.Args[3:], completedfile)
		os.Exit(1)
	}

	if *d {
		Delete(filename, os.Args[3:])
		fmt.Printf("Task deleted: %s\n", os.Args[3:])
		os.Exit(1)
	}

}

func Add(filename string, args []string) error {
	var i item
	w, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer w.Close()
	task := strings.Join(args, "")
	i.ID = len(ReadTableData(filename)) + 1
	i.Name = "," + task + ","
	i.Date = strings.Trim(time.Now().Format("01-02-2006 Monday"), " ") + ","
	i.Done = false
	_, err = fmt.Fprintln(w, i)
	fmt.Printf("Task added: %s\n", task)
	return err
}

func List(todo *flag.FlagSet, filename string) error {

	data := ReadTableData(filename)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Item", "Date"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render() // Send output

	fmt.Println("There are " + strconv.Itoa(len(data)) + " items in your todo list.")

	return nil
}

func Delete(filename string, id []string) error {
	ids, _ := strconv.Atoi(strings.Join(id, ""))
	var i item
	itemList := ReadItemStructData(filename)

	w, err := os.Create(filename + "_")

	if err != nil {
		return err
	}
	defer w.Close()

	for _, item := range itemList {
		if item.ID != ids {
			i.ID = item.ID
			i.Name = "," + item.Name + ","
			i.Date = strings.Trim(item.Date, " ") + ","
			i.Done = true
			_, err = fmt.Fprintln(w, i)
		}
	}
	w.Close()
	err = os.Remove(filename)
	if err != nil {
		return err
	}
	return os.Rename(filename+"_", filename)
}

func ListMarkedItem(filename string) error {
	data := ReadTableData(filename + "-")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Item", "Date"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render() // Send output

	fmt.Println("There are " + strconv.Itoa(len(data)) + " items in your todo list.")

	return nil
}

func Complete(filename string, id []string, completedfile string) error {
	ids, _ := strconv.Atoi(strings.Join(id, ""))
	itemList := ReadItemStructData(filename)
	w, err := os.Create(filename + "_")
	if err != nil {
		return err
	}
	defer w.Close()

	c, err := os.OpenFile(completedfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer c.Close()

	for _, element := range itemList {
		var i item
		if element.ID == ids {
			i.ID = element.ID
			i.Name = "," + element.Name + ","
			i.Date = strings.Trim(element.Date, " ") + ","
			i.Done = true
			//_, err = fmt.Fprintln(w, i) // if you want to keep marked item in the general list, just uncomment this line
			_, err = fmt.Fprintln(c, i)
		} else if element.ID != ids {
			i.ID = element.ID
			i.Name = "," + element.Name + ","
			i.Date = strings.Trim(element.Date, " ") + ","
			i.Done = element.Done
			_, err = fmt.Fprintln(w, i)
		}
	}
	w.Close()
	err = os.Remove(filename)
	if err != nil {
		return err
	}
	os.Rename(filename+"_", filename)
	return nil

}

func ReadTableData(filename string) [][]string {

	data := [][]string{}

	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer f.Close()
	br := bufio.NewReader(f)
	n := 1
	for {
		b, _, err := br.ReadLine()
		if err != nil {
			if err != io.EOF {
				return nil
			}
			break
		}
		line := string(b)

		index_id := strings.Index(line, ",")
		id := line[1 : index_id-1]

		rest := line[index_id+1:]

		index_name := strings.Index(rest, ",")
		name := rest[:index_name]

		restForDate := rest[index_name+1:]
		index_date := strings.Index(restForDate, ",")

		date := restForDate[:index_date]

		data = append(data, []string{id, name, date})

		n++

	}
	return data
}

func ReadItemStructData(filename string) []item {

	data := []item{}

	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer f.Close()
	br := bufio.NewReader(f)
	n := 1
	for {
		b, _, err := br.ReadLine()
		if err != nil {
			if err != io.EOF {
				return nil
			}
			break
		}
		line := string(b)
		index_id := strings.Index(line, ",")
		id := line[1 : index_id-1]

		rest := line[index_id+1:]

		index_name := strings.Index(rest, ",")
		name := rest[:index_name]

		restForDate := rest[index_name+1:]
		index_date := strings.Index(restForDate, ",")

		date := restForDate[:index_date]
		if_marked := restForDate[strings.LastIndex(restForDate, ",")+2 : len(restForDate)-1]

		var k item
		k.ID, _ = strconv.Atoi(id)
		k.Name = name
		k.Date = date
		if if_marked == "true" {
			k.Done = true
		} else if if_marked == "false" {
			k.Done = false
		}
		data = append(data, k)

		n++

	}
	return data
}
