# 鐭挎硥姘碅PI鍙戝竷涓庡洖婊氭墜鍐?
鏈墜鍐岀敤浜庢妸涓€娆℃敼鍔ㄤ粠鈥滄湰鍦拌兘璺戔€濇帹杩涘埌鈥滃彲鏀惰垂璇曡繍钀モ€濄€傛牳蹇冨師鍒欙細鍥哄畾鐗堟湰銆佸厛澶囦唤銆佸彲鍥炴粴銆佸啀鍏竷銆?
## 褰撳墠鍊欓€夌増鏈姸鎬侊紙2026-05-25锛?
- 鏈満楠岃瘉闀滃儚锛歚sub2api:kqs-api-prelaunch-20260525-2058`銆?- 鏈満瀹瑰櫒 `sub2api-gg` 宸茶繍琛岃鍥哄畾 tag锛宍/health` 杩斿洖 `200`銆?- Cloudflare Worker `sub2api-proxy` 宸查儴缃茬増鏈?`768ca6f6-93c9-4d6c-8e8a-1ab4cf9fadf6`銆?- 褰撳墠鍊欓€夊凡鍚敤 `SERVER_TRUSTED_PLATFORM=cloudflare`锛屽叕缃戝叆鍙ｇ殑鍙俊瀹㈡埛绔?IP 閾捐矾宸茬撼鍏ュ悗绔?`ClientIP()`銆?- 褰撳墠鍊欓€夊凡鍚敤 `SECURITY_URL_ALLOWLIST_ENABLED=true`锛屽苟鏄惧紡鏀捐 OpenAI/Codex 蹇呴渶涓婃父銆?- 褰撳墠浠ｇ爜宸ヤ綔鍖轰粛鏈夋湭鎻愪氦鏀瑰姩锛屽洜姝よ闀滃儚灞炰簬鍙戝竷鍊欓€夛紝涓嶇瓑鍚屼簬宸插喕缁?tag銆?- 姝ｅ紡鏍囪 `kqs-api-v0.2.0-beta` 鍓嶏紝蹇呴』鍏堟妸鍙戝竷鑼冨洿鏁寸悊鎴愬共鍑€鎻愪氦锛屽苟璁板綍瀵瑰簲 git SHA銆?- 鏈鍊欓€夐獙璇佽鎯呰 `docs/RELEASE_CANDIDATE_20260525_CN.md`銆?
## 1. 鍙戝竷鍓嶅喕缁?
1. 纭褰撳墠鐩爣鐗堟湰锛屼緥濡?`kqs-api-v0.2.0-beta`銆?2. 鏌ョ湅宸ヤ綔鍖猴細

```powershell
git status --short
git diff --name-only
```

3. 鍙繚鐣欐湰娆″彂甯冮渶瑕佺殑鏀瑰姩銆傛棤鍏虫敼鍔ㄨ涔堝崟鐙彁浜わ紝瑕佷箞鏄庣‘璁板綍涓轰笉鍙戝竷銆?4. 璁板綍褰撳墠绾夸笂鐗堟湰锛?
```powershell
docker inspect sub2api-gg --format '{{.Config.Image}}'
docker inspect sub2api-gg --format '{{.Image}}'
```

鎶婄粨鏋滃啓鍏ュ彂甯冭褰曪紝浣滀负鍥炴粴鐩爣銆?
## 2. 鍙戝竷鍓嶅浠?
鍦ㄥ悗鍙拌繘鍏?`绠＄悊鍚庡彴 -> 鏁版嵁澶囦唤`锛?
1. 纭 S3/R2 杩炴帴娴嬭瘯閫氳繃銆?2. 纭璁″垝澶囦唤寮€鍚紝寤鸿鍖椾含鏃堕棿姣忓ぉ 02:30锛屼繚鐣?14 澶╂垨 30 浠姐€?3. 鐐瑰嚮鎵嬪姩澶囦唤锛岀瓑寰呯姸鎬佸彉鎴?`completed`銆?4. 璁板綍澶囦唤 ID銆佹枃浠跺悕銆佸畬鎴愭椂闂村拰澶у皬銆?
绂佹鐢ㄧ敓浜у簱鍋氭仮澶嶆紨缁冦€傛仮澶嶆紨缁冨繀椤诲湪涓存椂鐜鎵ц銆?
## 3. 鏈湴楠岃瘉

