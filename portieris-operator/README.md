# Portieris Operator
Portieris is a Kubernetes admission controller for the enforcment of image security policies. Portieris can be deployed with operator. You can install portieris in a few simple steps.

## Installing Portieris Operator

0. move to portieris-operatar dir
```
cd portieris-operator
```
if you want to create local cluster, you can create a new kind cluster with this command
```
./dev-scripts/create-kind-cluster.sh
```
and you can delete cluster with the following command
```
kind delete cluster --name=portieris-cluster
```
1. build and push operator image
```
make docker-build IMG=localhost:5000/portieris-go-operator:0.1.5
make docker-push IMG=localhost:5000/portieris-go-operator:0.1.5
```
2. create namespace
```
oc create ns portieris-operator-system
```

3. deploy operator
```
make deploy IMG=localhost:5000/portieris-go-operator:0.1.5
```
4. check pod status
```
$ oc get all -n portieris-operator-system
NAME                                                         READY     STATUS    RESTARTS   AGE
pod/portieris-operator-controller-manager-7c6df5ffff-nwc2f   2/2       Running   0          2m6s

NAME                                                            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/portieris-operator-controller-manager-metrics-service   ClusterIP   10.96.239.150   <none>        8443/TCP   2m6s

NAME                                                    READY     UP-TO-DATE   AVAILABLE   AGE
deployment.apps/portieris-operator-controller-manager   1/1       1            1           2m6s

NAME                                                               DESIRED   CURRENT   READY     AGE
replicaset.apps/portieris-operator-controller-manager-7c6df5ffff   1         1         1         2m7s
```

5. check operator log
```
$ export PORTIERIS_NS=portieris-operator-system
$ make log
bash ./dev-scripts/log_operator.sh
2021-01-28T05:39:25.026Z	INFO	controller-runtime.metrics	metrics server is starting to listen	{"addr": "127.0.0.1:8080"}
2021-01-28T05:39:25.026Z	INFO	setup	starting manager
I0128 05:39:25.029883       1 leaderelection.go:242] attempting to acquire leader lease  portieris-operator-system/b463a212.portieris.io...
2021-01-28T05:39:25.031Z	INFO	controller-runtime.manager	starting metrics server	{"path": "/metrics"}
```

## Installing Portieris
### Custom Resource: Portieris
You can configure Portieris custom resource to define the configuration of Portieris.
#### Configuration of Portieris
`AllowAdmissionSkip`: Allow an annotation to be used to skip the webhook.  
`IBMContainerService`: If not running on IBM Cloud Container Service set to false.  
`securityContextConstraints`: If you deploy portieris in local cluster(minikub/kind etc.), please set to `false`.
```
apiVersion: apis.portieris.io/v1alpha1
kind: Portieris
metadata:
  name: portieris
spec:
  AllowAdmissionSkip: false
  IBMContainerService: false
  securityContextConstraints: false
```
#### Cluster Policy
```
  clusterPolicy:
  - name: '*'
```
### Deploy Custom Resource: Portieris

create Portieris CR
```
oc create -f config/samples/apis_v1alpha1_portieris.yaml -n portieris-operator-system
```

## Uninstalling Portieris

delete Portieris CR
```
oc delete -f config/samples/apis_v1alpha1_portieris.yaml -n portieris-operator-system
```
## Uninstalling Portieris Operator
```
make undeploy
```