package appmain

type GuessPlayerNode struct {
	//appgame property start here
	Guesstype int8
}

// //SyncOldPlayerData 同步
// func (g *GuessPlayerNode) SyncOldPlayerData(oldplayer base.IPlayerNode) {
// 	g.ClientNode.SyncOldPlayerData(oldplayer)
// 	guessPlayer := oldplayer.(*GuessPlayerNode)
// 	g.Guesstype = guessPlayer.Guesstype
// }
