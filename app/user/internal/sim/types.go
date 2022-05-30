package sim

import "time"

type simProvider = string // sim卡提供商
const (
	ChinaMobile  simProvider = "ChinaMobile"
	ChinaUnicom  simProvider = "ChinaUnicom"
	ChinaTelecom simProvider = "ChinaTelecom"
)

type Info struct {
	Id                     int
	SimNo                  string    `json:"sim_no"`
	Imsi                   string    `json:"imsi"`
	Iccid                  string    `json:"iccid"`
	Supplier               int       `json:"supplier"`
	TestPeriod             int       `json:"test_period"`
	QuietPeriod            int       `json:"quiet_period"`
	InnetDate              time.Time `json:"innet_date"`
	ActivatedTime          time.Time `json:"activated_time"`
	DefaultServiceDuration int       `json:"default_service_duration"`
	Status                 int       `json:"status"`
	Created                int64     `json:"created"`
	Updated                int64     `json:"updated"`
	Country                []string  `json:"country"`
	Remark                 string    `json:"remark"`
	Possess                int       `json:"possess"`
	LastActivated          time.Time `json:"last_activated"`
	Area                   int       `json:"area"`
	Service                int       `json:"service"`
}
