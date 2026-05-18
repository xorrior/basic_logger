package my_logger

import (
	"github.com/MythicMeta/MythicContainer/loggingstructs"
	"github.com/MythicMeta/MythicContainer/mythicrpc"
)

func Initialize() {
	myLoggerName := "my_logger"
	myLogger := loggingstructs.LoggingDefinition{
		Name:           myLoggerName,
		Description:    "basic stdout debug logger",
		LogToFilePath:  "mythic.log",
		LogLevel:       "debug",
		LogMaxSizeInMB: 20,
		//LogMaxBackups:  10,
		NewCallbackFunction: func(input loggingstructs.NewCallbackLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input)
		},
		NewTaskFunction: func(input loggingstructs.NewTaskLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewPayloadFunction: func(input loggingstructs.NewPayloadLog) {
			logger := loggingstructs.AllLoggingData.Get(myLoggerName)
			logger.LogInfo(input.Action, "data", input.Data)

			if input.Data.BuildPhase != "success" {
				return
			}

			searchResp, err := mythicrpc.SendMythicRPCPayloadSearch(mythicrpc.MythicRPCPayloadSearchMessage{
				PayloadUUID: input.Data.UUID,
			})
			if err != nil || !searchResp.Success || len(searchResp.PayloadConfigurations) == 0 {
				return
			}

			config := searchResp.PayloadConfigurations[0]

			if config.C2Profiles != nil {
				for _, c2 := range *config.C2Profiles {
					logger.LogInfo("payload_c2_config",
						"uuid", input.Data.UUID,
						"c2_profile", c2.Name,
						"parameters", c2.Parameters,
					)
				}
			}

			if config.AgentFileID != "" {
				fileResp, err := mythicrpc.SendMythicRPCFileSearch(mythicrpc.MythicRPCFileSearchMessage{
					AgentFileID: config.AgentFileID,
					IsPayload:   true,
				})
				if err == nil && fileResp.Success && len(fileResp.Files) > 0 {
					logger.LogInfo("payload_hashes",
						"uuid", input.Data.UUID,
						"md5", fileResp.Files[0].Md5,
						"sha1", fileResp.Files[0].Sha1,
					)
				}
			}
		},
		NewKeylogFunction: func(input loggingstructs.NewKeylogLog) {
			//loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewCredentialFunction: func(input loggingstructs.NewCredentialLog) {
			//loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewArtifactFunction: func(input loggingstructs.NewArtifactLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewFileFunction: func(input loggingstructs.NewFileLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewResponseFunction: func(input loggingstructs.NewResponseLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
	}
	loggingstructs.AllLoggingData.Get(myLoggerName).AddLoggingDefinition(myLogger)
}
