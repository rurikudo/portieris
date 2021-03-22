# Portieris Operator
Portieris is a Kubernetes admission controller for the enforcment of image security policies. Portieris can be deployed with operator. You can install portieris in a few simple steps.

## Prerequisites
You can access a cluster.
If you want to use local cluster, you can create local kind cluster and local registry with this command.
  ```
  ./dev-scripts/create-kind-cluster.sh
  ```
  and you can delete cluster with the following command
  ```
  kind delete cluster --name=portieris-cluster
  ```
## Installing Portieris Operator

0. Move to portieris-operatar dir
  ```
  cd portieris-operator
  ```
1. Set up parameters
```
$ export PORTIERIS_NS=portieris-operator-system
$ export PORTIERIS_REPO_ROOT=<absolute path>/portieris
```
1. Build and push operator image
```
make docker-build IMG=localhost:5000/portieris-operator:0.1.5
make docker-push IMG=localhost:5000/portieris-operator:0.1.5
```
2. Create namespace for portieris operator
```
oc create ns portieris-operator-system
```

3. Deploy portieris operator
```
make deploy IMG=localhost:5000/portieris-operator:0.1.5
```
4. Check if operator pod is running
```
$ oc get pod -n portieris-operator-system
NAME                                                         READY     STATUS    RESTARTS   AGE
pod/portieris-operator-controller-manager-7c6df5ffff-nwc2f   2/2       Running   0          2m6s
```

5. You can see operator log with `make log` command
```
$ make log
bash ./dev-scripts/log_operator.sh
2021-01-28T05:39:25.026Z	INFO	controller-runtime.metrics	metrics server is starting to listen	{"addr": "127.0.0.1:8080"}
2021-01-28T05:39:25.026Z	INFO	setup	starting manager
I0128 05:39:25.029883       1 leaderelection.go:242] attempting to acquire leader lease  portieris-operator-system/b463a212.portieris.io...
2021-01-28T05:39:25.031Z	INFO	controller-runtime.manager	starting metrics server	{"path": "/metrics"}
```

## Installing Portieris
### Custom Resource: Portieris
You can configure Portieris custom resource to define the configuration of Portieris. This CR has a similar structure to [values.yaml](https://github.com/IBM/portieris/blob/master/helm/portieris/values.yaml) in Helm chart.

#### Configuration of Portieris
+ allowAdmissionSkip: Allow an annotation to be used to skip the webhook  
+ IBMContainerService: If not running on IBM Cloud Container Service set to false   
+ securityContextConstraints: If you deploy portieris in local cluster(minikub/kind etc.), please set to `false`  
+ useCertManager: If using cert-manager to handle secrets, please set to true
+ skipSecretCreation: If managing portieris-certs secret externally, please set to true
```
apiVersion: apis.portieris.io/v1alpha1
kind: Portieris
metadata:
  name: portieris
spec:
  AllowAdmissionSkip: false
  IBMContainerService: false
  securityContextConstraints: false
  useCertManager: false
  skipSecretCreation: false
```
#### Allowed Repositories
This permissive policy allows all images in namespaces which do not have an ImagePolicy.
```
  allowedRepositories:
  - '*'
```
#### Finalizer
In Kubernetes 1.20 and later, a garbage collector ignore cluster scope children even if their owner is deleted. Please enable a finalizer for portieris when you deploy portieris into Kubernetes 1.20 and later.
```
metadata:
  name: portieris
  finalizers:
    - cleanup.finalizers.portieris.io
```

### Deploy Custom Resource: Portieris

1. Create Portieris CR
```
oc create -f config/samples/apis_v1alpha1_portieris.yaml -n portieris-operator-system
```
2. Check if portieris server pod is running
```
$ oc get pod -n portieris-operator-system
NAME                                                     READY     STATUS    RESTARTS   AGE
portieris-6f7c856876-z4pc6                               1/1       Running   0          88s
portieris-operator-controller-manager-7f85885977-mgtsh   2/2       Running   0          3m31s
```
## Uninstalling Portieris

Delete Portieris CR
```
oc delete -f config/samples/apis_v1alpha1_portieris.yaml -n portieris-operator-system
```
If CR is not removed when finalizer is available, please type the following command.
```
kubectl patch portieris.apis.portieris.io/portieris -p '{"metadata":{"finalizers":[]}}' --type=merge -n portieris-operator-system
```

## Uninstalling Portieris Operator
```
make undeploy
```

