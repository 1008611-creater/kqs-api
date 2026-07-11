# Cloudflare R2 澶囦唤閰嶇疆鎵嬪唽

褰撳墠鍙戝竷闂搁棬瑕佹眰澶囦唤鍙敤锛歋3/R2 閰嶇疆瀛樺湪銆佽繛鎺ユ祴璇曟垚鍔熴€佽鍒掑浠藉紑鍚€佹渶杩戜竴娆℃垚鍔熷浠藉皬浜?30 灏忔椂銆?
## 褰撳墠鐘舵€侊紙2026-05-25锛?
- Cloudflare R2 宸插惎鐢ㄣ€?- 绉佹湁 bucket 宸查厤缃负 `sub2api-backups`銆?- 鐭挎硥姘碅PI 鍚庡彴 S3/R2 杩炴帴宸叉祴璇曟垚鍔熴€?- 瀹氭椂澶囦唤宸插紑鍚細`30 2 * * *`锛屽嵆鍖椾含鏃堕棿姣忓ぉ 02:30銆?- 宸插畬鎴愪竴娆℃墜鍔ㄥ浠斤紝澶囦唤璁板綍 ID 涓?`0cf00f51`銆?
## 鏁版嵁杈圭晫涓庨殣绉?
褰撳墠鍚庡彴澶囦唤鍔熻兘鎵ц鐨勬槸 PostgreSQL 鍏ㄩ噺瀵煎嚭锛屽苟灏嗗帇缂╁悗鐨勬暟鎹簱澶囦唤涓婁紶鑷?R2銆傚畠涓嶄細涓诲姩璇诲彇鎴栦笂浼犵敤鎴风數鑴戜腑鐨勭湡瀹炴枃浠躲€佸浘鐗囨垨鏂囨。鐩綍銆?
鏁版嵁搴撳浠戒粛鍙兘鍖呭惈璐﹀彿鏍囪瘑銆佷綑棰濄€佸崱瀵嗗厬鎹㈣褰曘€丄PI Key 鍏冩暟鎹€佺珯鐐归厤缃互鍙婄敤閲?閿欒鏃ュ織锛屽洜姝や粛灞炰簬鏁忔劅杩愯惀鏁版嵁锛?
- bucket 蹇呴』淇濇寔绉佹湁锛屼笉閰嶇疆鍏紑璁块棶鍩熷悕銆?- R2 token 浠呮巿浜?`sub2api-backups` 鐨?`Object Read & Write` 鏉冮檺銆?- 涓嶆妸 Access Key銆丼ecret銆佸浠戒笅杞介摼鎺ュ啓鍏ヨ亰澶┿€佷粨搴撴垨鍏憡銆?- 鎭㈠婕旂粌鍙湪闅旂涓存椂鐜杩涜锛屼笉鎶婂浠戒氦缁欐櫘閫氱敤鎴枫€?
## 1. 鍚敤 R2

鎵撳紑 Cloudflare Dashboard锛?
```text
https://dash.cloudflare.com/?to=/:account/r2/overview
```

杩涘叆 `R2 object storage`锛屾寜椤甸潰鎻愮ず鍚敤 R2銆侰loudflare 瀹樻柟璇存槑锛氬繀椤诲厛璐拱/鍚敤 R2锛屼箣鍚庢墠鑳界敓鎴?R2 API token銆?
## 2. 鍒涘缓 Bucket

鍚敤鍚庡彲浠ュ湪 Dashboard 閲屽垱寤?bucket锛屼篃鍙互鐢?Wrangler锛?
```powershell
npx wrangler r2 bucket create sub2api-backups
npx wrangler r2 bucket list
```

Bucket 鍚嶇О寤鸿鍥哄畾涓猴細

```text
sub2api-backups
```

涓嶈鍏紑杩欎釜 bucket銆傚畠鍙敤浜庢暟鎹簱澶囦唤鏂囦欢銆?
## 3. 鍒涘缓 R2 API Token

鍦?`R2 object storage` 椤甸潰鎵惧埌 `Account Details`锛岀偣鍑?`API Tokens` 鏃佽竟鐨?`Manage`锛屽垱寤?token銆?
鎺ㄨ崘閰嶇疆锛?
- Token 绫诲瀷锛歚Create Account API token`锛屽鏋滈〉闈㈡潈闄愪笉瓒筹紝鍐嶇敤 `Create User API token`銆?- 鏉冮檺锛歚Object Read & Write`銆?- 鑼冨洿锛氬彧缁戝畾 `sub2api-backups` 杩欎釜 bucket銆?
鍒涘缓瀹屾垚鍚庯紝鍙鍒朵竴娆★細

- `Access Key ID`
- `Secret Access Key`

Secret 鍙湪 Cloudflare 椤甸潰鏄剧ず涓€娆★紝涓嶈鍙戠粰鐢ㄦ埛銆佷笉瑕佸啓杩涗粨搴撱€佷笉瑕佽创鍒拌亰澶╅噷銆?
## 4. 濉叆鐭挎硥姘碅PI鍚庡彴

杩涘叆绠＄悊鍛樺悗鍙帮細`璁剧疆 -> 鏁版嵁澶囦唤`锛屽～鍐欙細

```text
Endpoint: https://3b78db8f07bf8738064ca1606fe89cb5.r2.cloudflarestorage.com
Region: auto
Bucket: sub2api-backups
Prefix: backups/
Access Key ID: Cloudflare 鐢熸垚鐨?Access Key ID
Secret Access Key: Cloudflare 鐢熸垚鐨?Secret Access Key
Force Path Style: 鍏抽棴
```

淇濆瓨鍓嶅厛鐐?`娴嬭瘯杩炴帴`銆傛祴璇曢€氳繃鍚庡啀淇濆瓨銆?
## 5. 寮€鍚鍒掑浠?
鎺ㄨ崘閰嶇疆锛?
```text
鍚敤: 鏄?Cron: 30 2 * * *
淇濈暀澶╂暟: 14
淇濈暀浠芥暟: 30
```

璇存槑锛歚30 2 * * *` 琛ㄧず姣忓ぉ 02:30銆傚鍣ㄦ椂鍖哄綋鍓嶉厤缃负 `Asia/Shanghai`銆?
## 6. 鍙戝竷鍓嶉獙鏀?
鎵嬪姩鍒涘缓涓€娆″浠斤紝绛夊緟鐘舵€佸彉鎴?`completed`銆傜劧鍚庤繍琛岋細

```powershell
Set-Location D:\codex-work\gg\sub2api
.\scripts\prelaunch-readiness.ps1 -ApiBaseUrl https://api.cauai.fun
```

蹇呴』鐪嬪埌杩欎簺椤圭洰閫氳繃锛?
- `backup s3 config`
- `backup s3 connection`
- `backup schedule`
- `backup freshness`
- `backup failures`

濡傛灉 `backup freshness` 澶辫触锛屼笉瑕佹寮忓叕寮€鍙戝竷銆?
## 鍙傝€?
- Cloudflare R2 API token 鏂囨。锛歚https://developers.cloudflare.com/r2/api/tokens/`
- Cloudflare R2 鍒涘缓 bucket 鏂囨。锛歚https://developers.cloudflare.com/r2/buckets/create-buckets/`
- Cloudflare R2 S3 API 鍏煎鎬э細`https://developers.cloudflare.com/r2/api/s3/api/`
- Cloudflare R2 Go SDK 绀轰緥锛歚https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/`

