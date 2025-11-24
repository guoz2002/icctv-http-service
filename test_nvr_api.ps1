#!/usr/bin/env pwsh
# ICCTV HTTP Service - Complete API Test Script
# 测试所有 API 接口

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
    Write-Host "[$script:TOTAL] $Name" -ForegroundColor Cyan
    
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
        $script:PASSED++
        
        return $response
    }
    catch {
        Write-Host "    ✗ FAILED: $($_.Exception.Message)" -ForegroundColor Red
        $script:FAILED++
        return $null
    }
}

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "ICCTV HTTP Service - Complete API Test" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# ==================== 核心管理接口 ====================
Write-Host "========== 核心管理接口 ==========" -ForegroundColor Yellow
Write-Host ""

# 1. Health Check
$healthResp = Test-Endpoint "Health Check" "GET" "$BaseUrl/health"
Write-Host ""

# 3. Admin Login
$loginBody = @{
    username = $Username
    password = $Password
} | ConvertTo-Json

$loginResp = Test-Endpoint "Admin Login" "POST" "$BaseUrl/api/auth/login" @{} $loginBody

if ($null -eq $loginResp) {
    Write-Host ""
    Write-Host "❌ Login failed! Cannot continue with protected endpoints." -ForegroundColor Red
    exit 1
}

$jwtToken = $loginResp.data.accessToken
$authHeaders = @{"Authorization" = "Bearer $jwtToken" }
Write-Host ""

# ==================== 管理员账户管理 ====================
Write-Host "========== 管理员账户管理 (Adminer) ==========" -ForegroundColor Yellow
Write-Host ""

# 4. Query Admin List
$adminListResp = Test-Endpoint "Query Admin List" "GET" "$BaseUrl/api/admin" $authHeaders
Write-Host ""

# 5. Create Admin
$testId = Get-Random -Minimum 1000 -Maximum 9999
$adminBody = @{
    username = "test_admin_$testId"
    password = "testpass123"
} | ConvertTo-Json

$adminResp = Test-Endpoint "Create Admin" "POST" "$BaseUrl/api/admin" $authHeaders $adminBody
if ($adminResp) {
    $testAdminId = $adminResp.data.id
    Write-Host "    Created Admin ID: $testAdminId" -ForegroundColor Gray
}
Write-Host ""

# 6. Update Admin (if created)
if ($testAdminId) {
    $updateAdminBody = @{
        id       = $testAdminId
        username = "updated_admin_$testId"
    } | ConvertTo-Json
    
    Test-Endpoint "Update Admin" "PUT" "$BaseUrl/api/admin" $authHeaders $updateAdminBody
    Write-Host ""
}

# ==================== 建筑信息管理 ====================
Write-Host "========== 建筑信息管理 (Building) ==========" -ForegroundColor Yellow
Write-Host ""

# 12. Query Building List
$buildingListResp = Test-Endpoint "Query Building List" "GET" "$BaseUrl/api/building" $authHeaders
Write-Host ""

# 13. Create Building
$buildId = Get-Random -Minimum 10000 -Maximum 99999
$buildBody = @{
    ismartid = "test_building_$buildId"
    name     = "测试楼栋_$buildId"
    remark   = "自动化测试创建"
} | ConvertTo-Json

$buildResp = Test-Endpoint "Create Building" "POST" "$BaseUrl/api/building" $authHeaders $buildBody
if ($buildResp) {
    $testBuildId = $buildResp.data.id
    $testISmartId = $buildResp.data.ismartid
    Write-Host "    Created Building ID: $testBuildId" -ForegroundColor Gray
}
Write-Host ""

# 14. Update Building (if created)
if ($testBuildId) {
    $updateBuildBody = @{
        id     = $testBuildId
        name   = "更新后的楼栋_$buildId"
        remark = "已更新"
    } | ConvertTo-Json
    
    Test-Endpoint "Update Building" "PUT" "$BaseUrl/api/building" $authHeaders $updateBuildBody
    Write-Host ""
}

# ==================== OrangePi设备管理 ====================
Write-Host "========== OrangePi 设备管理 ==========" -ForegroundColor Yellow
Write-Host ""

