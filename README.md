# statistics
Basic statistical processor

## Notes

```golang
s := SampleSet()
s.Add([]int{1, 3, -5, 6, 15, 44})
s.Add([]float{33.0, 22.1, -3.5})

stdev  := s.PStdev()
mean   := s.Mean()
median := s.Median()
```