# statistics
Basic statistical processor

## Overview

Perform basic statistical functions on sample and population sets.  Objects are used to reduce duplication of functions
(e.g., sorting) necessary to compute different statistics.  Also seeks to be safe for concurrency.

## Examples

```golang
import (
    "fmt"
    stats "github.com/blorticus-go/statistics"
)

func main() {
    s, err := stats.MakeStatisticalSampleSetFrom([]float64{
        0, 3.4, -10.23, 85, 1, 3.4, 0, -5, 3.4,
    })

    if err != nil {
        panic(err)
    }

    q1, q3, iqr := s.InterQuartileRange()

    fmt.Printf("%f, %f, %f, %f, %f, %f, %f\n", s.Mean(), s.Median(), s.SampleStdev(), s.Range(), s.Minimum(), s.Maximum(), iqr)
}
```
