# Kubekadal

Talk to your infra using the Kubernetes [client-go](https://godoc.org/k8s.io/client-go/kubernetes) library. Inspired by kelsey hightower kube conference. [Youtube Link](https://youtu.be/6BYq6hNhceI?list=PL5IOR_b5llySpyG_2CE7dojO0ZVc2f5Ix&t=816)<br/> 

# About Kubekadal
Kubekadal is a simple golang app for making the API call to Kubernetes cluster to make changes such as routing the traffic and getting the cluster information. This app receives a request from the dialogflow webhooks. It parses the request parameters to understands the action. If it finds the action and intents, it will make an API call to the Kubernetes cluster to make those changes. <br/>
**For demo purpose only since voice commands are high volatile, don't use this in your prod/staging cluster setup.

### Before you begin
1) Install go runtime in your local machine. 
2) Create a Kubernetes cluster with Istio implementation.
3) Create a Istio virtual service with destination rules. Specify the virtual service name as "demo-app-vs". [reference](https://blog.webischia.com/2018/10/21/using-traffic-shifting-on-istio-to-make-blue-green-deployments-on-kubernetes/)
4) Place the kube config file in the HOME/.kuber folder (refer: [clientSet() func](./pkg/kubeclient.go))
5) Built a Dialogflow agent and intents to send the commands to this project. [reference](https://dialogflow.cloud.google.com)

### Getting started
1) Clone this project
2) Open [constant file](./app/utils/constant.go) and change the intent-id with your dialogflow intent ids.
3) Run this app locally
4) Use ngrok to expose the https domain url.

### Dialogflow
1) Open diagflow and create the agent
2) Replace the webhook URL with the ngrok url you have https://ngrokurl/app/dialogflow
3) Download the intents/entities from the [folder](./dialogflow/objects)
4) Import the entities/intents from that folder and import it to your dialogflow.
5) Note down the intent-id's you just created and modify the [constant file](./app/utils/constant.go) with those values
   - a welcome intent(replaces this intent-id for KubeCtlMain in constant file)<br/>
   - a serviceinfo intent(replaces this intent-id for KubeCtlServiceInfo in constant file)<br/>
   - a virtual service intent(replaces this intent-id for KubeCtlVirtualService in constant file)<br/>

### :punch: That's it. You are good to go...