# 8. Query Device List
$deviceListResp = Test-Endpoint "Query OrangePi List" "GET" "$BaseUrl/api/device" $authHeaders
Write-Host ""

# 9. Create Device (if building exists)
if ($testISmartId) {
    $devId = Get-Random -Minimum 100 -Maximum 999
    $devBody = @{
        ismartid                       = $testISmartId
        name                           = "TestOrangePi_$devId"
        icctv_auth_service_remote_port = (30000 + $devId)
        ssh_remote_port                = (20000 + $devId)
        is_active                      = $true
    } | ConvertTo-Json
    
    $devResp = Test-Endpoint "Create OrangePi" "POST" "$BaseUrl/api/device" $authHeaders $devBody
    if ($devResp) {
        $testDevId = $devResp.data.id
        Write-Host "    Created Device ID: $testDevId" -ForegroundColor Gray
    }
    Write-Host ""
    
    # 10. Update Device (if created)
    if ($testDevId) {
        $updateDevBody = @{
            id                             = $testDevId
            name                           = "UpdatedOrangePi_$devId"
            icctv_auth_service_remote_port = (30000 + $devId + 1)
        } | ConvertTo-Json
        
        Test-Endpoint "Update OrangePi" "PUT" "$BaseUrl/api/device" $authHeaders $updateDevBody
        Write-Host ""
    }
}

# ==================== NVR 管理 ====================
Write-Host "========== NVR (网络硬盘录像机) 管理 ==========" -ForegroundColor Yellow
Write-Host ""

# 查询 NVR 列表
$nvrListResp = Test-Endpoint "Query NVR List" "GET" "$BaseUrl/api/nvr" $authHeaders
Write-Host ""

# 创建 NVR (包含 RTSPUrls)
if ($testBuildId) {
    $nvrBody = @{
        name        = "测试NVR设备"
        url         = "192.168.1.100:8080"
        building_id = $testBuildId
        admin_user  = @{
            name     = "admin"
            password = "admin123"
        }
        users       = @(
            @{
                name     = "operator1"
                password = "pass123"
            },
            @{
                name     = "operator2"
                password = "pass456"
            }
        )
        rtsp_urls   = @(
            @{
                channel = 1
                url     = "rtsp://192.168.1.100:554/stream1"
            },
            @{
                channel = 2
                url     = "rtsp://192.168.1.100:554/stream2"
            },
            @{
                channel = 3
                url     = "rtsp://192.168.1.100:554/stream3"
            },
            @{
                channel = 4
                url     = "rtsp://192.168.1.100:554/stream4"
            }
        )
    } | ConvertTo-Json -Depth 10
    
    $nvrResp = Test-Endpoint "Create NVR with RTSPUrls" "POST" "$BaseUrl/api/nvr" $authHeaders $nvrBody
    if ($nvrResp) {
        $testNvrId = $nvrResp.data.id
        Write-Host "    Created NVR ID: $testNvrId" -ForegroundColor Gray
        Write-Host "    RTSP Channels: $($nvrResp.data.rtsp_urls.Count)" -ForegroundColor Gray
    }
    Write-Host ""
    
    # 查询 NVR 详情
    if ($testNvrId) {
        $nvrDetailResp = Test-Endpoint "Get NVR Details" "GET" "$BaseUrl/api/nvr?id=$testNvrId" $authHeaders
        if ($nvrDetailResp) {
            Write-Host "    Name: $($nvrDetailResp.data.name)" -ForegroundColor Gray
            Write-Host "    URL: $($nvrDetailResp.data.url)" -ForegroundColor Gray
            Write-Host "    RTSP Channels:" -ForegroundColor Gray
            foreach ($rtsp in $nvrDetailResp.data.rtsp_urls) {
                Write-Host "      CH$($rtsp.channel): $($rtsp.url)" -ForegroundColor Gray
            }
        }
        Write-Host ""
        
        # 更新 NVR (修改 RTSPUrls)
        $updateNvrBody = @{
            name      = "更新后的NVR设备"
            rtsp_urls = @(
                @{
                    channel = 1
                    url     = "rtsp://192.168.1.100:554/updated_stream1"
                },
                @{
                    channel = 2
                    url     = "rtsp://192.168.1.100:554/updated_stream2"
                },
                @{
                    channel = 5
                    url     = "rtsp://192.168.1.100:554/new_stream5"
                }
            )
        } | ConvertTo-Json -Depth 10
        
        $updateNvrResp = Test-Endpoint "Update NVR RTSPUrls" "PUT" "$BaseUrl/api/nvr?id=$testNvrId" $authHeaders $updateNvrBody
        if ($updateNvrResp) {
            Write-Host "    Updated RTSP Channels: $($updateNvrResp.data.rtsp_urls.Count)" -ForegroundColor Gray
        }
        Write-Host ""
    }
}

