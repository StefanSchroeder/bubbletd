package main

import (
	"github.com/StefanSchroeder/bubbletd"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	wantQuit := false
	btd := bubbletd.New()
	btd.ReadConfig()
	defer btd.WriteConfig()
	for scanner.Scan() {
		t := scanner.Text()
		switch {
		case t == "quit" || t == "q":
			fmt.Println("Good-bye!")
			wantQuit = true

		// Two argument commands
		case strings.HasPrefix(t, "done"):
			btd.MoveDone(t)
		case strings.HasPrefix(t, "calendar"):
			btd.MoveCalendarlist(t)
		case strings.HasPrefix(t, "queue"):
			btd.MoveQueuelist(t)
		case strings.HasPrefix(t, "quick"):
			btd.MoveQuicklist(t)
		case strings.HasPrefix(t, "add"):
			btd.AddTask(t)
		case strings.HasPrefix(t, "refer"):
			btd.Refer(t)
		case strings.HasPrefix(t, "trash"):
			btd.Trash(t)

		// Three argument commands
		case strings.HasPrefix(t, "delegate"):
			btd.Delegatetask(t)
		case strings.HasPrefix(t, "defer"):
			btd.Defertask(t)
		case strings.HasPrefix(t, "desc"):
			btd.Desc(t)

		// Print commands
		case t == "pl":
			btd.PrintFilter("Delegated")
		case t == "pu":
			btd.PrintFilter("Queue")
		case t == "pq":
			btd.PrintFilter("Quick")
		case t == "pt":
			btd.PrintFilter("Trash")
		case t == "pr":
			btd.PrintFilter("Reference")
		case t == "po":
			btd.PrintFilter("Done")
		case t == "pd":
			btd.PrintFilter("Defer")
		case t == "p":
			btd.PrintTasks()
		case t == "help":
			bubbletd.PrintHelp()

		default:
			fmt.Println("Don't know what to do.")
		}
		if wantQuit {
			break
		}
	}
	/*if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}*/
}


