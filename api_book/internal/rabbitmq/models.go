package rabbitmq

const (
	bookCreatedQueue          = "BookCreatedQueue"
	bookCreatedExchange       = "BookCreatedExchange"
	bookRatingChangedExchange = "BookRatingChangedExchange"
)

type RatingChangeModel struct {
	BookId       string
	Rating       float64
	TotalReviews int
}
