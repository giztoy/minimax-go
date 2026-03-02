# Voice List Example

`examples/voice/list` demonstrates voice list queries via `Voice.ListVoices`, including voice type filter and pagination parameters.

## Quick start

```bash
export MINIMAX_API_KEY="your_api_key"

go run ./examples/voice/list \
  -voice-type all \
  -page-size 20
```

Use `-json` to print the formatted typed response as JSON (`Raw` unknown fields are not included):

```bash
go run ./examples/voice/list \
  -voice-type system \
  -json
```

## Show all CLI options

```bash
go run ./examples/voice/list -h
```

## Common flags

- `-api-key`: Minimax API key (takes precedence over `MINIMAX_API_KEY`)
- `-base-url`: API endpoint (default: `https://api.minimax.io`)
- `-voice-type`: voice type filter (`system`, `voice_cloning`, `voice_generation`, `all`)
- `-page-size`: optional page size
- `-page-token`: optional next-page token
- `-timeout`: request timeout (default: `30s`)
- `-json`: print response as formatted JSON

## Environment variables

- `MINIMAX_API_KEY`
- `MINIMAX_BASE_URL`
- `MINIMAX_VOICE_TYPE`
- `MINIMAX_VOICE_PAGE_SIZE`
- `MINIMAX_VOICE_PAGE_TOKEN`
- `MINIMAX_VOICE_TIMEOUT`

## Official non-cloning voices (voice_type=system)

Snapshot time: 2026-03-02 00:53:40 UTC
Total official voices: 303

