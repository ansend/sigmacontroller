package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
        "strings"

	"github.com/ansend/sigmacontroller/pkg/controller/timerscale"
//        "github.com/ansend/sigmacontroller/pkg/kubeclient"
	log "github.com/golang/glog"
	"golang.org/x/net/context"
)

type CmdLineOpts struct {
	etcdEndpoints string
	etcdPrefix    string
	etcdKeyfile   string
	etcdCertfile  string
	etcdCAFile    string
	etcdUsername  string
	etcdPassword  string
	help          bool
	version       bool
}

var opts CmdLineOpts

func init() {
	flag.StringVar(&opts.etcdEndpoints, "etcd-endpoints", "http://127.0.0.1:4001,http://127.0.0.1:2379", "a comma-delimited list of etcd endpoints")
	flag.StringVar(&opts.etcdPrefix, "etcd-prefix", "/coreos.com/network", "etcd prefix")
	flag.StringVar(&opts.etcdKeyfile, "etcd-keyfile", "", "SSL key file used to secure etcd communication")
	flag.StringVar(&opts.etcdCertfile, "etcd-certfile", "", "SSL certification file used to secure etcd communication")
	flag.StringVar(&opts.etcdCAFile, "etcd-cafile", "", "SSL Certificate Authority file used to secure etcd communication")
	flag.StringVar(&opts.etcdUsername, "etcd-username", "", "Username for BasicAuth to etcd")
	flag.StringVar(&opts.etcdPassword, "etcd-password", "", "Password for BasicAuth to etcd")
	flag.BoolVar(&opts.help, "help", false, "print this message")
	flag.BoolVar(&opts.version, "version", false, "print version and exit")
}

func main() {

 //       client := kubeclient.Get()
  //      fmt.Println(client)
	// glog will log to tmp files by default. override so all entries
	// can flow into journald (if running under systemd)
	flag.Set("logtostderr", "true")

	// now parse command line args
	flag.Parse()

	if flag.NArg() > 0 || opts.help {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Register for SIGINT and SIGTERM
	log.Info("Installing signal handlers")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	//Compose the etcd registry
	cfg := &timerscale.EtcdConfig{
		//	Endpoints: []string{"http://127.0.0.1:4001", "http://127.0.0.1:2379"},
		Endpoints: strings.Split(opts.etcdEndpoints, ","),
		Prefix:    "/coreos.com/sigmacontroller",
	}

	r, err := timerscale.NewEtcdTimerScaleRegistry(cfg, nil)

	if err != nil {
		log.Info("Failed to Create Etcd Timescale Registry")
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {

			// Get all the configed timerscale
			_, err := r.GetTSConfig(ctx, "")
			if err != nil {
				log.Info("Fail to Get TS Config:", err)
			}

			_, _, errs := r.GetTSs(ctx, "timescale")

			if errs != nil {
				log.Info("Fail to Get Time Scaler spec:", err)
			}
			for _, val := range timerscale.GTimeScalerList {

				_, err := timerscale.IsValidTSSpec(val)
                                if err != nil {
                                     log.Info("Invalid Timer Scale config :", err)
                                } else{
 
                                     timerscale.ScalerRunRc(val)
                                }

			}
			time.Sleep(time.Minute * 1)
		}

		wg.Done()
	}()

	<-sigs
	// unregister to get default OS nuke behaviour in case we don't exit cleanly
	signal.Stop(sigs)

	log.Info("Exiting...")
	cancel()

	wg.Wait()

}
