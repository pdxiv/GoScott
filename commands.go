package main

// Skeleton of declaration for command logic
// Missing a lot of important stuff, still :)

type commandFunction func(int)

func initCommands() []commandFunction {
	command := []commandFunction{
		//  0 GETx
		func(actionId int) {

		},
		//  1 DROPx
		func(actionId int) {

		},
		//  2 GOTOy
		func(actionId int) {

		},
		//  3 x->RM0
		func(actionId int) {

		},
		//  4 NIGHT
		func(actionId int) {

		},
		//  5 DAY
		func(actionId int) {

		},
		//  6 SETz
		func(actionId int) {

		},
		//  7 x->RM0
		func(actionId int) {

		},
		//  8 CLRz
		func(actionId int) {

		},
		//  9 DEAD
		func(actionId int) {

		},
		// 10 x->y
		func(actionId int) {

		},
		// 11 FINI
		func(actionId int) {

		},
		// 12 DspRM
		func(actionId int) {

		},
		// 13 SCORE
		func(actionId int) {

		},
		// 14 INV
		func(actionId int) {

		},
		// 15 SET0
		func(actionId int) {

		},
		// 16 CLR0
		func(actionId int) {

		},
		// 17 FILL
		func(actionId int) {

		},
		// 18 CLS
		func(actionId int) {

		},
		// 19 SAVE
		func(actionId int) {

		},
		// 20 EXx,x
		func(actionId int) {

		},
		// 21 CONT
		func(actionId int) {

		},
		// 22 AGETx
		func(actionId int) {

		},
		// 23 BYx<-x
		func(actionId int) {

		},
		// 24 DspRM
		func(actionId int) {

		},
		// Newer post-1978 commands below
		// 25 CT-1
		func(actionId int) {

		},
		// 26 DspCT
		func(actionId int) {

		},
		// 27 CT<-n
		func(actionId int) {

		},
		// 28 EXRM0
		func(actionId int) {

		},
		// 29 EXm,CT
		func(actionId int) {

		},
		// 30 CT+n
		func(actionId int) {

		},
		// 31 CT-n
		func(actionId int) {

		},
		// 32 SAYw
		func(actionId int) {

		},
		// 33 SAYwCR
		func(actionId int) {

		},
		// 34 SAYCR
		func(actionId int) {

		},
		// 35 EXc,CR
		func(actionId int) {

		},
		// 36 DELAY
		func(actionId int) {

		},
	}
	command[2](1)
	return command
}

