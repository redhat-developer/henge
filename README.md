# henge
Transform multi container spec across providers: Docker Compose, Kubernetes, Openshift, etc.

## Usage

```
# This takes a Docker compose spec file and generates kubernetes artifacts inside output dir
henge --from compose --input compose.yml --to kubernetes --output kube/ 
```
