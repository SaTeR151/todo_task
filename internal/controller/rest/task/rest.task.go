package task

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/controller/rest/task/validation"
	"github.com/sater-151/todo-list/internal/entity"
)

func (c *TaskController) POST(ctx *gin.Context) {
	boardID := ctx.MustGet("board").(string)
	userID := ctx.MustGet("user").(string)

	var taskPOST dto.TaskPOST
	if err := ctx.ShouldBindJSON(&taskPOST); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	columns, err := c.columns.GetByBoardID(ctx, boardID)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	types, err := c.types.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	if err := validation.ValidateTaskCreate(taskPOST, columns, types); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	taskCreate := entity.TaskCreate{
		Label:       taskPOST.Label,
		Description: taskPOST.Description,
		ColumnID:    taskPOST.ColumnID,
		TypeID:      taskPOST.TypeID,
	}

	newTask, err := c.tasks.Create(ctx, boardID, taskCreate)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, newTask)
}

func (c *TaskController) LIST(ctx *gin.Context) {
	boardID := ctx.MustGet("board").(string)

	var taskGETQuery dto.TaskGETQuery
	if err := ctx.ShouldBindQuery(&taskGETQuery); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	taskGetOpts := entity.GetTasksOpts{
		ColumnID: taskGETQuery.ColumnID,
		TypeID:   taskGETQuery.TypeID,
	}

	tasks, err := c.tasks.Get(ctx, boardID, taskGetOpts)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(200, entity.Tasks{})
			return
		}
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, tasks)
}

func (c *TaskController) GET(ctx *gin.Context) {
	boardID := ctx.MustGet("board").(string)

	var taskGETuri dto.TaskGETUri
	if err := ctx.ShouldBindUri(&taskGETuri); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	task, err := c.tasks.GetByID(ctx, boardID, taskGETuri.TaskID)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(404, err.Error())
			return
		}
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, task)
}

func (c *TaskController) DELETE(ctx *gin.Context) {
	var taskDELETEQuery dto.TaskDELETEUri
	if err := ctx.ShouldBindUri(&taskDELETEQuery); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	err := c.tasks.Delete(ctx, taskDELETEQuery.TaskID)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(404, err.Error())
			return
		}
		ctx.JSON(500, err.Error())
		return
	}

	ctx.Status(200)
}

func (c *TaskController) MOVE(ctx *gin.Context) {
	boardID := ctx.MustGet("board").(string)

	var taskMOVEQuery dto.TaskMOVEUri
	if err := ctx.ShouldBindUri(&taskMOVEQuery); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	var taskMOVE dto.TaskMOVE
	if err := ctx.ShouldBindJSON(&taskMOVE); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	task, err := c.tasks.Move(ctx, boardID, taskMOVEQuery.TaskID, taskMOVE.ColumnDestenation)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(404, err.Error())
			return
		}
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, task)
}

func (c *TaskController) PATCH(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)
	boardID := ctx.MustGet("board").(string)

	var taskPATCHUri dto.TaskPATCHUri
	if err := ctx.ShouldBindUri(&taskPATCHUri); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	var taskPATCH dto.TaskPATCH
	if err := ctx.ShouldBindJSON(&taskPATCH); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	types, err := c.types.GetByUserID(ctx, userID)
	if err != nil && err != entity.ErrNotFound {
		ctx.JSON(500, err.Error())
		return
	}

	typeIDs := types.GetIDs()

	if err := validation.ValidateTaskUpdate(taskPATCH, typeIDs); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	taskUpdate := entity.TaskUpdate{
		ID:          taskPATCHUri.TaskID,
		Label:       taskPATCH.Label,
		Description: taskPATCH.Description,
		TypeID:      taskPATCH.TypeID,
	}

	task, err := c.tasks.Update(ctx, boardID, taskUpdate)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, task)
}
