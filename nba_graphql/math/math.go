package math

import "math"

func Mean(x []float64) float64 {
	var sum float64
	for _, v := range x {
		sum += v
	}
	return sum / float64(len(x))
}

func Mean_Online(x float64, before float64, nBefore int) float64 {
	return (x + before) / float64(nBefore+1)
}

func StdDev(x []float64) float64 {
	m := Mean(x)
	var sum float64
	for _, v := range x {
		sum += (v - m) * (v - m)
	}
	return RoundFloat(math.Sqrt(sum/float64(len(x))), 2)
}

func StdDevWithMean(x []float64, mean float64) float64 {
	var sum float64
	for _, v := range x {
		sum += (v - mean) * (v - mean)
	}
	return RoundFloat(math.Sqrt(sum/float64(len(x))), 2)
}

func StdDevMean(x []float64) (float64, float64) {
	m := Mean(x)
	stdDev := StdDev(x)
	return stdDev, m
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	x := math.Round(val*ratio) / ratio
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0
	}
	return x
}

func ZScore(x, mean, stdDev float64) float64 {
	zscore := (x - mean) / stdDev
	//round to 2 decimal places
	return RoundFloat(zscore, 2)
}
