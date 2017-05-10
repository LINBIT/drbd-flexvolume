# DRBD flexvolume plugin for Kubernetes

## Building

Requires Go 1.8 or higher and a configured GOPATH

`mkdir -p $GOPATH/src/github.com/`

`cd $GOPATH/src/github.com/`

`git clone git@drbd.org:drbd-flexvolume.git`

`cd drbd-flexvolume`

`make`

This will compile a binary targeting the local machine's architecture and
place it into the `_build` directory.

## Installing

Place the generated binary named `drbd` under the following path on the kubelet and
kube-controller-manager nodes: 

/usr/libexec/kubernetes/kubelet-plugins/volume/exec/linbit~drbd/

## Usage

Resources must be created before attachment with DRBD Manage before they are
available for Kubernetes.

The kube-controller-manager and all kubelets eligible to run containers must be
part of the same DRBD Manage cluster. Volumes will be attached to the kubelet
across the network via the DRBD Transport protocol, so they do not require local
storage.

kublet nodes names must match the output of `uname -n` exactly. If they do not,
this may be overridden via the kubelet `--hostname-override` parameter

`example.yaml`, located in the root of this project, contains an example
configuration that attaches a resource named `r0` to the container under the path
`/data`
