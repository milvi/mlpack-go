package main

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Le bilan humain distingue un calcul terminé d’une expiration sans partition.
func writeReport(out io.Writer, m Measurement) error {
	var text strings.Builder
	algorithm := "dual-tree / KD-tree"
	if m.Algorithm == "naive" {
		algorithm = "Lloyd naïf"
	}
	fmt.Fprintf(&text, "\nK-means — %s\n", algorithm)
	fmt.Fprintln(&text, "Points uniformes dans le carré [0, 1]²")
	row := func(label, value string) {
		fmt.Fprintf(&text, "  %-19s %s\n", label, value)
	}
	text.WriteByte('\n')
	state := "Convergence atteinte"
	if m.TimedOut {
		state = "Budget atteint — calcul interrompu"
	} else if !m.Converged {
		state = "Calcul terminé sans convergence"
	}
	row("État", state)
	text.WriteByte('\n')
	row("Points", grouped(uint64(m.N)))
	requested := " demandées"
	if m.K == 1 {
		requested = " demandée"
	}
	classes := grouped(uint64(m.K)) + requested
	if !m.TimedOut {
		obtained := " obtenues / "
		if m.KActual == 1 {
			obtained = " obtenue / "
		}
		classes = grouped(uint64(m.KActual)) + obtained + classes
	}
	row("Classes", classes)
	thread := "thread"
	if m.Threads != 1 {
		thread = "threads"
	}
	row("Exécution", fmt.Sprintf("%s %s · graine %s", grouped(uint64(m.Threads)), thread, grouped(m.Seed)))
	text.WriteByte('\n')
	row("Temps de calcul", duration(m.Seconds)+" / budget "+duration(m.DurationLimitSeconds))
	row("Tours terminés", grouped(uint64(m.Iterations)))
	row("Pic mémoire", memory(uint64(m.PeakRSSBytes))+" (processus entier)")
	row("Tas Go", memory(m.HeapAllocBytes)+" (après calcul)")
	if m.TimedOut {
		fmt.Fprintln(&text, "\nAucune partition disponible : le calcul a été interrompu.")
	} else {
		text.WriteByte('\n')
		row("Inertie", number(m.Inertia))
		row("Distances comptées", grouped(m.DistanceCalculations))
		text.WriteByte('\n')
	}
	_, err := io.WriteString(out, text.String())
	return err
}

// L’espace fine insécable regroupe les chiffres sans perdre leur précision.
func grouped(n uint64) string {
	s := strconv.FormatUint(n, 10)
	var out strings.Builder
	for i, digit := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out.WriteRune('\u202f')
		}
		out.WriteRune(digit)
	}
	return out.String()
}

func decimal(value float64, precision int) string {
	s := strconv.FormatFloat(value, 'f', precision, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}
	return strings.ReplaceAll(s, ".", ",")
}

func duration(seconds float64) string {
	switch {
	case seconds == 0:
		return "0 s"
	case seconds < 1e-6:
		return decimal(seconds*1e9, 2) + " ns"
	case seconds < 1e-3:
		return decimal(seconds*1e6, 2) + " µs"
	case seconds < 1:
		return decimal(seconds*1e3, 2) + " ms"
	default:
		return decimal(seconds, 3) + " s"
	}
}

func memory(bytes uint64) string {
	units := [...]string{"o", "Kio", "Mio", "Gio", "Tio", "Pio", "Eio"}
	value, unit := float64(bytes), 0
	for value >= 1024 && unit < len(units)-1 {
		value /= 1024
		unit++
	}
	return decimal(value, 2) + " " + units[unit]
}

func number(value float64) string {
	if a := math.Abs(value); a != 0 && (a < 1e-4 || a >= 1e9) {
		return strings.ReplaceAll(strconv.FormatFloat(value, 'g', 6, 64), ".", ",")
	}
	return decimal(value, 6)
}
