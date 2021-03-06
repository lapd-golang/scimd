package harness

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/fabbricadigitale/scimd/config"
	"github.com/fabbricadigitale/scimd/storage"
	"github.com/fabbricadigitale/scimd/storage/mongo"
	dc "github.com/fsouza/go-dockerclient"
	"github.com/ory/dockertest"
)

var pul *dockertest.Pool
var res *dockertest.Resource
var err error
var ada storage.Storer

func setup() {
	// It uses sensible defaults for windows (tcp/http) and linux/osx (socket)
	// Regarding darwin setting DOCKER_HOST environment variable is probably required
	pul, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	cwd, _ := os.Getwd()
	opt := dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "3.4",
		Env: []string{
			"MONGO_INITDB_DATABASE=scimd",
		},
		Mounts: []string{
			path.Clean(fmt.Sprintf("%s/../../testdata/initdb.d:/docker-entrypoint-initdb.d", cwd)),
		},
		PortBindings: map[dc.Port][]dc.PortBinding{
			"27017/tcp": {{HostIP: "", HostPort: strconv.Itoa(config.Values.Storage.Port)}},
		},
	}
	res, err = pul.RunWithOptions(&opt)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// The application in the container might not be ready to accept connections yet
	if err = pul.Retry(func() error {
		endpoint := fmt.Sprintf("localhost:%s", res.GetPort("27017/tcp"))
		ada, err = mongo.New(endpoint, "scimd", "resources")
		if err != nil {
			return err
		}

		return ada.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}

func teardown() {
	// Kill and remove the container
	if err := pul.Purge(res); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
