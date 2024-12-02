package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	max_threshold int8 = 3
)

var reports = make([][]int8, 1000)

func parse_input() {
	dat, err := os.ReadFile("input")
	if err != nil {
		panic(err)
	}

	for index, line := range strings.Split(string(dat), "\n") {
		fields := strings.Fields(line)
		report := make([]int8, len(fields))

		for index, field := range fields {
			num, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				panic(err)
			}

			report[index] = int8(num)
		}
		reports[index] = report
	}
}

func calculate_safe_reports(problem_dampen bool) (safe_reports uint16) {
	for _, report := range reports {
		if determine_report_safe(report, problem_dampen) {
			safe_reports++
		}
	}

	return
}

func determine_report_safe(report []int8, tolerate bool) bool {
	var (
		step      int8
		max_index = len(report) - 1
	)

	for index := 1; index <= max_index; index++ {
		prev_num := report[index-1]
		curr_num := report[index]

		// Make sure the previous and current number from the report is not equal and is below the allowed threshold.
		if diff := abs_int16(curr_num - prev_num); diff > max_threshold || diff == 0 {
			if tolerate {
				return check_report_with_tolerance(report, index)
			}
			return false
		}

		switch step {
		case 0:
			// If the step is zero, we haven't determined wether this report is increasing or decreasing.
			next_num := report[index+1]
			if prev_num > curr_num {
				step = -1
				if next_num > curr_num {
					if tolerate {
						// EDGE CASE: We need to check the current and previous index here. If we don't, we never attempt
						// to try and remove the first index for a possible valid report.
						return check_report_with_tolerance(report, index)
					}
					return false
				}
			} else {
				step = 1
				if next_num < curr_num {
					if tolerate {
						// EDGE CASE: We need to check the current and previous index here. If we don't, we never attempt
						// to try and remove the first index for a possible valid report.
						return check_report_with_tolerance(report, index)
					}
					return false
				}
			}
		case 1:
			// If the step is (positive) 1, we should be expecting an increasing report.
			if curr_num < prev_num {
				if tolerate {
					return check_report_with_tolerance(report, index)
				}
				return false
			}
		case -1:
			// If the step is (negative) 1, we should be expecting a decreasing report.
			if curr_num > prev_num {
				if tolerate {
					return check_report_with_tolerance(report, index)
				}
				return false
			}
		default:
			panic(fmt.Errorf("unexpected step value: %d", step))
		}
	}

	return true
}

func check_report_with_tolerance(report []int8, index int) bool {
	prev_index := index - 1
	copy_report_1 := make([]int8, len(report))
	copy_report_2 := make([]int8, len(report))
	copy(copy_report_1, report)
	copy(copy_report_2, report)

	// The commented out code below refuses to behave properly for some reason...
	// r1, r2 := append(report[:index], report[index+1:]...), append(report[:prev_index], report[prev_index+1:]...)

	s1, s2 := copy_report_1[:index], copy_report_1[index+1:]
	r1 := append(s1, s2...)

	s3, s4 := copy_report_2[:prev_index], copy_report_2[prev_index+1:]
	r2 := append(s3, s4...)

	return determine_report_safe(r1, false) || determine_report_safe(r2, false)
}

func abs_int16(v int8) int8 {
	if v < 0 {
		v *= -1
	}

	return v
}

func main() {
	parse_input()
	fmt.Println("safe reports w/o Problem Dampener:", calculate_safe_reports(false))
	fmt.Println("safe reports w/ Problem Dampener:", calculate_safe_reports(true))
}
