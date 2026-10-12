# Elmo Archive — detailed guide

[Quick start](../README.md) · [中文说明](README.zh-CN.md)

### Updating and existing accounts

Select **Sync records** after making more pulls; reconnect if the saved connection has expired. Use the same UID and identification method to keep history together.

The old generated-local-ID import option is no longer offered to new imports. Existing saved accounts remain readable, and saved legacy connections remain supported. Switching a legacy account to a real UID may create another account entry; automatic merging is not implemented.

## Capture, privacy, and troubleshooting

Capture temporarily uses a local HTTPS proxy and a certificate trusted for the current Windows user. Before installation, compare the name and thumbprint shown in the app with the Windows confirmation. Windows cannot independently validate this local CA; the warning remains visible. A matching thumbprint identifies the certificate, not the publisher. The app decrypts only explicitly listed game record hosts, plus the Darkwinter webshop API in webshop mode. Other HTTPS traffic is tunneled without decryption. Upstream TLS certificates are verified. These host limits are enforced by the app; Windows root trust itself is broader. The app shows the certificate folder during setup. `ca.crt` is the public certificate; `ca.key` is its private signing key. Keep the key private and do not attach it to bug reports.

**Remove after a successful import** is on by default. You can also open **Import records → Capture & certificate → Remove certificate trust** at any time. Removal stops capture, restores the previous proxy settings, and removes only the matching certificate from the current user’s root store. It keeps your records and local certificate files. A valid saved connection can sync without this certificate; a new capture requires installation again. If import fails or is cancelled, use the manual removal button. Other app copies or certificate stores may require separate removal through certmgr.msc. New certificates use **Elmo Archive Local Capture CA**; valid older **ElmoBeacon Custom CA** files retain their name and thumbprint.

Records are stored in `ElmoBeacon.db` in the app’s working folder. Connection files can contain session credentials. **Do not upload databases, game logs, capture URLs, tokens, `gacha_session_info.json`, `gacha_auth_details.json`, or `ca.key` in bug reports.** Crop or redact UIDs from screenshots.

If a crash leaves Windows using `127.0.0.1:8899` as its proxy, restore your previous proxy settings in Windows Settings → Network & internet → Proxy. Do not change a work or school proxy configuration without knowing its previous settings.

Enable **Show detailed import progress** to see capture stages, page counts, and record counts. This progress panel excludes tokens and request URLs; raw logs and technical errors still need review before sharing.

Automatic binary replacement from the original project’s update server is disabled. Download updates from this repository’s Releases page.

## Preview limitations

- Live capture still needs testing across Darkwinter and Haoplay regions. Compatibility paths are implemented; this is not a claim that every region has been tested.
- The game may no longer supply older records. The app cannot recover history absent from the server and your local database.
- Historical target choices are not supplied by the record API. Use **Edit target result** to correct labels you know are wrong; leave uncertain outcomes Unknown. A known loss followed by a guaranteed copy is counted as one recorded cost, even if you changed targets.
- The chart percentiles use interpolation between published thresholds, not an exact probability simulation or player leaderboard.
- English and Simplified/Traditional Chinese cover the current UI. Older Japanese/Korean translations may fall back to English. Item names depend on installed game data.

## Develop and build

Install Go **1.25 or later**, Node.js **22**, and the Windows requirements for [Wails v2](https://wails.io/docs/gettingstarted/installation/).

```powershell
npm --prefix frontend ci
go install github.com/wailsapp/wails/v2/cmd/wails@v2.10.1
node --test frontend/tests/*.test.mjs
go test ./...
wails dev
# Production build:
wails build -o ElmoArchive.exe
```

The executable is written to `build/bin/`. Run `node scripts/release-preflight.mjs` before publishing source. See [launch preparation](LAUNCH.md) for release packaging, repository discovery, and the remaining live checks.


