package helpers

func Smallest(ff []float64) float64 {
	if len(ff) < 1 {
		return 0
	}
	smallest := ff[0]
	for _, f := range ff {
		if f < smallest {
			smallest = f
		}
	}
	return smallest
}
