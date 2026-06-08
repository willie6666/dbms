# 專案開發硬碟空間清理指南 (Docker & WSL)

在進行 Docker Compose 專案開發時（如本專案反覆執行 `docker compose up --build`），經常會因為 Docker 的歷史映像檔、快取，以及 Windows WSL2 虛擬硬碟的特性，導致 C 槽空間急速減少幾十 GB。

若遇到空間不足的情況，請依照以下三個階段進行深度清理。

---

## 階段一：清理 Docker 無用快取與映像檔

當 Docker 容器處於「啟動中」狀態時，系統會鎖定快取與映像檔不讓您刪除。因此在清理前，必須先關閉專案。

1. **關閉正在運行的容器**
   在專案根目錄開啟終端機，停止並移除現有的容器：
   ```bash
   docker compose down
   ```

2. **執行終極清理指令**
   清理所有未被使用的容器、網路、懸空映像檔以及快取（這會釋放最多的內部空間）：
   ```bash
   docker system prune -a --volumes
   ```
   > 系統會詢問 `Are you sure you want to continue? [y/N]`，請輸入 `y` 並按下 Enter。

---

## 階段二：壓縮 WSL 2 虛擬硬碟 (Windows 專屬)

即使 Docker 刪除了幾十 GB 的檔案，Windows 底層的虛擬硬碟檔 (`.vhdx`) 預設**不會自動縮小**。我們必須手動壓縮它，才能把空間真正還給 C 槽。

### 1. 尋找肥大的 VHDX 檔案
開啟 **PowerShell** 並執行以下指令，找出 Docker 虛擬硬碟的精確位置：
```powershell
Get-ChildItem -Path C:\Users\HP\AppData\Local -Filter "*.vhdx" -Recurse -ErrorAction SilentlyContinue | Select-Object FullName, @{Name="Size(GB)";Expression={[math]::Round($_.Length / 1GB, 2)}}
```
*通常檔案會是 `C:\Users\HP\AppData\Local\Docker\wsl\disk\docker_data.vhdx`，且大小動輒超過 10GB 以上。*

### 2. 強制關閉 Docker 與 WSL
- **關閉 Docker Desktop**：在右下角系統匣對鯨魚圖示按右鍵，選擇 **Quit Docker Desktop**。
- **關閉 WSL**：打開新的 cmd 或 PowerShell 輸入：
  ```bash
  wsl --shutdown
  ```

### 3. 使用 DiskPart 進行壓縮
打開一個**以系統管理員身分執行**的命令提示字元 (cmd)，進入 `diskpart` 模式：
```cmd
diskpart
```
出現 `DISKPART>` 後，依序貼上以下指令（請把 `file="..."` 的路徑換成您在步驟 1 找到的路徑）：
```cmd
select vdisk file="C:\Users\HP\AppData\Local\Docker\wsl\disk\docker_data.vhdx"
attach vdisk readonly
compact vdisk
detach vdisk
exit
```
> **注意**：`compact vdisk` 步驟可能需要等待幾分鐘，請耐心等候它完成。

---

## 階段三：清理 Go 語言編譯快取 (選擇性)

Go 語言在編譯過程中也會產生快取，雖然佔用空間遠不及 Docker，但若空間非常緊繃也可順手清理。

1. **關閉 VS Code 或其他正在使用 Go 的 IDE**（避免檔案被鎖定而出現 `Access is denied`）。
2. 在終端機執行：
   ```bash
   go clean -cache -modcache -i -r
   ```
   *或您也可以手動刪除以下資料夾：*
   - `C:\Users\HP\AppData\Local\go-build`
   - `C:\Users\HP\go\pkg`
