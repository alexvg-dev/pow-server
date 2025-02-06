package services

type QuotesRepo interface {
	GetOneQuote() (string, error)
}

type QuoteProvider struct {
	Repo QuotesRepo
}

func NewQuoteProvider(repo QuotesRepo) *QuoteProvider {
	return &QuoteProvider{
		Repo: repo,
	}
}

func (q *QuoteProvider) GetQuote() (string, error) {
	quote, err := q.Repo.GetOneQuote()
	if err != nil {
		return "", err
	}

	return quote, nil
}
