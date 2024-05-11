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
  on lycos.fr/ under the name `nabokos` until September of 2008, along with
  several sokoban collections.  This author delighted in making a challenging
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
  responsive.  The latter half of these are multi-ban (some many-multi-ban!)
  and will not be included in the Hexaban collection (at least, not until an
  AI-assisted editor is implemented for solo-ban puzzles).

* **morehex.hsb** includes the puzzles which were hosted by DWS in his
  collections but authored by others.  It also contains two bonus puzzles I
  have recently authored.  I will be adding more but not through these text
  files.

## Puzzle-file (.hsb) format

Part of the difficulty interpreting this archive of puzzle collections is
that different authors include different amounts of non-puzzle information.
Added to that, there are peculiarities in the text representation of the
puzzle itself.

The one thing in common is that each file contains some top-level details
(author or collection name) and then each puzzle is in plain-text and
*separated by a double-newline*.  So, after reading the front matter up
to the first `"\n\n"` each puzzle can be partitioned out by looking for
the next double-newline.  Inspect the file to determine what non-puzzle
metadata is provided for that particular file (they are at least internally
consistent, typically).

![hsb-text](https://github.com/SymbolNotFound/hexaban/assets/1689/2eba85a5-ea59-47ff-b5b4-4c480c1776a9)

Each puzzle is then a line-oriented sequence of pairs of bytes: the first
byte represents the type of Tile at that coordinate (a wall, or a floor,
or a goal, or a crate to be moved, or the player).  These are represented
by "#" (hash), " " (space), "." (period), "$" (dollar) or "@" (at).  In the
case that a crate or player character is standing on a goal tile, then
"*" (asterisk) and "+" (plus) are used then.  Since nothing can coincide
with a wall tile and the player cannot stand on a crate, no other special
combinations are needed.

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

The first task, then, is to read in the definition for a single puzzle.  You
should be able to construct your own test case and use that to validate your
program works as expected for single puzzles, then proceed to parse the entire
collection it came from.  Then, having parsed one collection, extend your program
to be able to read all the collections in this directory.

There are many ways you could check your work -- you could do the conversion yourself
by hand first and store that in a program which checks tile-for-tile that they match.
You could opt for something simpler and just report the player's position and which
directions can be moved in, or which crate positions are reachable from the starting
position, and which of those can be moved.  You could remove all the walls and confirm
that they are easily derived from knowing just the floor tile positions.  You could
render it back as text or as symbols on a canvas, if you've gotten that far in your
language explorations, but it's best to keep the tests simple and easily verifiable.

I have included a reference implementation for this problem, it can be found in
/cmd/convert-levels and includes unit tests that verify a subset of each of the
collections.  It is in golang but I've included many comments that should also help.

For a great tutorial on hexagonal grid coordinates, refer to
[RedBlobGame's excellent coverage of the material](https://www.redblobgames.com/grids/hexagons/).