# ==================== Bind 接口测试 ====================
Write-Host "========== Bind 关联关系管理 ==========" -ForegroundColor Yellow
Write-Host ""

# Building-OrangePi 绑定
if ($testBuildId -and $testDevId) {
    $bindBody = @{
        building_id = $testBuildId
        orangepi_id = $testDevId
    } | ConvertTo-Json
    
    Test-Endpoint "Bind OrangePi to Building" "POST" "$BaseUrl/api/bind/building-orangepi" $authHeaders $bindBody
    Write-Host ""
    
    # 查询 Building 关联的 OrangePi
    Test-Endpoint "Get Building's OrangePis" "GET" "$BaseUrl/api/bind/building-orangepi/$testBuildId" $authHeaders
    Write-Host ""
}

# Building-NVR 绑定
if ($testBuildId -and $testNvrId) {
    $bindNvrBody = @{
        building_id = $testBuildId
        nvr_id      = $testNvrId
    } | ConvertTo-Json
    
    Test-Endpoint "Bind NVR to Building" "POST" "$BaseUrl/api/bind/building-nvr" $authHeaders $bindNvrBody
    Write-Host ""
    
    # 查询 Building 关联的 NVR
    Test-Endpoint "Get Building's NVRs" "GET" "$BaseUrl/api/bind/building-nvr/$testBuildId" $authHeaders
    Write-Host ""
}

# ==================== 清理测试数据 ====================
Write-Host "========== 清理测试数据 ==========" -ForegroundColor Yellow
Write-Host ""

$cleanup = Read-Host "是否删除测试数据? (y/n)"
if ($cleanup -eq "y" -or $cleanup -eq "Y") {
    # 解绑并删除 OrangePi
    if ($testDevId) {
        if ($testBuildId) {
            $unbindBody = @{
                building_id = $testBuildId
                orangepi_id = $testDevId
            } | ConvertTo-Json
            Test-Endpoint "Unbind OrangePi" "DELETE" "$BaseUrl/api/bind/building-orangepi" $authHeaders $unbindBody
            Write-Host ""
        }
        
        $delDevBody = @{id = $testDevId } | ConvertTo-Json
        Test-Endpoint "Delete OrangePi" "DELETE" "$BaseUrl/api/device" $authHeaders $delDevBody
        Write-Host ""
    }
    
    # 解绑并删除 NVR
    if ($testNvrId) {
        if ($testBuildId) {
            $unbindNvrBody = @{
                building_id = $testBuildId
                nvr_id      = $testNvrId
            } | ConvertTo-Json
            Test-Endpoint "Unbind NVR" "DELETE" "$BaseUrl/api/bind/building-nvr" $authHeaders $unbindNvrBody
            Write-Host ""
        }
        
        $delNvrBody = @{id = $testNvrId } | ConvertTo-Json
        Test-Endpoint "Delete NVR" "DELETE" "$BaseUrl/api/nvr" $authHeaders $delNvrBody
        Write-Host ""
    }
    
    # 删除 Building
    if ($testBuildId) {
        $delBuildBody = @{id = $testBuildId } | ConvertTo-Json
        Test-Endpoint "Delete Building" "DELETE" "$BaseUrl/api/building" $authHeaders $delBuildBody
        Write-Host ""
    }
    
    # 删除 Admin
    if ($testAdminId) {
        $delAdminBody = @{id = $testAdminId } | ConvertTo-Json
        Test-Endpoint "Delete Admin" "DELETE" "$BaseUrl/api/admin" $authHeaders $delAdminBody
        Write-Host ""
    }
}
else {
    Write-Host "测试数据保留:" -ForegroundColor Cyan
    if ($testAdminId) { Write-Host "  Admin ID: $testAdminId" -ForegroundColor Gray }
    if ($testBuildId) { Write-Host "  Building ID: $testBuildId" -ForegroundColor Gray }
    if ($testDevId) { Write-Host "  OrangePi ID: $testDevId" -ForegroundColor Gray }
    if ($testNvrId) { Write-Host "  NVR ID: $testNvrId" -ForegroundColor Gray }
    Write-Host ""
}

