
# BubbleTD

## A Gettings-Things-Done-Client 

This is a command line program implementing the GTD-method as advocated by David Allen. It allows you to track any tasks over time, freeing you from the mental load of having to remember anything. 

Design goals:

* Simplicity. It shall be possible to enter a new idea in a second without interruption.
* Easy to get back on track after losing focus.
* Few features. Essentials. No manual necessary.
* Keyboard navigation that allows blind entry.

## Functionality

The main screen is the INBOX-page. It has a text field for entering new ideas (COLLECT).
The cursor will by default be in that field. By pressing enter, the new entry is moved to the INBOX. 

The INBOX is a flat, scrollable list of all entered ideas, shown in the center of the INBOX page.

By pressing TAB the focus will move from the idea-entry field to the inbox, allow scrolling, editing
of existing tasks. One task will always be the CURRENT.

At the bottom of the INBOX-page is the "DESCRIPTION". When one task is selected, pressing TAB moves
to the DESCRIPTION field and allows entering additional information about the task.

At any time the program can be exited with ESC. The state of the program will be automatically saved in a Markdown file.
The tasks will be stored as H1 titles. The description will be standard text below the associated title.

In the tasks list, the currently selected entry can be "PROCESSED" by pressing "y" or "n". This answers
the question: "Is this task actionable?"

When pressing "n", a decision needs to be made, what to do with the task. There are three options.
TRASH, REFERENCE, POSTPONE. 

When selecting TRASH, the task will receive the TRASH tag and be removed from the INBOX.
When selecting REFERENCE, the task will receive the REFERENCE tag and be removed from the INBOX.
When selecting POSTPONE, the task will receive the POSTPONED flag and second question is asked: "How long?"
Answer options will be 1d, 1w, 1m, 1y. The task will be hidden from the INBOX until the time has passed.

When pressing "y", a decision needs to be made, what to do with the task. There are three options.
ACTION_QUICKLIST, ACTION_FREELIST, ACTION_CALLIST, ACTION_DELEGATE. 

When selecting ACTION_QUICKLIST, the task will receive the ACTION_QUICKLIST tag and be removed from the INBOX.

When selecting ACTION_FREELIST, the task will receive the ACTION_FREELIST tag and be removed from the INBOX.

When selecting ACTION_CALLIST, the task will receive the ACTION_CALLIST tag and be removed from the INBOX.

When selecting ACTION_DELEGATE, the task will receive the ACTION_TOBEDELEGATED tag and be removed from the INBOX.

### Other pages

The TRASH page shows all tasks with TRASH flag.

The REFERENCE page shows all tasks with TRASH flag.

The POSTPONE page shows all tasks with POSTPONED flag with target date. If the target date is reached, the POSTPONED flag is removed and the target date is removed. The task will show up again in the INBOX for re-evaluation.

The ACTION_QUICKLIST page shows all tasks with ACTION_QUICKLIST. It allows to apply a DONE tag to tasks. DONE tasks are no longer shown.

The ACTION_FREELIST page shows all tasks with ACTION_FREELIST. It allows to apply a DONE tag to tasks. DONE tasks are no longer shown.

The ACTION_CALLIST page shows all tasks with ACTION_CALLIST. It allows to apply a DONE tag to tasks. DONE tasks are no longer shown.

The ACTION_DELEGATE page shows all tasks with ACTION_TOBEDELEGATED. It allows to apply a WAITING tag to tasks. WAITING tasks are no longer shown.

A WAITING task has a target date. Tasks with expired WAITING time will be re-entered in the INBOX.

# Navigation

The main pages are navigated with the F-keys. F1 is the INBOX. FX are the other pages, Trash shall be obscured.

