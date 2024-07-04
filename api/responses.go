package api

import "hotel-reservation/types"

type ResourceResponse struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

type GenericResponse struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

type AuthResponse struct {
	User *types.User `json:"user"`
	Token string `json:"token"`
}