# ==================== 测试总结 ====================
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Total:  $TOTAL"
Write-Host "Passed: $PASSED" -ForegroundColor Green
Write-Host "Failed: $FAILED" -ForegroundColor Red
Write-Host ""

if ($FAILED -eq 0) {
    Write-Host "✅ All tests PASSED!" -ForegroundColor Green
    exit 0
}
else {
    Write-Host "❌ Some tests FAILED!" -ForegroundColor Red
    exit 1
}

  Write-Host "   ⚠ Building Creation Failed (may already exist)" -ForegroundColor Yellow
  $testBuildId = 1  # 使用默认ID
}

Write-Host ""

# 3. 创建 NVR (包含 RTSPUrls)
Write-Host "3. Create NVR with RTSPUrls" -ForegroundColor Cyan
$nvrBody = @{
  name        = "测试NVR设备"
  url         = "192.168.1.100:8080"
  building_id = $testBuildId
  admin_user  = @{
    name     = "admin"
    password = "admin123"
  }
  users       = @(
    @{
      name     = "operator1"
      password = "pass123"
    },
    @{
      name     = "operator2"
      password = "pass456"
    }
  )
  rtsp_urls   = @(
    @{
      channel = 1
      url     = "rtsp://192.168.1.100:554/stream1"
    },
    @{
      channel = 2
      url     = "rtsp://192.168.1.100:554/stream2"
    },
    @{
      channel = 3
      url     = "rtsp://192.168.1.100:554/stream3"
    },
    @{
      channel = 4
      url     = "rtsp://192.168.1.100:554/stream4"
    }
  )
} | ConvertTo-Json -Depth 10

try {
  $createResp = Invoke-RestMethod -Uri "$BaseUrl/api/nvr" -Method Post -Headers $authHeaders -ContentType "application/json" -Body $nvrBody
  $testNvrId = $createResp.data.id
  Write-Host "   ✓ NVR Created: ID=$testNvrId" -ForegroundColor Green
  Write-Host "   Response:" -ForegroundColor Gray
  Write-Host ($createResp | ConvertTo-Json -Depth 10) -ForegroundColor Gray
}
catch {
  Write-Host "   ✗ NVR Creation FAILED: $($_.Exception.Message)" -ForegroundColor Red
  if ($_.ErrorDetails.Message) {
    Write-Host "   Details: $($_.ErrorDetails.Message)" -ForegroundColor Red
  }
  exit 1
}

Write-Host ""

