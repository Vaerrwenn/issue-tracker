package controllers

import (
	"issue-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReplyCreateForm is a struct for Form binding.
type ReplyCreateForm struct {
	Body string `form:"description" binding:"required"`
}

// CreateReplyHandler handles Reply creation
//
// Needs "id" as param and "userID" as Header.
func CreateReplyHandler(c *gin.Context) {
	issueID, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, "Issue invalid.")
		return
	}

	userID, err := strconv.ParseUint(c.GetHeader("userID"), 10, 0)
	if err != nil {
		returnErrorAndAbort(c, http.StatusForbidden, "User invalid.")
	}

	var input ReplyCreateForm
	if err := c.ShouldBind(&input); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	var issue models.Issue
	_, err = issue.FindIssueByID(strconv.Itoa(int(issueID)))
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	reply := models.Reply{
		UserID:  uint(userID),
		IssueID: uint(issueID),
		Body:    input.Body,
	}
	err = reply.SaveReply()
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Reply added successfully!",
	})
	return
}
