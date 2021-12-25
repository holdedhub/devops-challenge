# Application deployment
In this folder you will find the steps to deploy the web application to GKE in 2 different ways, using the Bash shell and GitHub Actions

## Prerequisites
It is supposed that you have completed the [steps to deploy the infrastructure](../1_infrastructure). Then, you only need:
* The [Docker Engine](https://docs.docker.com/get-docker/) installed

## Manual deployment
The following commands deploy the application from a Bash shell. The CI/CD automated deployment is preferred, this deployment is provided only for learning purposes, to understand the basic steps to deploy the application.

### Build and push
We build the image of the [original application](../app) locally and push it to the Artifact Registry in Google Cloud
```
# Get some variables from Terraform state
REGION=$(terraform -chdir=../1_infrastructure output -raw region)
IMAGE=$(terraform -chdir=../1_infrastructure output -raw repository_name)/demo-app:v1
# Create the container image
docker build --build-arg SOURCE=https://raw.githubusercontent.com/mikakatua/devops-challenge/master/app/server.go -t $IMAGE .
```

This step is only requred if you have never pushed an image to GCP Artifact Registry. Then, you have to configure the credential helper to authenticate with the registry
```
gcloud auth configure-docker $REGION-docker.pkg.dev
```

Finally, push the image to the repository
```
docker push $IMAGE
```

### Kubernetes deployment
The following steps deploy the application in the `default` namespace and make it reachable from Internet using HTTPS. For this example we only use the `kubectl` tool to create the resources, instead of using YALM manifests.
```
# Create the deployment and the service
kubectl create deployment demo-app --image=$IMAGE --replicas=3
kubectl expose deployment demo-app --name=demo-app-svc --type=ClusterIP --port 8080 --target-port 8080
# Check that application Pods are running
kubectl get po -o wide -l app=demo-app
```
We have created a service of type `ClusterIP` to use [container-native load balancing](https://cloud.google.com/kubernetes-engine/docs/concepts/container-native-load-balancing)

Create a self-signed TLS certificate, using OpenSSL, to serve HTTPS requests
```
openssl req -x509 -newkey rsa:2048 -keyout demo-app.key -out demo-app.crt \
  -subj "/CN=app.example.com/O=Example" -days 365 -nodes

# Create a secret to allow Kubernetes use the TLS certificate
kubectl create secret tls demo-app-secret --key demo-app.key --cert demo-app.crt
```

We want to use the Ingress feature to redirect HTTP traffic to HTTPS. The `FrontendConfig` custom resource definition (CRD) allows us to further customize the load balancer
```
cat <<! | kubectl apply -f -
apiVersion: networking.gke.io/v1beta1
kind: FrontendConfig
metadata:
  name: http-to-https
spec:
  redirectToHttps:
    enabled: true
    responseCodeName: MOVED_PERMANENTLY_DEFAULT
!
```

Finally, create the Ingress resource specifying the service as a backend and the secret containing the TLS certificate
```
kubectl create ingress demo-app-ing --annotation=kubernetes.io/ingress.class=gce \
  --annotation=networking.gke.io/v1beta1.FrontendConfig=http-to-https \
  --rule="app.example.com/*=demo-app-svc:8080,tls=demo-app-secret"
```

This Ingress automatically deploys an external load balancer in GCP and it will take some minutes to be available. Once the Ingress is ready, we can get the load balancer external IP address
```
EXTERNAL_IP=$(kubectl get ing demo-app-ing -o jsonpath="{.status.loadBalancer.ingress[*].ip}")
```

We can check the health status of the Pod backends with the command
```
gcloud compute backend-services get-health "$(gcloud compute backend-services list \
  --filter="name~demo-app" --format="value(name)")" --global
```

To test the application we can request the endpoint url and check the Pod logs
```
curl -H "Host: app.example.com" http://$EXTERNAL_IP
curl -H "Host: app.example.com" -k https://$EXTERNAL_IP

kubectl logs deployment/demo-app --all-containers
```

Note: You will see multiple `Hello from my new fresh server` lines because the load balancer periodically performs a health check of the Pods

## CI/CD deployment

