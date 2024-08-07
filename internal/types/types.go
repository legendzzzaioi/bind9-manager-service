// Code generated by goctl. DO NOT EDIT.
package types

type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateRecord struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}

type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type GetUserResp struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	CreateAt string `json:"created_at"`
	UpdateAt string `json:"updated_at"`
}

type LoginLog struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Ip        string `json:"ip"`
	Operation string `json:"operation"`
	CreateAt  string `json:"created_at"`
}

type LoginLogResp struct {
	Logs  []LoginLog `json:"logs"`
	Total int        `json:"total"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type Message struct {
	Code    int    `json:"code"`
	Context string `json:"context"`
}

type OperationLog struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Operation string `json:"operation"`
	Context   string `json:"context"`
	CreateAt  string `json:"created_at"`
}

type OperationLogResp struct {
	Logs  []OperationLog `json:"logs"`
	Total int            `json:"total"`
}

type Record struct {
	Id     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}

type UpdateUserPassReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserLog struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Operation string `json:"operation"`
	Context   string `json:"context"`
	CreateAt  string `json:"created_at"`
}

type UserLogResp struct {
	Logs  []UserLog `json:"logs"`
	Total int       `json:"total"`
}

type Zone struct {
	Domain            string `json:"domain"`
	Ttl               int    `json:"ttl"`
	CacheTtl          int    `json:"cache_ttl"`
	Expire            int    `json:"expire"`
	MailAddress       string `json:"mail_address"`
	PrimaryNameServer string `json:"primary_name_server"`
	Refresh           int    `json:"refresh"`
	Retry             int    `json:"retry"`
	Serial            int    `json:"serial"`
}

type ZoneReq struct {
	Domain            string `json:"domain"`
	Ttl               int    `json:"ttl"`
	CacheTtl          int    `json:"cache_ttl"`
	Expire            int    `json:"expire"`
	MailAddress       string `json:"mail_address"`
	PrimaryNameServer string `json:"primary_name_server"`
	Refresh           int    `json:"refresh"`
	Retry             int    `json:"retry"`
}
