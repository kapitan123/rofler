package data

// AK TODO rename package
import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/kapitan123/telegrofler/config"
	"github.com/kapitan123/telegrofler/internal/data/model"

	log "github.com/sirupsen/logrus"
)

const roflersColName = "roflers"

// AK TODO need to separate in different data stores
type RoflersStore struct {
	ctx        *context.Context
	client     *firestore.Client
	roflersCol *firestore.CollectionRef
	postsCol   *firestore.CollectionRef
}

func NewRoflerStore() *RoflersStore {
	ctx := context.Background()
	client := createClient(ctx)
	return &RoflersStore{
		&ctx,
		client,
		client.Collection(roflersColName),
		client.Collection(postsColName),
	}
}

// actully I rarely need it. Only to store Wins and shit
func (rs *RoflersStore) GetAllRoflers() ([]model.Rofler, error) {
	docs, err := rs.roflersCol.Documents(*rs.ctx).GetAll()

	if err != nil {
		return nil, err
	}

	roflers := []model.Rofler{}
	for _, doc := range docs {
		r := model.Rofler{}
		doc.DataTo(&r)
		roflers = append(roflers, r)
	}

	return roflers, nil
}

// Tries to fetch a document by id
// Internally firestore throws an error if the document does not exist.
// We treat all errors on fetching as not found.
// Errors on convertion are treated in a regural way
func (rs *RoflersStore) GetByUserName(username string) (model.Rofler, bool, error) {
	var r model.Rofler
	doc := rs.roflersCol.Doc(username)
	snap, err := doc.Get(*rs.ctx)

	if err != nil {
		return r, false, nil
	}

	if err := snap.DataTo(&r); err != nil {
		return r, false, err
	}

	return r, true, nil
}

func (rs *RoflersStore) Upsert(r model.Rofler) error {
	doc := rs.roflersCol.Doc(r.UserName)
	_, err := doc.Set(*rs.ctx, r)

	return err
}

func (rs *RoflersStore) UpsertPost(p model.Post) error {
	doc := rs.postsCol.Doc(p.VideoId)
	_, err := doc.Set(*rs.ctx, p)

	return err
}

func (rs *RoflersStore) CreatPost(p model.Post) error {
	doc := rs.postsCol.Doc(p.VideoId)
	_, err := doc.Create(*rs.ctx, p)

	return err
}

func (rs *RoflersStore) Close() {
	rs.client.Close()
}

func createClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, config.ProjectId)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err) // AK TODO shouldn't be fatal
	}
	return client
}
