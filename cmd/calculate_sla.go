package cmd

import (
	"fmt"
	"log"
	"sort"
	"statuscakectl/helpers"
	"statuscakectl/statuscake"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/spf13/cobra"
)

var debug bool

var calculateSlaCmd = &cobra.Command{
	Use:   "calculate-sla",
	Short: "calculates sla for domain from to date",
	Example: `statuscakectl calculate-sla --domains testdomain.com
statuscakectl calculate-sla --domains foo.com,bar.org --from 2021-12-01 -to 2022-01-01
statuscakectl calculate-sla --domains foo.com,bar.org --maintenance-start-hour 0 --maintenance-finish-hour 2`,
	Run: func(cmd *cobra.Command, args []string) {
		api, _ := cmd.Flags().GetString("api")
		user, _ := cmd.Flags().GetString("user")
		key, _ := cmd.Flags().GetString("key")
		domains, _ := cmd.Flags().GetStringSlice("domains")
		fromString, _ := cmd.Flags().GetString("from")
		toString, _ := cmd.Flags().GetString("to")

		mStartHour, _ := cmd.Flags().GetInt("maintenance-start-hour")
		mEndHour, _ := cmd.Flags().GetInt("maintenance-finish-hour")

		from := parseDateToUnix(fromString)
		to := parseDateToUnix(toString)
		testID := 0
		var slas []float64

		if from > to && to != 0 {
			fmt.Println("to date should be after from")
			return
		}

		if len(domains) < 1 {
			fmt.Println("Please make sure to provide a valid domains flag")
			return
		}

		var allPeriods []statuscake.Period

		uptimeTests := statuscake.ListUptime(api, user, key)

		for _, domain := range domains {
			if domain != "" {
				for _, t := range uptimeTests {
					hostname, err := helpers.GetHostnameFromUrl(t.WebsiteURL)
					if err != nil {
						log.Printf("Can't get hostname from this URL:%v\n", t.WebsiteURL)
					}
					if hostname == domain {
						testID = t.TestID
						break
					}
				}
				if testID < 1 {
					fmt.Printf("Cannot find test for domain: %v\n", domain)
					return
				}
			}
			periods := statuscake.ListPeriods(api, user, key, testID)
			allPeriods = append(allPeriods, periods...)

			from, to, sla := calculateSLA(mStartHour, mEndHour, int(from), int(to), periods)
			slas = append(slas, sla)
			fmt.Printf("SLA for %v from:%v to:%v is %0.2f\n", domain, unixDateToString(from), unixDateToString(to), sla)
		}

		fmt.Printf("Worst SLA: %0.2f\n", helpers.Smallest(slas))
		meanSla, err := stats.Mean(slas)
		if err != nil {
			fmt.Printf("Can't calculate mean sla due to:%v\n", meanSla)
		} else {
			fmt.Printf("Mean SLA: %0.2f\n", meanSla)
		}

		medianSla, err := stats.Median(slas)
		if err != nil {
			fmt.Printf("Can't calculate median sla due to:%v\n", medianSla)
		} else {
			fmt.Printf("Median SLA: %0.2f\n", medianSla)
		}

		_, _, sla := calculateSLA(mStartHour, mEndHour, int(from), int(to), allPeriods)
		fmt.Printf("Combined down SLA: %0.2f\n", sla)

	},
}

func init() {
	// flags
	calculateSlaCmd.Flags().StringSliceP("domains", "d", []string{}, "Domain names to find periods for (eg. foo.com or foo.com,bar.com")
	calculateSlaCmd.Flags().String("from", "", "Date from which calculate SLA (Format: \"2006-01-02\"")
	calculateSlaCmd.Flags().String("to", "", "Date to which calculate SLA (Format: \"2006-01-02\"")
	calculateSlaCmd.Flags().Int("maintenance-start-hour", 0, "Periodic maintenance start hour in UTC (experimental)")
	calculateSlaCmd.Flags().Int("maintenance-finish-hour", 0, "Periodic maintenance end hour in UTC (experimental)")
}

