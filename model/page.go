package model

type Page struct {
	Num  int
	Next interface{} // Profile can set data for pagination
}
