# Limitations of the Scott Adams TRS-80 interpreter

## Introduction

The Scott Adams interpreter is a very flexible, compact and efficient system for running text adventure games, originally made in 1978 with TRS-80 Level II BASIC for the Tandy/Radio Shack TRS-80 Model I computer.

One of the most important innovations, was probably the separation of game interpreter and game data. This made it possible for games to be easily ported. Only the game interpreter needed to be rewritten to run the games on a different computer platform. Eventually, this meant that games were ported to at least 18 different computer platforms, during the lifetime of Scott Adams' company, Adventure International.

Another innovation, was representing an adventure game as a virtual machine, which made the games small and efficient.

With the obiquitous presence of the game interpreter on multiple platforms, came the unavoidable consequence of keeping the game data format the same, to maintain compatibility and portability. This meant that some of the original design decisions from the 1978 TRS-80 BASIC version seemed to have continued to limit the possibilities for new games, even when technology moved forward.

For the scope of this article, we will examine the limitations of the original format, and see what can be done to keep the original philosophy of the interpreter as much as possible, without limiting the size and complexity of the game world. While the possibility of running on resource-constrained systems (such as vintage 8-bit computers and microcontrollers) shouldn't be discarded entirely, the focus should be on what is possible with modern computer platforms, in terms of memory availability.

## Data file format

### Breakdown of limitations

- 100 text messages in action commands 2 and 4 (1 and 3 support addressing unlimited number of messages, as illustrated by the game "R" by "therealeasterbunny")
- 150 nouns (because of how verbs and nouns are stored truncated in actions)

### Possible improvement

Store data in JSON format instead.

## Interpreter

### Definition of "hard-coded behavior"

Some things happen in the game without being specified in the game data files. Examples of this are the GET and GO verbs, as well as light counter decrementing automatically when the light source (item 9) is being carried.

Other things depend on words being present in particular positions in the verb/noun lists, or are hard-coded to use specific flags or counters.

### Breakdown of hard-coded behavior

- GO is not implemented as a word action. references fixed nouns 1-6. Referenced by fixed verb position numbers. In later versions of the original interpreter directional noun text was specified in the interpreter code itself (otherwise, games such as #4 Voodoo Castle wouldn't look right).
- GET is not implemented as a word action. Referenced by fixed verb.
- Light references hard-coded counters and flags. Flag 15 and flag 16. Light source is hard-coded as object number 9. Time limit/light is hard coded counter 8.
- SCORE has some hard-coded behavior, related to objects being named in a certain way.

### Possible improvements

There are some enhancements that can be made to the "stock" interpreter to make it more flexible, and potentially handle games from other interpreters (like Quill), while still maintaining much of the original philosophy

- Turn "hard-coded" logic for "GET/DROP" and "GO" in the original interpreter into actions, with new underlying condition codes.
  - For "GET/DROP": Create a condition which takes a room number as an argument for resolving the noun word input to an object id, for an object, if it matches the word, and is in the supplied room id. The object id is stored in the "object" register, with object id of -1. If object not found, put -1 in the register.
  - For "GO": Create a condition which takes a noun and references it to the direction store for a room. If it matches, place the room id for the next room in the "-3" register. Otherwise, place "-1" in the "-3" register.
- Introduce the "-2" room register, which refers to the current room for the player.
- Turn remaining hard-coded events in the original interpreter into actions (light source, "SCORE", "DIE", "DSPRM" (etc?)).
  - This can possibly be done by introducing some sort of "GOSUB" command, which calls a numbered action entry.
  - Introduce "not a move" command code, which stops processing of the next set of auto actions, for creating informative error messages.
- Implement "SCORE" and treasure functionality by having an action dedicated for each treasure object, to check if the object is in the "treasure room". If the treasure object is located in the treasure room, increment counter 9. To show the score, commands for doing multiplication and division on the counters must be available. The correct way to calculate the score is 100\*counter/number_of_treasures. The new commands are "CT*n" (counter_multiply) and "CT/n" (counter_divide).

## Text parser

- Missing support for prepositions. The text parser uses a two word system for verb and noun. In some situations, a preposition is needed to adequately describe what must happen. For example: "Kill Dragon with Sword".
  - One way to solve this, may be to allow the user to write several verb-noun statements on one line. Prepositions could be triggered with a "fall-through" mechanism that keeps executing the next action, when the statement is matched.
- Missing support for spaces in verb phrases and noun phrases in the dictionary.
- Hard-coded "wildcard" on words which have an equal number of characters to the character limit.
  - For maximum flexibility, wildcard behavior could be explicitly specified on a per-word basis instead, so that `SWO*` could be a short form for `SWORD`.
