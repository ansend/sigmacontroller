package timerscale
import (
	//"encoding/json"
	//"errors"
	"fmt"
	"path"
	//"regexp"
	"sync"
        "time"

	etcd "github.com/coreos/etcd/client"
        "github.com/coreos/etcd/pkg/transport"
//	log "github.com/golang/glog"
	"golang.org/x/net/context"

)

func test() (error){

  return nil
}


type Registry interface {
	getTSConfig(ctx context.Context, network string) (string, error)
        getTSs(ctx context.Context, network string )  ([]string, uint64, error) 
}


type EtcdConfig struct {
	Endpoints []string
	Keyfile   string
	Certfile  string
	CAFile    string
	Prefix    string
	Username  string
	Password  string
}

type etcdNewFunc func(c *EtcdConfig) (etcd.KeysAPI, error)

type etcdTimerScaleRegistry struct {
	cliNewFunc   etcdNewFunc
	mux          sync.Mutex
	cli          etcd.KeysAPI
	etcdCfg      *EtcdConfig
}

func newEtcdClient(c *EtcdConfig) (etcd.KeysAPI, error) {
	tlsInfo := transport.TLSInfo{
		CertFile: c.Certfile,
		KeyFile:  c.Keyfile,
		CAFile:   c.CAFile,
	}

	t, err := transport.NewTransport(tlsInfo, 30 * time.Second)
	if err != nil {
		return nil, err
	}

	cli, err := etcd.New(etcd.Config{
		Endpoints: c.Endpoints,
		Transport: t,
		Username:  c.Username,
		Password:  c.Password,
	})
	if err != nil {
		return nil, err
	}

	return etcd.NewKeysAPI(cli), nil
}

func newEtcdTimerScaleRegistry(config *EtcdConfig, cliNewFunc etcdNewFunc) (Registry, error) {
	r := &etcdTimerScaleRegistry{
		etcdCfg:      config,
	}
	if cliNewFunc != nil {
		r.cliNewFunc = cliNewFunc
	} else {
		r.cliNewFunc = newEtcdClient
	}

	var err error
	r.cli, err = r.cliNewFunc(config)
	if err != nil {
		return nil, err
	}

	return r, nil
}


func (etr *etcdTimerScaleRegistry) getTSConfig(ctx context.Context, network string) (string, error) {
	key := path.Join(etr.etcdCfg.Prefix, network, "config")
	resp, err := etr.client().Get(ctx, key, nil)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

// getSubnets queries etcd to get a list of currently allocated leases for a given network.
// It returns the leases along with the "as-of" etcd-index that can be used as the starting
// point for etcd watch.
func (etr *etcdTimerScaleRegistry) getTSs(ctx context.Context, network string) ([]string, uint64, error) {
	key := path.Join(etr.etcdCfg.Prefix, network)
	resp, err := etr.client().Get(ctx, key, &etcd.GetOptions{Recursive: true})
	if err != nil {
		if etcdErr, ok := err.(etcd.Error); ok && etcdErr.Code == etcd.ErrorCodeKeyNotFound {
			// key not found: treat it as empty set
			return []string{}, etcdErr.Index, nil
		}
		return nil, 0, err
	}

	scalers := []string{}
	for _, node := range resp.Node.Nodes {
		//l, err := nodeToLease(node)
		/*if err != nil {
			log.Warningf("Ignoring bad subnet node: %v", err)
			continue
		} */
                s := node.Key
                fmt.Println("current key is :", s)
		scalers = append(scalers, s)
	}

	return scalers, resp.Index, nil
}

func (esr *etcdTimerScaleRegistry) client() etcd.KeysAPI {
	esr.mux.Lock()
	defer esr.mux.Unlock()
	return esr.cli
}

func (esr *etcdTimerScaleRegistry) resetClient() {
	esr.mux.Lock()
	defer esr.mux.Unlock()

	var err error
	esr.cli, err = newEtcdClient(esr.etcdCfg)
	if err != nil {
		panic(fmt.Errorf("resetClient: error recreating etcd client: %v", err))
	}
}
