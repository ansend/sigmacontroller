package timerscale

import (
	//	"encoding/json"
	//"errors"
	"fmt"
	//	"path"
	"regexp"
	//	"sync"
	//       "time"
	//      "strconv"
	"bytes"
	"os/exec"
	//	etcd "github.com/coreos/etcd/client"
	//       "github.com/coreos/etcd/pkg/transport"
	log "github.com/golang/glog"
	//	"golang.org/x/net/context"
)

const (
	MIN_PODS_NUM = 0
	MAX_PODS_NUM = 20
)

func IsValidTSSpec(tsspec *TSSpec) (bool, error) {

	err := validateTimeSpan(tsspec)

	if err != nil {
		return false, err
	}

	err = validateResource(tsspec)
	if err != nil {
		return false, err
	}

	err = validatePodNum(tsspec)
	if err != nil {
		return false, err
	}

	return true, nil
}

func validatePodNum(tsspec *TSSpec) error {

	if !(tsspec.DefaultNum >= MIN_PODS_NUM && tsspec.DefaultNum <= MAX_PODS_NUM) {
		return fmt.Errorf("%s Default Pods num should bettween %d and %d, current is %d ", "tsspec.Subresource", MIN_PODS_NUM, MAX_PODS_NUM, tsspec.DefaultNum)
	}

	for _, val := range tsspec.TimeSpanList {
		if !(val.Num >= MIN_PODS_NUM && val.Num <= MAX_PODS_NUM) {
			return fmt.Errorf("%s PODs num should bettween %d and %d, current is %d ", "tsspec.Subresource", MIN_PODS_NUM, MAX_PODS_NUM, val.Num)
		}
	}

	return nil
}

// Make sure no overlap  between any 2 of the timespans
// for 2 of the spans, following condition make sure no overlap
// startdate1 <=enddate2 and enddate1>=startdate2

func validateTimeSpan(tsspec *TSSpec) error {

	num := len(tsspec.TimeSpanList)

	for i := 0; i < num-1; i++ {
		for j := i + 1; j < num; j++ {

			if (tsspec.TimeSpanList[i].BeginTime.Unix() <= tsspec.TimeSpanList[j].EndTime.Unix()) &&
				(tsspec.TimeSpanList[i].EndTime.Unix() >= tsspec.TimeSpanList[j].BeginTime.Unix()) {

				log.Infof("There is Overlapped Timespan for Rc %s:", tsspec.SubResource)
				return fmt.Errorf("%s Overlapped Timespan", tsspec.SubResource)

			}

		}

	}

	return nil

}

// Check if the target resouce is a valid name to match the following rules.
// 1) it can be find in rc list and it's a valid rc
// 2) it can not be in the list of hpa since autoscale can not cowork with timer scale

func validateResource(tsspec *TSSpec) error {

	/*	re := regexp.MustCompile("resourcename")
			str := `fldjlafjldja
			        fdjklfjdlafdafdfdafa
		                fdjalfdjla resourcename
		                fdjlasfjd resourcename fdjla
		                fldlfjlda resourcename fdas `

			matches := re.FindAllString(str, -1)
			fmt.Println(re.FindAllString(str, -1))
			fmt.Println(len(matches))
	*/
        strns := fmt.Sprintf("--namespace=%s", tsspec.NameSpace)
	pattern := ".*" + tsspec.SubResource + ".*"
	cmd := exec.Command("kubectl", "get", "pods",strns)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	matched, err := regexp.MatchString(pattern, out.String())

	if matched {
		return nil
	} else {
		return fmt.Errorf("%s Not Exist in Rc List : %s", tsspec.SubResource,out.String())
	}

}
