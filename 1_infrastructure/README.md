# Infrastructure provisioning
This repo contains the instructions to deploy a GKE cluster using Terraform. The gcloud commands are provided but the same can be done from the GCP console

## Architecture
We will create a highly available [regional cluster](https://cloud.google.com/kubernetes-engine/docs/concepts/regional-clusters#regional) with the control plane and nodes replicated across three zones of a region. 

![](https://storage.googleapis.com/gweb-cloudblog-publish/original_images/gcp-google-kubernetes-engine-regional-clusterbcum.PNG)

## Prerequisites
Before trying to provision the infrastructure, you need:
* A [GCP account](https://console.cloud.google.com/) (free-tier is enough)
* The [gcloud SDK](https://cloud.google.com/sdk/docs/install) installed and configured with your GCP account
* The [kubectl CLI](https://kubernetes.io/docs/tasks/tools/) installed
* The [terraform CLI](https://www.terraform.io/downloads) installed

## Create a new project
This is not mandatory but it is good to keep the resources separate
```
gcloud projects create --name=devops-challenge --quiet
```
We save the project ID in a variable for later use
```
PROJECT_ID=$(gcloud projects list --filter="name: devops-challenge" --format="value(projectId)")
```

## Create the Terraform service account
These commands create the service account and grants the permissions
```
gcloud iam service-accounts create terraform-automation \
  --display-name="Terraform Service Account" \
  --project=$PROJECT_ID
```
```
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:terraform-automation@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/editor"
```
Note: For the sake of simplicity we granted the `Editor` role to the service account. This role allows create and delete most of the GCP resources.
From a security perspective, we should enforce the principle of least privilege and grant only the specific roles needed.

Finally, we download the service account private key that will be used later with the terraform CLI
```
gcloud iam service-accounts keys create ~/terraform-automation-key.json --iam-account="terraform-automation@$PROJECT_ID.iam.gserviceaccount.com"
```
