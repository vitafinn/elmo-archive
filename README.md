# Elmo Archive — GFL2 pull tracker

A Windows pull tracker for **Girls’ Frontline 2: Exilium**. View pity, doll and weapon copy costs, and outfit records.

**A fork of [ElmoBeacon](https://github.com/MatchaCabin/ElmoBeacon), originally maintained by MccWiki and contributors.** This fork is independently maintained. It is not an official game tool.

[Download](https://github.com/vitafinn/elmo-archive/releases) · [Report a bug](https://github.com/vitafinn/elmo-archive/issues) · [中文说明](docs/README.zh-CN.md)

## Install

Requires **Windows x64**, the **Windows game client**, and [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/). No separate proxy installation is needed.

1. Download the **windows-x64.zip** from Releases. Do not choose “Source code.” If Releases is empty, no public build is available yet.
2. Extract the ZIP to a folder you can write to. Run `ElmoArchive.exe`.

This is an early preview. The executable is unsigned. Darkwinter and Haoplay (Global, Asia, Japan, Korea) import paths are included; regional testing is not complete. China-client imports are not supported.

## Import records

1. Open the game and sign in.
2. In Elmo Archive, select **Import records → My game UID**. Enter the UID from your game profile.
3. Select **Connect to game history**. If certificate setup appears, compare the displayed name and thumbprint with the Windows prompt. Select **Yes** only if they match and you trust this app. If the prompt is hidden, select **Show Windows prompt**.
4. In the game, open a running banner → **Details → Access Record** side tab. If Access Record is already open, close and reopen it. No pull or purchase is needed.
5. When the app says **History connected**, select **Sync records**. Check the UID and region, then confirm. Wait for **Sync complete**.

Darkwinter users can also choose **Darkwinter webshop sign-in**. Follow the app’s steps and wait for **Webshop UID captured** before opening Access Record. Use the same account in the webshop and game; the app cannot verify that they match.

**Next time:** select **Sync records**. Reconnect only if the connection has expired. Keep using the same UID and import method.

## Certificate and saved data

Capture temporarily uses a local HTTPS proxy and a trusted root certificate. Root trust is a real security permission; read the setup explanation before accepting it. **Remove after a successful import** is on by default. You can also remove trust from **Import records → Capture & certificate**. A valid saved connection can still sync without reinstalling it.

Keep your app folder. Before updating, close the app and back up `ElmoBeacon.db`. Do not share your database, session files, game logs, or `ca.key` (the certificate’s private signing key). The app shows the certificate folder during setup.

[Certificate details and connection troubleshooting](docs/GUIDE.md#capture-privacy-and-troubleshooting)

## Read your stats

- Custom and Reunion banners have separate pity counters that expire. Outfit records are grouped by season.
- Select a result to see copies, preceding losses, and available doll/weapon costs, including V6R1 and V6R6.
- Use **Edit target result** only when a label is wrong. Choose **Unknown** if you do not know.
- Stats use imported history, not your live inventory. Luck values are estimates; missing records can affect them.

## Maintenance and credits

This fork focuses on bug fixes and game compatibility. No new features or regular release schedule are planned. English and Chinese interfaces are included.

Original code and copyright notices are retained under [GPL v3](LICENSE). Release downloads include corresponding source. Chart: 一棵平衡树; English translation: DiggusDBickus. Game names and assets belong to their owners.

[Build instructions](docs/GUIDE.md#develop-and-build) · [Contributing](CONTRIBUTING.md) · [Security](SECURITY.md)
