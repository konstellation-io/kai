#/bin/sh

. ./scripts/functions.sh

HOST="http://api-kre.local"
SIGNIN_URL="$HOST/api/v1/auth/signin"
DEV_EMAIL="dev@local.local"
ADMIN_API_POD=$(kubectl -n kre get pod -l app=kre-local-admin-api -o custom-columns=":metadata.name" --no-headers)

echo "calling api..."
curl -s $SIGNIN_URL \
  -H 'pragma: no-cache' -H 'cache-control: no-cache' \
  -H 'accept: application/json, text/plain, */*' \
  -H 'content-type: application/json;charset=UTF-8' -H "origin: $HOST" \
  -H 'sec-fetch-site: same-site' -H 'sec-fetch-mode: cors' \
  -H "referer: $HOST/login" \
  --data-binary "{\"email\":\"$DEV_EMAIL\"}" >/dev/null 2>&1

sleep 0.5

WATCH_FILE=$(mktemp)

echo "watching $WATCH_FILE"
echo "pod $ADMIN_API_POD"
kubectl -n kre logs $ADMIN_API_POD | tail -n 100 >$WATCH_FILE

# REVERSE CONTENT OF THE FILE
LINK=$(sed -n '1{h;T;};G;h;$p;' $WATCH_FILE | egrep -oh "http://.*/signin/([^\"]*)")

rm $WATCH_FILE

echo "\n Login done. Open your browser at: \n\n 🌎 $LINK\n"

# Open browser automacally
nohup xdg-open $LINK >/dev/null 2>&1 &

echo_green "\n✔️  Done.\n\n"
