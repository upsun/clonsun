package api

import (
	"github.com/upsun/clonsun/internal/logic"
	"github.com/upsun/lib-sun/entity"
)

func Clone(projectSrcContext entity.ProjectGlobal, projectDstContext entity.ProjectGlobal) {
	logic.Clone(projectSrcContext, projectDstContext)
}
