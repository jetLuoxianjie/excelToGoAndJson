package config
import (
    "encoding/json"
	"io/ioutil"
	"sync"
)

// 调用初始化读取配置
func InitConfig() {
  initGiftConfig()
  initGift2Config()
  initBuffConfig()

}

type GiftConfig struct {
	Lobby_GiftDesML int32 `json:"Lobby_GiftDesML"`
	GiftNum []int32 `json:"GiftNum"`
	GetSoldierNum []int32 `json:"GetSoldierNum"`
	RescueInjuryNum []int32 `json:"RescueInjuryNum"`
	Id int32 `json:"Id"`
}

var GiftListTable []*GiftConfig
var GiftListMap sync.Map
// 解析JSON数据到结构体
func initGiftConfig(){
  fileContent, err := ioutil.ReadFile("./jsons/Gift.json")
	if err != nil {
        panic("Gift Error reading file:" + err.Error())
		return
	}
	err = json.Unmarshal(fileContent, &GiftListTable)
	if err != nil {
		panic("Gift Error unmarshalling JSON:"+ err.Error())
		return
	}
	for _, v := range GiftListTable {
		GiftListMap.Store(v.Id, v)
	}
}
func GetGiftConfig(id float64) *GiftConfig {
	if v, ok := GiftListMap.Load(id); ok {
		return v.(*GiftConfig)
	}
    return nil
}


type Gift2Config struct {
	Id int32 `json:"Id"`
	Lobby_GiftDesML2 float64 `json:"Lobby_GiftDesML2"`
	GiftNum2 []float64 `json:"GiftNum2"`
	GetSoldierNum2 []float64 `json:"GetSoldierNum2"`
}

var Gift2ListTable []*Gift2Config
var Gift2ListMap sync.Map
// 解析JSON数据到结构体
func initGift2Config(){
  fileContent, err := ioutil.ReadFile("./jsons/Gift2.json")
	if err != nil {
        panic("Gift2 Error reading file:" + err.Error())
		return
	}
	err = json.Unmarshal(fileContent, &Gift2ListTable)
	if err != nil {
		panic("Gift2 Error unmarshalling JSON:"+ err.Error())
		return
	}
	for _, v := range Gift2ListTable {
		Gift2ListMap.Store(v.Id, v)
	}
}
func GetGift2Config(id float64) *Gift2Config {
	if v, ok := Gift2ListMap.Load(id); ok {
		return v.(*Gift2Config)
	}
    return nil
}


type BuffConfig struct {
	BuffType int32 `json:"BuffType"`
	BuffDurationType int32 `json:"BuffDurationType"`
	Int_Param int32 `json:"Int_Param"`
	Float_Param float64 `json:"Float_Param"`
	String_Param string `json:"String_Param"`
	Id int32 `json:"Id"`
	NameML int32 `json:"NameML"`
	IntArray_Param []int32 `json:"IntArray_Param"`
	FloatArray_Param []float64 `json:"FloatArray_Param"`
	StringArray_Param []string `json:"StringArray_Param"`
}

var BuffListTable []*BuffConfig
var BuffListMap sync.Map
// 解析JSON数据到结构体
func initBuffConfig(){
  fileContent, err := ioutil.ReadFile("./jsons/Buff.json")
	if err != nil {
        panic("Buff Error reading file:" + err.Error())
		return
	}
	err = json.Unmarshal(fileContent, &BuffListTable)
	if err != nil {
		panic("Buff Error unmarshalling JSON:"+ err.Error())
		return
	}
	for _, v := range BuffListTable {
		BuffListMap.Store(v.Id, v)
	}
}
func GetBuffConfig(id float64) *BuffConfig {
	if v, ok := BuffListMap.Load(id); ok {
		return v.(*BuffConfig)
	}
    return nil
}

