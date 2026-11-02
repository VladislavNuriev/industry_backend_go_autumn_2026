package main

type Stats struct {
	Count         int
	Sum, Min, Max int64
}

func Calc(nums []int64) Stats {
	n := len(nums)
	if n < 2 {
		return Stats{}
	}

	// Первая разность — инициализирует min и max.
	first := nums[1] - nums[0]
	stats := Stats{
		Count: n - 1,
		Sum:   first,
		Min:   first,
		Max:   first,
	}

	for i := 2; i < n; i++ {
		d := nums[i] - nums[i-1]
		stats.Sum += d
		if d < stats.Min {
			stats.Min = d
		}
		if d > stats.Max {
			stats.Max = d
		}
	}

	return stats
}
