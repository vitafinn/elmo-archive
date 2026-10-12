# Procurement categories and statistics

Custom Procurement and Reunion Procurement are separate banners sharing API history
types 6 (dolls) and 7 (weapons). The tracker keeps raw records unchanged and separates
views by banner identity. Reunion banners 252001/253001 and client name IDs
559138/559139 identify Reunion. Display-only categories -6/-7 are never sent to the
history API. New banners reusing those name IDs are discovered automatically.

Each Custom/Reunion banner period has an independent pull counter, even when targets
change within it. Completed five-star sequences determine the average pulls to a
five-star. All pulls, including unfinished sequences and banners with no five-stars,
contribute to pulls spent per five-star. Overall averages weight individual pulls and
drops, not banner averages. With zero five-stars the UI shows a dash, not a zero cost.
For 80 pulls with a five-star at 60, the completed average is 60 and spending is 80.
The remaining 20 are active pity until the server-specific end timestamp passes,
then become expired pity. Missing dates remain unknown, never guessed from inactivity.
Expiry does not alter imported history. The UI refreshes every 30 seconds and labels
counters as last-recorded values: sync near closing to confirm the final unused pity.

Target outcomes are distinct from pity. The first recorded five-star defaults to an
initial guarantee, assuming complete history. Further results stay unknown, regardless
of whether the item is a standard doll, until explicitly annotated. A known loss
makes the next automatic outcome guaranteed. Only annotated unguaranteed wins and
losses enter target win rates; no eligible attempts displays N/A. Users can override
outcomes, including the first record when history is incomplete, and restore Automatic.
Annotations are stored in the settings table by user and original record ID. They
never mutate raw records or pity. Historical selection changes cannot be reconstructed
from the fields currently imported. Weapon and doll outcomes are not assumed to share
the same base odds (75% weapons versus 50% dolls after the initial guarantee).

The sidebar reads client categories and the selected user's saved history. Unknown
positive API categories receive numeric labels and neutral record colours. New protocols
or banner families still require identification. Table parsing supports current and
legacy layouts and reads server-specific schedule fields 7/8 for expiry.

Validation: go test ./... and npm run build in frontend. ELMOBEACON_TEST_GAME_DATA enables
read-only installed-client tests; ELMOBEACON_TEST_DB enables a read-only history check.
Generate Wails bindings with wails generate module before building the frontend when
Go response types change, then use wails build -s to embed the frontend production build.

## Outfit Procurement

Outfit history types 9/10 have dedicated access-record views, each grouped by banner
(season) ID. Type 10 remains separate because it is a distinct API history category.
All accesses are displayed, including currency, normal rewards, and unknown item IDs.
Reward categories and the main outfit come from the region-specific GachaClothesData
fields 7/9/11. Season names and dates come from GachaData and the language table;
no fixed 365-day interval is imposed on server schedules. Concurrent seasons keep
independent access numbering. Main-outfit acquisition is shown by recorded access
number, without applying doll/weapon pity rules. Counts represent access/reward
occurrences, not quantities of items within a reward. Original records are unchanged.
Existing Synchronize Records already fetches these categories from the client type
list. Sync matching now also checks banner ID to avoid confusing same-time rewards
from different seasons. The UI shows all currently active seasons plus recorded
historical seasons and paginates the full access list in groups of 50.

## Five-star timing indicators

Result cards and copy details distinguish rarity timing from target outcomes. The
indicator uses an explicitly approximate linear soft-pity model: base 0.006 through
58 pulls, increasing linearly to 1 at 80 for standard/doll categories; base 0.007
through 50, increasing to 1 at 70 for weapon categories. CDF(n) = 1 - product(1-p(k)).
Only known ordinary categories 1/3/4/6/7/-6/-7 are rated. Invalid counts and special
or future types remain unrated. These are modeled probabilities, not account
percentiles or verified server per-pull rates. Imported sequences may be incomplete.

Early is CDF <= 25%; typical <= 75%; late <= 90%; very late > 90%. Labels accompany
colors. Target wins, losses, and guarantees remain separate. A cumulative target
fortification chart (including lost target rolls) cannot be applied to the sum of
sequences ending in a particular character, especially with first-target guarantees
and unknown target switches. The UI explains this instead of assigning a misleading
combined target-luck score.

