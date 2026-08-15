package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/unxed/kiwi-go"
	"golang.org/x/term"
)

func main() {
	fmt.Println("==========================================================================")
	fmt.Println("             KIWI-GO: DISCRETE CASSOWARY & TUI HINTING DEMO               ")
	fmt.Println("==========================================================================")
	fmt.Println()

	runContinuousDemo()
	runNaïveVsAutohintDemo()
	runTrueTypeRulesDemo()
	runTUIWindowLayoutDemo()

	if os.Getenv("KIWI_TESTING") == "1" {
		return
	}

	runInteractiveTUI()
}

func runContinuousDemo() {
	fmt.Println("--- 1. Continuous Cassowary Core Engine ---")

	solver := kiwi.NewSolver()
	left := kiwi.NewVariable("left")
	width := kiwi.NewVariable("width")
	right := kiwi.NewVariable("right")
	centerX := kiwi.NewVariable("centerX")

	// Subscribe to value change
	centerX.Subscribe(func(value, prev float64) {
		fmt.Printf("   [Event] CenterX updated: %.2f -> %.2f\n", prev, value)
	})

	_ = solver.AddEditVariable(left, kiwi.StrengthStrong)
	_ = solver.AddEditVariable(width, kiwi.StrengthStrong)

	_ = solver.SuggestValue(left, 10.5)
	_ = solver.SuggestValue(width, 40.0)

	// right == left + width
	_ = solver.AddConstraint(kiwi.NewConstraint(right, kiwi.OpEq, left.Plus(width)))
	// centerX == left + width / 2
	_ = solver.AddConstraint(kiwi.NewConstraint(centerX, kiwi.OpEq, left.Plus(width.Divide(2))))

	solver.UpdateVariables()

	fmt.Printf("   Left: %.2f | Width: %.2f | Right: %.2f | CenterX: %.2f\n",
		left.Value(), width.Value(), right.Value(), centerX.Value())

	fmt.Println("   Updating width to 80.0...")
	_ = solver.SuggestValue(width, 80.0)
	solver.UpdateVariables()

	fmt.Printf("   Left: %.2f | Width: %.2f | Right: %.2f | CenterX: %.2f\n\n",
		left.Value(), width.Value(), right.Value(), centerX.Value())
}

func runNaïveVsAutohintDemo() {
	fmt.Println("--- 2. Naïve Rounding vs. FreeType Autohinting (3 Equal Columns in 80 Chars) ---")

	c1 := kiwi.NewVariable("col1")
	c2 := kiwi.NewVariable("col2")
	c3 := kiwi.NewVariable("col3")

	vars := []*kiwi.Variable{c1, c2, c3}
	// Continuous solver output for 80 / 3 = 26.66666667
	floats := map[*kiwi.Variable]float64{
		c1: 26.66666667,
		c2: 26.66666667,
		c3: 26.66666667,
	}

	naive := kiwi.RoundValues(vars, floats)
	autohint := kiwi.ApportionSum(vars, floats, 80)

	fmt.Printf("   Continuous Values : %.4f, %.4f, %.4f (Sum: 80.0000)\n", floats[c1], floats[c2], floats[c3])
	fmt.Printf("   Naïve Rounding    : %d + %d + %d = %d (OVERFLOW by %d char!)\n",
		naive.Get(c1), naive.Get(c2), naive.Get(c3), naive.Sum(vars...), naive.Sum(vars...)-80)
	fmt.Printf("   Autohinting (FreeType): %d + %d + %d = %d (EXACT FIT, symmetric!)\n\n",
		autohint.Get(c1), autohint.Get(c2), autohint.Get(c3), autohint.Sum(vars...))

	renderBoxRow("Naïve (81 chars)", []int{naive.Get(c1), naive.Get(c2), naive.Get(c3)})
	renderBoxRow("Autohint (80 chars)", []int{autohint.Get(c1), autohint.Get(c2), autohint.Get(c3)})
	fmt.Println()
}

func runTrueTypeRulesDemo() {
	fmt.Println("--- 3. TrueType-Style Rule Directives (Grid Snapping & Clamping) ---")

	cjkCol := kiwi.NewVariable("cjkCol")
	smallBadge := kiwi.NewVariable("smallBadge")

	res := kiwi.DiscreteResult{
		cjkCol:     15, // Wants to snap to even grid (e.g. double-width CJK boundary)
		smallBadge: 1,  // Wants to clamp to min 3 chars for border [A]
	}

	fmt.Printf("   Before Hints : CJK Col = %d, Badge = %d\n", res.Get(cjkCol), res.Get(smallBadge))

	hinter := kiwi.NewRuleHinter()
	hinter.AddDirective(
		kiwi.SnapToGrid(cjkCol, 2),          // Snap to multiple of 2
		kiwi.ClampMinMax(smallBadge, 3, 10), // Clamp min 3 chars
	)
	hinter.Apply(res)

	fmt.Printf("   After Hints  : CJK Col = %d (snapped even), Badge = %d (clamped min 3)\n\n",
		res.Get(cjkCol), res.Get(smallBadge))
}

