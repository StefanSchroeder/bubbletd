// bubbletd is a package implementing a 
// Getting-Things-Done workflow.
// See README.md for details.
// Written by Stefan Schroeder
package bubbletd

import (
	"encoding/json"
	//"reflect"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

//              +-------- Waiting for some time --+ <--+
//              v                                 |    |
// Start -> State=Inbox        +--> State=Later --+    |
//              |              |                       |
//          Actionable --No--> +--> State=Reference    |
//              |              |                       |
//             Yes             +--> State=Trash -->GC  |
//              |                                      |
//              +--------------+--> State=Delegate ----+
//                             |
//                             +--> State=Calendar -> Done
//                             |
//                             +--> State=Queue ----> Done
//                             |
//                             +--> State=Quick ----> Done
type Task struct {
	AddTime         string
	Title           string
	Desc            string
	State		string
	//IsTrash         bool
	//IsReference     bool
	//IsDeferred      bool
	DeferredTime    string
	//IsQuicklist     bool
	//IsQueuelist     bool
	//IsCalendarlist  bool
	//IsToBeDelegated bool
	WaitingTime     string
	//IsDone          bool
}

type Bubbletd []Task

// Three part string
// Command-Name TaskID New Title
func (b Bubbletd) EditTitle(s string) {
	a := strings.SplitN(s, " ", 3)
	idx := b.GetTaskId(a[1])
	if idx == -1 {
		return
	}

	b[idx].Title = a[2]
}

func (b Bubbletd) Desc(s string) {
	a := strings.SplitN(s, " ", 3)
	idx := b.GetTaskId(a[1])
	if idx == -1 {
		return
	}

	b[idx].Desc = a[2]
}

func (b Bubbletd) Delegatetask(s string) {
	a := strings.SplitN(s, " ", 3)
	idx := b.GetTaskId(a[1])

	b[idx].State = "Delegate"
	b[idx].WaitingTime = a[2]
	fmt.Printf("Hibernating until %v\n", b[idx].WaitingTime)
}

func (b Bubbletd) Defertask(s string) {
	a := strings.SplitN(s, " ", 3)
	idx := b.GetTaskId(a[1])
	if idx == -1 {
		return
	}

	b[idx].State = "Later"
	now := time.Now()
	// default is to add one day
	nowplus := now.AddDate(0, 0, 1)
	switch a[2] {
	case "w":
		nowplus = now.AddDate(0, 0, 7)
	case "m":
		nowplus = now.AddDate(0, 1, 0)
	case "y":
		nowplus = now.AddDate(1, 0, 0)
	}
	fmt.Printf("Deferring until %v\n", nowplus)

	b[idx].DeferredTime = nowplus.Format("2006-01-02 15:04:05")
}

func New() *Bubbletd {
	return &Bubbletd{}
}

func (b *Bubbletd) GetDescriptions() []string {
	a := []string{}
	for _, j := range *b {
		a = append(a, j.Desc)
	}
	return a
}

func (b *Bubbletd) GetTitles() []string {
	a := []string{}
	for _, j := range *b {
		a = append(a, j.Title)
	}
	return a
}

func (b *Bubbletd) PrintTasks() {
	fmt.Println("********************")
	for i, j := range *b {
		fmt.Println(i, j)
	}
	fmt.Println("********************")
}

// printFilter generically looks up the boolean of the Tasks-struct and
// will print it only when the field is true.
func (b Bubbletd) PrintFilter(s string) {
	for i, j := range b {
		if j.State == s {
			fmt.Println(i, j)
		}
	}
}

func (b *Bubbletd) AddTask(s string) {
	a := strings.SplitN(s, " ", 2)
	if len(a) < 2 { return }
	title := a[1]
	now := time.Now()
	newtask := Task{
		now.String(), // AddTime
		title,            // Title
		"empty",      // Desc
		"Inbox",      // State
		//false,        // not actionable, IsTrash
		//false,        // not actionable, IsReference
		//false,        // not actionable, IsDeferred
		"none",       // not actionable, DeferredTime
		//false,        // actionable, IsQuicklist
		//false,        // actionable, IsQueuelist
		//false,        // actionable, IsCalendarlist
		//false,        // actionable, IsToBeDelegated
		"none",       // Waiting Time until delegation ends
		//false,        // isDone
	}
	*b = append(*b, newtask)
}

// getField uses reflection of find the boolean struct field
// given as argument.
/*func GetField(v *Task, field string) bool {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Bool()
}*/

func (b Bubbletd) GetTaskId(indexstring string) int {
	idx, err := strconv.Atoi(indexstring)
	if err != nil {
		fmt.Println("Error during conversion")
		return -1
	}
	if idx >= len(b) || idx < 0 {
		fmt.Println("Error. Index out of range")
		return -1
	}
	return idx
}

func (b Bubbletd) Refer(s string) {
	a := strings.SplitN(s, " ", 2)
	idx := b.GetTaskId(a[1])
	if idx == -1 {
		return
	}

	b[idx].State = "Reference"
}

func (b Bubbletd) MoveToBeDelegated(s string) {
	a := strings.SplitN(s, " ", 2)
	idx := b.GetTaskId(a[1])
	if idx == -1 {
		return
	}

	b[idx].State = "Delegate"
}

func (b Bubbletd) MoveDone(s string) {
	a := strings.SplitN(s, " ", 2)
	idx := b.GetTaskId(a[1])
	if idx == -1 {
		return
	}

	b[idx].State = "Done"
}

func (b Bubbletd) MoveCalendarlist(s string) {
	a := strings.SplitN(s, " ", 2)
	idx := b.GetTaskId(a[1])
	if idx == -1 {
		return
	}

	b[idx].State = "Calendar"
}

func (b Bubbletd) MoveQueuelist(s string) {
	a := strings.SplitN(s, " ", 2)
	idx := b.GetTaskId(a[1])
	if idx == -1 {
		return
	}

	b[idx].State = "Queue"
}

func (b Bubbletd) MoveQuicklist(s string) {
	a := strings.SplitN(s, " ", 2)
	idx := b.GetTaskId(a[1])
	if idx == -1 {
		return
	}

	b[idx].State = "Quick"
}

func (b Bubbletd) Trash(s string) {
	a := strings.SplitN(s, " ", 2)
	idx := b.GetTaskId(a[1])
	if idx == -1 {
		return
	}

	b[idx].State = "Trash"
}

func PrintHelp() {
	s := `
q		quit
add text	I am a task: add a task
desc X text	Set the description of task X
done X		Set the done flag of task X

trash X		Set the trash flag of task X
refer X		Set the reference flag of task X

defer X d	Defer the task X one day in the future	
defer X w	Defer the task X one week in the future	
defer X m	Defer the task X one month in the future	
defer X y	Defer the task X one year in the future	

quick X		Set the Quicklist flag to task X
queue X		Set the Queuelist flag to task X
delegate X Y	Delegate task X until Y

pq		print quicklist
pu		print queuelist
pt		print trash 
po		print done tasks
pr		print references
pd		print deferred tasks
pl		print delegated tasks
p		print all

`
	fmt.Printf(s)
}

func (b *Bubbletd) ReadConfig() {
	content, err := ioutil.ReadFile("./test.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, b)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
}

func (b *Bubbletd) Size() string {
	return fmt.Sprintf("Size <%v>", len(*b))
}

func (b *Bubbletd) WriteConfig() {
	fmt.Printf("Writing data...\n")
	file, _ := json.MarshalIndent(b, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0644)
	fmt.Printf("Done writing data.\n")
}

