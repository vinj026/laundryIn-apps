import urllib.request
import urllib.error
import json
import time
import subprocess
import os

BASE_URL = "http://localhost:8080/api/v1"

def print_result(test_name, expected, actual_status, actual_body):
    status_str = f"[{'SUCCESS' if expected == actual_status else 'FAILED'}]"
    print(f"{status_str} {test_name} - Expected: {expected}, Got: {actual_status}")
    if expected != actual_status:
        try:
            body_json = json.loads(actual_body)
            print(f"   => Response: {json.dumps(body_json, indent=2)}")
        except:
            print(f"   => Response: {actual_body}")

def do_post(url, payload, headers=None):
    if headers is None:
        headers = {}
    headers["Content-Type"] = "application/json"
    data = json.dumps(payload).encode('utf-8')
    req = urllib.request.Request(url, data=data, headers=headers, method='POST')
    try:
        with urllib.request.urlopen(req) as response:
            return response.status, response.read().decode('utf-8')
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode('utf-8')
    except Exception as e:
        return 0, str(e)

def do_patch(url, payload, headers=None):
    if headers is None:
        headers = {}
    headers["Content-Type"] = "application/json"
    data = json.dumps(payload).encode('utf-8')
    req = urllib.request.Request(url, data=data, headers=headers, method='PATCH')
    try:
        with urllib.request.urlopen(req) as response:
            return response.status, response.read().decode('utf-8')
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode('utf-8')
    except Exception as e:
        return 0, str(e)

def do_delete(url, headers=None):
    if headers is None:
        headers = {}
    req = urllib.request.Request(url, headers=headers, method='DELETE')
    try:
        with urllib.request.urlopen(req) as response:
            return response.status, response.read().decode('utf-8')
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode('utf-8')
    except Exception as e:
        return 0, str(e)

def do_get(url, headers=None):
    if headers is None:
        headers = {}
    req = urllib.request.Request(url, headers=headers, method='GET')
    try:
        with urllib.request.urlopen(req) as response:
            return response.status, response.read().decode('utf-8')
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode('utf-8')
    except Exception as e:
        return 0, str(e)

def test_auth_empty_fields():
    print("\n--- TEST: AUTH EMPTY FIELDS ---")
    payloads = [
        ("Missing Name", {"phone": "+628111222333", "password": "Password123!", "role": "owner"}),
        ("Empty Name", {"name": "", "phone": "+628111222333", "password": "Password123!", "role": "owner"}),
        ("Missing Phone", {"name": "Test", "password": "Password123!", "role": "owner"}),
        ("Missing Password", {"name": "Test", "phone": "+628111222333", "role": "owner"}),
        ("Empty Password", {"name": "Test", "phone": "+628111222333", "password": "", "role": "owner"}),
    ]

    for name, payload in payloads:
        status, body = do_post(f"{BASE_URL}/auth/register", payload)
        print_result(name, 400, status, body)

def test_auth_weak_password():
    print("\n--- TEST: AUTH WEAK PASSWORD ---")
    payload = {
        "name": "Weak Pass User",
        "phone": "+628999888777",
        "password": "weakpassword", 
        "role": "owner"
    }
    status, body = do_post(f"{BASE_URL}/auth/register", payload)
    print_result("Weak Password Registration", 400, status, body)

def test_valid_register_and_login():
    print("\n--- TEST: VALID REGISTER AND LOGIN ---")
    payload = {
        "name": "Valid Owner",
        "phone": "+628333444555",
        "password": "StrongPassword123!",
        "role": "owner"
    }
    status, body = do_post(f"{BASE_URL}/auth/register", payload)
    if status not in [201, 409]:
        print_result("Valid Register", 201, status, body)
        return ""
    
    # Login
    login_payload = {
        "phone": "+628333444555",
        "password": "StrongPassword123!"
    }
    status, body = do_post(f"{BASE_URL}/auth/login", login_payload)
    print_result("Valid Login", 200, status, body)
    if status == 200:
        return json.loads(body)["data"]["token"]
    return ""

def test_outlet_validations(token):
    print("\n--- TEST: OUTLET CREATION VALIDATIONS ---")
    headers = {"Authorization": f"Bearer {token}"}
    
    payloads = [
        ("Empty Name", {"name": "", "address": "Valid Address Here", "phone": "+628333444555"}),
        ("Empty Address", {"name": "Valid Name", "address": "", "phone": "+628333444555"}),
        ("Empty Phone", {"name": "Valid Name", "address": "Valid Address Here", "phone": ""}),
        ("Invalid Phone Format", {"name": "Valid Name", "address": "Valid Address Here", "phone": "081234abcd"}),
    ]

    for name, payload in payloads:
        status, body = do_post(f"{BASE_URL}/outlets", payload, headers)
        print_result(name, 400, status, body)

