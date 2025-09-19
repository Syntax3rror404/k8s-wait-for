# k8s-wait-for
Go cli app to wait for kubernetes pod/jobs to come ready 

Supports service account by default and if no one is found it's try to use /.kube/config

Example how to use it:

```yaml
kind: StatefulSet
metadata:
  name: myapp
  labels:
    app: myapp
    chart: example
  annotations:
    version: "0.1"
spec:
  serviceAccountName: waitfor
  selector:
    matchLabels:
      app: myapp
      chart: example
  serviceName: myapp
  template:
    metadata:
      labels:
        app: myapp
        chart: example
      annotations:
        version: "0.2.17"
    spec:
      initContainers:
        - name: wait-for-database
          image: ghcr.io/syntax3rror404/k8s-wait-for@sha256:7a58f56c216981117a196ba7a4949c19a2f2a9dba6933733b6e50e0539342a07
          imagePullPolicy: Always
          args:
            - "job"
            - "-n"
            - "myapp"
            - "-l"
            - "app=database"
      containers:
      - name: myapp
        image: ghcr.io/example/myapp:latest
        imagePullPolicy: Always
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: waitfor
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: waitfor-pod-reader
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: waitfor-pod-reader-binding
subjects:
  - kind: ServiceAccount
    name: waitfor
roleRef:
  kind: Role
  name: waitfor-pod-reader
  apiGroup: rbac.authorization.k8s.io
```

Example output while waiting for vault pods to be ready:
```
waitfor.exe pod -n vault -l app.kubernetes.io/instance=vault       
Info: Namespace set by user: "vault"
Info: Label set by user: "app.kubernetes.io/instance=vault"
Info: Trying in-cluster config...
Info: Error, falling back to kubeconfig...
Waiting for Pods in namespace vault with selector 'app.kubernetes.io/instance=vault'...
=== Fri, 19 Sep 2025 14:51:38 CEST ===========
State: Pod vault-0 --> Ready: true
State: Pod vault-1 --> Ready: true
State: Pod vault-2 --> Ready: true
==============================================
```

See help command for more information:
```
waitfor.exe -h                                    
This tool waits for kubernetes pods or jobs to be ready

a common usecase is to use it as init container to wait for other pods to be ready before starting the main application.
For example waiting for a database to be ready before starting the app to prevent errors.

Example:
  waitfor pod -n vault -l app.kubernetes.io/instance=vault
  waitfor job -n snipeit -l job=generate-app-key

Usage:
  waitfor [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  job         wait for a job to complete
  pod         wait for a pod to be ready

Flags:
  -h, --help               help for waitfor
  -l, --label string       Label to filter (required)
  -n, --namespace string   Namespace to use (default "default")
  -t, --timer int32        Wait time between checks (default 3)

Use "waitfor [command] --help" for more information about a command.
```