func unixDateToString(u int64) string {
	return time.Unix(u, 0).Format("2006-01-02")
}

func parseDateToUnix(s string) int64 {
	if s == "" {
		return 0
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Fatalln("Cannot parse given date:", s)
	}
	return t.Unix()
}

// calculateSLA takes from and to dates in unixseconds format and returns real from, real to and sla
func calculateSLA(mStart, mEnd int, from, to int, periods []statuscake.Period) (int64, int64, float64) {
	sort.Slice(periods, func(i, j int) bool {
		return periods[i].StartUnix < periods[j].StartUnix
	})

	first := getFirstTime(periods)
	if first > from {
		from = first
	}

	if to > int(time.Now().Unix()) || to < 1 {
		to = int(time.Now().Unix())
	}

	totalSeconds := to - from
	var downSeconds int
	for _, p := range periods {

		if p.EndUnix < from {
			continue
		}

		if p.StartUnix > to {
			continue
		}

		if p.Status == "Down" {
			if p.StartUnix < from {
				downSeconds = downSeconds + p.EndUnix - from
			} else if p.EndUnix > to {
				downSeconds = downSeconds + to - p.StartUnix
			} else {
				downSeconds = downSeconds + p.EndUnix - p.StartUnix
			}

			maintenanceWindowSeconds := secondsCoveredByMaintenanceWindow(p, mStart, mEnd)
			downSeconds = downSeconds - maintenanceWindowSeconds

		}

	}
	sla := 100 - ((100 * float64(downSeconds)) / float64(totalSeconds))
	return int64(from), int64(to), sla
}

// finds first time check
func getFirstTime(periods []statuscake.Period) int {
	if len(periods) < 1 {
		return 0
	}
	first := periods[0].StartUnix
	for _, p := range periods {
		if p.StartUnix < first {
			first = p.StartUnix
		}
	}
	return first
}

func secondsCoveredByMaintenanceWindow(p statuscake.Period, startHour, endHour int) int {
	startTime := time.Unix(int64(p.StartUnix), 0)
	endTime := time.Unix(int64(p.EndUnix), 0)

	var mww []maintenanceWindow

	for i := 0.0; i < endTime.Sub(startTime).Hours(); i = i + 24 {
		date := startTime.Add(time.Duration(i) * time.Hour)
		mww = append(mww, maintenanceWindowForDate(date, startHour, endHour))
	}

	var secondsCoveredByMaintenanceWindow int

	for _, mw := range mww {
		// if period is fully within maintenance window
		if startTime.After(mw.start) && endTime.Before(mw.end) {
			secondsCoveredByMaintenanceWindow = secondsCoveredByMaintenanceWindow + p.EndUnix - p.StartUnix
			continue
		}

		// if period starts within maintenance window
		if startTime.After(mw.start) && startTime.Before(mw.end) {
			secondsCoveredByMaintenanceWindow = secondsCoveredByMaintenanceWindow + int(mw.end.Unix()) - int(startTime.Unix())
			continue
		}

		// if period ends within maintenance window
		if endTime.After(mw.start) && endTime.Before(mw.end) {
			secondsCoveredByMaintenanceWindow = secondsCoveredByMaintenanceWindow + int(endTime.Unix()) - int(mw.start.Unix())
			continue
		}

		// if maintenance windows is within period
		if startTime.Before(mw.start) && endTime.After(mw.end) {
			secondsCoveredByMaintenanceWindow = secondsCoveredByMaintenanceWindow + int(mw.start.Unix()) - int(mw.end.Unix())
			continue
		}
	}

	return secondsCoveredByMaintenanceWindow
}

type maintenanceWindow struct {
	start time.Time
	end   time.Time
}

func maintenanceWindowForDate(date time.Time, startHour, endHour int) maintenanceWindow {
	start := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		startHour,
		0,
		0,
		0,
		time.UTC,
	)

	end := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		endHour,
		0,
		0,
		0,
		time.UTC,
	)

	return maintenanceWindow{
		start: start,
		end:   end,
	}

}
