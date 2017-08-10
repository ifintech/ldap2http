#!/bin/sh

exec ldap-auth -host="$HOST" -port="$PORT" -auth_url="$AUTH_URL" -auth_token="$AUTH_URL"
