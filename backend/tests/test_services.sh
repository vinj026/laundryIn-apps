#!/bin/bash

BASE="http://localhost:8080"
echo "Restarting server for clean run..."
pkill -9 -f "laundryin-server" 2>/dev/null
sleep 1
./laundryin-server &
SERVER_PID=$!
sleep 3

echo -e "\n1. Getting auth tokens..."
OWNER_A=$(curl -s -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"+6281111111111","password":"SecurePass123!"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

OWNER_B=$(curl -s -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"+6282222222222","password":"SecurePass123!"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

echo "Owner A Token: ${OWNER_A:0:20}..."
echo "Owner B Token: ${OWNER_B:0:20}..."

echo -e "\n2. Fetching an Outlet belonging to Owner A..."
OUTLET_A=$(curl -s "$BASE/api/v1/outlets" \
  -H "Authorization: Bearer $OWNER_A" \
  | python3 -c "import sys,json; d=json.load(sys.stdin)['data']['data']; print(d[0]['id'] if d else 'NONE')")
echo "Outlet A: $OUTLET_A"

if [ "$OUTLET_A" = "NONE" ]; then
    echo "No outlet found for Owner A, creating one..."
    OUTLET_A=$(curl -s -X POST "$BASE/api/v1/outlets" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $OWNER_A" \
      -d '{"name":"A Laundry","address":"Jl A No 1","phone":"+6281111111111"}' \
      | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])")
    echo "Created new Outlet A: $OUTLET_A"
fi

echo -e "\n3. Testing Service Creation (Owner A -> Outlet A) [EXPECT: 201]"
CREATE_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/v1/services" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_A" \
  -d "{\"outlet_id\":\"$OUTLET_A\",\"name\":\"Cuci Kiloan Biasa\",\"price\":5000.00,\"unit\":\"KG\"}")
echo "$CREATE_RESP"
SERVICE_ID=$(echo "$CREATE_RESP" | head -n 1 | python3 -c "import sys,json; print(json.load(sys.stdin).get('data', {}).get('id', 'ERROR'))")
echo "Service ID: $SERVICE_ID"

echo -e "\n4. Testing Anti-IDOR: Owner B trying to create service in Outlet A [EXPECT: 404 (Not Found / Denied)]"
curl -s -w "\n%{http_code}\n" -X POST "$BASE/api/v1/services" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_B" \
  -d "{\"outlet_id\":\"$OUTLET_A\",\"name\":\"Hacker Service\",\"price\":100,\"unit\":\"PCS\"}"

echo -e "\n5. Testing Input Validation: Negative Price [EXPECT: 400]"
curl -s -w "\n%{http_code}\n" -X POST "$BASE/api/v1/services" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_A" \
  -d "{\"outlet_id\":\"$OUTLET_A\",\"name\":\"Free\",\"price\":-10,\"unit\":\"KG\"}"

echo -e "\n6. Testing Input Validation: Invalid Unit ENUM [EXPECT: 400]"
curl -s -w "\n%{http_code}\n" -X POST "$BASE/api/v1/services" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_A" \
  -d "{\"outlet_id\":\"$OUTLET_A\",\"name\":\"Cuci Karpet\",\"price\":100000,\"unit\":\"LITER\"}"

echo -e "\n7. Fetching Services for Outlet A [EXPECT: 200]"
curl -s -w "\n%{http_code}\n" "$BASE/api/v1/outlets/$OUTLET_A/services" \
  -H "Authorization: Bearer $OWNER_A" | grep "Cuci Kiloan Biasa" >/dev/null && echo "Service Found!" || echo "Service Miss"

echo -e "\n8. Anti-IDOR Update: Owner B tries to update Service owned by A [EXPECT: 404]"
curl -s -w "\n%{http_code}\n" -X PUT "$BASE/api/v1/services/$SERVICE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_B" \
  -d "{\"outlet_id\":\"$OUTLET_A\",\"name\":\"Hacked Pwd\",\"price\":1,\"unit\":\"KG\"}"

echo -e "\n9. Anti-IDOR Delete: Owner B tries to delete Service owned by A [EXPECT: 404]"
curl -s -w "\n%{http_code}\n" -X DELETE "$BASE/api/v1/services/$SERVICE_ID" \
  -H "Authorization: Bearer $OWNER_B"

# cleanup
pkill -9 -P $$
kill -9 $SERVER_PID 2>/dev/null
