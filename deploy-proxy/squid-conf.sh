#
#kubectl create namespace external
#
#cat <<EOF > ./proxy.conf
#http_port 3128
#
#acl SSL_ports port 443
#acl CONNECT method CONNECT
#
#http_access deny CONNECT !SSL_ports
#http_access allow localhost manager
#http_access deny manager
#http_access allow all
#
#coredump_dir /var/spool/squid
#EOF
#
#kubectl create configmap proxy-configmap -n external --from-file=squid.conf=./proxy.conf

kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: squid
  namespace: external
spec:
  replicas: 1
  selector:
    matchLabels:
      app: squid
  template:
    metadata:
      labels:
        app: squid
    spec:
      volumes:
      - name: proxy-config
        configMap:
          name: proxy-configmap
      - name: log
        emptyDir: {}
      containers:
      - name: squid
        image: sameersbn/squid:3.5.27
        imagePullPolicy: IfNotPresent
        volumeMounts:
        - name: proxy-config
          mountPath: /etc/squid
          readOnly: true
        - name: log
          mountPath: /var/log/squid
EOF