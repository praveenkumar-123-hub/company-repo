# Test Script for Reducate Backend

$baseUrl = "http://localhost:8080"

# 1. Signup a User
echo "Signing up user..."
$signupUser = Invoke-RestMethod -Uri "$baseUrl/signup" -Method Post -Body '{"username":"johndoe", "password":"password123", "role":"User"}' -ContentType "application/json"
$signupUser | ConvertTo-Json

# 2. Signup an Admin
echo "Signing up admin..."
$signupAdmin = Invoke-RestMethod -Uri "$baseUrl/signup" -Method Post -Body '{"username":"admin", "password":"adminpassword", "role":"Admin"}' -ContentType "application/json"
$signupAdmin | ConvertTo-Json

# 3. Login as User
echo "Logging in as user..."
$loginUser = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body '{"username":"johndoe", "password":"password123"}' -ContentType "application/json"
$userToken = $loginUser.token
echo "User Token: $userToken"

# 4. Login as Admin
echo "Logging in as admin..."
$loginAdmin = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body '{"username":"admin", "password":"adminpassword"}' -ContentType "application/json"
$adminToken = $loginAdmin.token
echo "Admin Token: $adminToken"

# 5. Get User Profile
echo "Getting user profile..."
$headers = @{ "Authorization" = "Bearer $userToken" }
$profile = Invoke-RestMethod -Uri "$baseUrl/profile" -Method Get -Headers $headers
$profile | ConvertTo-Json

# 6. Get Admin Profile
echo "Getting admin profile..."
$headers = @{ "Authorization" = "Bearer $adminToken" }
$adminProfile = Invoke-RestMethod -Uri "$baseUrl/profile" -Method Get -Headers $headers
$adminProfile | ConvertTo-Json

# 7. Get All Users (Admin only) - Should Success
echo "Getting all users as admin..."
$users = Invoke-RestMethod -Uri "$baseUrl/users" -Method Get -Headers $headers
$users | ConvertTo-Json

# 8. Get All Users (User) - Should Fail
echo "Getting all users as user (expecting 403)..."
try {
    $headers = @{ "Authorization" = "Bearer $userToken" }
    Invoke-RestMethod -Uri "$baseUrl/users" -Method Get -Headers $headers
} catch {
    echo "Caught expected error: $($_.Exception.Message)"
}