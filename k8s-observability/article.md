Kubernetes observability with OpenTelemetry and Jaeger
***
// Entry paragraph on why topic is important
Having a Kubernetes cluster means hosting many compute resources inside. This amount quickly goes up, no matter which architecture solutions have. This is particularly the case with microservices. While a developer can maintain and troubleshoot any part of the solution in a playground, troubleshooting a composed solution deployed to Kubernetes is a tricky thing. Localizing the problem is a doable task with implemented observability in distributed solutions. This article shows a practical  approach to setting up monitoring solution and benefiting from it, even when it is hard/impossible to modify/change your microservices.
***
OpenTelemetry is a well-established Observability framework. It should be used to generate, gather, collect, and export your solution's telemetry data (metrics, traces, and logs). There are two approaches to extract telemetry from your running microservices, known as instrumentation:
- Code-based (or manual): when Developers modify their code to import OpenTelemetry libraries, add lines of code to produce custom telemetry, and DevOps rebuild and redeploy microservices.
- Zero-code (or auto-instrumentation): when DevOps enrich microservice deployments in Kubernetes with observability components, without changing source code at all:

Auto-instrumentation is a particularly promising approach, because it might be really long and difficult to change source code. This approach can and will help localize performance issues in an already working solution.
***
Jaeger is one of the possible Observability backends, which should import telemetry data from OpenTelemetry, store it and provide search and visualize (UI) capability. Jaeger is quite mature, as confirmed by its graduated status on CNCF. Even more interesting, is that after a few years of evolving in many releases of the first version, the new major Jeager v.2.1 (as of the end of 2024) is announced to go live soon. soon.
***
Altogether, the plan for the article is:
- use kind (Kubernetes-IN-Docker) to run Kubernetes cluster locally;
- install OpenTelemetry Operator and configure it to do Zero-code instrumentation (which is individual and tricky as custom services are implemented with different Programming languages);
- install Jaeger Operator in a simpler version and consume OpenTelemetry data;
- install the old-chap Vote microservice sample application as a Kubernetes deployment and configure its vote-part (Python-written) and result-part(Node.js-written) to be Zero-code Instrumented with OpenTelemetry. This code was created years ago when nobody heard about OpenTelemetry;
- install my hand-crafted backend service (Golang-written), which I used in my recent articles;
- run tests on two of the above services and view the telemetry records in the Jaeger UI;
It is important to stress, that I use both test Deployments without any change and code or re-build. And still receive some traces, which might give a hint of "how quick/reliable/healthy etc." arbitrary custom code behaves itself in Kubernetes. If this is not a enough, Code-based Instrumentation to the rescue.
.....
(see the rest on medium.com)