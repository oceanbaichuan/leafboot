package gamelogic

// func (f *FactoryGameLogic) WriteLoginRoomLog(loginlog *msg.Logingamelog) {
// 	log.BuringPoint("WriteLoginRoomLog", zap.Object("logingame", loginlog))
// }

// func (f *FactoryGameLogic) WriteLeaveRoomLog(leavelog *msg.Leavegamelog) {
// 	log.BuringPoint("WriteLeaveRoomLog", zap.Object("leavegame", leavelog))
// }
// func (f *FactoryGameLogic) WriteTableRoundLog(playlog *msg.Playgamelog) {
// 	log.BuringPoint("WriteTableRoundLog", zap.Object("tableroundlog", playlog))
// }
// func (f *FactoryGameLogic) WriteAttributionLog(attrlog *msg.AttributeChangelog) {
// 	log.BuringPoint("WriteAttributionLog", zap.Object("propchangelog", attrlog))
// }

func (f *FactoryGameLogic) WriteLoginRoomLog(loginlog interface{}) {
	//log.BuringPoint("WriteLoginRoomLog", zap.Object("logingame", loginlog))
}

func (f *FactoryGameLogic) WriteLeaveRoomLog(leavelog interface{}) {
	//log.BuringPoint("WriteLeaveRoomLog", zap.Object("leavegame", leavelog))
}
func (f *FactoryGameLogic) WriteTableRoundLog(playlog interface{}) {
	//log.BuringPoint("WriteTableRoundLog", zap.Object("tableroundlog", playlog))
}
func (f *FactoryGameLogic) WriteAttributionLog(attrlog interface{}) {
	//log.BuringPoint("WriteAttributionLog", zap.Object("propchangelog", attrlog))
}
