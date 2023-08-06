package models

/*

*/

type InvokeCloneOnDemandForm struct {
    Template string `json:"template" binding:"required"`
    SessionKey string `json:"jwtToken" binding:"required"`
}

type LoginForm struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type AuthForm struct {
    SessionKey string `json:"jwtToken" binding:"required"`
}