Mechanics reference: https://iopwiki.com/wiki/GFL2 (soft/hard pity thresholds).
The precise linear increase is an approximation, not a claimed official curve.
Run `node --test frontend/tests/luck.test.mjs` for model boundary checks.

Banner titles use the installed client language table when available. Read-only
inspection on this client resolved banner 2000011 to type 6, name ID 82050,
"Custom Procurement - Dolls". Reunion 252001 is also type 6 but name ID 559138,
"Reunion Procurement - Doll". Classification is not derived from the numeric banner
ID's appearance. Existing historical fallbacks remain in service/PoolCategory.go.

## Acquisition costs including preceding losses

Character details now retain both the five-star sequence count and an acquisition
cost that can include the immediately preceding loss. Normal limited categories
link a preceding non-featured result to the next featured result; verified metadata,
inferred outcomes, and cross-banner carried guarantees have distinct labels.
Custom/Reunion require an identified loss followed by a guaranteed result in the
same banner. Attribution to the same target is explicitly estimated because target
switches are unavailable. Unknown results are never silently counted as losses.
A preceding copy of the same character is not counted twice in that character's
total. Unfinished pity is excluded; incomplete history can understate costs.

The overview shows total and average recorded acquisition cost, direct versus loss
contributions, lowest/highest copy costs, and cumulative vertebra milestones. The
individual five-star timing indicator still uses only the final sequence count.
Tests: `go test ./handler` and `node --test frontend/tests/itemHistory.test.mjs`.

## Signature weapon costs and appearance

Character details read GunData field 24 for the corresponding weapon, validate both item types and membership in GunWeaponData, and resolve the localized weapon name. General recommended-weapon lists are not treated as signature links. Missing links, failed lookups, and absent imported weapon copies are distinct from zero cost.

The one-weapon total adds the earliest imported weapon acquisition (including its linked preceding loss) to the doll acquisition total. The V0–V6 selector can include this first weapon once. V6R6 appears only with at least seven doll and six weapon records, sums exactly the first seven and first six acquisition costs, and excludes excess copies. These are recorded costs, not live inventory or upgrade levels; existing attribution caveats apply. No raw pull records are changed.

Graphite & Jade is the default for new appearance preferences. Saved Slate/Deep Contrast choices remain intact; all three are available in Appearance.

Cost labels now distinguish loss-inclusive acquisition totals from individual five-star sequence lengths. The highest acquisition cost exposes the two sequences and its attribution caveat. VxR1 uses the first x+1 doll copies (capped at seven) plus the first weapon copy, not all excess copies. Both Simplified and Traditional Chinese include the new application messages; a key/parameter coverage test checks them against English.

## Import review and diagnostics

The UI prepares an import before syncing, displaying the exact UID as a decimal string, publisher/region, and identification source. A checkbox and explicit confirmation are required. The backend holds the reviewed connection snapshot for five minutes, consumes its ticket once, rejects a different UID or stale/replayed ticket, and prevents overlapping confirmed imports. Preparing never fetches/saves pull records. A Darkwinter webshop identity cannot be used with a Haoplay history connection. Webshop UID and captured game-history ownership cannot be cross-verified automatically; the review states this clearly. UID numeric decoding preserves large integers. Raw webshop response logging was removed.

Verbose import progress is optional and session-local (latest 300 entries). Events contain fixed stage names, category numbers, pages and counts only; no URLs, tokens, emails or response bodies. The appearance preference for displaying progress is saved. Copy detail dialogs open at the summary. The footer button locates the clicked record on request. A selected badge identifies that record.

## UI language and capture guidance

Use ASD-STE100 principles where useful: short action labels, direct commands, and consistent game terms. Use longer, natural sentences when explanations need context; clarity takes priority over strict controlled-language wording. Keep product names, UID, pity, V/R notation, and procurement names as domain terms. This is a writing policy, not a claim of full dictionary or standard conformance. Reference: https://www.asd-ste100.org/about_STE.html . Chinese translations preserve the same actions and meanings.

The webshop capture stage has an always-visible success state independent of verbose logging. It tells the user to open the game after the UID event. Capture start clears both stage flags. The copy locator is in the dialog footer so it remains available while the details scroll.

