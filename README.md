# GoScott
## Introduction
Project to make a Scott Adams advanture game interpreter in Go.
## Current status
Non-working, as a whole. Currently, loading the data in files into local data structures is working, but the work to create new custom actions for prepending "GO" and appending "GET" and "DROP" are not implemented in the loading routine.
## Scope
- Loading game data from files into internal data structures
- Turn all hard-coded events in the original engine into actions (light source, "GET", "DROP", "GO" (etc?))
- Saving of internal data structures back to original data files
- Save/Load functionality for saving all steps taken in the game, including the random number generator seed
- Be able to run games
