package carbon

import (
	"strconv"
	"strings"
	"time"
)

type Carbon struct {
	Time  time.Time
	Loc   *time.Location
	Error error
}

// Now 当前
func (c Carbon) Now() Carbon {
	c.Time = time.Now()
	return c
}

// Now 当前(默认时区)
func Now() Carbon {
	return SetTimezone(Local).Now()
}

// Tomorrow 明天
func (c Carbon) Tomorrow() Carbon {
	if c.Time.IsZero() {
		c.Time = time.Now().AddDate(0, 0, 1)
	} else {
		c.Time = c.Time.AddDate(0, 0, 1)
	}
	return c
}

// Tomorrow 明天(默认时区)
func Tomorrow() Carbon {
	return SetTimezone(Local).Tomorrow()
}

// Yesterday 昨天
func (c Carbon) Yesterday() Carbon {
	if c.Time.IsZero() {
		c.Time = time.Now().AddDate(0, 0, -1)
	} else {
		c.Time = c.Time.AddDate(0, 0, -1)
	}
	return c
}

// Yesterday 昨天(默认时区)
func Yesterday() Carbon {
	return SetTimezone(Local).Yesterday()
}

// CreateFromTimestamp 从时间戳创建Carbon实例
func (c Carbon) CreateFromTimestamp(timestamp int64) Carbon {
	ts := timestamp
	switch len(strconv.FormatInt(timestamp, 10)) {
	case 10:
		ts = timestamp
	case 13:
		ts = timestamp / 1e3
	case 16:
		ts = timestamp / 1e6
	case 19:
		ts = timestamp / 1e9
	default:
		ts = 0
	}
	c.Time = time.Unix(ts, 0)
	return c
}

// CreateFromTimestamp 从时间戳创建Carbon实例(默认时区)
func CreateFromTimestamp(timestamp int64) Carbon {
	return SetTimezone(Local).CreateFromTimestamp(timestamp)
}

// CreateFromDateTime 从年月日时分秒创建Carbon实例
func (c Carbon) CreateFromDateTime(year int, month int, day int, hour int, minute int, second int) Carbon {
	c.Time = time.Date(year, time.Month(month), day, hour, minute, second, 0, c.Loc)
	return c
}

// CreateFromDateTime 从年月日时分秒创建Carbon实例(默认时区)
func CreateFromDateTime(year int, month int, day int, hour int, minute int, second int) Carbon {
	return SetTimezone(Local).CreateFromDateTime(year, month, day, hour, minute, second)
}

// CreateFromDate 从年月日创建Carbon实例
func (c Carbon) CreateFromDate(year int, month int, day int) Carbon {
	hour, minute, second := time.Now().Clock()
	c.Time = time.Date(year, time.Month(month), day, hour, minute, second, 0, c.Loc)
	return c
}

// CreateFromDate 从年月日创建Carbon实例(默认时区)
func CreateFromDate(year int, month int, day int) Carbon {
	return SetTimezone(Local).CreateFromDate(year, month, day)
}

// CreateFromTime 从时分秒创建Carbon实例
func (c Carbon) CreateFromTime(hour int, minute int, second int) Carbon {
	year, month, day := time.Now().Date()
	c.Time = time.Date(year, month, day, hour, minute, second, 0, c.Loc)
	return c
}

// CreateFromTime 从时分秒创建Carbon实例(默认时区)
func CreateFromTime(hour int, minute int, second int) Carbon {
	return SetTimezone(Local).CreateFromTime(hour, minute, second)
}

// Parse 将标准格式时间字符串解析成 Carbon 实例
func (c Carbon) Parse(value string) Carbon {
	if c.Error != nil {
		return c
	}

	layout := DateTimeFormat

	if value == "" || value == "0" || value == "0000-00-00 00:00:00" || value == "0000-00-00" || value == "00:00:00" {
		return c
	}

	if len(value) == 10 && strings.Count(value, "-") == 2 {
		layout = DateFormat
	}

	if strings.Index(value, "T") == 10 {
		layout = RFC3339Format
	}

	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		switch len(value) {
		case 8:
			layout = ShortDateFormat
		case 14:
			layout = ShortDateTimeFormat
		}
	}

	tt, err := parseByLayout(value, layout)
	c.Time = tt
	c.Error = err
	return c
}

