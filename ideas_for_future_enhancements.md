# Ideas for future enhancements
## Introduction
There are some enhancements that can be made to the "stock" engine to make it more flexible, and potentially handle games from other engines (like Quill), while still maintaining much of the original philosophy
## Scope
- Turn "hard-coded" logic for "GET/DROP" and "GO" in the original engine into actions, with new underlying condition codes.
  - For "GET/DROP": Create a condition which takes a room number as an argument for resolving the noun word input to an object id, for an object, if it matches the word, and is in the supplied room id. The object id is stored in the "object" register, with object id of -1. If object not found, put -1 in the register.
  - For "GO": Create a condition which takes a noun and references it to the direction store for a room. If it matches, place the room id for the next room in the "-3" register. Otherwise, place "-1" in the "-3" register.
- Introduce the "-2" room register, which refers to the current room for the player.
- Turn remaining hard-coded events in the original engine into actions (light source, "SCORE", "DIE", "DSPRM" (etc?)).
  - This can possibly be done by introducing a "GOSUB" command, which calls a numbered action entry.
  - Introduce "not a move" command code, which stops processing of the next set of auto actions, for creating informative error messages.
- Implement "SCORE" and treasure functionality by having an action dedicated for each treasure object, to check if the object is in the "treasure room". If the treasure object is located in the treasure room, increment counter 9. To show the score, commands for doing multiplication and division on the counters must be available. The correct way to calculate the score is 100*counter/number_of_treasures. The new commands are "CT*n" (counter_multiply) and "CT/n" (counter_divide).
