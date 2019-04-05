# (Yet another) Dancing links (DLX) library for Go

## What is a DLX?

DLX stands for "**D**ancing **L**inks E**x**act [cover problem
solver]".  It's a cool algorithm/technique used to solve the _exact
cover problem_ (see below) in a memory-efficient manner.  Donald Knuth
came up with this algorithm, so you know it's good.

### Why are the links dancing (why is it called "dancing links")?

Basically, the implementation of the algorithm ends up involving a
bunch of linked list nodes that get removed and re-inserted over and
over again, so much so that Knuth decided the linked list nodes kind
of looked like they were dancing, or something.

## What is an exact cover problem?

The setup for an exact cover problem is as follows.  You are given a
(finite) set of _items_, along with a collection of some subsets of
the items, called _options_.  You are tasked to select some of the
options so that each item is included in, i.e. _covered by_, exactly
one of your selected options.  (It's worth noting that a given exact
cover problem can have multiple different valid solutions, or none at
all, depending on the available options.)

As it turns out, exact cover problems are NP-complete (or is it
NP-hard?  I don't know--you'll have to ask a real computer scientist
for that).  Basically, that means nobody in the world knows of an
"theoretically efficient" (i.e. polynomial time) algorithm to solve
the exact cover problem.  Instead, our best bet for solving problems
like these is basically by some sort of trial-and-error/backtracking.

The DLX algorithm is just one way to do this trial-and-error in a
memory/time-efficient manner by minimizing the amount of inefficient
overhead associated with copying and moving data around in memory,
etc.

### An example

You walk into a sushi restaurant.  You love eating fish, so you want
to try every single kind of fish that the restaurant has to offer.
The restaurant has the following 7 kinds of fish (the _items_):

1. Hibachi
1. Albacore
1. Salmon
1. Yellowtail
1. Tuna
1. Shrimp (is this a fish?  probably not)
1. Egg (this is most certainly not a fish, but you want to eat it
   anyway)

You have limited stomach space, so in order to save room in your
stomach for each kind of fish, you decide that you can only eat
exactly one piece of each fish.  Unfortunately, the restaurant menu
doesn't offer individual fish pieces.  Instead, they offer the
following sushi combo platters (the _options_):

1. Chef's choice: 1 salmon, 1 tuna
1. Deluxe platter: 1 hibachi, 1 yellowtail, 1 egg
1. Family favorites: 1 albacore, 1 salmon, 1 shrimp
1. Local specials: 1 hibachi, 1 yellowtail, 1 shrimp
1. Catch of the day: 1 albacore, 1 shrimp
1. Party platter: 1 yellowtail, 1 tuna, 1 egg

You want to order menu options so that you get to eat each kind of
fish and that you eat each kind of fish exactly once, and that you
don't order any more than you eat.

(This example is taken, with modification, from Knuth's 2018 Christmas
Lecture on DLX at Stanford.  I believe it is a "textbook" an exact
cover problem example taken from his book, _The Art of Computer
Programming_ (volume who-knows-what).)

#### An example solution

We will try to find a solution to this exact cover problem by a
recursive trial-and-error (i.e. backtracking) method.

To make it easier to see certain patterns, we can represent the
problem setup with a matrix of 0s and 1s.  We assign a column for each
_item_, a row for each _option_.  Each cell contains a 1 if its
associated option (row) covers its associated item (column), and 0
otherwise:

|                  | Hibachi | Albacore | Salmon | Yellowtail | Tuna | Shrimp | Egg |
|------------------|---------|----------|--------|------------|------|--------|-----|
| Chef's choice    | 0       | 0        | 1      | 0          | 1    | 0      | 0   |
| Deluxe platter   | 1       | 0        | 0      | 1          | 0    | 0      | 1   |
| Family favorites | 0       | 1        | 1      | 0          | 0    | 1      | 0   |
| Local specials   | 1       | 0        | 0      | 1          | 0    | 1      | 0   |
| Catch of the day | 0       | 1        | 0      | 0          | 0    | 0      | 1   |
| Party platter    | 0       | 0        | 0      | 1          | 1    | 0      | 1   |

We start by trying to cover _Hibachi_.  Only _Deluxe platter_ and
_Local specials_ cover _Hibachi_, so in order to cover _Hibachi_, we
must select one of _Deluxe platter_ and _Local specials_.  We consider
these choices separately.

1. Suppose we select _Deluxe platter_.  Since _Deluxe platter_ covers
   _Hibachi_, _Yellowtail_, and _Egg_, we cannot select any other
   option that contains _Yellowtail_ or _Egg_, since that would cause
   us to cover these items more than once.  Thus, we remove the
   covered items from our consideration (since they have now already
   been covered), and we also eliminate from our possible selection
   the following remaining options:
   
   - Conflicting on _Hibachi_: _Local specials_;
   - Conflicting on _Yellowtail_: _Local specials_, _Party platter_;
   - Conflicting on _Egg_: _Catch of the day_, _Party platter_.
   
   After doing so, these are the remaining items we have yet to cover:
   
   - _Albacore_,
   - _Salmon_,
   - _Tuna_,
   - _Shrimp_;
   
   and these are the remaining options we have to choose from:
   
   - _Chef's choice_,
   - _Family favorites_.
   
   |                  | Albacore | Salmon | Tuna | Shrimp |
   |------------------|----------|--------|------|--------|
   | Chef's choice    | 0        | 1      | 1    | 0      |
   | Family favorites | 1        | 1      | 0    | 1      |

   (At this point, we can see directly that there is no way to select
   an exact cover from the remaining options.  For sake of
   demonstrating the algorithm's process, we will continue
   systematically anyway until we hit an obvious end condition.)
   
   The next item to cover is _Albacore_.  Only _Family favorites_
   covers _Albacore_, so we are forced to select _Family favorites_.
   
   1. We select _Family favorites_.  Again, we go through and
      eliminate conflicting options again.  Note that _Chef's choice_
      is eliminated since it conflicts with _Family favorites_ on
      _Salmon_.  Therefore, after selecting _Family favorites_, the
      item _Tuna_ remains uncovered, but there are no more options
      left to cover this item.
      
      Thus there are no possible solutions down this decision path, so
      we backtrack to a previous decision (selecting _Deluxe platter_)
      and attempt to find a solution by taking a different decision.
  
2. Suppose we select _Local specials_.

   _Local specials_ covers _Hibachi_, _Yellowtail_, and _Shrimp_, so
   we eliminate the conflicting options:
   
   - Conflicting on _Hibachi_: _Deluxe platter_;
   - Conflicting on _Yellowtail_: _Deluxe platter_, _Party platter_;
   - Conflicting on _Shrimp_: _Family favorites_.
   
   We are left with the following items to cover:
   
   - _Albacore_,
   - _Salmon_,
   - _Tuna_,
   - _Egg_;
   
   and these are the remaining options to select from:
   
   - _Chef's choice_,
   - _Catch of the day_.
   
   |                  | Albacore | Salmon | Tuna | Egg |
   |------------------|----------|--------|------|-----|
   | Chef's choice    | 0        | 1      | 1    | 0   |
   | Catch of the day | 1        | 0      | 0    | 1   |

   (It is straightforward to see now that selecting all remaining
   options gives a valid solution to this exact cover problem.  For
   sake of illustrating the algorithm, we will continue systematically
   as before until we hit a clear end condition.)
   
   The next item to cover is _Albacore_.  Only _Catch of the day_
   covers _Albacore_, so we select it.
   
   1. We select _Catch of the day_, eliminating conflicting options
      (of which there are none).  The remaining items to cover are
      
      - _Salmon_,
      - _Tuna_;
      
      and the remaining options are
      
      - _Chef's choice_.
      
      |               | Salmon | Tuna |
      |---------------|--------|------|
      | Chef's choice | 1      | 1    |

      We cover the next item _Salmon_ by selecting the only available
      option _Chef's choice_.
      
      1. We select _Chef's choice_, leaving us with no remaining items
         left to cover (and also no remaining options to cover).  Thus
         all items have been covered, and we have constructed an exact
         cover by selecting the options
         
         - _Local specials_,
         - _Catch of the day_,
         - _Chef's choice_.
         
##### Applying the dancing links "technique"

### What are some other exact cover problems?

## What is this library?

# Design goals

- Flexible API.

## Comparison with other Go implementations

# How to use this library

