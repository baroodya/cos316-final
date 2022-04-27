package cache

import (
	"fmt"
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
 
 func TestPlotHits(t *testing.T) {
	 capacity := 1024
	 minVal := 256
	 maxVal := 768
	 lfu := NewLfu(capacity)
	 lru := NewLru(capacity)
	 log_lfu := NewLogLfu(capacity, 0.1, 0.1)
	 log_lfu_2 := NewLogLfu(capacity, 1.0, 0.01)
	 log_lfu_3 := NewLogLfu(capacity, 10.0, 1.0)
	 lfu_da := NewLFUDA(capacity, 0.001)
	 
	 trials := 100000

	 // choose trials random values between minVal and maxVal
	 lfu_hits := make([]opts.LineData, trials)
	 lru_hits := make([]opts.LineData, trials)
	 log_lfu_hits := make([]opts.LineData, trials)
	 log_lfu_2_hits := make([]opts.LineData, trials)
	 log_lfu_3_hits := make([]opts.LineData, trials)
	 lfu_da_hits := make([]opts.LineData, trials)
	 ideal_hits := make([]opts.LineData, trials)
	 xAxis:= make([]int, trials)
	 for i := 0; i < trials; i++ {
		xAxis[i] = i

		var randVal float64
		if i < trials / 4 {
			randVal = float64(minVal) * rand.Float64()
		} else {
			randVal = float64(minVal) + (float64(maxVal - minVal)) * rand.Float64()
		}
		key := fmt.Sprintf("%d", int(randVal))
		val := []byte(key)

		getLFUVal(t, lfu, key, val)
		getLRUVal(t, lru, key, val)
		getLogLFUVal(t, log_lfu, key, val)
		getLogLFUVal(t, log_lfu_2, key, val)
		getLogLFUVal(t, log_lfu_3, key, val)
		getLFUDAVal(t, lfu_da, key, val)

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
			log_lfu_2_hits[i] = opts.LineData{
				Value: 0.0,
			}
			log_lfu_3_hits[i] = opts.LineData{
				Value: 0.0,
			}
			lfu_da_hits[i] = opts.LineData{
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
			log_lfu_2_hits[i] = opts.LineData{
				Value: float64(log_lfu_2.stats.Hits) / float64(i),
			}
			log_lfu_3_hits[i] = opts.LineData{
				Value: float64(log_lfu_3.stats.Hits) / float64(i),
			}
			lfu_da_hits[i] = opts.LineData{
				Value: float64(lfu_da.stats.Hits) / float64(i),
			}
		}
		ideal_hits[i] = opts.LineData{
			Value: 0.85,
		}
	 }

	 // plot the results
	 // create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "LFU vs. Log LFU Cache Accuracy",
			Subtitle: "First 25000 accesses are 0-256, last 75000 are random 256-512",
		}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
		// charts.WithDataZoomOpts(opts.DataZoom{
		// 	Type:       "inside",
		// 	Start:      50,
		// 	End:        100,
		// 	XAxisIndex: []int{0},
		// }),
	)

	// Put data into instance
	line.SetXAxis(xAxis).
		AddSeries("LFU", lfu_hits).
		AddSeries("LRU", lru_hits).
		AddSeries("LogLFU", log_lfu_hits).
		// AddSeries("LogLFU 2", log_lfu_2_hits).
		// AddSeries("LogLFU 3", log_lfu_3_hits).
		AddSeries("LFU DA", lfu_da_hits).
		// AddSeries("Infinite Cache", ideal_hits).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create("line.html")
	line.Render(f)

	 // print stats for both caches
	 fmt.Printf("LFU hits: %d = %.3f\n", lfu.stats.Hits, float64(lfu.stats.Hits) / float64(trials))
	 fmt.Printf("Log LFU hits: %d = %.3f\n", log_lfu.stats.Hits, float64(log_lfu.stats.Hits) / float64(trials))

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
 