// Parse 将标准格式时间字符串解析成 Carbon 实例(默认时区)
func Parse(value string) Carbon {
	return SetTimezone(Local).Parse(value)
}

// ParseByFormat 将特殊格式时间字符串解析成 Carbon 实例
func (c Carbon) ParseByFormat(value string, format string) Carbon {
	if c.Error != nil {
		return c
	}
	layout := format2layout(format)
	return c.ParseByLayout(value, layout)
}

// ParseByFormat 将特殊格式时间字符串解析成 Carbon 实例(默认时区)
func ParseByFormat(value string, format string) Carbon {
	return SetTimezone(Local).ParseByFormat(value, format)
}

// ParseByLayout 将布局时间字符串解析成 Carbon 实例
func (c Carbon) ParseByLayout(value string, layout string) Carbon {
	if c.Error != nil {
		return c
	}
	tt, err := parseByLayout(value, layout)
	c.Time = tt
	c.Error = err
	return c
}

// ParseByLayout 将布局时间字符串解析成 Carbon 实例(默认时区)
func ParseByLayout(value string, layout string) Carbon {
	return SetTimezone(Local).ParseByLayout(value, layout)
}

// Time2Carbon 将 time.Time 转换成 Carbon
func Time2Carbon(tt time.Time) Carbon {
	loc, _ := time.LoadLocation(Local)
	return Carbon{Time: tt, Loc: loc}
}

// Carbon2Time 将 Carbon 转换成 time.Time
func (c Carbon) Carbon2Time() time.Time {
	return c.Time
}

// AddDurations 按照持续时间字符串增加时间
// 支持整数/浮点数和符号ns(纳秒)、us(微妙)、ms(毫秒)、s(秒)、m(分钟)、h(小时)的组合
func (c Carbon) AddDuration(duration string) Carbon {
	if c.Error != nil {
		return c
	}
	td, err := parseByDuration(duration)
	c.Time = c.Time.Add(td)
	c.Error = err
	return c
}

// SubDurations 按照持续时间字符串减少时间
// 支持整数/浮点数和符号ns(纳秒)、us(微妙)、ms(毫秒)、s(秒)、m(分钟)、h(小时)的组合
func (c Carbon) SubDuration(duration string) Carbon {
	if c.Error != nil {
		return c
	}
	td, err := parseByDuration("-" + duration)
	c.Time = c.Time.Add(td)
	c.Error = err
	return c
}

// AddCenturies N世纪后
func (c Carbon) AddCenturies(centuries int) Carbon {
	return c.AddYears(YearsPerCentury * centuries)
}

// AddCenturiesNoOverflow N世纪后(月份不溢出)
func (c Carbon) AddCenturiesNoOverflow(centuries int) Carbon {
	return c.AddYearsNoOverflow(centuries * YearsPerCentury)
}

// AddCentury 1世纪后
func (c Carbon) AddCentury() Carbon {
	return c.AddCenturies(1)
}

// AddCenturyNoOverflow 1世纪后(月份不溢出)
func (c Carbon) AddCenturyNoOverflow() Carbon {
	return c.AddCenturiesNoOverflow(1)
}

// SubCenturies N世纪前
func (c Carbon) SubCenturies(centuries int) Carbon {
	return c.SubYears(YearsPerCentury * centuries)
}

// SubCenturiesNoOverflow N世纪前(月份不溢出)
func (c Carbon) SubCenturiesNoOverflow(centuries int) Carbon {
	return c.SubYearsNoOverflow(centuries * YearsPerCentury)
}

// SubCentury 1世纪前
func (c Carbon) SubCentury() Carbon {
	return c.SubCenturies(1)
}

// SubCenturyNoOverflow 1世纪前(月份不溢出)
func (c Carbon) SubCenturyNoOverflow() Carbon {
	return c.SubCenturiesNoOverflow(1)
}

// AddYears N年后
func (c Carbon) AddYears(years int) Carbon {
	c.Time = c.Time.AddDate(years, 0, 0)
	return c
}

