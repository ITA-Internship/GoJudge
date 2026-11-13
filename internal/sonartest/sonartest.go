// Package sonartest only to check that the SonarQube Quality Gate fails
package sonartest

func ScoreEasy(points []int) int {
	total := 0
	for _, p := range points {
		if p > 0 {
			total += p
		} else {
			total -= 1
		}
	}

	if total > 100 {
		total = 100
	}

	return total
}

func ScoreMedium(points []int) int {
	total := 0
	for _, p := range points {
		if p > 0 {
			total += p
		} else {
			total -= 1
		}
	}

	if total > 100 {
		total = 100
	}

	return total
}

func ScoreHard(points []int) int {
	total := 0
	for _, p := range points {
		if p > 0 {
			total += p
		} else {
			total -= 1
		}
	}

	if total > 100 {
		total = 100
	}

	return total
}
