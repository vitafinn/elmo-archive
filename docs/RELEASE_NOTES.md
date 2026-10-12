# Elmo Archive 0.2.1

A Windows pull tracker for Girls’ Frontline 2: Exilium. A fork of [ElmoBeacon](https://github.com/MatchaCabin/ElmoBeacon).

## Changes

- Search records using English or Chinese names from installed game data.
- Simpler import and certificate setup, with technical details available when needed.
- Clearer buttons, readable Cancel states, and better spacing in Capture settings.
- Updated English and Chinese instructions.

Capture is built in: no Fiddler or separate capture tool is needed. Open Access Record inside the game. Windows asks before installing the capture certificate; removal after a successful import is on by default.

## Download

Choose **ElmoArchive-0.2.1-windows-x64.zip**, extract it, and run `ElmoArchive.exe`. Close the old app and back up `ElmoBeacon.db` before updating. If WebView2 is missing, the app offers Microsoft's installer.

The executable is unsigned. Live testing across all regions is not complete. Missing history and unknown target outcomes can affect statistics; luck values are estimates.

Matching source and SHA-256 checksums are included. GPL v3. Maintenance focuses on bug fixes and game compatibility.
