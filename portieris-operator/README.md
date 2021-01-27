## Installing Portieris Operator

1. build and push image
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

## Installing Portieris
prepare Portieris CR  
If you deploy portieris in local cluster(minikub/kind etc.), please set `securityContextConstraints` to `false`.
```
apiVersion: apis.portieris.io/v1alpha1
kind: Portieris
metadata:
  name: portieris
spec:
  securityContextConstraints: false
```
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