// AddYearsNoOverflow N年后(月份不溢出)
func (c Carbon) AddYearsNoOverflow(years int) Carbon {
	year := c.Time.Year() + years
	month := c.Time.Month()
	day := c.Time.Day()

	// 获取N年后本月的最后一天
	last := time.Date(year, month, 1, c.Time.Hour(), c.Time.Minute(), c.Time.Second(), c.Time.Nanosecond(), c.Loc).AddDate(0, 1, -1)

	if day > last.Day() {
		day = last.Day()
	}

	c.Time = time.Date(last.Year(), last.Month(), day, c.Time.Hour(), c.Time.Minute(), c.Time.Second(), c.Time.Nanosecond(), c.Loc)
	return c
}

// AddYear 1年后
func (c Carbon) AddYear() Carbon {
	return c.AddYears(1)
}

// AddYearNoOverflow 1年后(月份不溢出)
func (c Carbon) AddYearNoOverflow() Carbon {
	return c.AddYearsNoOverflow(1)
}

// SubYears N年前
func (c Carbon) SubYears(years int) Carbon {
	c.Time = c.Time.AddDate(-years, 0, 0)
	return c
}

// SubYearsWithOverflow N年前(月份不溢出)
func (c Carbon) SubYearsNoOverflow(years int) Carbon {
	return c.AddYearsNoOverflow(-years)
}

// SubYear 1年前
func (c Carbon) SubYear() Carbon {
	return c.SubYears(1)
}

// SubYearWithOverflow 1年前(月份不溢出)
func (c Carbon) SubYearNoOverflow() Carbon {
	return c.SubYearsNoOverflow(1)
}

// AddQuarters N季度后
func (c Carbon) AddQuarters(quarters int) Carbon {
	return c.AddMonths(quarters * MonthsPerQuarter)
}

// AddQuartersNoOverflow N季度后(月份不溢出)
func (c Carbon) AddQuartersNoOverflow(quarters int) Carbon {
	return c.AddMonthsNoOverflow(quarters * MonthsPerQuarter)
}

// AddQuarter 1季度后
func (c Carbon) AddQuarter() Carbon {
	return c.AddQuarters(1)
}

// NextQuarters 1季度后(月份不溢出)
func (c Carbon) AddQuarterNoOverflow() Carbon {
	return c.AddQuartersNoOverflow(1)
}

// SubQuarters N季度前
func (c Carbon) SubQuarters(quarters int) Carbon {
	return c.SubMonths(quarters * MonthsPerQuarter)
}

// SubQuartersNoOverflow N季度前(月份不溢出)
func (c Carbon) SubQuartersNoOverflow(quarters int) Carbon {
	return c.AddMonthsNoOverflow(-quarters * MonthsPerQuarter)
}

// SubQuarter 1季度前
func (c Carbon) SubQuarter() Carbon {
	return c.SubQuarters(1)
}

// SubQuarterNoOverflow 1季度前(月份不溢出)
func (c Carbon) SubQuarterNoOverflow() Carbon {
	return c.SubQuartersNoOverflow(1)
}

// AddMonths N月后
func (c Carbon) AddMonths(months int) Carbon {
	c.Time = c.Time.AddDate(0, months, 0)
	return c
}

// AddMonthsNoOverflow N月后(月份不溢出)
func (c Carbon) AddMonthsNoOverflow(months int) Carbon {
	year := c.Time.Year()
	month := c.Time.Month() + time.Month(months)
	day := c.Time.Day()

	// 获取N月后的最后一天
	last := time.Date(year, month, 1, c.Time.Hour(), c.Time.Minute(), c.Time.Second(), c.Time.Nanosecond(), c.Loc).AddDate(0, 1, -1)

	if day > last.Day() {
		day = last.Day()
	}

	c.Time = time.Date(last.Year(), last.Month(), day, c.Time.Hour(), c.Time.Minute(), c.Time.Second(), c.Time.Nanosecond(), c.Loc)
	return c
}

// AddMonth 1月后
func (c Carbon) AddMonth() Carbon {
	return c.AddMonths(1)
}

