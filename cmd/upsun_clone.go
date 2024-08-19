package main

import (
	"log"

	flag "github.com/spf13/pflag"
	app "github.com/upsun/clonsun"
	logic "github.com/upsun/clonsun/internal/logic"
	entity "github.com/upsun/lib-sun/entity"
	utils "github.com/upsun/lib-sun/utility"
)

const (
	APP_NAME = "clonsun"

	ARG_SRC_CLI = "src_provider"
	ARG_SRC_PRJ = "src_project"
	ARG_SRC_ENV = "src_env"
	ARG_DST_CLI = "dst_provider"
	ARG_DST_PRJ = "dst_project"
	ARG_DST_ENV = "dst_env"
	ARG_DST_ORG = "dst_organisation"
	ARG_DST_RGN = "dst_region"

	ARG_KEEP_DATA     = "keep_data"
	ARG_FORCE_REPO    = "force_psh_repo"
	ARG_RMT_ONLY      = "remote_only"
	ARG_VAR_SENSITIVE = "clone_sensitive"
)

func init() {
	flag.StringVarP(&app.ArgsM.SrcProvider, ARG_SRC_CLI, "", entity.PSH_PROVIDER, "Source provider CLI")
	flag.StringVarP(&app.ArgsM.SrcProjectID, ARG_SRC_PRJ, "", "", "Source Project ID")
	flag.StringVarP(&app.ArgsM.SrcEnvironment, ARG_SRC_ENV, "", "main", "Source Environment")
	flag.StringVarP(&app.ArgsM.DstProvider, ARG_DST_CLI, "", entity.UPS_PROVIDER, "Destination provider CLI")
	flag.StringVarP(&app.ArgsM.DstProjectID, ARG_DST_PRJ, "", "", "Destination Project ID")
	flag.StringVarP(&app.ArgsM.DstEnvironment, ARG_DST_ENV, "", "main", "Destination Environment")
	flag.StringVarP(&app.ArgsM.DstOrga, ARG_DST_ORG, "", "", "Destination Organisation")
	flag.StringVarP(&app.ArgsM.DstRegion, ARG_DST_RGN, "", "fr-4.platform.sh", "Destination Region")
	flag.BoolVarP(&app.ArgsM.NoData, "no_data", "", false, "Do not clone data (databases & mounts)")
	flag.BoolVarP(&app.ArgsM.OnlyData, "only_data", "", false, "Clone only data (databases & mounts)")
	flag.BoolVarP(&app.ArgsM.OnlyDb, "only_db", "", false, "Clone only databases")
	flag.BoolVarP(&app.ArgsM.OnlyMount, "only_mount", "", false, "Clone only mounts")
	flag.BoolVarP(&app.ArgsM.NoUsers, "no_users", "", false, "Do not clone user.")
	flag.BoolVarP(&app.ArgsM.Sensitive, ARG_VAR_SENSITIVE, "", false, "Clone sensitive value (variables)")
	flag.StringVarP(&app.ArgsC.TypeMount, "mount_type", "", "storage", "Change 'Local' mount to upsun compatible mode : 'storage' or 'instance'.")
	flag.StringVarP(&app.ArgsM.KeepData, ARG_KEEP_DATA, "", "", "Path where to keep all dumped data (WIP)")
	flag.BoolVarP(&app.ArgsM.PshRepo, ARG_FORCE_REPO, "", false, "Force to use provider repository")
	flag.BoolVarP(&app.ArgsM.NoLocal, ARG_RMT_ONLY, "", false, "Direct transfert data from source to destination without tempory local transfert (Not Implemented)")
	flag.BoolVarP(&app.Args.Verbose, "verbose", "v", false, "Enable verbose mode")

	if err := flag.CommandLine.MarkHidden(ARG_KEEP_DATA); err != nil {
		log.Println("Bad hidden flag !")
	}
	if err := flag.CommandLine.MarkHidden(ARG_FORCE_REPO); err != nil {
		log.Println("Bad hidden flag !")
	}
	if err := flag.CommandLine.MarkHidden(ARG_RMT_ONLY); err != nil {
		log.Println("Bad hidden flag !")
	}
	if err := flag.CommandLine.MarkHidden(ARG_VAR_SENSITIVE); err != nil {
		log.Println("Bad hidden flag !")
	}
	flag.CommandLine.SortFlags = false
	flag.Parse()
}

func main() {
	utils.InitLogger(APP_NAME)
	utils.Disclaimer(APP_NAME)
	utils.StartReporters(APP_NAME)

	// FROM
	app.ArgsM.SrcProvider = utils.RequireFlag(
		ARG_SRC_CLI,
		"[Source] Enter the CLI to which you want to clone (choose between upsun or platform) [%v]: ",
		app.ArgsM.SrcProvider,
		false)
	utils.ProviderCheck(app.ArgsM.SrcProvider)

	app.ArgsM.SrcProjectID = utils.RequireFlag(
		ARG_SRC_PRJ,
		"[Source] Enter the Project ID from which you want to clone [%v]: ",
		app.ArgsM.SrcProjectID,
		true)

	app.ArgsM.SrcEnvironment = utils.RequireFlag(
		ARG_SRC_ENV,
		"[Source] Enter the Environment Name from which you want to clone [%v]: ",
		app.ArgsM.SrcEnvironment,
		false)

	// TO
	app.ArgsM.DstProvider = utils.RequireFlag(
		ARG_DST_CLI,
		"[Destination] Enter the CLI to which you want to clone (choose between upsun or platform) [%v]: ",
		app.ArgsM.DstProvider,
		false)
	utils.ProviderCheck(app.ArgsM.DstProvider)

	app.ArgsM.DstProjectID = utils.RequireFlag(
		ARG_DST_PRJ,
		"[Destination] Enter the Project ID to which you want to clone [%v] (if empty, create new project): ",
		app.ArgsM.DstProjectID,
		false)

	app.ArgsM.DstOrga = utils.RequireFlag(
		ARG_DST_ORG,
		"[Destination] Enter the Organisation to which you want to clone [%v]: ",
		app.ArgsM.DstOrga,
		false)

	app.ArgsM.DstRegion = utils.RequireFlag(
		ARG_DST_RGN,
		"[Destination] Enter the Region to which you want to clone [%v]: ",
		app.ArgsM.DstRegion,
		false)

	// Init
	projectFrom := entity.MakeProjectContext(
		app.ArgsM.SrcProvider,
		app.ArgsM.SrcProjectID,
		app.ArgsM.SrcEnvironment,
	)

	projectTo := entity.MakeProjectContext(
		app.ArgsM.DstProvider,
		app.ArgsM.DstProjectID,
		app.ArgsM.DstEnvironment,
	)
	projectTo.OrgEmail = app.ArgsM.DstOrga
	projectTo.Region = app.ArgsM.DstRegion

	// Assert
	if !utils.IsAuthenticated(projectFrom) {
		log.Fatalf("You are not authenticated, please run: %v login\n", projectFrom.Provider)
	}
	if !utils.IsAuthenticated(projectTo) {
		log.Fatalf("You are not authenticated, please run: %v login\n", projectTo.Provider)
	}
	//TODO utils.HasSufficientRights(projectTo.Provider)
	//TODO utils.HasSufficientRights(projectFrom.Provider)

	// Process
	logic.Clone(projectFrom, projectTo)
}
