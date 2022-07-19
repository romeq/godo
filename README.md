# godo
simple TUI todo application with TermUI

## Installation
```sh
go get 
go build main.go
```
Binary should be now created at ./main

## Notes
Application does not include any proper error handling, 
and panics when one occurs. 

it creates a folder to `$XDG_DATA_HOME/godo` (if exists) or `$HOME/.local/godo`. 
This location contains the SQLite database used to save your todos.

## Known bugs
database lock is not handled, and will result in panic

