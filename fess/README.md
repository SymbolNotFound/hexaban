# FESS in golang for hexoban

Clean-room implementation of FESS to attempt reproducing results of Festival 
based on presentation of [FESS at CoG 2020] given by
[Yaron Shoham](https://festival-solver.site), and related details in the paper.

The work in progress is a solver applied to Hexoban, a hexagonal version of
Sokoban, though it should be possible to apply it to classical Sokoban by
changing the rules for legal inputs (traversable and input-move, pushable and
input-push).  I am designing it with eventual GDL or GGDL compatibility in mind
hence the abstractions around knowledge bases (IDB & EDB) and the planning in
terms of `init`, `base`, `legal`, `input`, `next`, `terminal`, `goal` and
game-specific relations as representation of state and available moves.

## Feature Selection Search

Fess is a new search approach for single-agent search applications, approaching
the state space in terms of higher-order features, thus becoming capable of
searching very large state spaces and very large branching factors.  It is the
first algorithm to solve all levels of the standard Sokoban benchmark.

[2020 Y. Shoham, J. Schaeffer "The FESS Algorithm"]

## Strategy & Tactics

A solution is guided by a path in a Feature Space where features are derived
from properties of the graph over state spaces.  Tactically, the search is able
to shift the value of one or more features.  Feature states are rewarded with
computate time, thus successful outcomes are allocated further search outcomes.

FESS will always find a solution, if one exists -- all states are visited at
least once, but not necessarily in order.  They are never pruned and the state
space is finite, so FESS would eventually visit the entire search space.  It
uses a transposition table to efficiently store positional state and keep the
search-visited representation from blowing up.

Fortunately, Sokoban solution paths are usually composed of subpaths that
themselves look good (crates may need to be shifted partially into a tunnel
but seldom up against a wall) so the search algorithm has an advantage such
that routes which get stuck or loopy are not given sufficient rewards to 
continue dominating the search time, but no matter how low the score even the
rare long unrewarding side path will eventually be visited.

This could be adapted to be an anytime search but the current application is
primarily interested in solving the available Hexoban puzzles within reasonable
resource allocations, and providing a baseline reasoning task that can be
translated into a game description language (albeit level-by-level in GDL, or
by passing the level into every action, state, object and function).

`Advisors` address the large branching factor by suggesting moves that lead to
feature space progress via priority values attached to feature space states.
Good Advisors (heuristics) speed the search by selecting likely good goals at
the strategy planning level.

Despite evaluating fewer states per second than the competing solvers, Festival
solves some puzzles that none of the other solvers were able to find.  The
published results show it is about two orders of magnitude slower!  But it is
able to find solution paths where none others (other than humans, of course!).
However, there are still some puzzles in the large sokoban corpus that are not
yet solved by Festival.

Presumably, this can be addressed by training smarter advisors.  I think perhaps
that is where the bitter lesson will be asserted.  If so, even better to have a
framework that alredy decomposes the problem into a feature-search goal driven
layer and rewards-are-resources allocation management.  Similarly, I wonder if
an Advisor that behaves like a h(x) function (as defined by a classic A* search)
would simplify the search well enough to present a "dumb" AI, the perfect
discriminant for difficulty assignment.  Then the same framework can be used for
different skill levels, even for General Game Playing, with appropriate use of
various Advisor roles.

## Puzzle Credits

I am very curious how this algorithm will do on the
~300 hexagonal variations created by these eight authors:

- David W. Skinner
- Erim SEVER
- François Marques
- Aymeric du Peloux
- LukaszM
- J. Kenneth Riviere
- David Holland
- Gerald Holler

Thank you for designing these puzzles!  and thank you kind reader for reading
this far.  I'll update with my results as soon as they are ready.

If you are an author of one of these puzzles and would like to see more
included, or would rather I not republish them this way, please let me know
(either as a github issue or via the site feedback where you find them
published) and I will update them as soon as I can.  As far as I could
determine these had all been authored over 20 years ago.



## _References_

- [FESS at CoG 2020] (video)

- [2020 Y. Shoham, J. Schaeffer "The FESS Algorithm"] (pdf)


[FESS at CoG 2020]: https://archive.org/details/fess-algorithm

[2020 Y. Shoham, J. Schaeffer "The FESS Algorithm"]: https://ieee-cog.org/2020/papers/paper_44.pdf