func runTUIWindowLayoutDemo() {
	fmt.Println("--- 4. Full TUI Window Layout Rendering (DiscreteSolver) ---")

	screenW := 80
	screenH := 12

	ds := kiwi.NewDiscreteSolver()
	solver := ds.Solver()

	totalW := kiwi.NewVariable("totalW")
	sidebarW := kiwi.NewVariable("sidebarW")
	mainW := kiwi.NewVariable("mainW")

	_ = solver.AddEditVariable(totalW, kiwi.StrengthStrong)
	_ = solver.SuggestValue(totalW, float64(screenW))

	// sidebar == 0.3 * totalW (24 chars)
	_ = solver.AddConstraint(kiwi.NewConstraint(sidebarW, kiwi.OpEq, totalW.Multiply(0.3)))
	_ = solver.AddConstraint(kiwi.NewConstraint(sidebarW.Plus(mainW), kiwi.OpEq, totalW))

	ds.SetMinSize(sidebarW, 15)
	ds.AddApportionGroup(kiwi.ApportionGroup{
		Vars:      []*kiwi.Variable{sidebarW, mainW},
		TargetVar: totalW,
	})
	ds.AddDirective(kiwi.SnapToGrid(sidebarW, 2)) // Snap sidebar to even width

	res := ds.SolveDiscrete()

	sw := res.Get(sidebarW)
	mw := res.Get(mainW)

	fmt.Printf("   Terminal Screen Size : %dx%d\n", screenW, screenH)
	fmt.Printf("   Sidebar Width        : %d chars (even grid)\n", sw)
	fmt.Printf("   Main Content Width   : %d chars\n", mw)
	fmt.Printf("   Total Layout Width   : %d chars\n\n", sw+mw)

	renderTUILayout(sw, mw, screenH-4)
}

func renderBoxRow(label string, widths []int) {
	fmt.Printf("   %-20s |", label)
	for i, w := range widths {
		text := fmt.Sprintf("Col %d (%d)", i+1, w)
		if len(text) > w {
			text = text[:w]
		}
		pad := w - len(text)
		padL := pad / 2
		padR := pad - padL
		fmt.Print(strings.Repeat(" ", padL) + text + strings.Repeat(" ", padR) + "|")
	}
	fmt.Println()
}

func renderTUILayout(sidebarW, mainW, contentHeight int) {
	totalW := sidebarW + mainW
	topBorder := "+" + safeRepeat("-", totalW-2) + "+"
	header := "| " + padCenter("KIWI-GO TUI AUTO-LAYOUT ENGINE", totalW-4) + " |"
	sepBar := "+" + safeRepeat("-", sidebarW-1) + "+" + safeRepeat("-", mainW-2) + "+"

	fmt.Println("   " + topBorder)
	fmt.Println("   " + header)
	fmt.Println("   " + sepBar)

	for i := 0; i < contentHeight; i++ {
		sideText := ""
		mainText := ""
		if i == 1 {
			sideText = "Tree View"
			mainText = "Welcome to Kiwi-Go TUI!"
		} else if i == 2 {
			sideText = " > src/"
			mainText = "Cassowary + FreeType Autohinting"
		} else if i == 3 {
			sideText = " > docs/"
			mainText = "Zero Gaps & Zero Overflows!"
		}
		line := "| " + padLeft(sideText, sidebarW-3) + " | " + padLeft(mainText, mainW-4) + " |"
		fmt.Println("   " + line)
	}

	botBorder := "+" + safeRepeat("-", sidebarW-1) + "+" + safeRepeat("-", mainW-2) + "+"
	footer := "| " + padLeft("Status: READY", sidebarW-3) + " | " + padLeft("Press Ctrl+C to exit", mainW-4) + " |"

	fmt.Println("   " + botBorder)
	fmt.Println("   " + footer)
	fmt.Println("   " + topBorder)
	fmt.Println()
}

func safeRepeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(s, count)
}

func padLeft(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if len(s) >= width {
		return s[:width]
	}
	return s + safeRepeat(" ", width-len(s))
}

