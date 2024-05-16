# Archive of hex-ban puzzles

### (most sources are no longer live and had to be fetched via the Wayback Machine)

> [!CAUTION]
> These files are NOT in a suitable format for directly interacting with the
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
  the deadlock challenges of Sokoban to the hexagonal grid, visible
  enthusiasm on the website.  Both collections are short but all of the
  puzzles are medium-hard challenges.

* **hexocet.hsb** contains puzzles from Aymeric du Peloux who had hosted them
  on lycos.fr/ under the name _nabokos_ until September of 2008, along with
  several Sokoban collections.  This author delighted in making a challenging
  puzzle, judging by the comments and hints provided for these puzzles.

* **all_E_Hex.hsb** has puzzles from Erim SEVER hosted on `erimsever[dot]com` 
  back before its site redesign (the puzzles are no longer found there).  There
  are other sites containing these puzzles but Erim has pointed out that those
  were unauthorized copies, and further raised the question of how to detect
  slightly-adjusted (but isomorphic, differing only in rotation or reflection
  or where the 'ban starting position is). \
  This is also the source of **test_level.hsb** rotations of 'dws002'.

* **lukaszm.hsb** can currently be found on play.fancade.com and have been
  transcribed into this file format to be included with this collection.
  The puzzles vary in size but all are easy-medium challenges.  Many of these
  would make great initial-world puzzles because of their approachability and
  the later ones will require some deeper reasoning.
  I'm curious how many of these would be solvable with a simple A* search, I
  suspect many of them are.

* **svenhex.hsb** has 20 puzzles authored by Sven Egevad of moderate to hard
  difficulty.  These had been hosted on telia.com but that site is no longer
  responsive.  Only six of them are not multi-ban (some are many-multi-ban!)
  and the other 14 will not be included in the Hexaban collection (at least,
  only after an AI-assisted editor is implemented for solo-hexaban puzzles).

* **morehex.hsb** a tiny collection, includes the puzzles which were hosted
  by DWS in his collections but authored by others.  It also contains two
  bonus puzzles by yours truly.  I will be adding more but not through
  these text files.

> [!Note]
> Attribution of authorship and source will be maintained
> through future representations of these puzzles.


## Puzzle-file (.hsb) format

Part of the difficulty of interpreting this archive of puzzle collections is
that different authors included different amounts of non-puzzle information.
Added to that, there are peculiarities in the text representation of the
puzzle itself.  Notably, there is no convention about whether the first line
is an even-aligned or odd-aligned line (this will become apparent later).

The one thing in common is that each file contains some top-level details
(author or collection name) and then each puzzle is in plain-text and
*separated by a double-newline* `"\n\n"`.  So, after reading up to the
first "\n\n" each puzzle can then be partitioned out by looking for
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
case that a crate or player character is standing on a goal tile,
"*" (asterisk) and "+" (plus) are used.  Since nothing can coincide
with a wall tile and the player cannot stand on a crate, no other special
combinations are needed.  The second byte is always a `" "` (space) or newline, 
or the end of the file.  That is, tiles are always spaced out.  Newlines are
always unix-style "\n", all Windows-style "\r\n" have been normalized.

These line-oriented tiles are thus all-even columns or all-odd columns.
This is called [double-height (offset) coordinates] and effectively means that each
row in the puzzle layout takes up two (2) lines in the text layout.  Determining
which coordinate is up-right or down-right depends on whether the current column
is even or odd.  Refer to the image above for clarification.

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
> \* A space " " byte that is outside the walls should not be counted as a floor tile.

## Writing a Converter

I've provided a reference implementation for reading from this text format, transforming
the double-height horizontal layout (row, column) coordinates into axial coordinates that
are centered on the player's initial position.  This has the advantage of being able to
translate and rotate freely without regard for the odd/even column of source and target
coordinates.  Being able to find rotation and reflection symmetries helps immensely when
detecting duplicate puzzle definitions.  I also provide a JSON encoder and decoder for
puzzle definitions and metadata (Author, Title, relative Difficulty if known, etc.).

This can be found in `/cmd/convert-levels/*` but I treated it as throwaway code so it has
what I consider the bare minimum of comments and organization.  Note that this means there
is still a lot of clarity provided in the comments there, and the data structures and
functionality are well organized.  I have a pretty high low-bar for these conventions.  But,
I didn't obsess over its presentation because I plan on removing it from the `main` branch
now that there is no more need for it.  Except, of course, as an illustrative aid for the
beginner programmer looking for a just-challenging-enough problem.  Is that you?

# Suggested problem definitions

These *DO NOT* assume that you have first written a parser for reading from a file
of a hexoban puzzle.  You can read from the JSON representation provided in `/levels/*`
of this repository.  If reading JSON from Python or JavaScript, there are functions in
the standard library for turning a file (or string) of JSON into a dict or list or value.

If you do write a parser, which I strongly recommend you do, you can trim out a
single one from any collection rather than parse an entire collection as a first task.
The parser itself will have some challenging sub-problems, especially w.r.t. handling
odd and even lines consistently, and choice of origin.

Writing the data format parser is a good opportunity for having to make code organization
decisions like what the structural and funcational abstractions are, file organization,
etc., and time spent in pondering the design.  Solving any of these problems would also
involve nearly as much attention to design & organization:

> [!Note]
> These problems have been annotated with their difficulty level
> from "[ * ]" approachable to "[***]" challenging, assuming a few
> weeks-months of programming introduction in a language like Python or JavaScript.

1. \[ * \] **Validate a puzzle definition** \
   Count the number of floor tiles in the puzzle.  Count the number of crates and goals, make sure they match.  Count the number of players, error or break if any of these checks fail.

2. \[ * \] **Legal Moves & Pushes** \
   Find the player position.  For each direction (the six cardinal directions, however you choose to enumerate them)
   * can the player move in this direction?
   * if a crate is in this direction, can the player push it?  Crates can only be pushed one at a time.
   * otherwise, indicate that it is a wall (can't move nor push).

3. \[* *\] **Reachability** \
  For any coordinate (i, j) in your chosen coordinate system, is there a path from the player's current position?  You will need to ignore tiles that were already inspected if you want to save yourself a lot
of redundant searching.

4. \[* *\] **Pushability** \
  From the player's current position, and considering all other positions reachable, list all of the stones
that can be pushed and in which directions they can be pushed.  Evaluate one of the pushes possible to create
a new puzzle state.

5. \[***\] **Terminal Program** \
  Provide the user a REPL: loop until done, print the current puzzle state, accept user input for which
move or push to perform, evaluate that input and continue looping.  Check the updated puzzle state for
the terminal condition (all goals are occupied by a crate).  Note: as in the original sokoban games, the
puzzle doesn't inform you that a crate is deadlocked, so without an undo some pushes are truly game over
but an agent may continue pushing the other blocks around -- deadlock is not a terminal condition.

6. \[* *\] **Valid Sequence** \
  Given a puzzle (terrain and initial conditions) and a sequence of moves, indicate (true/false) whether
the sequence of moves is valid for the provided puzzle.

7. \[***\] **Assignment Problem** \
  Consider each crate in isolation, and whether it can be pushed to each of the goals.  Is there an assignment
of crates to goals such that each goal is covered by a different crate?  You do not have to provide the optimal
solution, but cf. the Hungarian Algorithm for an approachable polynomial-time algorithm for finding an optimal
solution.  As a validity check it is still good to confirm that walls/hallways don't prohibit at least *some*
way of getting each crate to its own goal, ignoring the shuffling of other crates.

8. \[* *\] **JSON Encoding** \
  Using any structure coordinate system of your choice (it doesn't have to match the one here,
and it can be all one file, or file-per-collection), produce a JSON representation of each
puzzle and its metadata.  Python has a json package in its standard library that will help a
lot, you provide it with a dict or list and it does the encoding steps.  Once you have that,
you can see how easy it is to read the same definition out of its JSON representation.


# Why I think this is a good early project

This has enough structural complexity to justify an automated approach and at a scale
that is a little too onerous for humans, but a tiny amount for computation (~150 puzzles
would take weeks to months to convert coordinates by hand, and would likely include several 
errors).  It would take less time to write the automation and run it than it would to do the
entire task by hand.

Parsing a data format is something that is common enough to be practical, tricky enough to
provide a challenge, and yet (in this case) demands no exotic data structures or complicated
math.  It can be solved without relying on another library to do the work (and in this case,
defining a formal grammar and interpreting these puzzles would be more work than hand-rolling
a parser such as I did, so I can encourage directly writing a parser for it).  I did dig up one
example of a Python library for parsing hexobans (along with soko- and tri-) [sokoenginepy],
you could use it as a second reference implementation, but it is scarce on comments.

There are also some lessons that I think doing a larger project can teach you that no amount
of textbook exercises will, because some things emerge from handling the complexity and
underspecification.

## Lesson 1: Programming is Composition

Bonus points if you thought this was a musical reference, I have used that metaphor before.

What I mean, though, is the nesting-boxes kind of composition, and programming is a lot of that.
There are other forms of abstraction, like functional abstraction or file organization, naming
conventions, etc.  The abstraction that involves choosing the properties of an object and the
allocation of resources within an overall structure, that is at the core of program design.

Choose your composition/structure well and it will aid you and others as well.

This project involves a lot of composition -- from the nesting of tiles into columns, into rows,
in puzzles, in collections, to the metadata properties to associate with a puzzle, to the common
parser-related functionality like reading the next byte, reading until double-newline, etc.
Each puzzle collection has similar but different layouts, too, which requires figuring out
which parts are repeated often enough to be redundant (and thus factored out into a function)
and which parts can be duplicated anyway because they still slightly differ at each call site.

## Lesson 2: Programming is Communication

There is a lot of emphasis in introductory programming literature about how you're giving
instructions to the computer.  In effect, communicating your intent to the computer so it
knows what process to carry out.  But the communication that makes even more difference is
done between humans.  And yes, there's communication via the program between programmer and
user, but pay attention also to the communication between programmers reading and modifying
each others' code.  This even applies to the 6-months-from-now future you that will read code
that you've forgotton the details of.

A lot of students don't encounter this lesson until they start their first job.  Most coursework
is done independently so most communication is limited to (instructor -> student -> computer).
Even most internships involve very little attention to this kind of communication because intern
projects are meant to be isolated enough from critical-path project work.
But, designing and organizing your code so that it can be easily used and extended is a valuable
skill.  Documenting and evangelizing the code and its program artifacts is also very useful.

I've tried to provide an example here of the other layers of communication that come with
building a software project: comments, code structure, file organization, documentation.
There are aspects not shown here, such as conversations in Issues
and the unique kind of communication that happens in a code review, but you will have opportunity
to see those, too, I mention them here for you to be aware of it.


## Lesson 3: Programming is Iteration

There is a saying in programming circles that the way to write a program is

1. Make it

2. Make it work

3. Make it work \[better\]*

   > \* the original goes, "make it work faster" but it originates from a time of compute scarcity,
   > Other optimizations such as making it smaller, portable, reliable, etc. are also important.

This acknowledges that most programs, on first try, don't work.  Well, let's say "make it" includes
getting it to compile and perform the "happy flow", maybe a trivial example.  But to work it needs to
perform adequately to all acceptable user input.  And, you need to figure out how to check whether it
"works" for your definition of work.  Established best practice is to write enough tests
(enough is defined by "until fear \[becomes\] boredom" - Kent Beck),
but for most of the history of computers the notion of writing tests was rarely considered! and often inadequately executed!  These days, you could write multiple books on testing practices. 
You don't have to write tests first, but you should write enough tests to be confident
that you've met step 2 above.

The other takeaway is that optimizing your program before it fully works is doing things out of order.
Some people go years into their career before realizing this, especially if they have been in highly
competitive environments where knowing fancy algorithmic tricks for optimizing certain problems fetches
a lot of attention.  But, most of the time you aren't going to need that level of scalability, or at
least you won't until you've 10x your growth a few more times.  And there's considerable overhead in
some of the optimizations -- sometimes resources, often the mental burden of managing and maintaining
any code that's adjacent to the fancy faster-than-needed code.  Because: that code is a form of communication.

Also important in this is that the "it" here can grow in scope as you're building.  And it should.
This wasn't always considered standard practice but the idea of iteratively designing and building
software has not only proven to make sense (you can more dynamically react to what users prefer) it is
far easier to coordinate among coworkers when the scope of discussion is not a 300-page binder.

This project has many opportunities for iteration.  You can create small test puzzles with just walls for building the first limited-functionality version of a parser.  The parsing could be validated for entire
puzzles before attempting any coordinate conversions or metadata parsing.  The suggested projects
could be attempted without parsing entire collections.

> No project is too large when you've mastered composition, communication and iteration.


# See Also

For a great tutorial on hexagonal grid coordinates, refer to
[RedBlobGame's excellent coverage of the material](https://www.redblobgames.com/grids/hexagons/).

[double-height (offset) coordinates]: https://www.redblobgames.com/grids/hexagons/#coordinates-doubled

[sokoenginepy]: https://github.com/tadams42/sokoenginepy