// AddMonthNoOverflow 1月后(月份不溢出)
func (c Carbon) AddMonthNoOverflow() Carbon {
	return c.AddMonthsNoOverflow(1)
}

// SubMonths N月前
func (c Carbon) SubMonths(months int) Carbon {
	c.Time = c.Time.AddDate(0, -months, 0)
	return c
}

// SubMonthsNoOverflow N月前(月份不溢出)
func (c Carbon) SubMonthsNoOverflow(months int) Carbon {
	return c.AddMonthsNoOverflow(-months)
}

// SubMonth 1月前
func (c Carbon) SubMonth() Carbon {
	return c.SubMonths(1)
}

// SubMonthNoOverflow 1月前(月份不溢出)
func (c Carbon) SubMonthNoOverflow() Carbon {
	return c.SubMonthsNoOverflow(1)
}

// AddWeeks N周后
func (c Carbon) AddWeeks(weeks int) Carbon {
	return c.AddDays(weeks * DaysPerWeek)
}

// AddWeek 1天后
func (c Carbon) AddWeek() Carbon {
	return c.AddWeeks(1)
}

// SubWeeks N周后
func (c Carbon) SubWeeks(weeks int) Carbon {
	return c.SubDays(weeks * DaysPerWeek)
}

// SubWeek 1天后
func (c Carbon) SubWeek() Carbon {
	return c.SubWeeks(1)
}

// AddDays N天后
func (c Carbon) AddDays(days int) Carbon {
	c.Time = c.Time.AddDate(0, 0, days)
	return c
}

// AddDay 1天后
func (c Carbon) AddDay() Carbon {
	return c.AddDays(1)
}

// SubDays N天前
func (c Carbon) SubDays(days int) Carbon {
	c.Time = c.Time.AddDate(0, 0, -days)
	return c
}

// SubDay 1天前
func (c Carbon) SubDay() Carbon {
	return c.SubDays(1)
}

// AddHours N小时后
func (c Carbon) AddHours(hours int) Carbon {
	td := time.Duration(hours) * time.Hour
	c.Time = c.Time.Add(td)
	return c
}

// AddHour 1小时后
func (c Carbon) AddHour() Carbon {
	return c.AddHours(1)
}

// SubHours N小时前
func (c Carbon) SubHours(hours int) Carbon {
	td := time.Duration(hours) * -time.Hour
	c.Time = c.Time.Add(td)
	return c
}

// SubHour 1小时前
func (c Carbon) SubHour() Carbon {
	return c.SubHours(1)
}

// AddMinutes N分钟后
func (c Carbon) AddMinutes(minutes int) Carbon {
	td := time.Duration(minutes) * time.Minute
	c.Time = c.Time.Add(td)
	return c
}

// AddMinute 1分钟后
func (c Carbon) AddMinute() Carbon {
	return c.AddMinutes(1)
}

// SubMinutes N分钟前
func (c Carbon) SubMinutes(minutes int) Carbon {
	td := time.Duration(minutes) * -time.Minute
	c.Time = c.Time.Add(td)
	return c
}

// SubMinute 1分钟前
func (c Carbon) SubMinute() Carbon {
	return c.SubMinutes(1)
}

// AddSeconds N秒钟后
func (c Carbon) AddSeconds(seconds int) Carbon {
	td := time.Duration(seconds) * time.Second
	c.Time = c.Time.Add(td)
	return c
}

// AddSecond 1秒钟后
func (c Carbon) AddSecond() Carbon {
	return c.AddSeconds(1)
}

// SubSeconds N秒钟前
func (c Carbon) SubSeconds(seconds int) Carbon {
	td := time.Duration(seconds) * -time.Second
	c.Time = c.Time.Add(td)
	return c
}

// SubSecond 1秒钟前
func (c Carbon) SubSecond() Carbon {
	return c.SubSeconds(1)
}

// StartOfYear 本年开始时间
func (c Carbon) StartOfYear() Carbon {
	c.Time = time.Date(c.Time.Year(), 1, 1, 0, 0, 0, 0, c.Loc)
	return c
}

