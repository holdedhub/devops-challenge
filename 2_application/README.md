Create the registry to store our container image

gcloud artifacts repositories create devops-challenge-repo \
  --repository-format=docker \
  --location=europe-west6

Note: the first time you run the command it will ask you to enable the API `artifactregistry.googleapis.com`. Please answer Y

Build de container image
docker build -t europe-west6-docker.pkg.dev/${PROJECT_ID}/devops-challenge-repo/demo-app:v1 .

gcloud auth configure-docker europe-west6-docker.pkg.dev

Note: he first time you run the command it will ask you to add a credential helper to authenticate with the registry. Please answer Y

docker push europe-west6-docker.pkg.dev/${PROJECT_ID}/devops-challenge-repo/demo-app:v1

Deploy the application in the `default` namespace
kubectl create deployment demo-app --image=europe-west6-docker.pkg.dev/${PROJECT_ID}/devops-challenge-repo/demo-app:v1

Create a service of type `LoadBalancer` to expose the application to the internet. This automatically deploys an external load balancer in GCP
kubectl expose deployment demo-app --name=demo-app-svc --type=LoadBalancer --port 80 --target-port 8080

Wait until the `EXTERNAL-IP` column is populated
watch kubectl get svc demo-app-svc

EXTERNAL_IP=$(kubectl get svc demo-app-svc -o jsonpath="{.status.loadBalancer.ingress[*].ip}")

curl $EXTERNAL_IP

kubectl logs deployment/demo-app --all-containers

---

kubectl expose deployment demo-app --name=demo-app-svc --type=ClusterIP --port 8080 --target-port 8080

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout demo-app.key -out demo-app.crt -subj "/CN=app.example.com/O=Example"

kubectl create secret tls demo-app-secret --key demo-app.key --cert demo-app.crt

kubectl create ingress demo-app-ing --annotation=kubernetes.io/ingress.class=gce --rule="app.example.com/*=demo-app-svc:8080,tls=demo-app-secret"
watch kubectl get ing demo-app-ing

EXTERNAL_IP=$(kubectl get ing demo-app-ing -o jsonpath="{.status.loadBalancer.ingress[*].ip}")

curl -H "Host: app.example.com" http://$EXTERNAL_IP
curl -H "Host: app.example.com" -k https://$EXTERNAL_IP

kubectl logs deployment/demo-app --all-containers

