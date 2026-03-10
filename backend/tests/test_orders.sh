#!/bin/bash

BASE="http://localhost:8080"
echo "Restarting server for clean run..."
pkill -9 -f "laundryin-server" 2>/dev/null
sleep 1
./laundryin-server &
SERVER_PID=$!
sleep 2

echo -e "\n1. Getting auth tokens..."
OWNER_A=$(curl -s -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"+6281111111111","password":"SecurePass123!"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

CUST_A=$(curl -s -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"+6289999999999","password":"SecurePass123!"}' \
  | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data', {}).get('token', 'NONE'))")

if [ "$CUST_A" = "NONE" ]; then
    echo "Creating Customer A..."
    CUST_A=$(curl -s -X POST "$BASE/api/v1/auth/register" \
      -H "Content-Type: application/json" \
      -d '{"name":"Cust A", "phone":"+6289999999999", "password":"SecurePass123!", "role":"customer"}' \
      | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")
fi

echo "Owner A Token: ${OWNER_A:0:20}..."
echo "Cust A Token: ${CUST_A:0:20}..."

echo -e "\n2. Fetching Outlet A..."
OUTLET_A=$(curl -s "$BASE/api/v1/outlets" \
  -H "Authorization: Bearer $OWNER_A" \
  | python3 -c "import sys,json; d=json.load(sys.stdin)['data']['data']; print(d[0]['id'] if d else 'NONE')")

echo -e "\n3. Ensuring Service exists in Outlet A..."
SERVICE_A=$(curl -s "$BASE/api/v1/outlets/$OUTLET_A/services" \
  -H "Authorization: Bearer $OWNER_A" \
  | python3 -c "import sys,json; d=json.load(sys.stdin)['data']; print(d[0]['id'] if d else 'NONE')")

if [ "$SERVICE_A" = "NONE" ]; then
    SERVICE_A=$(curl -s -X POST "$BASE/api/v1/services" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $OWNER_A" \
      -d "{\"outlet_id\":\"$OUTLET_A\",\"name\":\"Cuci Decimal\",\"price\":5000.55,\"unit\":\"KG\"}" \
      | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])")
fi
echo "Service A: $SERVICE_A"

echo -e "\n4. Transaction: Happy Path Checkout (Decimal qty: 3.333)"
ORDER_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/v1/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CUST_A" \
  -d "{\"outlet_id\":\"$OUTLET_A\",\"items\":[{\"service_id\":\"$SERVICE_A\",\"qty\":3.333}]}")
echo "$ORDER_RESP"
ORDER_ID=$(echo "$ORDER_RESP" | head -n 1 | python3 -c "import sys,json; print(json.load(sys.stdin).get('data', {}).get('id', 'ERROR'))")
echo "Order ID: $ORDER_ID"

echo -e "\n5. Validation: State Machine (pending -> picked_up) [EXPECT 400 Invalid State]"
curl -s -w "\n%{http_code}\n" -X PATCH "$BASE/api/v1/orders/$ORDER_ID/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_A" \
  -d '{"status":"picked_up"}'

echo -e "\n6. Validation: State Machine (pending -> process) [EXPECT 200]"
curl -s -w "\n%{http_code}\n" -X PATCH "$BASE/api/v1/orders/$ORDER_ID/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_A" \
  -d '{"status":"process"}'

echo -e "\n7. Validation: Partial Failure & ACID Rollback"
# Create a fake nonexistent service ID
FAKE_SERVICE="11111111-2222-3333-4444-555555555555"
ROLLBACK_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/v1/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CUST_A" \
  -d "{\"outlet_id\":\"$OUTLET_A\",\"items\":[{\"service_id\":\"$SERVICE_A\",\"qty\":1}, {\"service_id\":\"$FAKE_SERVICE\",\"qty\":1}]}")
echo "$ROLLBACK_RESP"

echo -e "\n8. Confirming Rollback: Getting User Orders"
curl -s "$BASE/api/v1/orders" \
  -H "Authorization: Bearer $CUST_A" | grep -o '"id"' | wc -l > order_count.txt
echo "Total Orders found (should be 1 due to rollback of second): $(cat order_count.txt)"

# cleanup
pkill -9 -P $$
kill -9 $SERVER_PID 2>/dev/null
