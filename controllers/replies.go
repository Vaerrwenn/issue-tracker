package controllers

import (
	"issue-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReplyCreateUpdateForm is a struct for Form binding on Create and Update operation.
type ReplyCreateUpdateForm struct {
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

	var input ReplyCreateUpdateForm
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
		"replyID": reply.ID,
		"msg":     "Reply added successfully!",
	})
	return
}

// UpdateReplyHandler handles the Update request.
func UpdateReplyHandler(c *gin.Context) {
	userID, err := strconv.ParseUint(c.GetHeader("userID"), 10, 0)
	if err != nil {
		returnErrorAndAbort(c, http.StatusForbidden, "User invalid.")
		return
	}

	replyID, err := strconv.Atoi(c.Param("replyId"))
	if err != nil {
		returnErrorAndAbort(c, http.StatusNotFound, "No reply ID provided.")
		return
	}

	var reply models.Reply
	replySouce := reply.FindReplyByID(uint(replyID))
	if replySouce == nil {
		returnErrorAndAbort(c, http.StatusNotFound, "Reply not found.")
		return
	}

	if userID != uint64(replySouce.UserID) {
		returnErrorAndAbort(c, http.StatusForbidden, "This user is not allowed to update this Reply.")
		return
	}

	var input ReplyCreateUpdateForm
	if err := c.ShouldBind(&input); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	updateReply := models.Reply{
		Body: input.Body,
	}

	err = updateReply.UpdateReply(replySouce)
	if err != nil {
		returnErrorAndAbort(c, http.StatusNotAcceptable, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Data successfully updated.",
	})
	return
}
