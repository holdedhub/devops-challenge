# DevOps Technical Test

We think infrastructure is best represented as code, and provisioning of resources should be automated as much as possible.

We are testing your ability to implement modern automated infrastructure, as well as general knowledge of system administration and coding. In your solution you should emphasize readability, maintainability and DevOps methodologies.

To begin, create a GitHub repository and start adding your work. When you're finished, send us the URL to your repository.

You can use the following folder structure or create your own:

```
./
├─ 1_infrastructure
│  └─ <your answer>
├─ 2_application
│  └─ <your answer>
└─ 3_security
   └─ <your answer>
```

## 1. Infrastructure Test

Provide a diagram of the architecture in AWS or GCP to deploy our [Golang server](app/server.go) with HA that can be used in a repeatable way. 
The answer should contain the following:

* Use of Terraform.
* Use of Kubernetes.
* Clearly explaining why you're doing things a certain way.

## 2. Application (CI/CD)

Provide a PNG diagram of the CI/CD proposal using a tool of your choice such to automate the docker build and deploy of the [Golang server](app/server.go) that serves some static or dynamic content.
The answer should contain the following:

* Using Containers as part of your automation.
* Creating a CI pipeline, using a tool of your choice, that deploys the web server to a cloud environment of your choice.
* Serve traffic from 443 port with self-signed certificate would be highly appreciate.

## 3. Security

Provide an explanation about the security that you consider to be implemented to this stack always taking in consideration of the best practices
that should be introduced for a production deployment.
