package dsp

import (
	"log"
	"math"
	"sort"
)

// DataSet is a float64 slice with util functions.
type DataSet []float64

// Bounds returns the bounds for the dataset
func (d DataSet) Bounds() (float64, float64) {
	maxValue := math.Inf(-1)
	minValue := math.Inf(1)
	for i := 0; i < len(d); i++ {
		if d[i] > maxValue {
			maxValue = d[i]
		}
		if d[i] < minValue {
			minValue = d[i]
		}
	}
	return minValue, maxValue
}

// Range returns the data range for the data set
func (d DataSet) Range() float64 {
	min, max := d.Bounds()
	return max - min
}

// Len returns the length of the data set
func (d DataSet) Len() int {
	return len(d)
}

// Derivative returns the derivative
func (d DataSet) Derivative() DataSet {
	deriv := make([]float64, len(d))
	for i := 1; i < d.Len(); i++ {
		deriv[i] = d[i] - d[i-1]
	}
	return DataSet(deriv)
}

// MapRange maps the dataset onto [0, 1) after dividing the value by the entire range.
func (d DataSet) MapRange() DataSet {
	min, max := d.Bounds()
	dataRange := d.Range()
	return d.Do(func(v float64) float64 {
		return Map(v/dataRange, min/dataRange, max/dataRange, 0, 1)
	})
}

// Map maps the data from one range to another
func (d DataSet) Map(start1, stop1, start2, stop2 float64) DataSet {
	mappedValues := make([]float64, len(d))
	for i := 0; i < len(d); i++ {
		mappedValues[i] = Map(d[i], start1, stop1, start2, stop2)
	}
	return mappedValues
}

// Mult multiplies all the points by the given number and returns a new DataSet.
func (d DataSet) Mult(num float64) DataSet {
	return d.Do(Mult(num))
}

// Div divides all the points by the given number and returns a new DataSet.
func (d DataSet) Div(denom float64) DataSet {
	return d.Do(Div(denom))
}

// Add adds a number to all the points and returns a new DataSet.
func (d DataSet) Add(num float64) DataSet {
	return d.Do(Add(num))
}

// Sub subtracts a number from all the points and returns a new DataSet.
func (d DataSet) Sub(num float64) DataSet {
	return d.Do(Sub(num))
}

// Do performs a function on each point.
func (d DataSet) Do(fns ...MapFunc) DataSet {
	if len(fns) == 0 {
		log.Fatal("Do requires at least one function")
	}
	values := make([]float64, d.Len())
	for i := 1; i < d.Len(); i++ {
		val := fns[0](d[i])
		for j := 1; j < len(fns); j++ {
			val = fns[j](val)
		}
		values[i] = val
	}
	return DataSet(values)
}

// Reduce performs the reduce function on the data
func (d DataSet) Reduce(fn ReduceFunc) float64 {
	return fn(d)
}

// Min returns the minimum
func (d DataSet) Min() float64 {
	return d.Reduce(minReduce)
}

// Max returns the maximum
func (d DataSet) Max() float64 {
	return d.Reduce(maxReduce)
}

// Sum returns the sum of the data set
func (d DataSet) Sum() float64 {
	return d.Reduce(sumReduce)
}

// Mean returns the average of the data set
func (d DataSet) Mean() float64 {
	if d.Len() == 0 {
		return 0
	}
	return d.Reduce(sumReduce) / float64(d.Len())
}

// Var returns the variance
func (d DataSet) Var() float64 {
	if len(d) <= 1 {
		return 0.0
	}
	return d.Reduce(func(data []float64) float64 {
		var sum, ssq float64
		for i := 0; i < len(d); i++ {
			sum += d[i]
			ssq += d[i] * d[i]
			// val += math.Pow(d[i]-mean, 2) / float64(n)
		}
		n := float64(len(d))
		mean := sum / n
		return ssq/n - mean*mean
	})
}

// Stdev returns the standard deviation
func (d DataSet) Stdev() float64 {
	return math.Sqrt(d.Var())
}

// Sort copies and sorts the data
func (d DataSet) Sort() []float64 {
	s := make([]float64, len(d))
	copy(s, d)
	sort.Float64s(s)
	return s
}

// Median returns the median of the dataset
func (d DataSet) Median() float64 {
	s := d.Sort()
	half := len(s) / 2
	m := s[half]
	if len(s)%2 == 0 {
		m = (m + s[half-1]) / 2
	}
	return m
}

// MapFunc is a function that can be performed on a dataset
type MapFunc func(float64) float64

// Div returns a MapFunc which divides the point by a denominator.
func Div(denom float64) MapFunc {
	return func(v float64) float64 {
		return v / denom
	}
}

// Mult returns a MapFunc which multiplies the point by a numerator.
func Mult(num float64) MapFunc {
	return func(v float64) float64 {
		return v * num
	}
}

// Add returns a MapFunc which adds the given number to the to the point.
func Add(num float64) MapFunc {
	return func(v float64) float64 {
		return v + num
	}
}

// Sub returns a MapFunc which subtracts the given number from the point.
func Sub(num float64) MapFunc {
	return Add(-num)
}

// ReduceFunc reduces the float64 slice to a single float64 value
type ReduceFunc func([]float64) float64

// minReduce is a reduce function for finding the minimum.
func minReduce(data []float64) float64 {
	minValue := math.Inf(1)
	for i := 0; i < len(data); i++ {
		if data[i] < minValue {
			minValue = data[i]
		}
	}
	return minValue
}

// maxReduce is a reduce function for finding the maximum
func maxReduce(data []float64) float64 {
	maxValue := math.Inf(-1)
	for i := 0; i < len(data); i++ {
		if data[i] > maxValue {
			maxValue = data[i]
		}
	}
	return maxValue
}

// sumReduce is a reduce function for finding the sum.
func sumReduce(data []float64) float64 {
	value := 0.0
	for i := 0; i < len(data); i++ {
		value += data[i]
	}
	return value
}

// Map maps a number from one range to another.
func Map(value, start1, stop1, start2, stop2 float64) float64 {
	return start2 + (stop2-start2)*((value-start1)/(stop1-start1))
}
