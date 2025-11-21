#!/usr/bin/env pwsh
# OrangePi HTTP Service - Complete API Test Script

param(
    [string]$BaseUrl = "http://127.0.0.1:8080",
    [string]$Username = "admin",
    [string]$Password = "123456"
)

$PASSED = 0
$FAILED = 0
$TOTAL = 0

function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Uri,
        [hashtable]$Headers = @{},
        [string]$Body = $null
    )
    
    $script:TOTAL++
    Write-Host "[$TOTAL] $Name" -ForegroundColor Cyan
    
    try {
        $params = @{
            Uri         = $Uri
            Method      = $Method
            ContentType = "application/json"
            ErrorAction = "Stop"
        }
        
        if ($Headers -and $Headers.Count -gt 0) {
            $params["Headers"] = $Headers
        }
        
        if ($Body) {
            $params["Body"] = $Body
        }
        
        $response = Invoke-RestMethod @params
        Write-Host "    ✓ SUCCESS" -ForegroundColor Green
        Write-Host "    Response:" -ForegroundColor Gray
        Write-Host ($response | ConvertTo-Json -Depth 10) -ForegroundColor Gray
        $script:PASSED++
        
        return $response
    }
    catch {
        Write-Host "    ✗ FAILED: $($_.Exception.Message)" -ForegroundColor Red
        $script:FAILED++
        return $null
    }
}

Write-Host "=======================================" -ForegroundColor Cyan
Write-Host "OrangePi HTTP Service - API Test" -ForegroundColor Cyan
Write-Host "=======================================" -ForegroundColor Cyan
Write-Host ""

# 1. Health Check
Test-Endpoint "1. Health Check" "Get" "$BaseUrl/health"

# 2. Generate Video Token
Write-Host ""
$videoBody = @{
    building_id = "ismart_001"
    channels    = @("channel1", "channel2", "channel3")
} | ConvertTo-Json

$videoResp = Test-Endpoint "2. Generate Video Token" "Post" "$BaseUrl/api/auth/public" @{} $videoBody

# 3. Admin Login
Write-Host ""
$loginBody = @{
    username = $Username
    password = $Password
} | ConvertTo-Json

$loginResp = Test-Endpoint "3. Admin Login" "Post" "$BaseUrl/api/auth/login" @{} $loginBody

if ($null -eq $loginResp) {
    Write-Host ""
    Write-Host "Login failed! Cannot continue with protected endpoints." -ForegroundColor Red
    exit 1
}

$jwtToken = $loginResp.data.accessToken
$authHeaders = @{"Authorization" = "Bearer $jwtToken" }

# 4. Query Admin List
Write-Host ""
Test-Endpoint "4. Query Admin List" "Get" "$BaseUrl/api/admin" $authHeaders

# 5. Query Building List
Write-Host ""
Test-Endpoint "5. Query Building List" "Get" "$BaseUrl/api/building" $authHeaders

# 6. Query Device List
Write-Host ""
Test-Endpoint "6. Query Device List" "Get" "$BaseUrl/api/device" $authHeaders

# 7. Get Device Info
Write-Host ""
Test-Endpoint "7. Get Device Info" "Get" "$BaseUrl/api/device/info" $authHeaders

# 8. Create Admin
Write-Host ""
$testId = Get-Random -Minimum 1000 -Maximum 9999
$adminBody = @{
    username = "test_$testId"
    password = "testpass123"
} | ConvertTo-Json

$adminResp = Test-Endpoint "8. Create Admin" "Post" "$BaseUrl/api/admin" $authHeaders $adminBody
if ($adminResp) {
    $testAdminId = $adminResp.data.id
}

# 9. Create Building
Write-Host ""
$buildId = Get-Random -Minimum 10000 -Maximum 99999
$buildBody = @{
    ismartid = "test_ismart_$buildId"
    name     = "Test Building"
    remark   = "Test"
} | ConvertTo-Json

$buildResp = Test-Endpoint "9. Create Building" "Post" "$BaseUrl/api/building" $authHeaders $buildBody
if ($buildResp) {
    $testBuildId = $buildResp.data.id
    $testISmartId = $buildResp.data.ismartid
}

