package latvis

// Configures a LatVis server which uses appengine services (e.g. blob storage,
// http client, etc.).
//
// All AppEngine-specific code should be completely encapsulated inside this package.
//
// Run it locally with:
// $ dev_appserver.py .
// From the root latvis directory.
//
// Also works in a deployed appengine instance.
import (
	"github.com/mrjones/latvis"

	"appengine"
	"appengine/datastore"
	"appengine/taskqueue"
	"appengine/urlfetch"

	"net/http"
	"net/url"
)

const (
	LATVIS_OUTPUT_DATATYPE = "latvis-output"
)

func init() {
//	config := latvis.NewConfig(
//		&AppengineBlobStoreProvider{},
//		&AppengineHttpClientProvider{},
//		&latvis.InMemoryOauthSecretStoreProvider{},
//		&AppengineUrlTaskQueueProvider{})
//	latvis.Setup(config)
	latvis.Setup(AppengineEnvironmentFactory{})
}

type AppengineEnvironmentFactory struct { }

func (fac AppengineEnvironmentFactory) ForRequest(request *http.Request) *latvis.Environment {
	context := appengine.NewContext(request)

	return latvis.NewEnvironment(
		&AppengineBlobStore{context: context},
		&AppengineUrlTaskQueue{context: context},
		&AppengineLogger{context: context},
		&urlfetch.Transport{Context: context},
		)
}

//

type AppengineLogger struct {
	context appengine.Context
}

func (l *AppengineLogger) Errorf(format string, args ...interface{}) {
	l.context.Errorf(format, args)
}

//
// TASK QUEUE
//

type AppengineUrlTaskQueue struct {
	context appengine.Context
}

func (q *AppengineUrlTaskQueue) Enqueue(url string, params *url.Values) error {
	t := taskqueue.NewPOSTTask(url, *params)
	_, err := taskqueue.Add(q.context, t, "")
	return err
}

//
// BLOB STORAGE
//

type AppengineBlobStore struct {
	context appengine.Context
}

func (store *AppengineBlobStore) Store(handle *latvis.Handle, blob *latvis.Blob) error {
	store.context.Infof("Storing blob with handle: '%s'", handle.String())

	datastore.Put(store.context, keyFromHandle(store.context, handle), blob)
	return nil
}

func (store *AppengineBlobStore) Fetch(handle *latvis.Handle) (*latvis.Blob, error) {
	store.context.Infof("Looking up blob with handle: '%s'", handle.String())

	blob := new(latvis.Blob)
	if err := datastore.Get(store.context, keyFromHandle(store.context, handle), blob); err != nil {
		return nil, err
	}
	return blob, nil
}

func keyFromHandle(c appengine.Context, h *latvis.Handle) *datastore.Key {
	return datastore.NewKey(c, LATVIS_OUTPUT_DATATYPE, h.String(), 0, nil)
}
