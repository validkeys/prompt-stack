# Technical Specification: Building Bubbletea Applications in Go

**Source Context:** Based on the architectural patterns and best practices described in "Building Bubbletea Programs."
**Target Audience:** AI Agent coding in Go (using the `bubbletea` package).
**Goal:** Generate strictly structured, event-driven Terminal User Interface (TUI) applications.

---

## 1. Architectural Philosophy: The Elm Architecture

Bubbletea is based on **The Elm Architecture**. The application is a single loop consisting of three distinct parts:

1.  **Model:** The application state.
2.  **Update:** The logic to handle incoming messages (events) and update the state.
3.  **View:** A way to render the state to the terminal.

The AI must strictly separate these three concepts. No logic should exist inside the View, and no rendering should happen inside the Update.

---

## 2. Core Interfaces and Types

You must utilize the `tea` package from `github.com/charmbracelet/bubbletea`.

### The `Model` Interface
Any struct serving as the application state must implement the following interface:

```go
type Model interface {
    // Init returns the initial command to run.
    Init() Cmd

    // Update handles incoming messages and returns the updated model and command.
    Update(Msg) (Model, Cmd)

    // View renders the model to a string.
    View() string
}
```

### `Msg` (Messages)
A message is any data that triggers an update. It is an empty interface:

```go
type Msg interface{}
```

Common built-in messages include:
*   `tea.KeyMsg`: Keypress events.
*   `tea.WindowSizeMsg`: Terminal resize events.
*   `tea.MouseMsg`: Mouse events (if enabled).

### `Cmd` (Commands)
A command is a function that runs asynchronously (outside the main loop) and returns a `Msg` when complete.

```go
type Cmd func() Msg
```

---

## 3. The Application Lifecycle

When initializing a Bubbletea program, follow this exact sequence:

1.  **Initialization:** The `Model` is created. The `Init()` function is called.
2.  **Command Execution:** Any `Cmd` returned by `Init()` is executed (e.g., waiting for user input or a timer).
3.  **Event Loop:**
    *   A `Msg` is received (from a Cmd or system event).
    *   `Update(msg)` is called.
    *   `Update` returns a **new** state (Model) and potentially a new `Cmd`.
    *   `View()` is called with the new state and rendered to the terminal.
4.  **Termination:** The program exits when a `tea.Quit` command is returned.

---

## 4. Implementation Specification

### A. Defining the Model
The model should be a struct containing all necessary state variables (counters, text inputs, loading states).

```go
type myModel struct {
    choices  []string         // Items in a list
    cursor   int              // Which item is selected
    selected map[int]struct{} // Which items are marked
    quitting bool
}
```

### B. The `Init` Function
This function sets up the initial state. It usually returns `nil` (waiting for input) or a command to fetch initial data.

```go
func (m myModel) Init() tea.Cmd {
    // Just return nil, meaning we want to wait for input (no I/O right now)
    return nil
}
```

### C. The `Update` Function
This is the brain of the application. It must be a pure function (same input = same output) regarding logic, though it triggers side effects via Commands.

**Rules for AI:**
*   Always type-assert `msg` to determine the message type.
*   Always return the modified `m` (model) and a `tea.Cmd`.
*   Return `m, tea.Quit` to exit the application.

```go
func (m myModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    
    // Handle Key Presses
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            m.quitting = true
            return m, tea.Quit

        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        }

    // Handle Window Resizing
    case tea.WindowSizeMsg:
        // Usually handled by sub-components, but can update layout logic here
        // if you are doing manual layout calculations.
    }

    return m, nil
}
```

### D. The `View` Function
This function converts the state into a string (using `fmt.Sprintf` or styling libraries like `lipgloss`).

**Rules for AI:**
*   **Do not** update state here.
*   **Do not** perform I/O here.
*   It is common to use `\n` for new lines.
*   If `m.quitting` is true, you might want to print a "Goodbye!" message and then exit logic handles the rest.

```go
func (m myModel) View() string {
    // The header
    s := "What should we buy at the market?\n\n"

    // Iterate over choices
    for i, choice := range m.choices {
        // Is the cursor pointing at this choice?
        cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

        // Is this choice selected?
        checked := " "
        if _, ok := m.selected[i]; ok {
            checked = "x"
        }

        // Render the row
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }

    // The footer
    s += "\nPress q to quit.\n"

    return s
}
```

---

## 5. Running the Program

To execute the Bubbletea loop, use `tea.NewProgram`. This is typically done in the `main` function.

```go
func main() {
    // Initialize the model with default data
    initialModel := myModel{
        choices:  []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
        selected: make(map[int]struct{}),
    }

    // Start the program
    // tea.WithAltScreen() creates a fullscreen UI (cleaner)
    // tea.WithMouseCellMotion() enables mouse support
    p := tea.NewProgram(initialModel, tea.WithAltScreen())

    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
```

---

## 6. Handling Async I/O (Advanced)

When the application needs to perform I/O (like reading a file or an HTTP request), use **Commands**.

1.  **Define a custom message struct:**
    ```go
    type fileReadedMsg struct{ content string }
    type fileReadErrMsg struct{ err error }
    ```
2.  **Create a command function:**
    ```go
    func readFile() tea.Msg {
        content, err := os.ReadFile("myfile.txt")
        if err != nil {
            return fileReadErrMsg{err}
        }
        return fileReadedMsg{string(content)}
    }
    ```
3.  **Return the command from `Init` or `Update`:**
    ```go
    func (m model) Init() tea.Cmd {
        return readFile
    }
    ```
4.  **Handle the result in `Update`:**
    ```go
    case fileReadedMsg:
        m.content = msg.content
        return m, nil
    
    case fileReadErrMsg:
        m.err = msg.err
        return m, nil
    ```

---

## 7. Best Practices Checklist

1.  **Immutability:** Treat the model as immutable. `Update` should return a modified copy of the model or modify the fields of the pointer receiver and return it, ensuring the old state is never used accidentally after the tick.
2.  **String Building:** Use `strings.Builder` or `lipgloss` Join functions for complex views to avoid performance bottlenecks in string concatenation.
3.  **Cleanup:** If you start a goroutine or open a resource in a Command, ensure you handle the teardown (usually via `tea.Quit`).
4.  **Styling:** Use `github.com/charmbracelet/lipgloss` for styling. Do not manually add ANSI color codes in the `View` string.