# 10. Create Device (if building exists)
Write-Host ""
if ($buildResp -and $testISmartId) {
    $devId = Get-Random -Minimum 100 -Maximum 999
    $devBody = @{
        ismartid                       = $testISmartId
        name                           = "TestDevice"
        icctv_auth_service_remote_port = (30000 + $devId)
        ssh_remote_port                = (20000 + $devId)
        is_active                      = $true
    } | ConvertTo-Json
    
    $devResp = Test-Endpoint "10. Create Device" "Post" "$BaseUrl/api/device" $authHeaders $devBody
    if ($devResp) {
        $testDevId = $devResp.data.id
    }
}

# 11. Update Public Net Config
Write-Host ""
$configBody = @{
    external_ip = "192.168.1.100"
} | ConvertTo-Json

Test-Endpoint "11. Update Public Net Config" "Put" "$BaseUrl/api/publicnet/config" $authHeaders $configBody

# 12. Test Building-OrangePi Binding Management
Write-Host ""
Write-Host "Testing Building-OrangePi Binding Management..." -ForegroundColor Yellow

# 13.1 Bind OrangePi to Building
if ($testBuildId -and $testDevId) {
    $bindBody = @{
        building_id = $testBuildId
        orangepi_id = $testDevId
    } | ConvertTo-Json
    
    Test-Endpoint "13.1 Bind OrangePi to Building" "Post" "$BaseUrl/api/building/bind" $authHeaders $bindBody
}

# 13.2 Update Binding (if we have multiple buildings)
if ($testDevId) {
    # Try to bind to building ID 1 (assuming it exists)
    $updateBindBody = @{
        orangepi_id     = $testDevId
        new_building_id = 1
    } | ConvertTo-Json
    
    Test-Endpoint "13.2 Update OrangePi Binding" "Put" "$BaseUrl/api/building/bind" $authHeaders $updateBindBody
}

# 14. Test OrangePi Remote Management
Write-Host ""
Write-Host "Testing OrangePi Remote Management..." -ForegroundColor Yellow

if ($targetDeviceId) {
    # 14.1 Remote Health Check
    Test-Endpoint "14.1 Remote Health Check" "Get" "$BaseUrl/api/orangepi/remote/health?id=$targetDeviceId" $authHeaders

    # 14.2 Remote Get Info
    Test-Endpoint "14.2 Remote Get Device Info" "Get" "$BaseUrl/api/orangepi/remote/info?id=$targetDeviceId" $authHeaders

    # 14.3 Remote Update Ports
    $remotePortBody = @{
        id                             = $targetDeviceId
        ssh_remote_port                = 20001
        icctv_auth_service_remote_port = 30001
    } | ConvertTo-Json

    Test-Endpoint "14.3 Remote Update Ports" "Post" "$BaseUrl/api/orangepi/remote/ports" $authHeaders $remotePortBody
}

# Clean up - Delete test data
Write-Host ""
Write-Host "Cleaning up test data..." -ForegroundColor Yellow

# Unbind device before deletion
if ($testDevId) {
    $unbindBody = @{orangepi_id = $testDevId } | ConvertTo-Json
    Test-Endpoint "Unbind Test Device" "Post" "$BaseUrl/api/building/unbind" $authHeaders $unbindBody
    
    $delDevBody = @{id = $testDevId } | ConvertTo-Json
    Test-Endpoint "Delete Test Device" "Delete" "$BaseUrl/api/device" $authHeaders $delDevBody
}

if ($testBuildId) {
    $delBuildBody = @{id = $testBuildId } | ConvertTo-Json
    Test-Endpoint "Delete Test Building" "Delete" "$BaseUrl/api/building" $authHeaders $delBuildBody
}

if ($testAdminId) {
    $delAdminBody = @{id = $testAdminId } | ConvertTo-Json
    Test-Endpoint "Delete Test Admin" "Delete" "$BaseUrl/api/admin" $authHeaders $delAdminBody
}

# Summary
Write-Host ""
Write-Host "=======================================" -ForegroundColor Cyan
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "=======================================" -ForegroundColor Cyan

Write-Host ""
Write-Host "Total:  $TOTAL"
Write-Host "Passed: $PASSED" -ForegroundColor Green
Write-Host "Failed: $FAILED" -ForegroundColor Red

if ($FAILED -eq 0) {
    Write-Host ""
    Write-Host "All tests PASSED! ✓" -ForegroundColor Green
    exit 0
}
else {
    Write-Host ""
    Write-Host "Some tests FAILED! ✗" -ForegroundColor Red
    exit 1
}

