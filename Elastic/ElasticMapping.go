package Elastic

import shared "okey101/Shared"

const (
	okey101LoginMappings = ` 
{ 
	"aliases":{
		"okey101_login": {}
	},
	"settings":{ 
		"number_of_shards": 3,
		"number_of_replicas": 1
	}, 
	"mappings":{ 
		"properties":{ 
			"PostDate": { "type": "date" }, 
			"UserID": { "type": "long" }, 
			"Username": { "type": "text" } 
		} 
	} 
}
`

	okey101AuditMappings = ` 
{ 
	"aliases":{
		"okey101_audit": {}
	},
	"settings":{ 
		"number_of_shards": 3,
		"number_of_replicas": 1
	}, 
	"mappings":{ 
		"properties":{ 
			"DateTime": { "type": "date" },
			"TimeStamp": { "type": "date", "format": "strict_date_optional_time||epoch_millis" },
			"OrderID": { "type": "long" },
			"UserName": { "type": "keyword" },
			"UserID": { "type": "long" },
			"ActionType": { "type": "integer" },
			"ActionName": { "type": "keyword" },
			"Message": { "type": "text" },
			"ModuleName": { "type": "keyword" },
			"GameID": { "type": "keyword" },
			"RoomID": { "type": "keyword" },
			"Tiles": { 
				"type": "nested", 
				"properties":{ 
					"ID": { "type": "integer" },
					"Number": { "type": "integer" },
					"Color": { "type": "integer" },
					"IsJoker": { "type": "boolean" },
					"IsOkey": { "type": "boolean" }
				}
			},
			"PenaltyReasonID": { "type": "integer" },
      		"PenaltyReason": { "type": "keyword" },
      		"PenaltyMultiplier": { "type": "float" },
      		"PenaltyPoints": { "type": "integer" },
			"HadOkeyTile": { "type": "boolean" },
      		"OpenedFivePairsButLost": { "type": "boolean" },
      		"OkeyUsedInFinish": { "type": "boolean" },
			"ReconnectDelaySeconds": { "type": "float" },
			"GameDurationSeconds": { "type": "float" },
            "PlayerReactionTimeSeconds": { "type": "float" }, 
			"IPAddress": { "type": "ip" },
			"Browser": { "type": "keyword" },
			"Device": { "type": "keyword" },
			"Platform": { "type": "keyword" },
			"ErrorCode": { "type": "integer" },
			"ExtraData": { "type": "object", "enabled": true }
		} 
	} 
}
`

	okey101ErrorMappings = ` 
{ 
	"aliases":{
		"okey101_error": {}
	},
	"settings":{ 
		"number_of_shards": 3,
		"number_of_replicas": 1
	}, 
	"mappings":{ 
		"properties":{ 
			"DateTime": { "type": "date" }, 
			"UserName": { "type": "text" }, 
			"Message": { "type": "text" },
			"ModuleName": { "type": "text" }, 
			"ActionName": { "type": "text" },
			"ErrorCode": { "type": "integer" } 
		} 
	} 
}
`
)

var ElasticMaps = map[string]string{
	shared.Config.ELASTICLOGININDEX: okey101LoginMappings,
	shared.Config.ELASTICAUDITINDEX: okey101AuditMappings,
	shared.Config.ELASTICERRORINDEX: okey101ErrorMappings,
}
