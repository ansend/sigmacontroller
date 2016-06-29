package timerscale

import (
	//	"encoding/json"
	//	"errors"
	"fmt"
	//	"path"
	//	"regexp"
	//	"sync"
	//	"time"
	"testing"

	//	etcd "github.com/coreos/etcd/client"
	//	log "github.com/golang/glog"
	"golang.org/x/net/context"
)

func newTestEtcdRegistry(t *testing.T) Registry {
	cfg := &EtcdConfig{
		Endpoints: []string{"http://127.0.0.1:4001", "http://127.0.0.1:2379"},
		Prefix:    "/coreos.com/network",
	}

	// if a func is passed to 'NewEtcdTimerScaleRegistry' a customized etcd.KeysAPI instance
	// aka etcd client will returned such as 'newMockEtcd' will returned. otherwise, a etcd
	// provided standard client will be returned.

	r, err := NewEtcdTimerScaleRegistry(cfg, nil)

	if err != nil {
		t.Fatal("Failed to create etcd subnet registry")
	}

	return r
}

func TestEtcdRegistry(t *testing.T) {
	r := newTestEtcdRegistry(t)

	ctx, _ := context.WithCancel(context.Background())

	conf, err := r.GetTSConfig(ctx, "")

	if err != nil {
		t.Fatal("Failed to get networks config")
	}

	fmt.Println("return value is :", conf)

	scalers, index, err := r.GetTSs(ctx, "")
	fmt.Println("return etcd resp index :", index)
	fmt.Println("return Ts is :", len(scalers))
	// TODO: watchSubnet and watchNetworks
}

func TestScaleAll(t *testing.T) {
	r := newTestEtcdRegistry(t)

	ctx, _ := context.WithCancel(context.Background())

	conf, err := r.GetTSConfig(ctx, "")

	if err != nil {
		t.Fatal("Failed to get networks config")
	}

	fmt.Println("return value is :", conf)
	scalers, index, err := r.GetTSs(ctx, "")

	fmt.Println("return etcd resp index :", index)
	fmt.Println("return Ts is :", len(scalers))
	for _, val := range GTimeScalerList {

		validateTimeSpan(val)
		validateResource(val)
		ScalerRunAll(val)
	}

}

func TestValidate(t *testing.T) {

	str := "test"
	/*	num, err := getCurrentPodNum()
		fmt.Println(str, num, err) */

	fmt.Println(str)

}