// EndOfYear 本年结束时间
func (c Carbon) EndOfYear() Carbon {
	c.Time = time.Date(c.Time.Year(), 12, 31, 23, 59, 59, 0, c.Loc)
	return c
}

// StartOfMonth 本月开始时间
func (c Carbon) StartOfMonth() Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), 1, 0, 0, 0, 0, c.Loc)
	return c
}

// EndOfMonth 本月结束时间
func (c Carbon) EndOfMonth() Carbon {
	t := time.Date(c.Time.Year(), c.Time.Month(), 1, 23, 59, 59, 0, c.Loc)
	c.Time = t.AddDate(0, 1, -1)
	return c
}

// StartOfWeek 本周开始时间
func (c Carbon) StartOfWeek() Carbon {
	days := c.Time.Weekday()
	if days == 0 {
		days = DaysPerWeek
	}
	t := time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), 0, 0, 0, 0, c.Loc)
	c.Time = t.AddDate(0, 0, int(1-days))
	return c
}

// EndOfWeek 本周结束时间
func (c Carbon) EndOfWeek() Carbon {
	days := c.Time.Weekday()
	if days == 0 {
		days = DaysPerWeek
	}
	t := time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), 23, 59, 59, 0, c.Loc)
	c.Time = t.AddDate(0, 0, int(DaysPerWeek-days))
	return c
}

// StartOfDay 本日开始时间
func (c Carbon) StartOfDay() Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), 0, 0, 0, 0, c.Loc)
	return c
}

// EndOfDay 本日结束时间
func (c Carbon) EndOfDay() Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), 23, 59, 59, 0, c.Loc)
	return c
}

// StartOfHour 小时开始时间
func (c Carbon) StartOfHour() Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), c.Time.Hour(), 0, 0, 0, c.Loc)
	return c
}

// EndOfHour 小时结束时间
func (c Carbon) EndOfHour() Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), c.Time.Hour(), 59, 59, 0, c.Loc)
	return c
}

// StartOfMinute 分钟开始时间
func (c Carbon) StartOfMinute() Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), c.Time.Hour(), c.Time.Minute(), 0, 0, c.Loc)
	return c
}

// EndOfMinute 分钟结束时间
func (c Carbon) EndOfMinute() Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), c.Time.Hour(), c.Time.Minute(), 59, 0, c.Loc)
	return c
}

// Timezone 设置时区
func SetTimezone(name string) Carbon {
	loc, err := getLocationByTimezone(name)
	return Carbon{Loc: loc, Error: err}
}

// Timezone 设置时区
func (c Carbon) SetTimezone(name string) Carbon {
	loc, err := getLocationByTimezone(name)
	c.Loc = loc
	c.Error = err
	return c
}

// SetYear 设置年
func (c Carbon) SetYear(year int) Carbon {
	c.Time = time.Date(year, c.Time.Month(), c.Time.Day(), c.Time.Hour(), c.Time.Minute(), c.Time.Second(), c.Time.Nanosecond(), c.Loc)
	return c
}

// SetMonth 设置月
func (c Carbon) SetMonth(month int) Carbon {
	c.Time = time.Date(c.Time.Year(), time.Month(month), c.Time.Day(), c.Time.Hour(), c.Time.Minute(), c.Time.Second(), c.Time.Nanosecond(), c.Loc)
	return c
}

// SetDay 设置日
func (c Carbon) SetDay(day int) Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), day, c.Time.Hour(), c.Time.Minute(), c.Time.Second(), c.Time.Nanosecond(), c.Loc)
	return c
}

// SetHour 设置时
func (c Carbon) SetHour(hour int) Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), hour, c.Time.Minute(), c.Time.Second(), c.Time.Nanosecond(), c.Loc)
	return c
}

// SetMinute 设置分
func (c Carbon) SetMinute(minute int) Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), c.Time.Hour(), minute, c.Time.Second(), c.Time.Nanosecond(), c.Loc)
	return c
}

// SetSecond 设置秒
func (c Carbon) SetSecond(second int) Carbon {
	c.Time = time.Date(c.Time.Year(), c.Time.Month(), c.Time.Day(), c.Time.Hour(), c.Time.Minute(), second, c.Time.Nanosecond(), c.Loc)
	return c
}