def test_service_validations(token, outlet_id):
    print("\n--- TEST: SERVICE CREATION VALIDATIONS ---")
    headers = {"Authorization": f"Bearer {token}"}
    
    payloads = [
        ("Empty Name", {"outlet_id": outlet_id, "name": "", "price": "100.5", "unit": "KG"}),
        ("Empty Price", {"outlet_id": outlet_id, "name": "Valid", "price": "", "unit": "KG"}),
        ("Invalid Price Format", {"outlet_id": outlet_id, "name": "Valid", "price": "abcd", "unit": "KG"}),
        ("Empty Unit", {"outlet_id": outlet_id, "name": "Valid", "price": "100.5", "unit": ""}),
        ("Wrong Unit Case", {"outlet_id": outlet_id, "name": "Valid", "price": "100.5", "unit": "kg"}), # Should fail because allowed is KG
        ("Negative Price", {"outlet_id": outlet_id, "name": "Valid", "price": "-5", "unit": "KG"}),
        ("Invalid UUID format ID", {"outlet_id": "not-a-uuid", "name": "Valid", "price": "10", "unit": "KG"}),
    ]

    for name, payload in payloads:
        status, body = do_post(f"{BASE_URL}/services", payload, headers)
        print_result(name, 400, status, body)

def test_order_validations(token, outlet_id, service_id):
    print("\n--- TEST: ORDER & ZERO-TRUST VALIDATIONS ---")
    headers = {"Authorization": f"Bearer {token}"}
    
    payloads = [
        ("Empty Items List", {"outlet_id": outlet_id, "items": []}),
        ("Missing Required Qty", {"outlet_id": outlet_id, "items": [{"service_id": service_id, "qty": ""}]}),
        ("Negative Qty", {"outlet_id": outlet_id, "items": [{"service_id": service_id, "qty": "-2"}]}),
        ("Zero Qty", {"outlet_id": outlet_id, "items": [{"service_id": service_id, "qty": "0"}]}),
        ("Invalid Outlet ID in Order", {"outlet_id": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", "items": [{"service_id": service_id, "qty": "1"}]}),
    ]

    for name, payload in payloads:
        status, body = do_post(f"{BASE_URL}/orders", payload, headers)
        # Even if structure valid, negative/0 qty and invalid outlet should be 400
        print_result(name, 400, status, body)

def test_cross_outlet_idor(token):
    print("\n--- TEST: CROSS-OUTLET IDOR ---")
    headers = {"Authorization": f"Bearer {token}"}
    
    status_a, body_a = do_post(f"{BASE_URL}/outlets", {"name": "Outlet A", "address": "Valid Address", "phone": "+628111111111"}, headers)
    if status_a != 201: return print_result("Create Outlet A failed", 201, status_a, body_a)
    outlet_a_id = json.loads(body_a)["data"]["id"]
    
    status_b, body_b = do_post(f"{BASE_URL}/outlets", {"name": "Outlet B", "address": "Valid Address", "phone": "+628222222222"}, headers)
    if status_b != 201: return print_result("Create Outlet B failed", 201, status_b, body_b)
    outlet_b_id = json.loads(body_b)["data"]["id"]
    
    status_s, body_s = do_post(f"{BASE_URL}/services", {"outlet_id": outlet_b_id, "name": "Service B", "price": "1000", "unit": "KG"}, headers)
    if status_s != 201: return print_result("Create Service B failed", 201, status_s, body_s)
    service_b_id = json.loads(body_s)["data"]["id"]
    
    payload = {"outlet_id": outlet_a_id, "items": [{"service_id": service_b_id, "qty": "1"}]}
    status_o, body_o = do_post(f"{BASE_URL}/orders", payload, headers)
    print_result("Cross-Outlet Order Creation (Using service from other outlet)", 400, status_o, body_o)

def test_fsm_brute_force(token, outlet_id, service_id):
    print("\n--- TEST: FSM BRUTE FORCE ---")
    headers = {"Authorization": f"Bearer {token}"}
    
    payload = {"outlet_id": outlet_id, "items": [{"service_id": service_id, "qty": "1"}]}
    status_o, body_o = do_post(f"{BASE_URL}/orders", payload, headers)
    if status_o != 201: return print_result("FSM - Create Order failed", 201, status_o, body_o)
    order_id = json.loads(body_o)["data"]["id"]
    
    status, body = do_patch(f"{BASE_URL}/orders/{order_id}/status", {"status": "picked_up"}, headers)
    print_result("pending -> picked_up (Invalid Jump)", 400, status, body)
    
    status, body = do_patch(f"{BASE_URL}/orders/{order_id}/status", {"status": "process"}, headers)
    print_result("pending -> process (Valid)", 200, status, body)
    
    status, body = do_patch(f"{BASE_URL}/orders/{order_id}/status", {"status": "completed"}, headers)
    print_result("process -> completed (Valid)", 200, status, body)
    
    status, body = do_patch(f"{BASE_URL}/orders/{order_id}/status", {"status": "cancelled"}, headers)
    print_result("completed -> cancelled (Invalid Backwards)", 400, status, body)

def test_orphan_soft_delete(token, outlet_id):
    print("\n--- TEST: ORPHAN SOFT DELETE ---")
    headers = {"Authorization": f"Bearer {token}"}
    
    status_s, body_s = do_post(f"{BASE_URL}/services", {"outlet_id": outlet_id, "name": "Transient Service", "price": "1000", "unit": "KG"}, headers)
    if status_s != 201: return
    service_id = json.loads(body_s)["data"]["id"]
    
    payload = {"outlet_id": outlet_id, "items": [{"service_id": service_id, "qty": "1"}]}
    status_o, body_o = do_post(f"{BASE_URL}/orders", payload, headers)
    if status_o != 201: return
    
    status_d, body_d = do_delete(f"{BASE_URL}/services/{service_id}", headers)
    if status_d != 200: return
        
    status_get, body_get = do_get(f"{BASE_URL}/orders", headers)
    print_result("GET Orders after Service Deleted shouldn't crash", 200, status_get, "Success" if status_get == 200 else body_get)

def test_pagination_limit_abuse(token):
    print("\n--- TEST: PAGINATION LIMIT ABUSE ---")
    headers = {"Authorization": f"Bearer {token}"}
    status, body = do_get(f"{BASE_URL}/outlets?page=1&limit=999999", headers)
    print_result("Limit 999999 blocked by validator", 400, status, body)

def test_customer_role_block():
    print("\n--- TEST: CUSTOMER ROLE BLOCK ---")
    payload = {"name": "Iseng Customer", "phone": "+628555555555", "password": "StrongPassword123!", "role": "customer"}
    status, body = do_post(f"{BASE_URL}/auth/register", payload)
    if status == 409:
        status, body = do_post(f"{BASE_URL}/auth/login", {"phone": "+628555555555", "password": "StrongPassword123!"})
    if status not in [200, 201]: return print_result("Customer Auth failed", 200, status, body)
        
    token = json.loads(body)["data"]["token"]
    headers = {"Authorization": f"Bearer {token}"}
    
    status_out, body_out = do_post(f"{BASE_URL}/outlets", {"name": "Hacker Outlet", "address": "X", "phone": "+628111111111"}, headers)
    print_result("Customer prohibited from POST /outlets", 403, status_out, body_out)

def test_reports(token, outlet_id):
    print("\n--- TEST: ANALYTICS & REPORTS ---")
    headers = {"Authorization": f"Bearer {token}"}
    
    status, body = do_get(f"{BASE_URL}/reports/omzet", headers)
    print_result("Get Total Omzet", 200, status, "Success" if status == 200 else body)

    status, body = do_get(f"{BASE_URL}/reports/orders/summary", headers)
    print_result("Get Order Summary", 200, status, "Success" if status == 200 else body)

    status, body = do_get(f"{BASE_URL}/reports/services/top", headers)
    print_result("Get Top Services", 200, status, "Success" if status == 200 else body)

    # Test with full filters map
    status, body = do_get(f"{BASE_URL}/reports/omzet?outlet_id={outlet_id}&start_date=2023-01-01&end_date=2030-12-31", headers)
    print_result("Get Total Omzet with Filters", 200, status, "Success" if status == 200 else body)

def run_all_tests():
    # Attempt to wait for server to be ready
    try:
        urllib.request.urlopen("http://localhost:8080/ping")
        server_running_externally = True
    except:
        server_running_externally = False
        print("Starting server...")
        proc = subprocess.Popen(["go", "run", "cmd/api/main.go"], stdout=subprocess.PIPE, stderr=subprocess.PIPE, cwd="/home/vin/vin/Projects/laundryIn-app")
        time.sleep(3)
    
    test_auth_empty_fields()
    test_auth_weak_password()
    token = test_valid_register_and_login()
    if token:
        test_outlet_validations(token)
        test_pagination_limit_abuse(token)
        test_customer_role_block()
        test_cross_outlet_idor(token)
        
        # Test creating a valid outlet to get an ID
        status, body = do_post(f"{BASE_URL}/outlets", {"name": "TestOutlet", "address": "Valid Address", "phone": "+628444555666"}, {"Authorization": f"Bearer {token}"})
        outlet_id = ""
        if status == 201:
            outlet_id = json.loads(body)["data"]["id"]
            test_service_validations(token, outlet_id)
            
            # Test creating a valid service to get an ID
            s_status, s_body = do_post(f"{BASE_URL}/services", {"outlet_id": outlet_id, "name": "Valid Service", "price": "1000", "unit": "KG"}, {"Authorization": f"Bearer {token}"})
            if s_status == 201:
                service_id = json.loads(s_body)["data"]["id"]
                test_order_validations(token, outlet_id, service_id)
                test_fsm_brute_force(token, outlet_id, service_id)
                test_orphan_soft_delete(token, outlet_id)
            
            # Lastly check reports to ensure logic parsed our FSM jump created orders.
            test_reports(token, outlet_id)
    
    if not server_running_externally:
        proc.terminate()
    print("\nTests completed.")

if __name__ == "__main__":
    run_all_tests()
