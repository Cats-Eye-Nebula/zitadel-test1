package mock

//go:generate mockgen -package mock -destination ./repository.mock.go github.com/zitadel/zitadel/v2/internal/eventstore Querier,Pusher