// Working Perl code below, for reference - for now
/*
# Code for all the action commands
my @command_function = (

    #  0 GETx
    sub {
        my $action_id = shift;
        $carried_objects = 0;

        foreach my $location (@object_location) {
            if ( $location == $ROOM_INVENTORY ) {
                $carried_objects++;
            }
        }
        if ( $carried_objects >= $max_objects_carried ) {
            print "I've too much too carry. try -take inventory-\n" or croak;

            # Stop processing later commands if this one fails
            my $continue = shift;
            ${$continue} = $FALSE;
        }

        get_command_parameter($action_id);
        $object_location[$command_parameter] = $ROOM_INVENTORY;
    },

    #  1 DROPx
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        $object_location[$command_parameter] = $current_room;
    },

    #  2 GOTOy
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        $current_room = $command_parameter;
    },

    #  3 x->RM0
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        $object_location[$command_parameter] = 0;
    },

    #  4 NIGHT
    sub {
        my $action_id = shift;
        $status_flag[$FLAG_NIGHT] = $TRUE;
    },

    #  5 DAY
    sub {
        my $action_id = shift;
        $status_flag[$FLAG_NIGHT] = $FALSE;
    },

    #  6 SETz
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        $status_flag[$command_parameter] = 1;
    },

    #  7 x->RM0
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        $object_location[$command_parameter] = 0;
    },

    #  8 CLRz
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        $status_flag[$command_parameter] = 0;
    },

    #  9 DEAD
    sub {
        my $action_id = shift;
        print "I'm dead...\n" or croak;
        $current_room = $number_of_rooms;
        $status_flag[$FLAG_NIGHT] = $FALSE;
        show_room_description();
    },

    # 10 x->y
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        my $temporary_1 = $command_parameter;
        get_command_parameter($action_id);
        $object_location[$temporary_1] = $command_parameter;
    },

    # 11 FINI
    sub {
        my $action_id = shift;
        exit 0;
    },

    # 12 DspRM
    sub {
        my $action_id = shift;
        show_room_description();
    },

    # 13 SCORE
    sub {
        my $action_id = shift;
        $stored_treasures = 0;
        {
            my $object = 0;
            foreach my $location (@object_location) {
                if ( $location == $treasure_room_id ) {
                    if ( substr( $object_description[$object], 0, 1 ) eq q{*} )
                    {
                        $stored_treasures++;
                    }
                }
                $object++;
            }
        }

        print "I've stored $stored_treasures treasures. "
          . "ON A SCALE OF 0 TO $PERCENT_UNITS THAT RATES A "
          . int( $stored_treasures / $number_of_treasures * $PERCENT_UNITS )
          . "\n"
          or croak;
        if ( $stored_treasures == $number_of_treasures ) {
            print "Well done.\n" or croak;
            exit 0;
        }
    },

    # 14 INV
    sub {
        my $action_id = shift;
        print "I'm carrying:\n" or croak;
        my $carrying_nothing_text = 'Nothing';
        my $object_text;
        {
            my $object = 0;
            foreach my $location (@object_location) {
                if ( $location != $ROOM_INVENTORY ) {
                    $object++;
                    next;
                }
                else {
                    $object_text = strip_noun_from_object_description($object);
                }
                print "$object_text. " or croak;
                $carrying_nothing_text = q{};
                $object++;
            }
        }
        print "$carrying_nothing_text\n\n" or croak;
    },

    # 15 SET0
    sub {
        my $action_id = shift;
        $command_parameter = 0;
        $status_flag[$command_parameter] = 1;

    },

    # 16 CLR0
    sub {
        my $action_id = shift;
        $command_parameter = 0;
        $status_flag[$command_parameter] = 0;
    },

    # 17 FILL
    sub {
        my $action_id = shift;
        $alternate_counter[$COUNTER_TIME_LIMIT] = $time_limit;
        $object_location[$LIGHT_SOURCE_ID]      = $ROOM_INVENTORY;
        $status_flag[$FLAG_LAMP_EMPTY]          = $FALSE;
    },

    # 18 CLS
    sub {
        my $action_id = shift;
        cls();

    },

    # 19 SAVE
    sub {
        my $action_id = shift;
        save_game();
    },

    # 20 EXx,x
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        my $temporary_1 = $command_parameter;
        get_command_parameter($action_id);
        my $temporary_2 = $object_location[$command_parameter];
        $object_location[$command_parameter] = $object_location[$temporary_1];
        $object_location[$temporary_1]       = $temporary_2;
    },

    # 21 CONT
    sub {
        $cont_flag = 1;
    },

    # 22 AGETx
    sub {
        my $action_id = shift;
        $carried_objects = 0;
        get_command_parameter($action_id);
        $object_location[$command_parameter] = $ROOM_INVENTORY;
    },

    # 23 BYx<-x
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        my $first_object_location = $object_location[$command_parameter];
        get_command_parameter($action_id);
        $object_location[$command_parameter] = $first_object_location;
    },

    # 24 DspRM
    sub {
        my $action_id = shift;
        show_room_description();
    },

    # Newer post-1978 commands below

    # 25 CT-1
    sub {
        my $action_id = shift;
        $counter_register--;
    },

    # 26 DspCT
    sub {
        my $action_id = shift;
        print "$counter_register" or croak;

    },

    # 27 CT<-n
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        $counter_register = $command_parameter;
    },

    # 28 EXRM0
    sub {
        my $action_id = shift;
        my $temp      = $current_room;
        $current_room = $alternate_room[0];
        $alternate_room[0] = $temp;
    },

    # 29 EXm,CT
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        my $temp = $counter_register;
        $counter_register = $alternate_counter[$command_parameter];
        $alternate_counter[$command_parameter] = $temp;
    },

    # 30 CT+n
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        $counter_register += $command_parameter;
    },

    # 31 CT-n
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        $counter_register -= $command_parameter;

        # According to ScottFree source, the counter has a minimum value of -1
        if ( $counter_register < $MINIMUM_COUNTER_VALUE ) {
            $counter_register = $MINIMUM_COUNTER_VALUE;
        }
    },

    # 32 SAYw
    sub {
        my $action_id = shift;
        print $global_noun or croak;
    },

    # 33 SAYwCR
    sub {
        my $action_id = shift;
        print "$global_noun\n" or croak;
    },

    # 34 SAYCR
    sub {
        my $action_id = shift;
        print "\n" or croak;
    },

    # 35 EXc,CR
    sub {
        my $action_id = shift;
        get_command_parameter($action_id);
        my $temp = $current_room;
        $current_room = $alternate_room[$command_parameter];
        $alternate_room[$command_parameter] = $temp;
    },

    # 36 DELAY
    sub {
        my $action_id = shift;
        sleep 1;
    },

);
*/
