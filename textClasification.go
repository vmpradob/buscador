package main 

type TextClasification struct {
  Time       int64      `json:"time"`      
  Categories []Category `json:"categories"`
  Lang       string     `json:"lang"`      
  Timestamp  string     `json:"timestamp"` 
}

type Category struct {
  Name  string  `json:"name"` 
  Score float64 `json:"score"`
}
