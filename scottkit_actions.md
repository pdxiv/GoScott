## Conditions
Mnemonic | Parameter | Description
-------- | --------- | -----------
at | room | True if the player's current room is ROOM, which must be the name of a room defined somewhere in the ScottKit file.
carried | item | True if the player is carrying ITEM, which must be the name of an item defined somewhere in the ScottKit file.
here | item | True if ITEM is in the player's current room.
present | item | True if ITEM is either being carried by the player or in the player's current room (i.e. if either carried ITEM or here ITEM is true.)
exists | item | True if ITEM is in the game (i.e. is not "nowhere").
moved | item | True if ITEM has been moved from its original location. 
loaded | n/a | True if the player is carrying at least one item.
flag | num | True if flag number NUM is set.
not_at | room | True if the player's current room is not ROOM, which must be the name of a room defined somewhere in the ScottKit file.
not_carried | item | True if the player is not carrying ITEM, which must be the name of an item defined somewhere in the ScottKit file.
not_here | item | True if ITEM is not in the player's current room.
not_present | item | True if ITEM is either not being carried by the player or not in the player's current room (i.e. if either carried ITEM or here ITEM is true.)
not_exists | item | True if ITEM is not in the game (i.e. is not "nowhere").
not_moved | item | True if ITEM has not been moved from its original location. 
not_loaded | n/a | True if the player is not carrying at least one item.
not_flag | num | True if flag number NUM is not set.
counter_eq | num | True if the current counter's value is NUM. (A different counter may be nominated as "current" by the select_counter action.)
counter_le | num | True if the current counter's value is NUM or less.
counter_gt | num | True if the current counter's value is greater than NUM. Note the asymmetry here: you can check for less-than-or-equal, or strictly-greater-than; but not for strictly-less-than or greater-than-or-equal.

## Commands
Mnemonic | Parameter(s) | Description
-------- | ------------ | -----------
print | string | Prints the specified string. Within that string, \n sequences are interpreted as newlines, and \t sequences as tabs. Since double-quotes are used to enclose the string, they may not appear within it. So backquotes (`) are replaced by double quotes when they are printed.
goto | room | Moves to the specified room and displays its description.
look | | Redisplays the description of the current room, the obvious exits and any visible items. (In a future version, this will be done automatically whenever the player moves (with the goto action), gets an item from the current room, or drops an item. Then it will only need to be done explicitly when changing the value of the darkness flag.)
look2 | | Exactly the same as look, but implemented using a different op-code in the compiled game file. (Why are both of these supported? So that when decompiling a game that uses the latter and then recompiling it, it remains the same.)
get | item | The specified item is put in the player's inventory, unless too many items are already being carried (Cf. the superget action). This works even with items that can't be picked up and dropped otherwise.
superget | item | The specified item is put in the player's inventory, even if too many items are already being carried. This can be used to give the player things he doesn't want, such as the chigger bites in Adventureland.
drop | item | The specified item is put in the player's current location, irrespective of whether it was previous carried, there, elsewhere or nowhere (out of the game). This is the standard way to bring into the game items which begin nowhere.
put | item room | Puts the specified item in the specified room.
put_with | item1 item2 | Puts the first-specified item into the same location as the second.
swap | item1 item2 | Exchanges the two specified items, so that each occupies the location previously occupied by the other. This can be used to switch one object out of the game while bringing another in, as well as for swapping objects that are already in the game.
destroy | item | Removes the specified item from the game, irrespective of whether it was previously carried, in the current location, elsewhere or already out of the game (in which case it's a no-op).
destroy2 | item | Exactly the same as destroy, but implemented using a different op-code in the compiled game file.
inventory | | Lists the items that the player carrying.
score | | Prints the current score, expressed as a mark out of 100, based on how many treasures have been stored in the treasury location. This causes a division-by-zero error if there are no treasures in the game - i.e. items whose descriptions begin with an asterisk (*). So games without treasures, such as Scott Adams's Impossible Mission, should not provide an action with this result.
die | | Implements death by printing an "I am dead" message, clearing the darkness flag and moving to the last defined room, which is conventionally a "limbo" room, as in Adventureland's "Find right exit and live again!" This is not a proper, permanent death: for that, you need the game_over action.
game_over | | Prints "The game is now over", waits five seconds and exits.
print_noun | | Prints the noun that the user just typed.
println_noun | | Prints the noun that the user just typed, followed by a newline.
println | | Emits a newline (i.e. moves to the beginning of the next line).
clear | | Clears the screen. Who could have guessed?
pause | | Waits for two seconds. Useful before clearing the screen.
refill_lamp | | Refills the lightsource object so that it is reset to give light for the initial number of turns, as specified by lighttime.
save_game | | Initiates the save-game diaglogue, allowing the player to save the state of the game to a file. (Unfortunately, there is no corresponding load_game action, so the only way to use a saved game is to restart the interpreter, providing the name of the saved-game file on the command-line.)
set_flag | number | Sets flag number. In general, this is useful only so that subsequent actions and occurrences can check the value of the flag, so there are no pre-defined meanings to the flags. The only flag with a "built-in" meaning is number 15 (darkness).
clear_flag | number | Clears flag number.
set_dark | | Sets flag 15, which indicates darkness. Exactly equivalent to set_flag 15.
clear_dark | | Clears flag 15, which indicates darkness. Exactly equivalent to clear_flag 15.
set_flag0 | | Sets flag 0. Exactly equivalent to set_flag 0.
clear_flag0 | | Clears flag 0. Exactly equivalent to clear_flag 0.
set_counter | number | Sets the value of the currently selected counter to the specified value. Negative values will not be honoured. Do not confuse this with the similarly named select_counter action!
print_counter | | Prints the value of the currently selected counter. Apparently some drivers can't print values greater than 99, so if you're designing your games for maximum portability, you should avoid using numbers higher than this.
dec_counter | | Decreases the value of the currently selected counter by one. The value cannot be decreased below zero. Surprisingly, there is no corresponding increase_counter action, but you can use add_to_counter 1.
add_to_counter | number | Increases the value of the currently selected counter by the specified number.
subtract_from_counter | number | Decreases the value of the currently selected counter by the specified number.
select_counter | number | Chooses which of the sixteen counters is the current one. Subsequent dec_counter, print_counter, etc., actions will operate on the nominated counter. (Initially, counter 0 is used.)
swap_room | | Swaps the player between the current location and a backup location. The backup location is initially undefined, so the first use of this should be immediately followed by a goto to a known room; the next use will bring the player back where it was first used.
swap_specific_room | number | Like swap_room but works with one of a sixteen numbered backup locations, nominated by number. Swaps the current location with backup location number, so that subsequently doing swap_specific_room again with the same argument will result in returning to the original place. This can be used to implement vehicles.
draw number | | Performs a "special action" that is dependent on the driver. For some drivers, it draws a picture specified but the number. In ScottKit (as in ScottFree), this does nothing.
continue | | Never use this action. It is used internally to allow a sequence of actions that is too long to fit into a single action slot, but there is no reason at all why you would ever explicitly use it: in fact, this kind of low-level detail is precisely what ScottKit is supposed to protect you from. I don't know why I'm even mentioning it.
