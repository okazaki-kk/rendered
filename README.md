# rendered

`helm template`や`helm install --dryrun --debug`の結果を、指定したディレクトリに保存するツール

## Example
```
$ helm template nginx-stable/nginx-ingress | rendered

$ tree output  
output
└── nginx-ingress
    └── templates
        ├── clusterrole.yaml
        ├── clusterrolebinding.yaml
        ├── controller-configmap.yaml
        ├── controller-deployment.yaml
        ├── controller-ingress-class.yaml
        ├── controller-leader-election-configmap.yaml
        ├── controller-lease.yaml
        ├── controller-role.yaml
        ├── controller-rolebinding.yaml
        ├── controller-service.yaml
        └── controller-serviceaccount.yaml

3 directories, 11 files
```
