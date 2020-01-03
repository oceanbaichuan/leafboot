package appmain

import (
	"github.com/hudgit2019/leafboot/msg"
)

/*主玩法通信示例*/
func init() {
	msg.Processor.Register(&GuessReq{})
}

type GuessStartNotice struct {
}
type GuessReq struct {
	Guesstype int8
}

type GuessRes struct {
	Guesstype int8
}

type GuessResult struct {
	Guesstype [2]int8
	Socres    [2]int
}
