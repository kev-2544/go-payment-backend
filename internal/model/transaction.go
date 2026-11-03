package model

type Transaction struct {
    ID       string
    UserID   string
    Amount   int64
    Fee      int64
    Currency string
    Status   string
}