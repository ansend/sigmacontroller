package timerscale

import (
	"encoding/json"
	//"errors"
	"fmt"
	"path"
        //"regexp"
	//	"bytes"
	//	"os/exec"
	//	"strconv"
	"sync"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/coreos/etcd/pkg/transport"
	log "github.com/golang/glog"
	"golang.org/x/net/context"
)

type Registry interface {
	GetTSConfig(ctx context.Context, network string) (string, error)
	GetTSs(ctx context.Context, network string) ([]string, uint64, error)
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

type EtcdTimerScaleRegistry struct {
	cliNewFunc etcdNewFunc
	mux        sync.Mutex
	cli        etcd.KeysAPI
	etcdCfg    *EtcdConfig
}

type TimeSpan struct {
	Num       uint
	Begin     string
	End       string
	BeginTime time.Time
	EndTime   time.Time
}

type TSSpec struct {
	NameSpace    string `json:"omitempty"`
	DefaultNum   uint   `json:"omitempty"`
	SubResource  string
	TimeSpanList []TimeSpan
}

var GTimeScalerList []*TSSpec

func newEtcdClient(c *EtcdConfig) (etcd.KeysAPI, error) {
	tlsInfo := transport.TLSInfo{
		CertFile: c.Certfile,
		KeyFile:  c.Keyfile,
		CAFile:   c.CAFile,
	}

	t, err := transport.NewTransport(tlsInfo, 30*time.Second)
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

func NewEtcdTimerScaleRegistry(config *EtcdConfig, cliNewFunc etcdNewFunc) (Registry, error) {
	r := &EtcdTimerScaleRegistry{
		etcdCfg: config,
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

func (etr *EtcdTimerScaleRegistry) GetTSConfig(ctx context.Context, network string) (string, error) {
	key := path.Join(etr.etcdCfg.Prefix, network, "config")
	resp, err := etr.client().Get(ctx, key, nil)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

// getTSs queries etcd to get a list of currently configed timer scaler resource.
func (etr *EtcdTimerScaleRegistry) GetTSs(ctx context.Context, network string) ([]string, uint64, error) {
	key := path.Join(etr.etcdCfg.Prefix, network)
	resp, err := etr.client().Get(ctx, key, &etcd.GetOptions{Recursive: true})
	if err != nil {
		if etcdErr, ok := err.(etcd.Error); ok && etcdErr.Code == etcd.ErrorCodeKeyNotFound {
			// key not found: treat it as empty set
			return []string{}, etcdErr.Index, nil
		}
		return nil, 0, err
	}
	// clear the list of scaler
	GTimeScalerList = nil

	scalers := []string{}
	for _, node := range resp.Node.Nodes {

		spec, err := nodeToTSSpec(node)

		if err == nil {
		
	               GTimeScalerList = append(GTimeScalerList, spec)

		} else {

			log.Warningf("Ignoring bad spec node: %v", err)
		}

		s := node.Key
		log.Info("Current Key is : ", s)
		scalers = append(scalers, s)
	}

	fmt.Println("current global timer list:",GTimeScalerList,len(GTimeScalerList))

	return scalers, resp.Index, nil
}

func nodeToTSSpec(node *etcd.Node) (*TSSpec, error) {

	attrs := &TSSpec{NameSpace: "default"}

	if err := json.Unmarshal([]byte(node.Value), attrs); err != nil {

		return nil, err
	}

	var err error

	timeForm := "2006-01-02 15:04:05"
	for inx, val := range attrs.TimeSpanList {

		//val.BeginTime, err  = time.Parse(timeForm, val.Begin)
		//val.EndTime, err = time.Parse(timeForm, val.End)

		//in go , the range syntx make "val" a copy of the origin slice
		// so mordification of the val did not change the origin slice value
		// at all, we need to new a pointer (ref) to the origin slice with index.

		valref := &attrs.TimeSpanList[inx]

		valref.BeginTime, err = time.ParseInLocation(timeForm, val.Begin, time.Local)
		valref.EndTime, err = time.ParseInLocation(timeForm, val.End, time.Local)

		if err != nil {
			fmt.Println("Error to parse the time")
			return nil, err
		}

	}

	//.fmt.Println("print the json text ")
	//fmt.Println(attrs)

	return attrs, nil
}

func (esr *EtcdTimerScaleRegistry) client() etcd.KeysAPI {
	esr.mux.Lock()
	defer esr.mux.Unlock()
	return esr.cli
}

func (esr *EtcdTimerScaleRegistry) resetClient() {
	esr.mux.Lock()
	defer esr.mux.Unlock()

	var err error
	esr.cli, err = newEtcdClient(esr.etcdCfg)
	if err != nil {
		panic(fmt.Errorf("resetClient: error recreating etcd client: %v", err))
	}
}