鍚庣锛?
```powershell
Set-Location D:\codex-work\gg\sub2api
$backend = (Resolve-Path .\backend).Path
docker run --rm -v "${backend}:/app/backend" -w /app/backend `
  -e GOPROXY=https://goproxy.cn,direct `
  -e GOSUMDB=sum.golang.google.cn `
  golang:1.26.3-alpine sh -lc `
  'export PATH=/usr/local/go/bin:$PATH; apk add --no-cache git ca-certificates tzdata >/dev/null; go test ./...'
```

璇存槑锛氬綋鍓?Windows 瀹夸富鏈烘病鏈夊畨瑁?Go toolchain锛屽彂甯冮獙璇佺粺涓€浣跨敤鍥哄畾 Go 瀹瑰櫒鎵ц鍚屼竴浠芥簮鐮佺殑鍏ㄩ噺娴嬭瘯銆?
鍓嶇锛?
```powershell
Set-Location D:\codex-work\gg\sub2api\frontend
npm run typecheck
npm run build
```

鍙戝竷闂搁棬锛?
```powershell
Set-Location D:\codex-work\gg\sub2api
.\scripts\prelaunch-readiness.ps1 -ApiBaseUrl https://api.cauai.fun
```

鑴氭湰浼氫紭鍏堜娇鐢?`-AdminToken`锛屽惁鍒欎細灏濊瘯璇诲彇鏈満 `deploy/.env` 鐨勭鐞嗗憳璐﹀彿锛屽湪鏈満鐧诲綍鍚庡彧鎶婁复鏃?token 鐢ㄤ簬澶囦唤鍋ュ悍妫€鏌ワ紝涓嶄細鎵撳嵃瀵嗙爜鎴?token銆?
濡傛灉鍦?CI 鎴栬繙绋嬫満鍣ㄤ笂娌℃湁 `deploy/.env`锛屾樉寮忎紶鍏ョ鐞嗗憳 token锛?
```powershell
.\scripts\prelaunch-readiness.ps1 -ApiBaseUrl https://api.cauai.fun -AdminToken $env:SUB2API_ADMIN_TOKEN
```

濡傛灉鍏綉鍩熷悕鍙敤浜?`/health`锛屼絾绠＄悊鎺ュ彛甯屾湜璧版湰鏈猴紝鍙樉寮忔寚瀹氾細

```powershell
.\scripts\prelaunch-readiness.ps1 -ApiBaseUrl https://api.cauai.fun -AdminApiBaseUrl http://localhost:18080
```

澶囦唤闂搁棬蹇呴』鍚屾椂閫氳繃锛歋3/R2 閰嶇疆瀛樺湪銆佽繛鎺ユ祴璇曟垚鍔熴€佽鍒掑浠藉紑鍚€佹渶杩戜竴娆?completed 澶囦唤灏忎簬 30 灏忔椂銆?
## 4. 鏋勫缓鍥哄畾闀滃儚

涓嶈鎶婃寮忓彂甯冪粦瀹氬埌 `latest`銆傚彂甯冮暅鍍忓繀椤诲寘鍚増鏈拰 git 鐭?SHA銆?
PowerShell 绀轰緥锛?
```powershell
$sha = git rev-parse --short HEAD
$image = "sub2api:kqs-api-v0.2.0-beta-$sha"
bash .\deploy\build_image.sh $image
```

濡傛灉 Docker Hub 鎷夊彇鍩虹闀滃儚澶辫触锛屽厛瑙ｅ喅鍩虹闀滃儚鎷夊彇鎴栭厤缃暅鍍忔簮锛涗笉瑕佹妸涓存椂 `docker cp` 浜岃繘鍒惰鐩栧綋浣滄寮忓彂甯冦€?
## 5. 涓婄嚎

1. 淇敼 `deploy/.env`锛?
```env
SUB2API_IMAGE=sub2api:kqs-api-v0.2.0-beta-<sha>
```

2. 鍚姩鍥哄畾鐗堟湰锛?
```powershell
docker compose --env-file .\deploy\.env -f .\deploy\docker-compose.yml up -d
```

3. 楠岃瘉瀹瑰櫒锛?
```powershell
docker ps --filter "name=sub2api-gg"
docker logs --tail 120 sub2api-gg
Invoke-RestMethod https://api.cauai.fun/health
```

4. 娴忚鍣ㄧ儫娴嬶細

- 鏃犵櫥褰曠姸鎬佽闂?`/register`锛岃兘鐪嬪埌閭€璇风爜杈撳叆鍜屽洓浠藉崗璁叆鍙ｃ€?- 鏃犵櫥褰曠姸鎬佽闂?`/legal/terms`銆乣/legal/privacy`銆乣/legal/refund`銆乣/legal/usage-policy`锛屾鏂囧潎鍙鍙栥€?- `/dashboard` 鑳界湅鍒版柊鎵嬪紩瀵煎拰鏀寔鍏ュ彛銆?- `/redeem` 鑳界湅鍒板崱瀵嗚鏄庡拰鏀寔鍏ュ彛銆?- `/keys` 鏂板缓 Key 榛樿鏈夐搴﹀拰閫熺巼淇濇姢锛屼笖鍙鍒?Codex 閰嶇疆銆?
## 6. Cloudflare 缂撳瓨

