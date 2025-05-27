param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("start", "stop")]
    [string]$Action
)

$services = @(
    @{
        Name = "auth-service"
        Path = "backend/auth-service/cmd"
        Port = 8080
    },
    @{
        Name = "profile-service"
        Path = "backend/profile-service/cmd"
        Port = 8081
    },
    @{
        Name = "orchestrator-service"
        Path = "backend/orchestrator-service/cmd"
        Port = 8082
    }
)

function Start-Services {
    Write-Host "Запуск микросервисов..." -ForegroundColor Green
    
    foreach ($service in $services) {
        $serviceName = $service.Name
        $servicePath = $service.Path
        
        # Создаем новое окно PowerShell для каждого сервиса
        Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PSScriptRoot/../$servicePath'; go run main.go" -WindowStyle Normal
        
        Write-Host "Сервис $serviceName запущен на порту $($service.Port)" -ForegroundColor Cyan
    }
    
    Write-Host "Все сервисы успешно запущены!" -ForegroundColor Green
}

function Stop-Services {
    Write-Host "Остановка микросервисов..." -ForegroundColor Yellow
    
    foreach ($service in $services) {
        $port = $service.Port
        $serviceName = $service.Name
        
        # Находим процесс по порту и завершаем его
        $processInfo = netstat -ano | findstr ":$port" | findstr "LISTENING"
        
        if ($processInfo) {
            $processId = ($processInfo -split ' ')[-1]
            Stop-Process -Id $processId -Force
            Write-Host "Сервис $serviceName остановлен (PID: $processId)" -ForegroundColor Cyan
        } else {
            Write-Host "Сервис $serviceName не запущен или не найден на порту $port" -ForegroundColor Yellow
        }
    }
    
    Write-Host "Все сервисы остановлены!" -ForegroundColor Green
}

# Запуск функции в зависимости от параметра
switch ($Action) {
    "start" { Start-Services }
    "stop" { Stop-Services }
} 