# Archive of hex-ban puzzles

### (most sources are no longer live and had to be fetched via the Wayback Machine)

> [!CAUTION]
> These files are not in a suitable format for directly interacting with the
> golang and typescript code in this repository.  The puzzles in these text
> files are slightly more convenient for hand-editing, if all you have is a
> text editor, but they are in "double-height" offset coordinates and hexaban
> code will assume an axial coordinate system.  This README only pertains to
> the format of these files, for as long as they exist, for instructional
> purposes.  The coordinate system will be explained below.

## Puzzles in this collection

* **dwshex.hsb** has puzzles by the inventor of Hexoban, David W. Skinner. 
  Previously hosted on bentonrea.com/~sasquatch but no longer up.  These 
  puzzles range in difficulty but are all at least a little challenging.
  Typical puzzle size is 20-40 tiles with some larger ones.  As far as I
  can tell, the file format was also invented by DWS, it is clearly adapted
  from the sokoban text description files which DWS was also familiar with.

* **heloban.hsb** and **heroban.hsb** are both by François Marques who
  really enjoyed making new variations on the puzzles and applying some of
  the deadlock challenges of Sokoban to the hexagonal grid.  These are
  short collections but all of the puzzles are medium-hard challenges.

* **hexocet.hsb** contains puzzles from Aymeric du Peloux who had hosted them
  on lycos.fr/ under the name _nabokos_ until September of 2008, along with
  several Sokoban collections.  This author delighted in making a challenging
  puzzle, judging by the comments and hints provided for these puzzles.

* **all_E_Hex.hsb** has puzzles from Erim SEVER hosted on `erimsever[dot]com` 
  but that site looks like it's been taken over by a casino now.  The very
  old wayback versions of the site point to a lot of these sokoban collections,
  including about 50 of Erim's hex-ban puzzles.

* **lukaszm.hsb** can currently be found on play.fancade.com and have been
  transcribed into this file format to be included with this collection.
  The puzzles vary in size but all are easy-medium challenges.  Many of these
  would make great initial-world puzzles because of their approachability and
  the later ones still make you have to reason some.  I am curious how many
  of these would still be solvable with a simple A* search.

* **svenhex.hsb** has 20 puzzles authored by Sven Egevad of moderate to hard
  difficulty.  These had been hosted on telia.com but that site is no longer
  responsive.  Only six of them are not multi-ban (some are many-multi-ban!)
  and will not be included in the Hexaban collection (at least, not until an
  AI-assisted editor is implemented for solo-hexaban puzzles).

* **morehex.hsb** a tiny collection, includes the puzzles which were hosted
  by DWS in his collections but authored by others.  It also contains two
  bonus puzzles I have recently authored.  I will be adding more but not through
  these text files.

## Puzzle-file (.hsb) format

Part of the difficulty interpreting this archive of puzzle collections is
that different authors include different amounts of non-puzzle information.
Added to that, there are peculiarities in the text representation of the
puzzle itself.  Notably, there is no convention about whether the first line
is an even-aligned or odd-aligned line (this will become apparent later).

The one thing in common is that each file contains some top-level details
(author or collection name) and then each puzzle is in plain-text and
*separated by a double-newline*.  So, after reading the front matter up
to the first `"\n\n"` each puzzle can be partitioned out by looking for
the next double-newline.  Inspect the file to determine what non-puzzle
metadata is provided for that particular file (there is usually consistency
that a property appears for all puzzles, but this is not always true, you
will have to inspect the file to know whether to ignore its absence).

