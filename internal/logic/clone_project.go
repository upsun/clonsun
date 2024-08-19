package logic

import (
	"log"

	app "github.com/upsun/clonsun"
	converter "github.com/upsun/convsun/api"
	cmd "github.com/upsun/lib-sun/cmd"
	entity "github.com/upsun/lib-sun/entity"
	utils "github.com/upsun/lib-sun/utility"
)

func Clone(projectSrcContext entity.ProjectGlobal, projectDstContext entity.ProjectGlobal) {
	// Define workspace
	var ws utils.PathTmp
	if utils.IsKeep() {
		ws = utils.BuildPersistWorkspace()
	} else {
		ws = utils.BuildTemporyWorkspace()
		defer ws.CleanUp() // clean up on exit
	}

	// Collect Meta-model & Data
	err := cmd.ProjectRead(&projectSrcContext)
	if err == nil {
		cmd.VariablesRead(projectSrcContext)
		cmd.UsersRead(projectSrcContext)
		cmd.ServicesMountsRead(projectSrcContext)

		if !app.ArgsM.NoData && !app.ArgsM.NoLocal {
			if !app.ArgsM.OnlyMount {
				cmd.ServicesExport(projectSrcContext, ws)
			}
			if !app.ArgsM.OnlyDb {
				cmd.MountsExport(projectSrcContext, ws)
			}
		}

		// Inject Meta-model from source to destination.
		projectDstContext.Copy(projectSrcContext) // Sync meta-model
		projectDstContext.Name = projectDstContext.Name + " - " + utils.TimeStamp()

		// Initialize Dst project
		if !app.ArgsM.OnlyData {
			cmd.ProjectCreate(&projectDstContext)

			// Push Meta-model to Dst projecrt
			cmd.ProjectWrite(projectDstContext)
			cmd.VariablesWrite(projectDstContext)
			if !app.ArgsM.NoUsers {
				cmd.UsersWrite(projectDstContext)
			}

			// Push All (code base and data (service & mount))
			// Clone repo
			cmd.ImportRepository(projectSrcContext, ws.Repo)

			// Upsun convert
			if projectSrcContext.Provider != entity.UPS_PROVIDER && projectDstContext.Provider == entity.UPS_PROVIDER {
				upsunPath, upsunConfigPath := utils.MakeUpsunConfigPath(ws.Repo)

				if !utils.IsExist(upsunConfigPath) {
					log.Print("New upsun project detected, need to convert it !\n WARNING this make a new commit on your project...")
					converter.Convert(ws.Repo, upsunPath)
					cmd.CommitOnRepository(ws.Repo, upsunPath)
				} else {
					log.Printf("WARNING: Upsun config file already exist. Do nothing !")
				}
			}

			// Push to new repo
			cmd.ExportRepository(projectDstContext, ws.Repo)
		}

		//ImportDatabase()
		if !app.ArgsM.NoData {
			if !app.ArgsM.OnlyMount {
				cmd.ServicesImport(projectDstContext)
			}
			if !app.ArgsM.OnlyDb {
				cmd.MountsImport(projectDstContext)
			}
		}

		// Check & Warning
		log.Printf("Project %v is cloned into %v", projectSrcContext.ID, projectDstContext.ID)
		cmd.DisplaySensitiveVariables(projectSrcContext)
		cmd.DisplayUnsupportedServices(projectSrcContext)
		cmd.DisplayUnsupportedIntegrations(projectSrcContext)

		utils.LinkToProject(projectDstContext)
	}

}
