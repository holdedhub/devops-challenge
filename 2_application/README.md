Build de container image
REGION=$(terraform -chdir=../1_infrastructure output -raw region)
IMAGE=$(terraform -chdir=../1_infrastructure output -raw repository_name)/demo-app:v1
docker build -f Dockerfile.v1 -t $IMAGE .

gcloud auth configure-docker $REGION-docker.pkg.dev

Note: The first time you run the command it will ask you to add a credential helper to authenticate with the registry. Please answer Y

docker push $IMAGE

Deploy the application in the `default` namespace
kubectl create deployment demo-app --image=$IMAGE

Create a service of type `ClusterIP` to use container-native load balancing
kubectl expose deployment demo-app --name=demo-app-svc --type=ClusterIP --port 8080 --target-port 8080

Create a self-signed SSL certificate using OpenSSL to serve HTTPS requests
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout demo-app.key -out demo-app.crt -subj "/CN=app.example.com/O=Example"

kubectl create secret tls demo-app-secret --key demo-app.key --cert demo-app.crt

Expose the application to the internet. This automatically deploys an external load balancer in GCP
kubectl create ingress demo-app-ing --annotation=kubernetes.io/ingress.class=gce --rule="app.example.com/*=demo-app-svc:8080,tls=demo-app-secret"
watch kubectl get ing demo-app-ing

watch gcloud compute backend-services get-health $(gcloud compute backend-services list --filter="name~demo-app" --format="value(name)") --global

EXTERNAL_IP=$(kubectl get ing demo-app-ing -o jsonpath="{.status.loadBalancer.ingress[*].ip}")

curl -H "Host: app.example.com" http://$EXTERNAL_IP
curl -H "Host: app.example.com" -k https://$EXTERNAL_IP

kubectl logs deployment/demo-app --all-containers

kubectl scale deployment demo-app --replicas=3

kubectl get po -o wide

