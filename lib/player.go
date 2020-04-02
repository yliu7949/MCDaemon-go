package lib

type Player interface {
	GetPosition(playerName string) (string, string, string)     //获取玩家当前的坐标
	GetDim(playerName string) string 							//获取玩家所在的维度，"0":"主世界","-1":"地狱","1":"末地"

}
