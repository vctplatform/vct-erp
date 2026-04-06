param(
    [string]$action,
    [string]$arg1,
    [string]$arg2,
    [string]$arg3
)

$TaskFile = "d:\VCT PLATFORM\.tasks.json"
$DockerTestFile = "d:\VCT PLATFORM\docker-compose.test.yml"
$KnowledgeDir = "d:\VCT PLATFORM\.knowledge"

if (-not (Test-Path $TaskFile)) {
    Set-Content -Path $TaskFile -Value '[]'
}

$tasks = Get-Content $TaskFile | ConvertFrom-Json
if ($null -eq $tasks) {
    $tasks = @()
}
if ($tasks.GetType().Name -ne 'Object[]') {
    $tasks = @($tasks)
}

if ($action -eq "create") {
    $newTask = [PSCustomObject]@{
        id = ($tasks.Count + 1)
        name = $arg1
        assignee = $arg2
        status = "todo"
    }
    $tasks += $newTask
    $tasks | ConvertTo-Json -Depth 5 | Set-Content $TaskFile
    Write-Host "Task created: $($newTask.id) - $($newTask.name) assigned to $($newTask.assignee)"
} elseif ($action -eq "list") {
    if ($tasks.Count -gt 0) {
        $tasks | Format-Table -Property id, name, assignee, status
    } else {
        Write-Host "No tasks currently in the system."
    }
} elseif ($action -eq "ask") {
    if (-not $arg1) {
        Write-Host "Usage: vct.cmd ask `'<search term>`'"
        exit 1
    }
    Write-Host "Searching Knowledge Base for '$arg1'..."
    $results = Get-ChildItem -Path $KnowledgeDir -Recurse -Filter *.md | Select-String -Pattern $arg1 -Context 0
    if ($results) {
        foreach ($match in $results) {
            Write-Host "--- Found in $($match.Filename) ---"
            Write-Host $match.Line
        }
    } else {
        Write-Host "No context found. Proceed with max level constraints."
    }
} elseif ($action -eq "complete") {
    $id = $arg1
    $team = $arg2
    
    if (-not $id) {
        Write-Host "Usage: vct.cmd complete <id> [team]"
        exit 1
    }

    if ($team) {
        $container = "test-vct-$team"
        Write-Host ">> System-Enforced Docker Validation for $container..."
        $process = Start-Process docker-compose -ArgumentList "-f", "`"$DockerTestFile`"", "run", "--rm", $container -PassThru -Wait -NoNewWindow
        
        if ($process.ExitCode -ne 0) {
            $failFile = "d:\VCT PLATFORM\.vct_fail_count_$id"
            $failCount = 1
            if (Test-Path $failFile) {
                $failCount = [int](Get-Content $failFile) + 1
            }
            if ($failCount -ge 3) {
                Write-Host "`n[FATAL SYSTEM LOCK] Agent has failed 3 times in a row." -ForegroundColor Red
                Write-Host ">> EXECUTING (V12) SELF-HEALING ROLLBACK: git reset --hard HEAD" -ForegroundColor Red
                git reset --hard HEAD
                Remove-Item -Path $failFile -ErrorAction SilentlyContinue
                
                foreach ($task in $tasks) {
                    if ($task.id -match "^$id$") {
                        $task.status = "failed"
                    }
                }
                $tasks | ConvertTo-Json -Depth 5 | Set-Content $TaskFile
                Write-Host "Task $id explicitly marked as FAILED. LLM locked out." -ForegroundColor Red
                exit 1
            } else {
                Set-Content -Path $failFile -Value $failCount
                Write-Host "`n[ERROR][FATAL] Docker Validation Failed (Attempt $failCount/3)." -ForegroundColor Red
                Write-Host "Task completion REJECTED. Fix your code and try again." -ForegroundColor Yellow
                exit 1
            }
        }
        $failFile = "d:\VCT PLATFORM\.vct_fail_count_$id"
        if (Test-Path $failFile) { Remove-Item -Path $failFile -ErrorAction SilentlyContinue }
        Write-Host "`n[SUCCESS] Tests passed for $team. (V12 Enforcer Authorized)" -ForegroundColor Green
    }

    $found = $false
    foreach ($task in $tasks) {
        if ($task.id -match "^$id$") {
            $task.status = "done"
            $found = $true
        }
    }
    if ($found) {
        $tasks | ConvertTo-Json -Depth 5 | Set-Content $TaskFile
        Write-Host "Task $id marked as done."
    } else {
        Write-Host "Task $id not found."
    }
} elseif ($action -eq "sync") {
    Write-Host ">> Synchronizing VCT Platform Repositories..." -ForegroundColor Cyan
    $repos = Get-ChildItem -Path "d:\VCT PLATFORM\vct-platform" -Directory
    foreach ($repo in $repos) {
        if (Test-Path "$($repo.FullName)\.git") {
            Write-Host "Syncing $($repo.Name)..." -ForegroundColor Yellow
            git -C $repo.FullName pull
            git -C $repo.FullName push
        }
    }
    Write-Host ">> All repos synchronized." -ForegroundColor Green
} elseif ($action -eq "status") {
    Write-Host ">> VCT Platform Cluster Status:" -ForegroundColor Cyan
    $repos = Get-ChildItem -Path "d:\VCT PLATFORM\vct-platform" -Directory
    foreach ($repo in $repos) {
        $gitStatus = "No Git"
        if (Test-Path "$($repo.FullName)\.git") {
            $gitStatus = "Git Ready"
        }
        Write-Host "[$gitStatus] $($repo.Name)"
    }
} else {
    Write-Host "Usage: vct.cmd [create <name> <assignee> | list | complete <id> <team> | ask <query> | sync | status]"
}
