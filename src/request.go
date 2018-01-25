package main

type Request struct {
	Instance  string  `json:"instance"`
	IdApp       int  `json:"idApp"`
	IdUser      int  `json:"idUser"`
}

type Requests []Request
