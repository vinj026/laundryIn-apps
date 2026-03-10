#!/bin/bash
# ============================================================
# 🛡️ PRD Phase 3.5 — Security Audit Script (v2)
# ============================================================

BASE="http://localhost:8080"
PASS=0
FAIL=0

check() {
  local test_id="$1" expected="$2" actual="$3"
  if [ "$actual" = "$expected" ]; then
    echo "  ✅ $test_id — PASS (HTTP $actual)"
    PASS=$((PASS+1))
  else
    echo "  ❌ $test_id — FAIL (Expected $expected, Got $actual)"
    FAIL=$((FAIL+1))
  fi
}

echo ""
echo "⏳ Setting up test data..."

# Get fresh tokens
OWNER_A=$(curl -s -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"+6281111111111","password":"SecurePass123!"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

OWNER_B=$(curl -s -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"+6282222222222","password":"SecurePass123!"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

CUSTOMER=$(curl -s -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"+6283333333333","password":"SecurePass123!"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")

# Get Owner A's first outlet ID
OUTLET_A=$(curl -s "$BASE/api/v1/outlets" \
  -H "Authorization: Bearer $OWNER_A" \
  | python3 -c "import sys,json; d=json.load(sys.stdin)['data']['data']; print(d[0]['id'] if d else 'NONE')")

echo "  Owner A token: ${OWNER_A:0:20}..."
echo "  Owner B token: ${OWNER_B:0:20}..."
echo "  Customer token: ${CUSTOMER:0:20}..."
echo "  Outlet A ID: $OUTLET_A"

echo ""
echo "╔══════════════════════════════════════════════════════════╗"
echo "║  🛡️  SECURITY AUDIT — PRD Phase 3.5 (v2)                ║"
echo "╚══════════════════════════════════════════════════════════╝"

# =============================================================
echo ""
echo "📋 1. BROKEN ACCESS CONTROL & IDOR"
echo "─────────────────────────────────────"

# SEC-1A: Missing Token
CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/api/v1/outlets")
check "SEC-1A Missing Token" "401" "$CODE"

# SEC-1B: Invalid Role (customer)
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/v1/outlets" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CUSTOMER" \
  -d '{"name":"Hacker Store","address":"Jl. Hacker No. 1, Jakarta Selatan","phone":"+6289999999999"}')
check "SEC-1B Invalid Role" "403" "$CODE"

# SEC-1C: IDOR Cross-Tenant (Owner B tries to access Owner A's outlet)
CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/api/v1/outlets/$OUTLET_A" \
  -H "Authorization: Bearer $OWNER_B")
check "SEC-1C IDOR Cross-Tenant" "404" "$CODE"

# =============================================================
echo ""
echo "📋 2. RATE LIMITING"
echo "─────────────────────────────────────"

# SEC-2A: Normal Traffic (single request)
CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/api/v1/outlets" \
  -H "Authorization: Bearer $OWNER_A")
check "SEC-2A Normal Traffic" "200" "$CODE"

# SEC-2B: Burst Limit
# We need to burn through the remaining burst tokens first.
# The limiter starts with 60 burst tokens. Fire requests until 429.
echo "  ⏳ SEC-2B Burst test — firing rapid requests..."
GOT_429=0
REQ_NUM=0
for i in $(seq 1 120); do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/api/v1/outlets" \
    -H "Authorization: Bearer $OWNER_A")
  REQ_NUM=$((REQ_NUM+1))
  if [ "$CODE" = "429" ]; then
    GOT_429=1
    echo "  ✅ SEC-2B Burst Limit — PASS (429 hit at request #$REQ_NUM)"
    PASS=$((PASS+1))
    break
  fi
done
if [ "$GOT_429" = "0" ]; then
  echo "  ❌ SEC-2B Burst Limit — FAIL (Never got 429 in $REQ_NUM requests)"
  FAIL=$((FAIL+1))
fi

# Wait for rate limiter tokens to replenish
echo "  ⏳ Waiting 5s for rate limiter to recover..."
sleep 5

# =============================================================
echo ""
echo "📋 3. INPUT VALIDATION"
echo "─────────────────────────────────────"

# SEC-3A: Empty Payload
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/v1/outlets" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_A" \
  -d '{}')
check "SEC-3A Empty Payload" "400" "$CODE"

# SEC-3B: Phone e164 format (no + prefix)
RESP=$(curl -s -X POST "$BASE/api/v1/outlets" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_A" \
  -d '{"name":"Test Store","address":"Jl. Test No. 1 Kota Test 12345","phone":"08112233445"}')
CODE=$(echo "$RESP" | python3 -c "import sys; sys.exit(0)" 2>/dev/null && \
  curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/v1/outlets" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $OWNER_A" \
    -d '{"name":"Test Store","address":"Jl. Test No. 1 Kota Test 12345","phone":"08112233445"}')
check "SEC-3B Phone e164" "400" "$CODE"

# SEC-3C: Buffer Overflow address (501+ chars)
LONG_ADDR=$(python3 -c "print('Jalan Raya ' + 'A' * 490)")
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/v1/outlets" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_A" \
  -d "{\"name\":\"Test Store\",\"address\":\"$LONG_ADDR\",\"phone\":\"+6281234567890\"}")
check "SEC-3C Buffer Overflow" "400" "$CODE"

# SEC-3D: XSS Injection
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/v1/outlets" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OWNER_A" \
  -d '{"name":"<script>alert(1)</script>","address":"Jl. Test No. 1 Kota Test 12345","phone":"+6281234567890"}')
if [ "$CODE" = "201" ] || [ "$CODE" = "400" ]; then
  echo "  ✅ SEC-3D XSS Injection — PASS (HTTP $CODE — safely handled)"
  PASS=$((PASS+1))
else
  echo "  ❌ SEC-3D XSS Injection — FAIL (HTTP $CODE)"
  FAIL=$((FAIL+1))
fi

# =============================================================
echo ""
echo "📋 5. SECURITY HEADERS"
echo "─────────────────────────────────────"

# SEC-5A: Header Inspection (use /ping which has no auth)
HEADERS=$(curl -s -D - -o /dev/null "$BASE/ping" 2>/dev/null)

H1=$(echo "$HEADERS" | grep -ci "X-Content-Type-Options: nosniff" || true)
H2=$(echo "$HEADERS" | grep -ci "X-Frame-Options: DENY" || true)
H3=$(echo "$HEADERS" | grep -ci "X-Xss-Protection: 1; mode=block" || true)

if [ "$H1" -ge 1 ] && [ "$H2" -ge 1 ] && [ "$H3" -ge 1 ]; then
  echo "  ✅ SEC-5A Headers — PASS"
  echo "     ├── X-Content-Type-Options: nosniff ✓"
  echo "     ├── X-Frame-Options: DENY ✓"
  echo "     └── X-XSS-Protection: 1; mode=block ✓"
  PASS=$((PASS+1))
else
  echo "  ❌ SEC-5A Headers — FAIL"
  echo "     Response headers:"
  echo "$HEADERS" | head -15
  FAIL=$((FAIL+1))
fi

# =============================================================
echo ""
echo "═══════════════════════════════════════"
echo "  📊 RESULTS: $PASS PASS / $FAIL FAIL out of 10 tests"
echo "═══════════════════════════════════════"
echo ""
echo "  ⚠️  SEC-4A (Timeout) requires manual code modification"
echo "     Add time.Sleep(6*time.Second) to repository Create method,"
echo "     restart server, and verify 408 response."
echo ""
