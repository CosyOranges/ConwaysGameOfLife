[![CICD](https://github.com/CosyOranges/ConwaysGameOfLife/actions/workflows/cicd.yml/badge.svg)](https://github.com/CosyOranges/ConwaysGameOfLife/actions/workflows/cicd.yml)
# ConwaysGameOfLife
---

Implementing Conways game of life to get a feel for how Go's Interface feature works.

Probably has no practical use... but it's pretty to watch...

---
## User Interaction

As per the description of Conway's Game of Life, it is a 0 person game where the user has no interaction besides providing the initial state.

You can interact by either stating the size of the grid, the number of generations to run for, the frames per second (fps), the initial starting state (If randomly chosen you can choose the percentage of the random population that will start out as alive).

### Future Improvements
Alternatively you can provide just the number of generations to run for, the frames persecond, **AND** a simple `.txt` file with the initial state representated as a 2D matrix of 1's and 0's.

- This `.txt` MUST have the number of rows and number of columns seperated by a `tab` space on the **FIRST ROW** of the `.txt` file, followed immediately by the grid below.

**e.g. for a Grid with 4 rows and 9 columns**:
```
4   9
0   0   0   0   0   0   0   0   0
0   1   0   1   0   0   1   1   0
0   0   1   1   1   0   0   0   0
0   0   0   0   0   1   0   0   0
```