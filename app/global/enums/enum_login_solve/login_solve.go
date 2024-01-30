package enum_login_solve

import "goskeleton/app/global/enums"

// LoginSolveConfig 口令破解配置
type LoginSolveConfig struct {
	Switch    int `json:"switch"`
	Gap       int `json:"gap"`
	Threshold int `json:"threshold"`
}

func (c *LoginSolveConfig) Format() {
	if c.Switch == 0 {
		c.Switch = enums.SwitchOff
	}
	if c.Gap == 0 {
		c.Gap = 3
	}
	if c.Threshold == 0 {
		c.Threshold = 5
	}
}
