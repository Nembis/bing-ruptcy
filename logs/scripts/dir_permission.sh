#!/bin/sh

set -e

ENV="$1"

echo "Running permission fix for environment: ${ENV:-prod}"

# always run permission for grafana
echo "Setting /tmp/grafana -> 472:0..."
chown -R 472:0 /tmp/grafana

# always run permission for loki
echo "Setting /tmp/loki -> 10001:10001..."
chown -R 10001:10001 /tmp/loki

# run only if in dev
if [ "$ENV" = "dev" ]; then 
    echo "Setting /tmp/pgadmin_data -> 5050:0..."
    chown -R 5050:0 /tmp/pgadmin_data
fi

echo "Successfully set all permissions"
exit 0