鍓嶇闈欐€佽祫婧愭垨 logo 鏇存柊鍚庯紝鍦?Cloudflare 鎵ц缂撳瓨鍒锋柊锛?
1. 杩涘叆 `cauai.fun`銆?2. 鎵撳紑 `Caching -> Configuration -> Purge Cache`銆?3. 浼樺厛 Purge URL锛歚https://api.cauai.fun/`銆乣https://api.cauai.fun/assets/*`銆?4. 鑻ラ〉闈粛鏃э紝鍐?Purge Everything銆?
鍒锋柊鍚庣敤鏃犵棔绐楀彛璁块棶 `https://api.cauai.fun/home` 鍜?`/login`銆?
## 7. 鍥炴粴

鍥炴粴鍙垏闀滃儚锛屼笉鏀规暟鎹簱锛岄櫎闈炴湰娆″彂甯冨寘鍚笉鍙吋瀹硅縼绉汇€?
```powershell
# 鎶?deploy/.env 鐨?SUB2API_IMAGE 鏀瑰洖鍙戝竷鍓嶈褰曠殑闀滃儚
docker compose --env-file .\deploy\.env -f .\deploy\docker-compose.yml up -d
Invoke-RestMethod https://api.cauai.fun/health
```

濡傛灉鏁版嵁搴撹縼绉婚€犳垚涓嶅彲閫嗛棶棰橈紝浣跨敤鍙戝竷鍓嶆墜鍔ㄥ浠藉湪涓存椂鐜鍏堥獙璇侊紝鍐嶅喅瀹氭槸鍚︽仮澶嶇敓浜с€備笉瑕佸湪娌℃湁楠岃瘉鐨勬儏鍐典笅鐩存帴瑕嗙洊鐢熶骇搴撱€?
## 8. 姝ｅ紡鍏竷鏉′欢

鍙湁鍏ㄩ儴婊¤冻鎵嶅彂甯冪粰瀹㈡埛锛?
- 宸ヤ綔鍖哄凡鍐荤粨鎴愬彲杩芥函鎻愪氦锛屽彂甯?tag 涓庨暅鍍?tag 鍧囪褰?git SHA銆?- 鍥哄畾闀滃儚 tag 宸蹭笂绾匡紝绾夸笂涓嶆槸鏈褰曠殑 `latest` 婕傜Щ鐗堟湰銆?- 鍙戝竷鍓嶆墜鍔ㄥ浠藉畬鎴愶紝骞朵笖鏈€杩戣嚜鍔ㄥ浠藉皬浜?30 灏忔椂銆?- 鐧诲綍鏉℃銆侀€€娆捐鍒欍€侀殣绉佽鏄庛€佸彲鎺ュ彈浣跨敤鏀跨瓥鍙闂€?- 鏂扮敤鎴锋祦绋嬪畬鎴愶細娉ㄥ唽銆佸厬鎹€佸垱寤?Key銆佸鍒堕厤缃€乣/responses` 鎴愬姛銆?- 閿欒鎻愮ず瑕嗙洊锛氶噸澶嶅崱瀵嗐€侀敊璇崱瀵嗐€佷綑棰濅笉瓒炽€並ey 鏈垎缁勩€佹棤鏁?Key銆?- 宸茶褰曞洖婊氶暅鍍忓拰鍥炴粴鍛戒护銆?
