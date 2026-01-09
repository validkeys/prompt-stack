---
description: Systematic debugging guidelines for tests
---

** Be systematic about debugging **

First: State the current problem -- whether it's an error or otherwise.

Then:

1. Isolate a failing test with .only
2. Add comprehensive diagnostic console logging to the entire call stack
3. Run the test
4. Trace the diagnostic logging
5. Determine the root cause.
6. Present the root cause to the user along with a proposed solution. If there are several options, present the options for solving.