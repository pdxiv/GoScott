# Limitations of the Scott Adams TRS-80 engine

## Introduction

The Scott Adams engine is a very flexible, compact and efficient system for running text adventure games, originally made in 1978 with TRS-80 Level II BASIC for the Tandy/Radio Shack TRS-80 Model I computer.

One of the most important innovations, was probably the separation of game engine and game data. This made it possible for games to be easily ported. Only the game engine needed to be rewritten to run the games on a different computer platform. Eventually, this meant that games were ported to at least 18 different computer platforms, during the lifetime of Scott Adams' company, Adventure International.

Another innovation, was representing an adventure game as a virtual machine, which made the games small and efficient.

With the obiquitous presence of the game engine on multiple platforms, came the unavoidable consequence of keeping the game data format the same, to maintain compatibility and portability. This meant that some of the original design decisions from the 1978 TRS-80 BASIC version seemed to have continued to limit the possibilities for new games, even when technology moved forward.

For the scope of this article, this has the unfortunate consequence of limiting the size and complexity of the game world 

## Data file format

### Breakdown of limitations

- 100 text messages
- 150 verbs
- 150 nouns

### Possible improvement

Store data in JSON format instead.

## Engine

### Definition of "hard-coded behavior"

Some things happen in the game without being specified in the game data files. Examples of this are the GET and GO verbs, as well as light counter decrementing automatically when the light source (item 9) is being carried.

Other things depend on words being present in particular positions in the verb/noun lists, or are hard-coded to use specific flags or counters.

### Breakdown of hard-coded behavior

- GO is not implemented as a word action. references fixed nouns 1-6. Referenced by fixed verb.
- GET is not implemented as a word action. Referenced by fixed verb.
- Light references hard-coded counters and flags. Flag 15 and flag 16. Light source is hard-coded as object number 9. Time limit/light is hard coded counter 8.
- SCORE has some hard-coded behavior, related to objects being named in a certain way.

### Possible improvements

There are some enhancements that can be made to the "stock" engine to make it more flexible, and potentially handle games from other engines (like Quill), while still maintaining much of the original philosophy

- Turn "hard-coded" logic for "GET/DROP" and "GO" in the original engine into actions, with new underlying condition codes.
  - For "GET/DROP": Create a condition which takes a room number as an argument for resolving the noun word input to an object id, for an object, if it matches the word, and is in the supplied room id. The object id is stored in the "object" register, with object id of -1. If object not found, put -1 in the register.
  - For "GO": Create a condition which takes a noun and references it to the direction store for a room. If it matches, place the room id for the next room in the "-3" register. Otherwise, place "-1" in the "-3" register.
- Introduce the "-2" room register, which refers to the current room for the player.
- Turn remaining hard-coded events in the original engine into actions (light source, "SCORE", "DIE", "DSPRM" (etc?)).
  - This can possibly be done by introducing some sort of "GOSUB" command, which calls a numbered action entry.
  - Introduce "not a move" command code, which stops processing of the next set of auto actions, for creating informative error messages.
- Implement "SCORE" and treasure functionality by having an action dedicated for each treasure object, to check if the object is in the "treasure room". If the treasure object is located in the treasure room, increment counter 9. To show the score, commands for doing multiplication and division on the counters must be available. The correct way to calculate the score is 100\*counter/number_of_treasures. The new commands are "CT*n" (counter_multiply) and "CT/n" (counter_divide).

## Text parser

- Missing support for prepositions. The text parser uses a two word system for verb and noun. In some situations, a preposition is needed to adequately describe what must happen. For example: "Kill Dragon with Sword".
- Missing support for spaces in verb phrases and noun phrases in the dictionary.