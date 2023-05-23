package app

import 
(
	"context"
	"cloud.google.com/go/storage"
)

type Application struct {
	cloudStorage *storage.Client
}

func NewApplication(ctx context.Context) Application {
	newClient, err := storage.NewClient(ctx)
	
	if err != nil {
		panic(err)
	}

	return Application{
		cloudStorage : newClient,
	}
}

func (app *Application) SaveVideoToStorage() {

}

// These wrappers can be put in a separate repo, but it seem like an overkill
func (app *Application) GetVideo() {
	// if it is not present tell the link is expired
}

func (app *Application) save() {

}

func (app *Application) download() {

}