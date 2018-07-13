package model

type Process struct {
	Id      int64  `json:"ID" xorm:"default 0 BIGINT(21)"`
	User    string `json:"USER" xorm:"not null default '' VARCHAR(32)"`
	Host    string `json:"HOST" xorm:"not null default '' VARCHAR(64)"`
	Db      string `json:"DB" xorm:"VARCHAR(64)"`
	Command string `json:"COMMAND" xorm:"not null default '' VARCHAR(16)"`
	Time    int    `json:"TIME" xorm:"not null default 0 INT(7)"`
	State   string `json:"STATE" xorm:"VARCHAR(64)"`
	Info    string `json:"INFO" xorm:"LONGTEXT"`
}

func (p Process) TableName() string {
	return "PROCESSLIST"
}
