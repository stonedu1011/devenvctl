#!/bin/bash

export CONSUL_ADDR="http://consul:8500"

initial_delay=3s
timeout=20s
backoff=1s

# wait for consul to accept HTTP requests
echo "Waiting Consul to accept HTTP requests..."
sleep $initial_delay
timeout $timeout bash << EOF || echo "Consul is not accepting HTTP requests after $timeout seconds."
while ! curl -f -s -X GET $CONSUL_ADDR/v1/agent/services 2>&1 1>/dev/null; do
  echo "Retry in $backoff..."
  sleep $backoff
done
EOF

if [ -z "$PROPERTIES" ]; then
  # try to deregister all services
  services=(`curl -f -s -X GET $CONSUL_ADDR/v1/agent/services | jq -r '.[]["ID"]'`)
  echo "Deregistering ${#services[@]} service instances ..."
  for id in ${services[@]}; do
    consul services deregister -http-addr=$CONSUL_ADDR -id=$id
  done
fi

if [ "$PROPERTIES" = "load" ]; then
  echo "load consul properties"
  curl --request PUT -g -k -v  --data 'usermanagementgoservice' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication/remoteservice.usermanagementservice.service 2>&1 1>/dev/null
  curl --request PUT -g -k -v  --data 'authservice' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication/remoteservice.authservice.service 2>&1 1>/dev/null
  curl --request PUT -g -k -v  --data 'http://localhost:8900/auth' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication/swagger.security.sso.baseurl 2>&1 1>/dev/null

  curl --request PUT -g -k -v  --data 'authservice' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/integration.security.server.name 2>&1 1>/dev/null
  curl --request PUT -g -k -v  --data 'auth' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/integration.security.server.contextPath 2>&1 1>/dev/null
  curl --request PUT -g -k -v  --data 'authservice' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/security.compatibility.auth-service-id 2>&1 1>/dev/null
  curl --request PUT -g -k -v  --data '${integration.security.client.client-id}' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/security.opaque-token.client-id 2>&1 1>/dev/null
  curl --request PUT -g -k -v  --data '${integration.security.client.client-secret}' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/security.opaque-token.client-secret 2>&1 1>/dev/null
  curl --request PUT -g -k -v  --data 'http://${security.compatibility.auth-service-id}/auth' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/security.opaque-token.introspection-base-url 2>&1 1>/dev/null
  curl --request PUT -g -k -v  --data 'dev-0' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/security.keys.jwt.id 2>&1 1>/dev/null
  curl --request PUT -g -k -v  --data 'usermanagementgoservice' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/remoteservice.usermanagementservice.service 2>&1 1>/dev/null
  curl --request PUT -g -k -v  --data 'usermanagementgoservice' $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/remoteservice.idm.service 2>&1 1>/dev/null
elif [ "$PROPERTIES" = "unload" ]; then
  echo "unload consul properties"
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication/remoteservice.usermanagementservice.service 2>&1 1>/dev/null
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication/remoteservice.authservice.service 2>&1 1>/dev/null
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication/swagger.security.sso.baseurl 2>&1 1>/dev/null

  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/integration.security.server.name 2>&1 1>/dev/null
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/integration.security.server.contextPath 2>&1 1>/dev/null
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/security.compatibility.auth-service-id 2>&1 1>/dev/null
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/security.opaque-token.client-id 2>&1 1>/dev/null
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/security.opaque-token.client-secret 2>&1 1>/dev/null
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/security.opaque-token.introspection-base-url 2>&1 1>/dev/null
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/security.keys.jwt.id 2>&1 1>/dev/null
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/remoteservice.usermanagementservice.service 2>&1 1>/dev/null
  curl --request DELETE $CONSUL_ADDR/v1/kv/userviceconfiguration/defaultapplication,indepauthservice/remoteservice.idm.service 2>&1 1>/dev/null
fi
