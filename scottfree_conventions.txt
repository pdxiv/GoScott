// Conditions:
// at          ROOM -- True if the player's current room is ROOM, which must be the name of a room defined somewhere in the ScottKit file.
// carried     ITEM -- True if the player is carrying ITEM, which must be the name of an item defined somewhere in the ScottKit file.
// here        ITEM -- True if ITEM is in the player's current room.
// present     ITEM -- True if ITEM is either being carried by the player or in the player's current room (i.e. if either carried ITEM or here ITEM is true.)
// exists      ITEM -- True if ITEM is in the game (i.e. is not "nowhere").
// moved       ITEM -- True if ITEM has been moved from its original location. 
// loaded      N/A  -- True if the player is carrying at least one item.
// flag        NUM  -- True if flag number NUM is set.
// not_at      ROOM -- True if the player's current room is not ROOM, which must be the name of a room defined somewhere in the ScottKit file.
// not_carried ITEM -- True if the player is not carrying ITEM, which must be the name of an item defined somewhere in the ScottKit file.
// not_here    ITEM -- True if ITEM is not in the player's current room.
// not_present ITEM -- True if ITEM is either not being carried by the player or not in the player's current room (i.e. if either carried ITEM or here ITEM is true.)
// not_exists  ITEM -- True if ITEM is not in the game (i.e. is not "nowhere").
// not_moved   ITEM -- True if ITEM has not been moved from its original location. 
// not_loaded  N/A  -- True if the player is not carrying at least one item.
// not_flag    NUM  -- True if flag number NUM is not set.
// counter_eq  NUM  -- True if the current counter's value is NUM. (A different counter may be nominated as "current" by the select_counter action.)
// counter_le  NUM  -- True if the current counter's value is NUM or less.
// counter_gt  NUM  -- True if the current counter's value is greater than NUM. Note the asymmetry here: you can check for less-than-or-equal, or strictly-greater-than; but not for strictly-less-than or greater-than-or-equal.

// Commands:
