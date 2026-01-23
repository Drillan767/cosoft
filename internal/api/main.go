package api

type Api struct{}

var (
	apiUrl     = "https://hub612.cosoft.fr/v2/api/api"
	spaceId    = "a4928a70-38c1-42b9-96f9-b2dd00db5b02"
	categoryId = "7f1e5757-b9b9-4530-84ad-b2dd00db5f0f"
)

func NewApi() *Api {
	return &Api{}
}
