package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/romeq/godo/database"
	"github.com/romeq/godo/todos"
	"github.com/romeq/godo/utils"
)

func renderScreen(widget ui.Drawable) {
    screen_width, screen_height := ui.TerminalDimensions()
    if screen_width > 120 && screen_height > 35 {
        widget.SetRect(screen_width/4, screen_height/4, screen_width/4*3, screen_height/4*3)
    } else if screen_width > 80 && screen_height > 35 {
        widget.SetRect(8, screen_height / 4, screen_width-8, screen_height / 4*3)
    } else if screen_width < 80 && screen_height > 35 {
        widget.SetRect(4, screen_height/4, screen_width-4, screen_height/4*3)
    } else if screen_width > 80 && screen_height < 35 {
        widget.SetRect(screen_width/10*2, 2, screen_width/10*8, screen_height-2)
    } else {
        widget.SetRect(4, 2, screen_width-4, screen_height-2)
    }

    ui.Render(widget)
}


func main() {
    utils.Check(ui.Init())
    defer ui.Close()

    err := database.Init()
    utils.Check(err)
    
    uicolor := ui.ColorMagenta

    todoListWidget := widgets.NewList()
    todoListWidget.Title = "Loading your To-Dos..."
    todoListWidget.Rows = []string{}
    todoListWidget.SelectedRowStyle = ui.NewStyle(uicolor)
    todoListWidget.TitleStyle = ui.NewStyle(uicolor)
    todoListWidget.BorderStyle = ui.NewStyle(ui.ColorBlack)
    renderScreen(todoListWidget)

    screen_width, _ := ui.TerminalDimensions()
    input_prefix := " > "
    inputWidget := widgets.NewParagraph()
    inputWidget.SetRect(2, 1, screen_width-2, 4)
    inputWidget.Text = input_prefix
    inputWidget.Title = "Add a To-Do"
    inputWidget.BorderStyle = todoListWidget.BorderStyle

    // load tasks in background
    var tasks todos.Todos
    go func() {
        db_tasks, err := database.GetTodos()
        utils.Check(err)

        for _, todo := range db_tasks {
            tasks = append(tasks, todo)
            todoListWidget.Rows = append(todoListWidget.Rows, todo.Display())
        }

        todoListWidget.Title = "To-Dos"
        renderScreen(todoListWidget)
    }()

    events := ui.PollEvents()
    inputLoop:for {
        e := <-events
        switch e.ID {
        case "q", "<C-c>", "<C-q>":
            break inputLoop
        case "j", "<Down>":
            if len(todoListWidget.Rows) == 0 { break }
            todoListWidget.ScrollDown()
        case "k", "<Up>":
            if len(todoListWidget.Rows) == 0 { break }
            todoListWidget.ScrollUp() 
        case "<C-d>":
            if len(todoListWidget.Rows) == 0 { break }
            todoListWidget.ScrollHalfPageDown()
        case "<C-u>":
            if len(todoListWidget.Rows) == 0 { break }
            todoListWidget.ScrollHalfPageUp()
        case "<C-f>":
            if len(todoListWidget.Rows) == 0 { break }
            todoListWidget.ScrollPageDown() 
        case "<C-b>":
            if len(todoListWidget.Rows) == 0 { break }
            todoListWidget.ScrollPageUp()
        case "g", "<Home>":
            if len(todoListWidget.Rows) == 0 { break }
            todoListWidget.ScrollTop()
        case "G", "<End>":
            if len(todoListWidget.Rows) == 0 { break }
            todoListWidget.ScrollBottom() 
        case "d":
            if len(todoListWidget.Rows) == 0 { break }
            // Get current task from tasks array and delete todo from database
            currentTask := &tasks[todoListWidget.SelectedRow]
            utils.Check(database.RemoveTodoById(currentTask.ID))

            // Remove task from both lists
            listRows := todoListWidget.Rows
            selectedRow := todoListWidget.SelectedRow
            tasks = append(tasks[:selectedRow], tasks[selectedRow+1:]...)
            todoListWidget.Rows = append(listRows[:selectedRow], listRows[selectedRow+1:]...)

            // Scroll up if deleted todo was not last one 
            if len(todoListWidget.Rows) == selectedRow { 
                todoListWidget.ScrollUp() 
            }

        case "<Space>": // Toggle done
            if len(todoListWidget.Rows) == 0 { break }

            // Get current task and toggle it's done property
            currentTask := &tasks[todoListWidget.SelectedRow]
            currentTask.Done = !currentTask.Done

            // Update the todo in database and screen
            todoListWidget.Rows[todoListWidget.SelectedRow] = currentTask.Display()
            utils.Check(database.UpdateDoneById(currentTask.ID, currentTask.Done))

        case "<Enter>": // Create a new todo
            var input_text string
            renderScreen(inputWidget)
            input_evs := ui.PollEvents()
            newtask_prompt_loop:for {
                ev := <-input_evs
                switch ev.ID {
                case "<Enter>":
                    if len(input_text) > 0 {
                        todo_id, err := database.NewTodo(input_text, 0, 0)
                        utils.Check(err)
                        todo := todos.New(int(todo_id), input_text, 0, 0)
                        tasks = append(tasks, todo)
                        todoListWidget.Rows = append(todoListWidget.Rows, todo.Display())
                        inputWidget.Text = input_prefix
                    }
                    break newtask_prompt_loop
                case "<Backspace>":
                    if len(input_text) > 0 {
                        input_text = input_text[:len(input_text)-1]
                    }
                    inputWidget.Text = input_prefix + input_text
                    renderScreen(inputWidget)
                case "<Escape>", "<C-c>":
                    inputWidget.Text = input_prefix
                    break newtask_prompt_loop
                default:
                    x := ev.ID
                    if x == "<Space>" {
                        x = " "
                    }
                    if len(x) != 1 { break }

                    input_text += x
                    inputWidget.Text = input_prefix + input_text
                    renderScreen(inputWidget)
                }
            }
        case "<Resize>":
            ui.Clear()
        }
        renderScreen(todoListWidget)

    }

}
