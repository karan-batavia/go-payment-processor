package dto

type Payment struct {
	Id string `json:"id"`
}

func NewPayment(id string) *Payment {
	return &Payment{
		Id: id,
	}
}
