## 部署

```
kubectl create -f ./ -n kube-system
kubectl label ns <空间名> pod-admission-webhook-injection=enabled  // 开启对应空间注入，通过 label 实现
```