![hsb converter](https://github.com/SymbolNotFound/hexaban/assets/1689/43faf985-6bd2-45db-ab43-b96d1ca0d0f5)

> [!Note]
> This example has its initial row as an odd row.  Some puzzles begin with
> an even row; both need to be handled consistently.

Each puzzle is then a line-oriented sequence of pairs of bytes: the first
byte represents the type of Tile at that coordinate (a wall, or a floor,
or a goal, or a crate to be moved, or the player).  These are represented
by "#" (hash), " " (space), "." (period), "$" (dollar) or "@" (at).  In the
case that a crate or player character is standing on a goal tile, then
"*" (asterisk) and "+" (plus) are used then.  Since nothing can coincide
with a wall tile and the player cannot stand on a crate, no other special
combinations are needed.  The second byte is always a `" "` (space) but may
be a newline or the end of the file.  That is, tiles are always spaced out.

This is called [double-height (offset) coordinates] and effectively means that each
row in the puzzle layout takes up two lines in the text layout (the odd-valued columns
on one line and the even-valued columns on the other).  It 

| byte |  represents  |
| ---- | ------------ |
| `#`  | Wall tile    |
| ` `  | Floor tile*  |
| `.`  | Empty goal   |
| `@`  | Player start |
| `+`  | Player (standing on goal) |
| `$`  | Crate on floor |
| `*`  | Crate on goal |

> [!Note]
> A space " " byte that is outside the walls should not be counted as a floor tile.

## Writing a Converter

I've provided a reference implementation for reading from this text format, transforming
the double-height horizontal layout (row, column) coordinates into axial coordinates that
are centered on the player's initial position.  This has the advantage of being able to
translate and rotate freely without regard for the odd/even column of source and target
coordinates.  Being able to find rotation and reflection symmetries helps immensely when
detecting duplicate puzzle definitions.  I also provide a JSON encoder and decoder for
puzzle definitions and metadata (Author, Title, relative Difficulty if known, etc.).

This can be found in `/cmd/convert/levels/*` but I treated it as throwaway code so it has
what I consider the bare minimum of comments and organization.  Note that this means there
is still a lot of clarity provided in the comments there, and the data structures and
functionality are well organized.  I have a pretty high low-bar for these conventions.  But,
I didn't obsess over its presentation because I plan on removing it from the `main` branch
now that there is no more need for it.  Except, of course, as an illustrative aid for the
beginner programmer looking for a just-challenging-enough problem.  Is that you?

# Suggested problem definitions

These assume that you have first written a parser for reading in the definition
of a hexoban puzzle.  You can trim out a single one from any collection or choose
to parse entire collections before tackling these.  The parser itself will have
some challenging sub-problems, especially w.r.t. handling odd and even lines,
and choice of origin.  It is a good opportunity for having to make code organization
decisions like what the structural and funcational abstractions are, file organization,
etc., and even amount of deliberation.

> [!Note]
> These problems have been annotated with their difficulty level
> from "[ * ]" approachable to "[***]" challenging, assuming
> a few weeks of programming introduction in, e.g. Python or JavaScript.

1. \[ * \] **Validate a puzzle definition**
   Count the number of floor tiles in the puzzle.  Count the number of crates and goals, make sure they match.  Count the number of players, error or break if any of these checks fail.

2. \[ * \] **Legal Moves & Pushes**
   Find the player position.  For each direction (the six cardinal directions, however you choose to enumerate them)
   * can the player move in this direction?
   * if a crate is in this direction, can the player push it?  Crates can only be pushed one at a time.
   * otherwise, indicate that it is a wall (can't move nor push).

3. \[* *\] **Reachability**
  For any coordinate (i, j) in your chosen coordinate system, is there a path from the player's current position?  You will need to ignore tiles that were already inspected if you want to save yourself a lot
of redundant searching.

4. \[* *\] **Pushability**
  From the player's current position, and considering all other positions reachable, list all of the stones
that can be pushed and in which directions they can be pushed.  Evaluate one of the pushes possible to create
a new puzzle state.

5. \[***\] **Terminal Program**
  Provide the user a REPL: loop until done, print the current puzzle state, accept user input for which
move or push to perform, evaluate that input and continue looping.  Check the updated puzzle state for
the terminal condition (all goals are occupied by a crate).  Note: as in the original sokoban games, the
puzzle doesn't inform you that a crate is deadlocked, so without an undo some pushes are truly game over
but an agent may continue pushing the other blocks around -- deadlock is not a terminal condition.

6. \[***\] **Assignment Problem**
  Consider each crate in isolation, and whether it can be pushed to each of the goals.  Is there an assignment
of crates to goals such that each goal is covered by a different crate?  You do not have to provide the optimal
solution, but cf. the Hungarian Algorithm for an approachable polynomial-time algorithm for finding an optimal
solution.  As a validity check it is still good to confirm that walls/hallways don't prohibit at least *some*
way of getting each crate to its own goal.

7. \[* *\] **JSON Encoding**
  Using any structure of your choice (it doesn't have to match mine), and the coordinate system of your
choice, produce a JSON representation of each puzzle and its metadata.  Python has a json package in its
standard library that will help a lot, you provide it with a dict or list and it does the encoding steps.
Once you have that, you can see how easy it is to read the same definition out of its JSON representation.


# Why I think this is a good early project

This has enough structural complexity to justify an automated approach and at a scale
that is a little too onerous for humans and a tiny amount for computation (~150 puzzles
would take months to convert coordinates by hand, and would likely include several errors).

Parsing a data format is something that is common enough to be practical, tricky enough to
provide a challenge, and yet (in this case) demands no exotic data structures or complicated
math.  It can be solved without relying on another library to do the work (and in this case,
defining a formal grammar and interpreting that would be more work than writing a hand-rolled
parser such as I did, so I can encourage directly writing a parser for it).  Indeed, I don't
think there is a Python library for reading hexoban puzzles, because it is relatively rare.

There are also some lessons that I think doing a larger project can teach you that no amount
of textbook exercises will, because some things emerge from handling the complexity.

## Lesson 1: Programming is Composition

... not a music reference, but pun intended

## Lesson 2: Programming is Communication

... obviously communicating with the machine, but importantly communicating with each other

## Lesson 3: Programming is Iteration

... make it, make it work, make it work \[_better_\]


There are many ways you could check your work -- you could do the conversion yourself
by hand first and store that in a program which checks tile-for-tile that they match.
You could opt for something simpler and just report the player's position and which
directions can be moved in, or which crate positions are reachable from the starting
position, and which of those can be moved.  You could remove all the walls and confirm
that they are easily derived from knowing just the floor tile positions.  You could
render it back as text or as symbols on a canvas, if you've gotten that far in your
language explorations, but it's best to keep the tests simple and easily verifiable.


# See Also

For a great tutorial on hexagonal grid coordinates, refer to
[RedBlobGame's excellent coverage of the material](https://www.redblobgames.com/grids/hexagons/).

[double-height (offset) coordinates]: https://www.redblobgames.com/grids/hexagons/#coordinates-doubled
