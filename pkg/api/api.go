package api

import (
	"encoding/json"
	"fmt"
	"linbit/drbd-flexvolume/pkg/drbd"
)

type flexAPIErr struct {
	message string
}

func (e flexAPIErr) Error() string {
	return fmt.Sprintf("DRBD Flexvoume API: %s", e.message)
}

type response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type attachResponse struct {
	response
	Device string `json:"device"`
}

type isAttachedResponse struct {
	response
	Attached string `json:"attached"`
}

type getVolNameResponse struct {
	response
	VolumeName string `json:"volumeName"`
}

type options struct {
	FsType    string `json:"kubernetes.io/fsType"`
	Readwrite string `json:"kubernetes.io/readwrite"`
	Resource  string `json:"resource"`
}

func parseOptions(s string) (options, error) {
	opts := options{}
	err := json.Unmarshal([]byte(s), &opts)
	if err != nil {
		return opts, flexAPIErr{fmt.Sprintf("couldn't parse options from %s", s)}
	}

	return opts, nil
}

type FlexVolumeApi struct {
}

func (api FlexVolumeApi) Call(s []string) (string, int) {
	if len(s) < 1 {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: "No driver action! Valid actions are: init, attach, detach, mountdevice, unmountdevice, getvolumename, isattached",
		})
		return string(res), 2
	}
	switch s[0] {
	case "init":
		return api.init()
	case "attach":
		return api.attach(s)
	case "waitforattach":
		return api.waitForAttach(s)
	case "detach":
		return api.detach(s)
	case "mountdevice":
		return api.mountDevice(s)
	case "unmountdevice":
		return api.unmountDevice(s)
	case "getvolumename":
		return api.getVolumeName(s)
	case "isattached":
		return api.isAttached(s)
	default:
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: fmt.Sprintf("Unsupported driver action: %q", s[0]),
		})
		return string(res), 2
	}
}

func (api FlexVolumeApi) init() (string, int) {
	res, _ := json.Marshal(response{Status: "Success"})
	return string(res), 0
}

func (api FlexVolumeApi) attach(s []string) (string, int) {
	if len(s) < 2 {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: flexAPIErr{fmt.Sprintf("attach: too few arguments passed: %s", s)}.Error(),
		})
		return string(res), 2
	}

	opts, err := parseOptions(s[0])
	if err != nil {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: err.Error(),
		})
		return string(res), 2
	}

	resource := drbd.Resource{Name: opts.Resource, NodeName: s[1]}

	_, err = drbd.AssignRes(resource)
	if err != nil {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: flexAPIErr{fmt.Sprintf("attach: failed to assign resource %q", resource.Name)}.Error(),
		})
		return string(res), 1
	}

	path, err := drbd.WaitForDevPath(resource, 4)
	if err != nil {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: flexAPIErr{fmt.Sprintf("attach: unable to find device path for resource %q", resource.Name)}.Error(),
		})
		return string(res), 1
	}

	res, _ := json.Marshal(attachResponse{
		Device: path,
		response: response{
			Status: "Success",
		},
	})
	return string(res), 0
}

func (api FlexVolumeApi) waitForAttach(s []string) (string, int) {
	res, _ := json.Marshal(response{Status: "Success"})
	return string(res), 0
}

func (api FlexVolumeApi) detach(s []string) (string, int) {
	if len(s) < 2 {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: flexAPIErr{fmt.Sprintf("detach: too few arguments passed: %s", s)}.Error(),
		})
		return string(res), 2
	}

	opts, err := parseOptions(s[0])
	if err != nil {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: err.Error(),
		})
		return string(res), 2
	}

	resource := drbd.Resource{Name: opts.Resource, NodeName: s[1]}

	err = drbd.UnassignRes(resource)
	if err != nil {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: err.Error(),
		})
		return string(res), 2
	}

	res, _ := json.Marshal(response{Status: "Success"})
	return string(res), 0
}

func (api FlexVolumeApi) mountDevice(s []string) (string, int) {
	res, _ := json.Marshal(response{Status: "Success"})
	return string(res), 0
}

func (api FlexVolumeApi) unmountDevice(s []string) (string, int) {
	res, _ := json.Marshal(response{Status: "Success"})
	return string(res), 0
}

func (api FlexVolumeApi) getVolumeName(s []string) (string, int) {
	if len(s) < 1 {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: flexAPIErr{fmt.Sprintf("getvolumename: too few arguments passed: %s", s)}.Error(),
		})
		return string(res), 2
	}

	opts, err := parseOptions(s[0])
	if err != nil {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: err.Error(),
		})
		return string(res), 2
	}

	res, _ := json.Marshal(getVolNameResponse{
		VolumeName: opts.Resource,
		response: response{
			Status: "Success",
		},
	})
	return string(res), 0
}

func (api FlexVolumeApi) isAttached(s []string) (string, int) {
	if len(s) < 2 {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: flexAPIErr{fmt.Sprintf("isattached: too few arguments passed: %s", s)}.Error(),
		})
		return string(res), 2
	}

	opts, err := parseOptions(s[0])
	if err != nil {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: err.Error(),
		})
		return string(res), 2
	}

	resource := drbd.Resource{Name: opts.Resource, NodeName: s[1]}

	ok, err := drbd.WaitForAssignment(resource, 4)
	if err != nil {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: err.Error(),
		})
		return string(res), 2
	}

	if !ok {
		res, _ := json.Marshal(response{
			Status:  "Failure",
			Message: flexAPIErr{fmt.Sprintf("resource %q not attached", resource.Name)}.Error(),
		})
		return string(res), 2
	}

	res, _ := json.Marshal(isAttachedResponse{
		Attached: "true",
		response: response{Status: "Success"},
	})
	return string(res), 0
}
