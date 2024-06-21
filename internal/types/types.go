// Code generated by goctl. DO NOT EDIT.
package types

type CreateRecord struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}

type Message struct {
	Code    int    `json:"code"`
	Context string `json:"context"`
}

type Record struct {
	Id     int    `json:"id"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Value  string `json:"value"`
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