# Setup Kubernetes cluster

The following manual explains how to install a Kubernetes cluster with `kubeadm`.

## Setup VM
Setup a server and make sure that you can SSH into the server.
In the following examples we assume that our server is named `k1`.

For this we use an Ubuntu image:
* https://cloud-images.ubuntu.com/minimal/releases/focal/release-20220810/

Before you proceed you should be able to SSH into the server:
```
ssh k1
```

### Setup VM with vu (optional)
How you setup your server is up to you. If you are on a Linux machine with KVM and Libvirt an easy way to quickly setup a VM is using [vu](https://github.com/dvob/vu).
The tool `vu` provides an easy way to setup a virtual machine based on a cloud-init image.
```
vu image add https://cloud-images.ubuntu.com/minimal/releases/focal/release-20220810/ubuntu-20.04-minimal-cloudimg-amd64.img
vu create --disk-size 100GB --memory 4G --cpu 2 ubuntu-20.04-minimal-cloudimg-amd64.img k1
```

## Setup Kubernetes
If the server is up and running, we first run the [setup script](./setup.sh) on the server to install kubelet, kubeadm, containerd, etc.
The setup script basically performs the tasks described under https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/.

Run the setup script on the server:
```
ssh k1 <setup.sh
```

Then you can bootstrap the cluster with kubeadm:
```
ssh k1 sudo kubeadm init --pod-network-cidr=10.244.0.0/16
```

### Join Nodes
To test the WASM extension a single node is enough.

If you still want more nodes in your cluster, you can repeat the actions above to add more nodes to a cluster (setup VM, run setup script).
Instead of `kubeadm init` you have to run the appropriate `kubeadm join` command which is shown after `kubeadm init`.

You also can generate a new join command like this:
```
token=$( ssh k1 sudo kubeadm token generate )
join_cmd=$( ssh k1 sudo kubeadm token create $token --print-join-command )
```

Then run the join command on the additional node `k2`:
```
ssh k2 sudo $join_cmd
```

### Install Network Plugin
To test the features of the WASM extension you don't need a network plugin.
If you want to deploy workload in your cluster you have to install a network plugin like Flannel:
```
ssh k1 sudo kubectl --kubeconfig /etc/kubernetes/admin.conf apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
```

## Configure kubectl
To connect to your test cluster from your local machine you have to configure `kubectl` accordingly.

The easiest way to do so is to copy the admin configuration from your server:
```
ssh k1 sudo cat /etc/kubernetes/admin.conf > ~/.kube/test_config
```

To not overwrite our existing configuration we put the configuration in a different location and then set the location of the config using the `KUBECONFIG` environment variable.
```
export KUBECONFIG=~/.kube/test_config
```

A `kubectl get nodes` should now list the available nodes.

## Run a custom build of the API server
Publish the kube-apiserver Docker image of your custom build to a Docker registry.
In the following example I use the image `dvob/kube-apiserver:dev`.

Login to the `k1` server and update the manifest of the API server `/etc/kubernetes/manifests/kube-apiserver.yaml` under `spec.containers[0].image` and set it to `dvob/kube-apiserver:dev`.

Alternatively, you can run this command:
```
ssh k1 'sudo kubectl set image -f /etc/kubernetes/manifests/kube-apiserver.yaml kube-apiserver=dvob/kube-apiserver:dev --local -o yaml | sudo sponge /etc/kubernetes/manifests/kube-apiserver.yaml'
```

Usually it takes some time until the kubelet picks up the new manifest and restarts the API server.
To speed up the restart of the API server you can login to the server `k1` and force a restart as follows:
```
ssh k1
sudo su -
crictl ps --name kube-apiserver | awk 'NR != 1{ print $1}' | xargs crictl stop
systemctl restart kubelet
```
