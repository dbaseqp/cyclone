package models

/*

*/

type InvokeCloneOnDemandForm struct {
    Template string `json:"template" binding:"required"`
}