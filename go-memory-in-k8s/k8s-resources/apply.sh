kubectl create ns ns1 
kubectl apply -f resource_quota.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml 