package timerscale

import (
	"fmt"
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
	"time"
	log "github.com/golang/glog"

        //"encoding/json"
        //"errors"
        //"path"
        //"sync"
)

func ScalerRunRc(tsspec *TSSpec) { 

        targetResource := tsspec.SubResource
	
        found := false

        now := time.Now()

        currNum, err := getCurrentPodNum(tsspec)

        if err != nil {
             log.Warning("Fail to Get Current Pod number for : ", targetResource, err)
             return
        }


	for _, val := range tsspec.TimeSpanList {

		if now.After(val.BeginTime) && now.Before(val.EndTime) {

                        log.Info("Start to Run Scaler for RC :", targetResource)
			
                        if currNum != val.Num {

				err := runScaler(tsspec.NameSpace, targetResource, val.Num)
                                if err != nil {
			       	    log.Warning("Fail to Scale for :", targetResource, err)
                                 }else{
                                    log.Infof("Success to Scale %s from %d to %d", targetResource,currNum, val.Num)
                                 }
                                 
                                 
			}

                        found = true

			break
		}
	}

        // No time span located in current list, scale to default num instance.

        if !found && (currNum != tsspec.DefaultNum){
            
             err := runScaler(tsspec.NameSpace, targetResource, tsspec.DefaultNum)
             if err != nil {
                   log.Warning("Fail to Scale for :", targetResource, err)
             }else{
                   log.Infof("Success to Scale to Default Num  %s from %d to %d", targetResource,currNum, tsspec.DefaultNum)
             }

        }

	return 

}

// After find a timespan for target rc, scale the instance number

func runScaler(namespace string, resource string, num uint) error {

	//strnum := strconv.Itoa(int(num))

	strrep := fmt.Sprintf("--replicas=%d", num)
	strns := fmt.Sprintf("--namespace=%s", namespace)

	log.Info("Start to Run ", strrep)
	cmd := exec.Command("kubectl", "-s", KUBE_LOCAL_APISERVER, "scale", strrep, "replicationcontrollers", strns, resource)

	strcmd := "kubectl" + "scale" + strrep + "replicationcontrollers" + strns + resource
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Info("Fail", strcmd, ":", out.String())
		return err
	}


	return nil
}

//Get current Pod instance number within the target RC

func getCurrentPodNum(tsspec *TSSpec) (uint, error) {

	strns := fmt.Sprintf("--namespace=%s", tsspec.NameSpace)
	cmd := exec.Command("kubectl", "-s", KUBE_LOCAL_APISERVER, "get", "rc", tsspec.SubResource, strns)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	numMatcher := regexp.MustCompile("(?m)" + `^` + tsspec.SubResource + `\s+\S+\s+\S+\s+\S+\s+([0-9]+)`)
	result := numMatcher.FindStringSubmatch(out.String())

	if result == nil {
	        return 0, fmt.Errorf("No Rc %s Found in Rc Result : %s", tsspec.SubResource, out.String())
	}

	v1, err := strconv.Atoi(result[1])
	if err != nil {
		return 0, err
	}

	return uint(v1), nil

}
