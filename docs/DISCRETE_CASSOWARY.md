# Discrete Cassowary & Grid Hinting for TUI Auto Layout

## Overview

Standard Cassowary solvers solve linear constraints in continuous floating-point space ($\mathbb{R}$). However, Terminal User Interfaces (TUI) render on a discrete character grid ($\mathbb{Z}^2$).

Naïve rounding of continuous Cassowary solutions (such as `int(math.Round(val))` or truncation) leads to severe layout defects on discrete character grids:
- **Gaps & Overflows**: Fractional column widths like `26.666` for 3 equal columns summing to `80` characters truncate to `26 + 26 + 26 = 78` (a 2-character gap) or round to `27 + 27 + 27 = 81` (a 1-character screen overflow).
- **Edge Misalignment**: Adjacent boxes (`box1.x + box1.width == box2.x`) become disjoint or overlap due to independent rounding.
- **Stem Collapse**: Small elements (borders, icons, status indicators) round to `0` characters and disappear.
- **Asymmetry Distortion**: Equal-priority sibling components get unequal sizes arbitrarily depending on variable evaluation order.

`kiwi-go` solves this by introducing a **Discrete Cassowary Grid Hinting Engine** adapting concepts from font hinting.

---

## Typography Hinting Analogies

### 1. FreeType Autohinting Paradigm (Automatic Heuristics)
FreeType uses automated heuristics to analyze vector glyph outlines and fit them to discrete pixel grids without requiring manual bytecode:
- **Hamilton / Hare-Niemeyer Apportionment**: Distributes integer remainders across layout variables so that $\sum \text{Discrete}(w_i) = \text{Target}$ exactly, while minimizing maximum deviation $|\text{Discrete}(w_i) - \text{Continuous}(w_i)| < 1$.
- **Symmetry Preservation**: Uses distance-from-center tie-breaking so symmetric columns or margins receive balanced integer remainders.
- **Stem Protection / Minimum Size Clamping**: Guarantees that visible UI elements maintain a minimum character dimension (e.g. minimum 1 character width, or 3 characters for bordered boxes).

### 2. TrueType Bytecode Paradigm (Explicit Directives)
TrueType uses explicit stack bytecode instructions attached to control points:
- **Grid Snapping (`SnapToGrid`)**: Snaps coordinates/widths to multiples of $N$ (e.g., 2 characters for double-byte CJK character grids).
- **Group Equalization (`EqualizeGroup`)**: Forces sibling elements to share identical integer dimensions before remainder allocation.
- **Min/Max Clamping (`ClampMinMax`)**: Restricts discrete coordinates to strict integer ranges.
- **Edge Alignment (`AlignEdges`)**: Forces two distinct variables to evaluate to the exact same discrete integer grid line.

---

## Mathematical Formulation

Given continuous variable values $v_1, v_2, \ldots, v_n \in \mathbb{R}$ subject to a sum constraint $\sum v_i = T \in \mathbb{Z}$:

1. **Base Allocation**:
   $$d_i^0 = \max\left(\lfloor v_i \rfloor, \text{minSize}_i\right)$$
2. **Deficit Calculation**:
   $$K = T - \sum_{i=1}^n d_i^0$$
3. **Remainder & Tie-Breaking**:
   $$r_i = v_i - d_i^0$$
   Variables are sorted by $r_i$ descending. Ties are broken symmetrically based on distance from the group midpoint.
4. **Apportionment**:
   The top $K$ variables receive $+1$ character unit:
   $$d_i = \begin{cases} d_i^0 + 1 & \text{if } \text{rank}(r_i) \le K \\ d_i^0 & \text{otherwise} \end{cases}$$

This guarantees $\sum d_i = T$ with zero gaps or overflows.