# 4. 查询 NVR 详情
Write-Host "4. Get NVR Details" -ForegroundColor Cyan
try {
  $getResp = Invoke-RestMethod -Uri "$BaseUrl/api/nvr?id=$testNvrId" -Method Get -Headers $authHeaders
  Write-Host "   ✓ NVR Details Retrieved" -ForegroundColor Green
  Write-Host "   Name: $($getResp.data.name)" -ForegroundColor Gray
  Write-Host "   URL: $($getResp.data.url)" -ForegroundColor Gray
  Write-Host "   RTSP Channels: $($getResp.data.rtsp_urls.Count)" -ForegroundColor Gray
  Write-Host "   Response:" -ForegroundColor Gray
  Write-Host ($getResp | ConvertTo-Json -Depth 10) -ForegroundColor Gray
}
catch {
  Write-Host "   ✗ Get NVR FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# 5. 更新 NVR (修改 RTSPUrls)
Write-Host "5. Update NVR RTSPUrls" -ForegroundColor Cyan
$updateBody = @{
  rtsp_urls = @(
    @{
      channel = 1
      url     = "rtsp://192.168.1.100:554/updated_stream1"
    },
    @{
      channel = 2
      url     = "rtsp://192.168.1.100:554/updated_stream2"
    },
    @{
      channel = 5
      url     = "rtsp://192.168.1.100:554/new_stream5"
    },
    @{
      channel = 6
      url     = "rtsp://192.168.1.100:554/new_stream6"
    }
  )
} | ConvertTo-Json -Depth 10

try {
  $updateResp = Invoke-RestMethod -Uri "$BaseUrl/api/nvr?id=$testNvrId" -Method Put -Headers $authHeaders -ContentType "application/json" -Body $updateBody
  Write-Host "   ✓ NVR Updated" -ForegroundColor Green
  Write-Host "   New RTSP Channels: $($updateResp.data.rtsp_urls.Count)" -ForegroundColor Gray
  Write-Host "   Response:" -ForegroundColor Gray
  Write-Host ($updateResp | ConvertTo-Json -Depth 10) -ForegroundColor Gray
}
catch {
  Write-Host "   ✗ Update NVR FAILED: $($_.Exception.Message)" -ForegroundColor Red
  if ($_.ErrorDetails.Message) {
    Write-Host "   Details: $($_.ErrorDetails.Message)" -ForegroundColor Red
  }
}

Write-Host ""

# 6. 再次查询验证更新
Write-Host "6. Verify Update" -ForegroundColor Cyan
try {
  $verifyResp = Invoke-RestMethod -Uri "$BaseUrl/api/nvr?id=$testNvrId" -Method Get -Headers $authHeaders
  Write-Host "   ✓ Verification SUCCESS" -ForegroundColor Green
  Write-Host "   RTSP Channels After Update:" -ForegroundColor Gray
  foreach ($rtsp in $verifyResp.data.rtsp_urls) {
    Write-Host "      Channel $($rtsp.channel): $($rtsp.url)" -ForegroundColor Gray
  }
}
catch {
  Write-Host "   ✗ Verification FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# 7. 查询所有 NVR 列表
Write-Host "7. List All NVRs" -ForegroundColor Cyan
try {
  $listResp = Invoke-RestMethod -Uri "$BaseUrl/api/nvr" -Method Get -Headers $authHeaders
  Write-Host "   ✓ NVR List Retrieved" -ForegroundColor Green
  Write-Host "   Total NVRs: $($listResp.data.Count)" -ForegroundColor Gray
}
catch {
  Write-Host "   ✗ List NVRs FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""

# 8. 清理测试数据
Write-Host "8. Cleanup Test Data" -ForegroundColor Yellow
$cleanup = Read-Host "Do you want to delete test data? (y/n)"
if ($cleanup -eq "y") {
  # 删除 NVR
  try {
    $delNvrBody = @{id = $testNvrId } | ConvertTo-Json
    Invoke-RestMethod -Uri "$BaseUrl/api/nvr" -Method Delete -Headers $authHeaders -ContentType "application/json" -Body $delNvrBody
    Write-Host "   ✓ NVR Deleted: ID=$testNvrId" -ForegroundColor Green
  }
  catch {
    Write-Host "   ✗ Delete NVR FAILED: $($_.Exception.Message)" -ForegroundColor Red
  }

  # 删除 Building
  if ($buildResp) {
    try {
      $delBuildBody = @{id = $testBuildId } | ConvertTo-Json
      Invoke-RestMethod -Uri "$BaseUrl/api/building" -Method Delete -Headers $authHeaders -ContentType "application/json" -Body $delBuildBody
      Write-Host "   ✓ Building Deleted: ID=$testBuildId" -ForegroundColor Green
    }
    catch {
      Write-Host "   ✗ Delete Building FAILED: $($_.Exception.Message)" -ForegroundColor Red
    }
  }
}
else {
  Write-Host "   Test data preserved:" -ForegroundColor Cyan
  Write-Host "   NVR ID: $testNvrId" -ForegroundColor Gray
  Write-Host "   Building ID: $testBuildId" -ForegroundColor Gray
}

Write-Host ""
Write-Host "=======================================" -ForegroundColor Cyan
Write-Host "Test Complete!" -ForegroundColor Cyan
Write-Host "=======================================" -ForegroundColor Cyan
