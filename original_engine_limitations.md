# Limitations of the Scott Adams TRS-80 engine

## Data file format

### Introduction

The Scott Adams engine is a very flexible, compact and efficient system for running text adventure games, originally made in 1978 with TRS-80 Level II BASIC for the Tandy/Radio Shack TRS-80 Model I computer.

One of the most important innovations, was probably the separation of game engine and game data. This made it possible for games to be easily ported. Only the game engine needed to be rewritten to run the games on a different computer platform. Eventually, this meant that games were ported to at least 18 different computer platforms, during the lifetime of Scott Adams' company, Adventure International.

Another innovation, was representing an adventure game as a virtual machine, which made the games small and efficient.

With the obiquitous presence of the game engine on multiple platforms, came the unavoidable consequence of keeping the game data format the same, to maintain compatibility and portability. This meant that some of the original design decisions from the 1978 TRS-80 BASIC version seemed to have continued to limit the possibilities for new games, even when technology moved forward.

For the scope of this article, this has the unfortunate consequence of limiting the size and complexity of the game world 

### Breakdown of limitations

- 100 text messages
- 150 verbs
- 150 nouns

### Proposal

JSON

## Engine

### Definition of "hard-coded behavior"

Some things happen in the game without being specified in the game data files. Examples of this are the GET and GO verbs, as well as light counter decrementing automatically when the light source (item 9) is being carried.

Other things depend on words being present in particular positions in the verb/noun lists, or are hard-coded to use specific flags or counters.

### Breakdown of hard-coded behavior

- GO is not implemented as a word action. references fixed nouns 1-6. Referenced by fixed verb.
- GET is not implemented as a word action. Referenced by fixed verb.
- Light references hard-coded counters and flags. Flag 15 and flag 16. Light source is hard-coded as object number 9. Time limit/light is hard coded counter 8.