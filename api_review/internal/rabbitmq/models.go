package rabbitmq

const (
	bookRatingChangedExchange = "BookRatingChangedExchange"
)

type RatingChangeModel struct {
	BookId       string
	Rating       float64
	TotalReviews int
}
