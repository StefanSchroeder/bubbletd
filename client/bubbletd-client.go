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
			btd.PrintFilter("IsToBeDelegated")
		case t == "pu":
			btd.PrintFilter("IsQueuelist")
		case t == "pq":
			btd.PrintFilter("IsQuicklist")
		case t == "pt":
			btd.PrintFilter("IsTrash")
		case t == "pr":
			btd.PrintFilter("IsReference")
		case t == "po":
			btd.PrintFilter("IsDone")
		case t == "pd":
			btd.PrintFilter("IsDeferred")
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