func padCenter(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if len(s) >= width {
		return s[:width]
	}
	pad := width - len(s)
	padL := pad / 2
	padR := pad - padL
	return safeRepeat(" ", padL) + s + safeRepeat(" ", padR)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func runInteractiveTUI() {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return
	}

	fmt.Println("\n--- 5. Interactive Full TUI Window Layout ---")
	fmt.Println("   Press ENTER to start interactive mode (Resize window to test! Press any key to exit)...")
	buf := make([]byte, 1)
	_, _ = os.Stdin.Read(buf)

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Printf("   [Failed to enter raw mode: %v]\n", err)
		return
	}
	defer term.Restore(fd, oldState)

	ds := kiwi.NewDiscreteSolver()
	solver := ds.Solver()

	totalW := kiwi.NewVariable("totalW")
	sidebarW := kiwi.NewVariable("sidebarW")
	mainW := kiwi.NewVariable("mainW")

	_ = solver.AddEditVariable(totalW, kiwi.StrengthStrong)
	_ = solver.AddConstraint(kiwi.NewConstraint(sidebarW, kiwi.OpEq, totalW.Multiply(0.3)))
	_ = solver.AddConstraint(kiwi.NewConstraint(sidebarW.Plus(mainW), kiwi.OpEq, totalW))

	ds.SetMinSize(sidebarW, 15)
	ds.AddApportionGroup(kiwi.ApportionGroup{
		Vars:      []*kiwi.Variable{sidebarW, mainW},
		TargetVar: totalW,
	})
	ds.AddDirective(kiwi.SnapToGrid(sidebarW, 2))

	keyCh := make(chan struct{})
	go func() {
		b := make([]byte, 1)
		_, _ = os.Stdin.Read(b)
		keyCh <- struct{}{}
	}()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	lastW, lastH := -1, -1

	for {
		select {
		case <-keyCh:
			fmt.Print("\r\033[H\033[2J") // Clear screen
			fmt.Print("Exited interactive demo.\r\n")
			return
		case <-ticker.C:
			w, h, err := term.GetSize(fd)
			if err != nil {
				w, h = 80, 24
			}
			if w != lastW || h != lastH {
				lastW, lastH = w, h
				_ = solver.SuggestValue(totalW, float64(w))
				res := ds.SolveDiscrete()

				sw := res.Get(sidebarW)
				mw := res.Get(mainW)

				fmt.Print("\r\033[H\033[2J") // Clear screen
				renderInteractiveTUILayout(sw, mw, h, w)
			}
		}
	}
}

func renderInteractiveTUILayout(sidebarW, mainW, screenH, screenW int) {
	if sidebarW+mainW > screenW {
		mainW = screenW - sidebarW
	}
	if mainW < 2 {
		mainW = 2
	}
	if sidebarW < 2 {
		sidebarW = 2
	}

	totalW := sidebarW + mainW
	topBorder := "+" + safeRepeat("-", totalW-2) + "+"
	headerText := "KIWI-GO INTERACTIVE TUI (Press any key to exit)"
	if totalW-4 < len(headerText) {
		headerText = headerText[:max(0, totalW-4)]
	}
	header := "| " + padCenter(headerText, totalW-4) + " |"
	sepBar := "+" + safeRepeat("-", sidebarW-1) + "+" + safeRepeat("-", mainW-2) + "+"

	fmt.Print(topBorder + "\r\n")
	fmt.Print(header + "\r\n")
	fmt.Print(sepBar + "\r\n")

	contentHeight := screenH - 6
	if contentHeight < 1 {
		contentHeight = 1
	}

	for i := 0; i < contentHeight; i++ {
		sideText := ""
		mainText := ""
		if i == 1 {
			sideText = "Tree View"
			mainText = "Resize the terminal window!"
		} else if i == 2 {
			sideText = " > src/"
			mainText = fmt.Sprintf("Current Size: %dx%d", screenW, screenH)
		} else if i == 3 {
			sideText = " > docs/"
			mainText = "Sidebar is 30% (snapped to even)"
		}
		line := "| " + padLeft(sideText, sidebarW-3) + " | " + padLeft(mainText, mainW-4) + " |"
		fmt.Print(line + "\r\n")
	}

	botBorder := "+" + safeRepeat("-", sidebarW-1) + "+" + safeRepeat("-", mainW-2) + "+"
	footerText1 := "Status: RESIZING"
	if sidebarW-3 < len(footerText1) {
		footerText1 = footerText1[:max(0, sidebarW-3)]
	}
	footerText2 := fmt.Sprintf("W:%d H:%d", screenW, screenH)
	if mainW-4 < len(footerText2) {
		footerText2 = footerText2[:max(0, mainW-4)]
	}
	footer := "| " + padLeft(footerText1, sidebarW-3) + " | " + padLeft(footerText2, mainW-4) + " |"

	fmt.Print(botBorder + "\r\n")
	fmt.Print(footer + "\r\n")
	fmt.Print(topBorder + "\r\n")
}
