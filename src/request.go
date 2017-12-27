package main

type Request struct {
	InstanceID  string  `json:"instanceID"`
	AppID       int  `json:"appID"`
	UserID      int  `json:"userID"`
}

type Requests []Request
