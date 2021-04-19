# Cosign verifier 
Cosign verifier is a tool for cosign verification at Admission controller.  
Cosign verifier works as sidecar of portieris.

## Enable cosign verifier
You can enable cosign verifier from portieris operator.

0. Move to `portieris` dir.
```
$ pwd
/home/repo/portieris
```
1. Prepare portieris image to support cosign verifier.
```
make image
docker tag portieris:v0.10.1 localhost:5000/portieris:v0.10.1
docker push localhost:5000/portieris:v0.10.1
docker push 
```
2. Build cosign verifier image
```
cd cosign-verify
docker build -t localhost:5000/cosign-verifier:0.0.1 .
docker push localhost:5000/cosign-verifier:0.0.1
```
3.  Move to `portieris-operator` dir and prepare operator images.
```
make docker-build IMG=localhost:5000/portieris-operator:0.0.1
docker push localhost:5000/portieris-operator:0.0.1
```
4. Edit image and `enableCosignVerify` fields in `config/samples/apis_v1alpha1_portieris.yaml`.
```
spec:
  enableCosignVerify: true
  cosignVerifierImage: localhost:5000/cosign-verifier:0.0.1
  image:
  host: localhost:5000
  image: portieris
  pullPolicy: Always
  pullSecret: []
  tag: v0.10.1
```
3. Install operator. please check this [document](https://github.com/rurikudo/portieris/blob/master/portieris-operator/README.md) for operator installation.