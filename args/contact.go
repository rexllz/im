package args

//frontend bind struct
type ContactArg struct {
	PageArg
	Userid int64 	`json:"userid" form:"userid"`
	Dstid int64 	`json:"dstid" form:"dstid"`
}