| # | voice_id | voice_name |
|---:|---|---|
| 1 | `male-qn-qingse` | 青涩青年音色 |
| 2 | `male-qn-jingying` | 精英青年音色 |
| 3 | `male-qn-badao` | 霸道青年音色 |
| 4 | `male-qn-daxuesheng` | 青年大学生音色 |
| 5 | `female-shaonv` | 少女音色 |
| 6 | `female-yujie` | 御姐音色 |
| 7 | `female-chengshu` | 成熟女性音色 |
| 8 | `female-tianmei` | 甜美女性音色 |
| 9 | `male-qn-qingse-jingpin` | 青涩青年音色-beta |
| 10 | `male-qn-jingying-jingpin` | 精英青年音色-beta |
| 11 | `male-qn-badao-jingpin` | 霸道青年音色-beta |
| 12 | `male-qn-daxuesheng-jingpin` | 青年大学生音色-beta |
| 13 | `female-shaonv-jingpin` | 少女音色-beta |
| 14 | `female-yujie-jingpin` | 御姐音色-beta |
| 15 | `female-chengshu-jingpin` | 成熟女性音色-beta |
| 16 | `female-tianmei-jingpin` | 甜美女性音色-beta |
| 17 | `clever_boy` | 聪明男童 |
| 18 | `cute_boy` | 可爱男童 |
| 19 | `lovely_girl` | 萌萌女童 |
| 20 | `cartoon_pig` | 卡通猪小琪 |
| 21 | `bingjiao_didi` | 病娇弟弟 |
| 22 | `junlang_nanyou` | 俊朗男友 |
| 23 | `chunzhen_xuedi` | 纯真学弟 |
| 24 | `lengdan_xiongzhang` | 冷淡学长 |
| 25 | `badao_shaoye` | 霸道少爷 |
| 26 | `tianxin_xiaoling` | 甜心小玲 |
| 27 | `qiaopi_mengmei` | 俏皮萌妹 |
| 28 | `wumei_yujie` | 妩媚御姐 |
| 29 | `diadia_xuemei` | 嗲嗲学妹 |
| 30 | `danya_xuejie` | 淡雅学姐 |
| 31 | `Santa_Claus` | Santa Claus |
| 32 | `Grinch` | Grinch |
| 33 | `Rudolph` | Rudolph |
| 34 | `Arnold` | Arnold |
| 35 | `Charming_Santa` | Charming Santa |
| 36 | `Charming_Lady` | Charming Lady |
| 37 | `Sweet_Girl` | Sweet Girl |
| 38 | `Cute_Elf` | Cute Elf |
| 39 | `Attractive_Girl` | Attractive Girl |
| 40 | `Serene_Woman` | Serene Woman |
| 41 | `Chinese (Mandarin)_Reliable_Executive` | 沉稳高管 |
| 42 | `Chinese (Mandarin)_News_Anchor` | 新闻女声 |
| 43 | `Chinese (Mandarin)_Mature_Woman` | 傲娇御姐 |
| 44 | `Chinese (Mandarin)_Unrestrained_Young_Man` | 不羁青年 |
| 45 | `Arrogant_Miss` | 嚣张小姐 |
| 46 | `Robot_Armor` | 机械战甲 |
| 47 | `Chinese (Mandarin)_Kind-hearted_Antie` | 热心大婶 |
| 48 | `Chinese (Mandarin)_HK_Flight_Attendant` | 港普空姐 |
| 49 | `Chinese (Mandarin)_Humorous_Elder` | 搞笑大爷 |
| 50 | `Chinese (Mandarin)_Gentleman` | 温润男声 |
| 51 | `Chinese (Mandarin)_Warm_Bestie` | 温暖闺蜜 |
| 52 | `Chinese (Mandarin)_Male_Announcer` | 播报男声 |
| 53 | `Chinese (Mandarin)_Sweet_Lady` | 甜美女声 |
| 54 | `Chinese (Mandarin)_Southern_Young_Man` | 南方小哥 |
| 55 | `Chinese (Mandarin)_Wise_Women` | 阅历姐姐 |
| 56 | `Chinese (Mandarin)_Gentle_Youth` | 温润青年 |
| 57 | `Chinese (Mandarin)_Warm_Girl` | 温暖少女 |
| 58 | `Chinese (Mandarin)_Kind-hearted_Elder` | 花甲奶奶 |
| 59 | `Chinese (Mandarin)_Cute_Spirit` | 憨憨萌兽 |
| 60 | `Chinese (Mandarin)_Radio_Host` | 电台男主播 |
| 61 | `Chinese (Mandarin)_Lyrical_Voice` | 抒情男声 |
| 62 | `Chinese (Mandarin)_Straightforward_Boy` | 率真弟弟 |
| 63 | `Chinese (Mandarin)_Sincere_Adult` | 真诚青年 |
| 64 | `Chinese (Mandarin)_Gentle_Senior` | 温柔学姐 |
| 65 | `Chinese (Mandarin)_Stubborn_Friend` | 嘴硬竹马 |
| 66 | `Chinese (Mandarin)_Crisp_Girl` | 清脆少女 |
| 67 | `Chinese (Mandarin)_Pure-hearted_Boy` | 清澈邻家弟弟 |
| 68 | `Chinese (Mandarin)_Soft_Girl` | 软软女孩 |
| 69 | `Cantonese_ProfessionalHost（F)` | 专业女主持 |
| 70 | `Cantonese_GentleLady` | 温柔女声 |
| 71 | `Cantonese_ProfessionalHost（M)` | 专业男主持 |
| 72 | `Cantonese_PlayfulMan` | 活泼男声 |
| 73 | `Cantonese_CuteGirl` | 可爱女孩 |
| 74 | `Cantonese_KindWoman` | 善良女声 |
| 75 | `English_Trustworthy_Man` | Trustworthy Man |
| 76 | `English_Graceful_Lady` | Graceful Lady |
| 77 | `English_Aussie_Bloke` | Aussie Bloke |
| 78 | `English_Whispering_girl` | Whispering girl |
| 79 | `English_Diligent_Man` | Diligent Man |
| 80 | `English_Gentle-voiced_man` | Gentle-voiced man |
| 81 | `Japanese_IntellectualSenior` | Intellectual Senior |
| 82 | `Japanese_DecisivePrincess` | Decisive Princess |
| 83 | `Japanese_LoyalKnight` | Loyal Knight |
| 84 | `Japanese_DominantMan` | Dominant Man |
| 85 | `Japanese_SeriousCommander` | Serious Commander |
| 86 | `Japanese_ColdQueen` | Cold Queen |
| 87 | `Japanese_DependableWoman` | Dependable Woman |
| 88 | `Japanese_GentleButler` | Gentle Butler |
| 89 | `Japanese_KindLady` | Kind Lady |
| 90 | `Japanese_CalmLady` | Calm Lady |
| 91 | `Japanese_OptimisticYouth` | Optimistic Youth |
| 92 | `Japanese_GenerousIzakayaOwner` | Generous Izakaya Owner |
| 93 | `Japanese_SportyStudent` | Sporty Student |
| 94 | `Japanese_InnocentBoy` | Innocent Boy |
| 95 | `Japanese_GracefulMaiden` | Graceful Maiden |
| 96 | `Dutch_kindhearted_girl` | Kind-hearted girl |
| 97 | `Dutch_bossy_leader` | Bossy leader |
| 98 | `Vietnamese_kindhearted_girl` | Kind-hearted girl |
| 99 | `Korean_SweetGirl` | Sweet Girl |
| 100 | `Korean_CheerfulBoyfriend` | Cheerful Boyfriend |
| 101 | `Korean_EnchantingSister` | Enchanting Sister |
| 102 | `Korean_ShyGirl` | Shy Girl |
| 103 | `Korean_ReliableSister` | Reliable Sister |
| 104 | `Korean_StrictBoss` | Strict Boss |
| 105 | `Korean_SassyGirl` | Sassy Girl |
| 106 | `Korean_ChildhoodFriendGirl` | Childhood Friend Girl |
| 107 | `Korean_PlayboyCharmer` | Playboy Charmer |
| 108 | `Korean_ElegantPrincess` | Elegant Princess |
| 109 | `Korean_BraveFemaleWarrior` | Brave Female Warrior |
| 110 | `Korean_BraveYouth` | Brave Youth |
| 111 | `Korean_CalmLady` | Calm Lady |
| 112 | `Korean_EnthusiasticTeen` | Enthusiastic Teen |
| 113 | `Korean_SoothingLady` | Soothing Lady |
| 114 | `Korean_IntellectualSenior` | Intellectual Senior |
| 115 | `Korean_LonelyWarrior` | Lonely Warrior |
| 116 | `Korean_MatureLady` | Mature Lady |
| 117 | `Korean_InnocentBoy` | Innocent Boy |
| 118 | `Korean_CharmingSister` | Charming Sister |
| 119 | `Korean_AthleticStudent` | Athletic Student |
| 120 | `Korean_BraveAdventurer` | Brave Adventurer |
| 121 | `Korean_CalmGentleman` | Calm Gentleman |
| 122 | `Korean_WiseElf` | Wise Elf |
| 123 | `Korean_CheerfulCoolJunior` | Cheerful Cool Junior |
| 124 | `Korean_DecisiveQueen` | Decisive Queen |
| 125 | `Korean_ColdYoungMan` | Cold Young Man |
| 126 | `Korean_MysteriousGirl` | Mysterious Girl |
| 127 | `Korean_QuirkyGirl` | Quirky Girl |
| 128 | `Korean_ConsiderateSenior` | Considerate Senior |
| 129 | `Korean_CheerfulLittleSister` | Cheerful Little Sister |
| 130 | `Korean_DominantMan` | Dominant Man |
| 131 | `Korean_AirheadedGirl` | Airheaded Girl |
| 132 | `Korean_ReliableYouth` | Reliable Youth |
| 133 | `Korean_FriendlyBigSister` | Friendly Big Sister |
| 134 | `Korean_GentleBoss` | Gentle Boss |
| 135 | `Korean_ColdGirl` | Cold Girl |
| 136 | `Korean_HaughtyLady` | Haughty Lady |
| 137 | `Korean_CharmingElderSister` | Charming Elder Sister |
| 138 | `Korean_IntellectualMan` | Intellectual Man |
| 139 | `Korean_CaringWoman` | Caring Woman |
| 140 | `Korean_WiseTeacher` | Wise Teacher |
| 141 | `Korean_ConfidentBoss` | Confident Boss |
| 142 | `Korean_AthleticGirl` | Athletic Girl |
| 143 | `Korean_PossessiveMan` | Possessive Man |
| 144 | `Korean_GentleWoman` | Gentle Woman |
| 145 | `Korean_CockyGuy` | Cocky Guy |
| 146 | `Korean_ThoughtfulWoman` | Thoughtful Woman |
| 147 | `Korean_OptimisticYouth` | Optimistic Youth |
| 148 | `Spanish_SereneWoman` | Serene Woman |
| 149 | `Spanish_MaturePartner` | Mature Partner |
| 150 | `Spanish_CaptivatingStoryteller` | Captivating Storyteller |
| 151 | `Spanish_Narrator` | Narrator |
| 152 | `Spanish_WiseScholar` | Wise Scholar |
| 153 | `Spanish_Kind-heartedGirl` | Kind-hearted Girl |
| 154 | `Spanish_DeterminedManager` | Determined Manager |
| 155 | `Spanish_BossyLeader` | Bossy Leader |
| 156 | `Spanish_ReservedYoungMan` | Reserved Young Man |
| 157 | `Spanish_ConfidentWoman` | Confident Woman |
| 158 | `Spanish_ThoughtfulMan` | Thoughtful Man |
| 159 | `Spanish_Strong-WilledBoy` | Strong-willed Boy |
| 160 | `Spanish_SophisticatedLady` | Sophisticated Lady |
| 161 | `Spanish_RationalMan` | Rational Man |
| 162 | `Spanish_AnimeCharacter` | Anime Character |
| 163 | `Spanish_Deep-tonedMan` | Deep-toned Man |
| 164 | `Spanish_Fussyhostess` | Fussy hostess |
| 165 | `Spanish_SincereTeen` | Sincere Teen |
| 166 | `Spanish_FrankLady` | Frank Lady |
| 167 | `Spanish_Comedian` | Comedian |
| 168 | `Spanish_Debator` | Debator |
| 169 | `Spanish_ToughBoss` | Tough Boss |
| 170 | `Spanish_Wiselady` | Wise Lady |
| 171 | `Spanish_Steadymentor` | Steady Mentor |
| 172 | `Spanish_Jovialman` | Jovial Man |
| 173 | `Spanish_SantaClaus` | Santa Claus |
| 174 | `Spanish_Rudolph` | Rudolph |
| 175 | `Spanish_Intonategirl` | Intonate Girl |
| 176 | `Spanish_Arnold` | Arnold |
| 177 | `Spanish_Ghost` | Ghost |
| 178 | `Spanish_HumorousElder` | Humorous Elder |
| 179 | `Spanish_EnergeticBoy` | Energetic Boy |
| 180 | `Spanish_WhimsicalGirl` | Whimsical Girl |
| 181 | `Spanish_StrictBoss` | Strict Boss |
| 182 | `Spanish_ReliableMan` | Reliable Man |
| 183 | `Spanish_SereneElder` | Serene Elder |
| 184 | `Spanish_AngryMan` | Angry Man |
| 185 | `Spanish_AssertiveQueen` | Assertive Queen |
| 186 | `Spanish_CaringGirlfriend` | Caring Girlfriend |
| 187 | `Spanish_PowerfulSoldier` | Powerful Soldier |
| 188 | `Spanish_PassionateWarrior` | Passionate Warrior |
| 189 | `Spanish_ChattyGirl` | Chatty Girl |
| 190 | `Spanish_RomanticHusband` | Romantic Husband |
| 191 | `Spanish_CompellingGirl` | Compelling Girl |
| 192 | `Spanish_PowerfulVeteran` | Powerful Veteran |
| 193 | `Spanish_SensibleManager` | Sensible Manager |
| 194 | `Spanish_ThoughtfulLady` | Thoughtful Lady |
| 195 | `Portuguese_SentimentalLady` | Sentimental Lady |
| 196 | `Portuguese_BossyLeader` | Bossy Leader |
| 197 | `Portuguese_Wiselady` | Wise lady |
| 198 | `Portuguese_Strong-WilledBoy` | Strong-willed Boy |
| 199 | `Portuguese_Deep-VoicedGentleman` | Deep-voiced Gentleman |
| 200 | `Portuguese_UpsetGirl` | Upset Girl |
| 201 | `Portuguese_PassionateWarrior` | Passionate Warrior |
| 202 | `Portuguese_AnimeCharacter` | Anime Character |
| 203 | `Portuguese_ConfidentWoman` | Confident Woman |
| 204 | `Portuguese_AngryMan` | Angry Man |
| 205 | `Portuguese_CaptivatingStoryteller` | Captivating Storyteller |
| 206 | `Portuguese_Godfather` | Godfather |
| 207 | `Portuguese_ReservedYoungMan` | Reserved Young Man |
| 208 | `Portuguese_SmartYoungGirl` | Smart Young Girl |
| 209 | `Portuguese_Kind-heartedGirl` | Kind-hearted Girl |
| 210 | `Portuguese_Pompouslady` | Pompous lady |
| 211 | `Portuguese_Grinch` | Grinch |
| 212 | `Portuguese_Debator` | Debator |
| 213 | `Portuguese_SweetGirl` | Sweet Girl |
| 214 | `Portuguese_AttractiveGirl` | Attractive Girl |
| 215 | `Portuguese_ThoughtfulMan` | Thoughtful Man |
| 216 | `Portuguese_PlayfulGirl` | Playful Girl |
| 217 | `Portuguese_GorgeousLady` | Gorgeous Lady |
| 218 | `Portuguese_LovelyLady` | Lovely Lady |
| 219 | `Portuguese_SereneWoman` | Serene Woman |
| 220 | `Portuguese_SadTeen` | Sad Teen |
| 221 | `Portuguese_MaturePartner` | Mature Partner |
| 222 | `Portuguese_Comedian` | Comedian |
| 223 | `Portuguese_NaughtySchoolgirl` | Naughty Schoolgirl |
| 224 | `Portuguese_Narrator` | Narrator |
| 225 | `Portuguese_ToughBoss` | Tough Boss |
| 226 | `Portuguese_Fussyhostess` | Fussy hostess |
| 227 | `Portuguese_Dramatist` | Dramatist |
| 228 | `Portuguese_Steadymentor` | Steady Mentor |
| 229 | `Portuguese_Jovialman` | Jovial Man |
| 230 | `Portuguese_CharmingQueen` | Charming Queen |
| 231 | `Portuguese_SantaClaus` | Santa Claus |
| 232 | `Portuguese_Rudolph` | Rudolph |
| 233 | `Portuguese_Arnold` | Arnold |
| 234 | `Portuguese_CharmingSanta` | Charming Santa |
| 235 | `Portuguese_CharmingLady` | Charming Lady |
| 236 | `Portuguese_Ghost` | Ghost |
| 237 | `Portuguese_HumorousElder` | Humorous Elder |
| 238 | `Portuguese_CalmLeader` | Calm Leader |
| 239 | `Portuguese_GentleTeacher` | Gentle Teacher |
| 240 | `Portuguese_EnergeticBoy` | Energetic Boy |
| 241 | `Portuguese_ReliableMan` | Reliable Man |
| 242 | `Portuguese_SereneElder` | Serene Elder |
| 243 | `Portuguese_GrimReaper` | Grim Reaper |
| 244 | `Portuguese_AssertiveQueen` | Assertive Queen |
| 245 | `Portuguese_WhimsicalGirl` | Whimsical Girl |
| 246 | `Portuguese_StressedLady` | Stressed Lady |
| 247 | `Portuguese_FriendlyNeighbor` | Friendly Neighbor |
| 248 | `Portuguese_CaringGirlfriend` | Caring Girlfriend |
| 249 | `Portuguese_PowerfulSoldier` | Powerful Soldier |
| 250 | `Portuguese_FascinatingBoy` | Fascinating Boy |
| 251 | `Portuguese_RomanticHusband` | Romantic Husband |
| 252 | `Portuguese_StrictBoss` | Strict Boss |
| 253 | `Portuguese_InspiringLady` | Inspiring Lady |
| 254 | `Portuguese_PlayfulSpirit` | Playful Spirit |
| 255 | `Portuguese_ElegantGirl` | Elegant Girl |
| 256 | `Portuguese_CompellingGirl` | Compelling Girl |
| 257 | `Portuguese_PowerfulVeteran` | Powerful Veteran |
| 258 | `Portuguese_SensibleManager` | Sensible Manager |
| 259 | `Portuguese_ThoughtfulLady` | Thoughtful Lady |
| 260 | `Portuguese_TheatricalActor` | Theatrical Actor |
| 261 | `Portuguese_FragileBoy` | Fragile Boy |
| 262 | `Portuguese_ChattyGirl` | Chatty Girl |
| 263 | `Portuguese_Conscientiousinstructor` | Conscientious Instructor |
| 264 | `Portuguese_RationalMan` | Rational Man |
| 265 | `Portuguese_WiseScholar` | Wise Scholar |
| 266 | `Portuguese_FrankLady` | Frank Lady |
| 267 | `Portuguese_DeterminedManager` | Determined Manager |
| 268 | `French_Male_Speech_New` | Level-Headed Man |
| 269 | `French_Female_News Anchor` | Patient Female Presenter |
| 270 | `French_CasualMan` | Casual Man |
| 271 | `French_MovieLeadFemale` | Movie Lead Female |
| 272 | `French_FemaleAnchor` | Female Anchor |
| 273 | `French_MaleNarrator` | Male Narrator |
| 274 | `Indonesian_SweetGirl` | Sweet Girl |
| 275 | `Indonesian_ReservedYoungMan` | Reserved Young Man |
| 276 | `Indonesian_CharmingGirl` | Charming Girl |
| 277 | `Indonesian_CalmWoman` | Calm Woman |
| 278 | `Indonesian_ConfidentWoman` | Confident Woman |
| 279 | `Indonesian_CaringMan` | Caring Man |
| 280 | `Indonesian_BossyLeader` | Bossy Leader |
| 281 | `Indonesian_DeterminedBoy` | Determined Boy |
| 282 | `Indonesian_GentleGirl` | Gentle Girl |
| 283 | `German_FriendlyMan` | Friendly Man |
| 284 | `German_SweetLady` | Sweet Lady |
| 285 | `German_PlayfulMan` | Playful Man |
| 286 | `Russian_HandsomeChildhoodFriend` | Handsome Childhood Friend |
| 287 | `Russian_BrightHeroine` | Bright Queen |
| 288 | `Russian_AmbitiousWoman` | Ambitious Woman |
| 289 | `Russian_ReliableMan` | Reliable Man |
| 290 | `Russian_CrazyQueen` | Crazy Girl |
| 291 | `Russian_PessimisticGirl` | Pessimistic Girl |
| 292 | `Russian_AttractiveGuy` | Attractive Guy |
| 293 | `Russian_Bad-temperedBoy` | Bad-tempered Boy |
| 294 | `Italian_BraveHeroine` | Brave Heroine |
| 295 | `Italian_Narrator` | Narrator |
| 296 | `Italian_WanderingSorcerer` | Wandering Sorcerer |
| 297 | `Italian_DiligentLeader` | Diligent Leader |
| 298 | `Arabic_CalmWoman` | Calm Woman |
| 299 | `Arabic_FriendlyGuy` | Friendly Guy |
| 300 | `Turkish_CalmWoman` | Calm Woman |
| 301 | `Turkish_Trustworthyman` | Trustworthy man |
| 302 | `Ukrainian_CalmWoman` | Calm Woman |
| 303 | `Ukrainian_WiseScholar` | Wise Scholar |
