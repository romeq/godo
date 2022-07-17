# godo
simpel TUI todo application with TermUI

## installation
```sh
go get 
go build main.go
```
the binary should be now created at ./main

## notes
application does not include any proper error handling, 
and panics when one occurs. 

it creates a folder to $HOME/.local/godo which contains the
database used to save todos.

## known bugs
database lock is not handled, and will result in panic of this application

