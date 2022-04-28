package cache

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	// "github.com/go-echarts/go-echarts/v2/types"
)

/******************************************************************************/
/*                                Constants                                   */
/******************************************************************************/
// Constants can go here

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/
 
// stores 3 byte keys and 1 byte values
 func TestPlotHits(t *testing.T) {
	 capacity := 1024
	 inf_capacity := int(math.Exp2(32))
	 minVal := 0
	 maxVal := 2048
	 lfu := NewLfu(capacity)
	 lru := NewLru(capacity)
	 log_lfu := NewLogLfu(capacity, 0.1, 10.0)
	 lin_lfu := NewLinearLfu(capacity, 0.5)
	 exp_lfu := NewExpLfu(capacity, 0.1, 0.5)
	 lfu_da := NewLFUDA(capacity)
	 ideal := NewLfu(inf_capacity)
	 
	 trials := 100000

	 // choose trials random values between minVal and maxVal
	 lfu_hits := make([]opts.LineData, trials)
	 lru_hits := make([]opts.LineData, trials)
	 log_lfu_hits := make([]opts.LineData, trials)
	 lin_lfu_hits := make([]opts.LineData, trials)
	 exp_lfu_hits := make([]opts.LineData, trials)
	 lfu_da_hits := make([]opts.LineData, trials)
	 ideal_hits := make([]opts.LineData, trials)
	 xAxis:= make([]int, trials)
	 for i := 0; i < trials; i++ {
		xAxis[i] = i

		var randVal float64
		// if i < trials / 4 {
		// 	randVal = float64((minVal)) * rand.Float64()
		// } else if i < trials / 2 {
		// 	randVal = float64(i % minVal + minVal)
		// } else {
		// 	randVal = float64(minVal) + (float64(maxVal - minVal)) * rand.Float64()
		// }
		randVal = float64(minVal) + (float64(maxVal - minVal)) * math.Exp((-10 * math.Pow(rand.Float64(),2)))
		var key string
		if randVal < 10 {
			key = fmt.Sprintf("__%d", int(randVal))
		} else if randVal < 100 {
			key = fmt.Sprintf("_%d", int(randVal))
		}
		key = fmt.Sprintf("%d", int(randVal))
		val := []byte(key)

		getLFUVal(t, lfu, key, val)
		getLRUVal(t, lru, key, val)
		getLogLFUVal(t, log_lfu, key, val)
		getLinearLFUVal(t, lin_lfu, key, val)
		getExpLFUVal(t, exp_lfu, key, val)
		getLFUDAVal(t, lfu_da, key, val)
		getLFUVal(t, ideal, key, val)

		if i == 0 {
			lfu_hits[i] = opts.LineData{
				Value: 0.0,
			}
			lru_hits[i] = opts.LineData{
				Value: 0.0,
			}
			log_lfu_hits[i] = opts.LineData{
				Value: 0.0,
			}
			lin_lfu_hits[i] = opts.LineData{
				Value: 0.0,
			}
			exp_lfu_hits[i] = opts.LineData{
				Value: 0.0,
			}
			lfu_da_hits[i] = opts.LineData{
				Value: 0.0,
			}
			ideal_hits[i] = opts.LineData{
				Value: 0.0,
			}
		} else {
			lfu_hits[i] = opts.LineData{
				Value: float64(lfu.stats.Hits) / float64(i),
			}
			lru_hits[i] = opts.LineData{
				Value: float64(lru.stats.Hits) / float64(i),
			}
			log_lfu_hits[i] = opts.LineData{
				Value: float64(log_lfu.stats.Hits) / float64(i),
			}
			lin_lfu_hits[i] = opts.LineData{
				Value: float64(lin_lfu.stats.Hits) / float64(i),
			}
			exp_lfu_hits[i] = opts.LineData{
				Value: float64(exp_lfu.stats.Hits) / float64(i),
			}
			lfu_da_hits[i] = opts.LineData{
				Value: float64(lfu_da.stats.Hits) / float64(i),
			}
			ideal_hits[i] = opts.LineData{
				Value: float64(ideal.stats.Hits) / float64(i),
			}
		}
	 }

	 // plot the results
	 // create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Hit Rate for Cache Algorithms",
			// Subtitle: "First 25000 accesses are randomly chosen between 0 and 256, last 75000 are randomly chosen between 256 and 512", 
			Subtitle: "Accesses are random between 0 and 2048, according to the PDF: e^(-10 * x^2)",
		}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
		charts.WithColorsOpts(opts.Colors{"blue", "red", "green", "orange", "purple"}),
		// charts.WithDataZoomOpts(opts.DataZoom{
		// 	Type:       "inside",
		// 	Start:      100,
		// 	End:        150,
		// 	YAxisIndex: []int{0},
		// }),
	)

	// Put data into instance
	line.SetXAxis(xAxis).
		// AddSeries("LFU", lfu_hits).
		// AddSeries("LRU", lru_hits).
		AddSeries("LogLFU", log_lfu_hits).
		AddSeries("LinLFU", lin_lfu_hits).
		AddSeries("ExpLFU", exp_lfu_hits).
		AddSeries("LFU DA", lfu_da_hits).
		// AddSeries("Infinite Cache", ideal_hits).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create("line.html")
	line.Render(f)
 }

 func getLFUVal(t *testing.T, cache *LFU, key string, val []byte) {
	_, ok := cache.Get(key)
	if !ok {
		ok = cache.Set(key, val)
		if !ok {
			fmt.Printf("Failed to add binding to lfu with key: %s\n", key)
			t.FailNow()
		}
	}
 }

 func getLRUVal(t *testing.T, cache *LRU, key string, val []byte) {
	_, ok := cache.Get(key)
	if !ok {
		ok = cache.Set(key, val)
		if !ok {
			fmt.Printf("Failed to add binding to lfu with key: %s\n", key)
			t.FailNow()
		}
	}
 }

 func getLogLFUVal(t *testing.T, cache *LogLFU, key string, val []byte) {
	_, ok := cache.Get(key)
	if !ok {
		ok = cache.Set(key, val)
		if !ok {
			fmt.Printf("Failed to add binding to lfu with key: %s\n", key)
			t.FailNow()
		}
	}
 }

 func getLFUDAVal(t *testing.T, cache *LFUDA, key string, val []byte) {
	_, ok := cache.Get(key)
	if !ok {
		ok = cache.Set(key, val)
		if !ok {
			fmt.Printf("Failed to add binding to lfu with key: %s\n", key)
			t.FailNow()
		}
	}
 }

 func getLinearLFUVal(t *testing.T, cache *LinearLFU, key string, val []byte) {
	_, ok := cache.Get(key)
	if !ok {
		ok = cache.Set(key, val)
		if !ok {
			fmt.Printf("Failed to add binding to lfu with key: %s\n", key)
			t.FailNow()
		}
	}
 }

 func getExpLFUVal(t *testing.T, cache *ExpLFU, key string, val []byte) {
	_, ok := cache.Get(key)
	if !ok {
		ok = cache.Set(key, val)
		if !ok {
			fmt.Printf("Failed to add binding to lfu with key: %s\n", key)
			t.FailNow()
		}
	}
 }
 