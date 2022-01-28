package date

import (
	"math"
	"time"
)

// 时间格式 常量
const (
	Format01 = "2006-01-02"
	Format02 = "2006-01-02 15:04:05"
	Format03 = "2006/01/02"
	Format04 = "2006/01/02 15:04:05"
	Format05 = "20060102"
	Format06 = "20060102150405"
)

const roundEpsilon = 1e-9
const nanosInADay = float64((24 * time.Hour) / time.Nanosecond)

/**
*  ParseTime
*  @Description: 解析时间字符串
*  @param value 时间
*  @param layout 时间格式 可以自定义 和 使用 标准库自带
*  @return t
*  @return err
 */
func ParseTime(value string, layout string) (time.Time, error) {
	return time.Parse(layout, value)
}

/**
*  SubMonth
*  @Description: 两个时间 相差多个月 自然月
*  @param time1
*  @param time2
*  @return int
 */
func SubMonth(time1, time2 time.Time) int {
	y1 := time1.Year()
	y2 := time2.Year()
	m1 := int(time1.Month())
	m2 := int(time2.Month())
	d1 := time1.Day()
	d2 := time2.Day()

	diffYear := y1 - y2
	// 如果 d1的 月-日 小于 d2的 月-日 那么 diffYear-- 这样就得到了相差的年数
	if m1 < m2 || m1 == m2 && d1 < d2 {
		diffYear--
	}
	diffMonth := (m1 + 12) - m2
	if d1 < d2 {
		diffMonth--
	}
	diffMonth %= 12

	return diffYear*12 + diffMonth
}

/**
*  IsLeapYear
*  @Description: 是否是闰年
*  @param year
*  @return bool
 */
func IsLeapYear(year int) bool {
	if year == year/400*400 {
		return true
	}
	if year == year/100*100 {
		return false
	}
	return year == year/4*4
}

// 参考或copy了 excelize 的时间处理部分
// excel的时间部分是从1900-01-01起算
// 也有excel的时间部分是从1904-01-01起算
const OFFSET1900 = 15018.0
const OFFSET1904 = 16480.0
const MJD0 float64 = 2400000.5

var (
	daysInMonth           = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	excel1900Epoc         = time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)
	excel1904Epoc         = time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC)
	excelMinTime1900      = time.Date(1899, time.December, 31, 0, 0, 0, 0, time.UTC)
	excelBuggyPeriodStart = time.Date(1900, time.March, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)
)

func ExcelDateToTime(excelTime float64) time.Time {
	return timeFromExcelDate(excelTime, false)
}

func ExcelDateToTime1904(excelTime float64) time.Time {
	return timeFromExcelDate(excelTime, false)
}

// timeFromExcelTime provides a function to convert an excelTime
// representation (stored as a floating point number) to a time.Time.
func timeFromExcelDate(excelTime float64, date1904 bool) time.Time {
	var days = int(excelTime)
	var date time.Time
	if days <= 61 { // 没懂 ！！！！
		if date1904 {
			date = julianDateToGregorianTime(MJD0, excelTime+OFFSET1900)
		} else {
			date = julianDateToGregorianTime(MJD0, excelTime+OFFSET1904)
		}
		return date
	}
	var floats = excelTime - float64(days) + roundEpsilon
	durationPart := time.Duration(nanosInADay * floats)
	if date1904 {
		date = excel1904Epoc
	} else {
		date = excel1900Epoc
	}

	return date.AddDate(0, 0, days).Add(durationPart).Truncate(time.Second)
}

// julianDateToGregorianTime provides a function to convert julian date to
// gregorian time.
func julianDateToGregorianTime(part1, part2 float64) time.Time {
	part1I, part1F := math.Modf(part1)
	part2I, part2F := math.Modf(part2)
	julianDays := part1I + part2I
	julianFraction := part1F + part2F
	julianDays, julianFraction = shiftJulianToNoon(julianDays, julianFraction)
	day, month, year := doTheFliegelAndVanFlandernAlgorithm(int(julianDays))
	hours, minutes, seconds, nanoseconds := fractionOfADay(julianFraction)
	return time.Date(year, time.Month(month), day, hours, minutes, seconds, nanoseconds, time.UTC)
}

// shiftJulianToNoon provides a function to process julian date to noon.
func shiftJulianToNoon(julianDays, julianFraction float64) (float64, float64) {
	switch {
	case -0.5 < julianFraction && julianFraction < 0.5:
		julianFraction += 0.5
	case julianFraction >= 0.5:
		julianDays++
		julianFraction -= 0.5
	case julianFraction <= -0.5:
		julianDays--
		julianFraction += 1.5
	}
	return julianDays, julianFraction
}

// doTheFliegelAndVanFlandernAlgorithm; By this point generations of
// programmers have repeated the algorithm sent to the editor of
// "Communications of the ACM" in 1968 (published in CACM, volume 11, number
// 10, October 1968, p.657). None of those programmers seems to have found it
// necessary to explain the constants or variable names set out by Henry F.
// Fliegel and Thomas C. Van Flandern.  Maybe one day I'll buy that jounal and
// expand an explanation here - that day is not today.
func doTheFliegelAndVanFlandernAlgorithm(jd int) (day, month, year int) {
	l := jd + 68569
	n := (4 * l) / 146097
	l = l - (146097*n+3)/4
	i := (4000 * (l + 1)) / 1461001
	l = l - (1461*i)/4 + 31
	j := (80 * l) / 2447
	d := l - (2447*j)/80
	l = j / 11
	m := j + 2 - (12 * l)
	y := 100*(n-49) + i + l
	return d, m, y
}

// fractionOfADay provides a function to return the integer values for hour,
// minutes, seconds and nanoseconds that comprised a given fraction of a day.
// values would round to 1 us.
func fractionOfADay(fraction float64) (hours, minutes, seconds, nanoseconds int) {

	const (
		c1us  = 1e3
		c1s   = 1e9
		c1day = 24 * 60 * 60 * c1s
	)

	frac := int64(c1day*fraction + c1us/2)
	nanoseconds = int((frac%c1s)/c1us) * c1us
	frac /= c1s
	seconds = int(frac % 60)
	frac /= 60
	minutes = int(frac % 60)
	hours = int(frac / 60)
	return
}
