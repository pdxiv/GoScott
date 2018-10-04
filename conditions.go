package main

// Skeleton of declaration for condition logic
// Missing a lot of important stuff, still :)
import "fmt"

type conditionFunction func(int) bool

func initConditions() []conditionFunction {
	condition := []conditionFunction{
		//   0 Par
		func(parameter int) bool {
			return parameter == 1
		},
		//   1 HAS
		func(parameter int) bool {
			return parameter == 1
		},
		//   2 IN/W
		func(parameter int) bool {
			return parameter == 1
		},
		//   3 AVL
		func(parameter int) bool {
			return parameter == 1
		},
		//   4 IN
		func(parameter int) bool {
			return parameter == 1
		},
		//   5 -IN/W
		func(parameter int) bool {
			return parameter == 1
		},
		//   6 -HAVE
		func(parameter int) bool {
			return parameter == 1
		},
		//   7 -IN
		func(parameter int) bool {
			return parameter == 1
		},
		//   8 BIT
		func(parameter int) bool {
			return parameter == 1
		},
		//   9 -BIT
		func(parameter int) bool {
			return parameter == 1
		},
		//  10 ANY
		func(parameter int) bool {
			return parameter == 1
		},
		//  11 -ANY
		func(parameter int) bool {
			return parameter == 1
		},
		//  12 -AVL
		func(parameter int) bool {
			return parameter == 1
		},
		//  13 -RM0
		func(parameter int) bool {
			return parameter == 1
		},
		//  14 RM0
		func(parameter int) bool {
			return parameter == 1
		},
		// Newer post-1978 conditions below
		//  15 CT<=
		func(parameter int) bool {
			return parameter == 1
		},
		//  16 CT>
		func(parameter int) bool {
			return parameter == 1
		},
		//  17 ORIG
		func(parameter int) bool {
			return parameter == 1
		},
		//  18 -ORIG
		func(parameter int) bool {
			return parameter == 1
		},
		//  19 CT=
		func(parameter int) bool {
			return parameter == 1
		},
	}
	fmt.Printf("\n%v\n", condition[2](1))
	return condition
}

// Working Perl code below, for reference - for now
/*
# Code for all the action conditions
my @condition_function = (

    #  0 Par
    sub {
        my $parameter = shift;
        return $TRUE;
    },

    #  1 HAS
    sub {
        my $parameter = shift;
        my $result    = $FALSE;
        if ( $object_location[$parameter] == $ROOM_INVENTORY ) {
            $result = $TRUE;
        }
        return $result;
    },

    #  2 IN/W
    sub {
        my $parameter = shift;
        return ( $object_location[$parameter] == $current_room );
    },

    #  3 AVL
    sub {
        my $parameter = shift;
        my $result;
        $result = ( $object_location[$parameter] == $ROOM_INVENTORY );
        $result = $result
          || ( $object_location[$parameter] == $current_room );
        return $result;
    },

    #  4 IN
    sub {
        my $parameter = shift;
        return ( $current_room == $parameter );
    },

    #  5 -IN/W
    sub {
        my $parameter = shift;
        return ( $object_location[$parameter] != $current_room );
    },

    #  6 -HAVE
    sub {
        my $parameter = shift;
        my $result    = $FALSE;
        if ( $object_location[$parameter] != $ROOM_INVENTORY ) {
            $result = $TRUE;
        }
        return $result;
    },

    #  7 -IN
    sub {
        my $parameter = shift;
        return ( $current_room != $parameter );
    },

    #  8 BIT
    sub {
        my $parameter = shift;
        my $result;

        $result = $status_flag[$parameter];
        return $result;
    },

    #  9 -BIT
    sub {
        my $parameter = shift;
        my $result;

        $result = ( !$status_flag[$parameter] );
        return $result;
    },

    # 10 ANY
    sub {
        my $parameter = shift;
        my $result    = $FALSE;
        foreach my $location (@object_location) {
            if ( $location == $ROOM_INVENTORY ) {
                $result = $TRUE;
            }
        }
        return $result;
    },

    # 11 -ANY
    sub {
        my $parameter = shift;
        my $result    = $FALSE;
        foreach my $location (@object_location) {
            if ( $location == $ROOM_INVENTORY ) {
                $result = $TRUE;
            }
        }
        return ( !$result );
    },

    # 12 -AVL
    sub {
        my $parameter = shift;
        my $result;
        $result = ( $object_location[$parameter] == $ROOM_INVENTORY );
        $result = $result
          || ( $object_location[$parameter] == $current_room );
        return ( !$result );
    },

    # 13 -RM0
    sub {
        my $parameter = shift;
        return ( $object_location[$parameter] != $ROOM_STORE );
    },

    # 14 RM0
    sub {
        my $parameter = shift;
        return ( $object_location[$parameter] == $ROOM_STORE );
    },

    # Newer post-1978 conditions below

    # 15 CT<=
    sub {
        my $parameter = shift;
        return $counter_register <= $parameter;
    },

    # 16 CT>
    sub {
        my $parameter = shift;
        return $counter_register > $parameter;
    },

    # 17 ORIG
    sub {
        my $parameter = shift;
        return $object_original_location[$parameter] ==
          $object_location[$parameter];
    },

    # 18 -ORIG
    sub {
        my $parameter = shift;
        return !( $object_original_location[$parameter] ==
            $object_location[$parameter] );
    },

    # 19 CT=
    sub {
        my $parameter = shift;
        return $counter_register == $parameter;
    },
);
*/
