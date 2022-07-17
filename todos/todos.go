package todos

import (
	"fmt"
)

type Todo struct {
    ID int
    Task string
    Done bool
    Deadline int
    Date int
}


func New(id int, task string, deadline int, date int) Todo {
    return Todo {
        ID: id,
        Task: task,
        Done: false,
        Deadline: deadline,
        Date: date,
    }
}

func (t *Todo) Display() string {
    checkmark := " "
    if t.Done {
        checkmark = "*"
    }

    return fmt.Sprintf("[%s] %s", checkmark, t.Task)
}

func (t *Todo) Update(new_task string, done bool, deadline int) {
    if new_task != "" {
        t.Task = new_task
    }
    if deadline > 0 {
        t.Deadline = deadline
    }
    t.Done = done
}

type Todos []Todo

func (t *Todos) SortByDone() Todos {
    var sortedArr Todos
    for _, el := range *t {
        if !el.Done {
            sortedArr = append(sortedArr, el)
        }
    }
    for _, el := range *t {
        if el.Done {
            sortedArr = append(sortedArr, el)
        }
    }

    return sortedArr
}

