package dawg

import (
	"strconv"
	"testing"
)

var testStrings []string = []string{
	"abacaba",
	"Словарь",
	"тест-42",
	"word-ФЫ",
	"Azərbaycanca",
	"Čeština",
	"עברית",
	"日本語",
	"Қазақша",
	"Об'єднаних",
	"а б в г д е ё",
	"а,б,в,г,д,е,ё",
	"中ab文01文",
	"威尔士三",
	"表面か",
	"面には10",
	"　♪リンゴ可愛いや可愛いやリンゴ。半世紀も前に流行した「リンゴの歌」がぴったりするかもしれない。米アップルコンピュータ社のパソコン「マック（マッキントッシュ）」を、こよなく愛する人たちのことだ。「アップル信者」なんて言い方まである。",
	"Tiếng Việt",
	"الأخبار",
	"naïve rèsumé",
	"ܐܠܦ ܒܝܬ ܣܘܪܝܝܐ",
	"সাঁওতালি বাংলা সমশব্দ অভিধান",
	"య్ఢయ్ణ	య్తయ్థయ్ద	య్ధయ్న	య్ప",
	"ნუსხური",
	"እሱ ኢትዮጵያዊ ነው",
	"₠₡	₢	₣	10₤	₥	₦	₧	₨	₩	₪",
	"ᠮᠠᠲᠠᠭᠠᠷ ᠰᠢᠯᠪᠢ",
	"🏴‍☠️🇬🇷❤️🔥",
	"T̪̩̼h̥̫̪͔̀e̫̯͜ ̨N̟e҉͔̤zp̮̭͈̟é͉͈ṛ̹̜̺̭͕d̺̪̜͇͓i̞á͕̹̣̻n͉͘ ̗͔̭͡h̲͖̣̺̺i͔̣̖̤͎̯v̠̯̘͖̭̱̯e̡̥͕-m͖̭̣̬̦͈i͖n̞̩͕̟̼̺͜d̘͉ ̯o̷͇̹͕̦f̰̱ ̝͓͉̱̪̪c͈̲̜̺h̘͚a̞͔̭̰̯̗̝o̙͍s͍͇̱͓.̵͕̰͙͈ͅ ̯̞͈̞̱̖Z̯̮̺̤̥̪̕a͏̺̗̼̬̗ḻg͢o̥̱̼.̺̜͇͡ͅ ̴͓͖̭̩͎̗	̧̪͈̱̹̳͖͙H̵̰̤̰͕̖e̛ ͚͉̗̼̞w̶̩̥͉̮h̩̺̪̩͘ͅọ͎͉̟ ̜̩͔̦̘ͅW̪̫̩̣̲͔̳a͏͔̳͖i͖͜t͓̤̠͓͙s̘̰̩̥̙̝ͅ ̲̠̬̥Be̡̙̫̦h̰̩i̛̫͙͔̭̤̗̲n̳͞d̸ ͎̻͘T̛͇̝̲̹̠̗ͅh̫̦̝ͅe̩̫͟ ͓͖̼W͕̳͎͚̙̥ą̙l̘͚̺͔͞ͅl̳͍̙̤̤̮̳.̢	̟̺̜̙͉Z̤̲̙̙͎̥̝A͎̣͔̙͘L̥̻̗̳̻̳̳͢G͉̖̯͓̞̩̦O̹̹̺!̙͈͎̞̬ *",

	"Тан Турă Амăш турăшĕ — Турамăш Христос пепкене аллинче тытса ларнине сăнарлакан турăш, ăна Елеуса (Юмартлăх) мелĕпе çырнă. Турăш икĕ енлĕ, хыçалти енĕнче Турă Амăш Вилнине ӳкернĕ. Халăх сăмахĕпе (Тан мăнастăрĕн 1692 çулхи кĕнекинче çырнипе), Сиротин хулинчен тан каcакĕсем мускав кнеçне, Тан Дмитрийне Куликово çапăçăвĕ умĕн (1380 çул) парнеленĕ. Вырăс Чиркĕвĕнче асамлă турăш шутланса хисеплĕ вырăн йышăнать.",
	"Jacques Marie Émile Lacan (pronunciación en francés: /ʒak lakɑ̃/; París, 13 de abril de 1901-ibídem, 9 de septiembre de 1981) fue un psiquiatra y psicoanalista francés conocido por los aportes teóricos que hizo al psicoanálisis, sobre la base de la experiencia analítica y en la lectura de Sigmund Freud, combinada con elementos de la filosofía, el estructuralismo, la lingüística estructural y las matemáticas.",
	"Světla velkoměsta (anglicky City Lights) je americký němý film studia United Artists z roku 1931. Snímek režíroval Charlie Chaplin a sám si i zahrál hlavní roli. Film sleduje životní peripetie typické Chaplinovy postavy Tuláka, který se zamiluje do slepé dívky a získá poněkud vrtkavé přátelství alkoholického milionáře. Chaplin trval na konceptu němého filmu, přestože byl v období příprav film zvukový již na svém vzestupu. Natáčení začalo v prosinci 1928 a trvalo až do září 1930. Za pouhých šest týdnů k filmu vznikla doprovodná hudba, kterou poprvé napsal sám režisér, a to ve spolupráci se skladatelem Arthurem Johnstonem.",
	"বিষ্ণুপ্রিয়া মণিপুরী ঠার এহান ভারতর অসম, ত্রিপুরা, মণিপুর বারো বাংলাদেশ, মায়ানমার বাদেউ আরাকউ দেশ কতহাত অতারতারা, ঠার এহান ইন্দো-আর্যর ঠার বাংলা, অহমীয়া, ওড়িয়া ঠারেত্তউ তঙাল। বিষ্ণুপ্রিয়া মণিপুরী ঠার এহান মুলত ভারতর মণিপুর বারো মণিপুর রাজ্যর হমবুকে আসে লগতাকর চারিয়বারেদে আসে লয়া অতাত হঙসেহান বারো মুঙবারেসেহান। ঠারহানর বারাদে হাবির পয়লাকা বা হাব্বিত্ত পুরানা তথ্য উৎসহান ১৮শ শতাব্দীত ইকরিসি পণ্ডিত নবখেন্দ্র শর্ম্মার 'খুমল পুরান' বুলতারা লেরিক এহাত পেয়ার। আরতা উল্লেখযোগ্য উৎস পেয়ারতা মেজর মেককুলাকর An account of the valley of Manipore, ই.টি. ডালটনর Descriptive Ethnology of Bengal বারো স্যার জি.এ. গ্রিয়ার্সনর Linguistic Survey of India লেরিক এহানিত মাতেসিতা ঠার এহান ১৯শ শতাব্দীতর আগে মণিপুরে আসিল। ড. গিয়ার্সন গিরকে ঠার এহানরে বিষ্ণুপুরিয়া মণিপুরী বুলিয়া মাতেসে অন্যতায় যেপাগা হুদ্দা বিষ্ণুপ্রিয়া বুলিয়া মাতেসি। মুলত ঠার এহান মণিপুরর খাঙাবুক, হেইরুক, মিয়াং ইম্ফল, বিষ্ণুপুর খুনৌ, নংথৌখং, ঙাইখং বারো থামানপকপি লয়াত্ত চলিয়া আহেরহান।",
	"イオには400個を超える火山があり、太陽系内で最も地質学的に活発な天体である[5][6]。この極端な地質活動は、木星と他のガリレオ衛星であるエウロパ、ガニメデとの重力相互作用に伴うイオ内部での潮汐加熱の結果である[7]。いくつかの火山は硫黄と二酸化硫黄の噴煙を発生させており、その高さは表面から 500 km にも達する。イオの表面には100以上の山も見られ、イオの岩石地殻の底部における圧縮によって持ち上げられ形成されたと考えられる。これらのうちいくつかはエベレストよりも高い[8]。大部分が水の氷からなる大部分の太陽系遠方の衛星とは異なり、イオの主成分は岩石であり、溶けた鉄もしくは硫化鉄の核を岩石が取り囲んだ構造をしている。イオの表面の大部分は、硫黄と二酸化硫黄の霜で覆われた広い平原からなっている。",
	"Халыкара хатын-кызлар көне — хатын-кызларның ирләр белән хокукый тигезлеге, кешеләр җәмгыятендә хатын-кыз хокукларын киңәйтүгә юнәлгән көрәштә халыкара сәяси чаралар үткәрү өчен билгеләнгән көн, һәр елның 8 март көненә төшә.",
	"پروژه‌نین تاریخی ۱۹۹۹-جۇ ایلدن باشلاییر. پروژه‌نین باش رداکتوْرو، تشکیلاتچیسی لاری سنقر (Larri Senqer) و بوْمیس (Bomis) کوْمپانییاسی‌نین ایجراچی دیرکتوْرو، لاییحه‌نی مالیه‌لشدیرن جیمی اۇلس (Cimmi Uels) ویکی تکنوْلوْژیسی اساسیندا آنلاین بیلیکلیک یاراتماق قرارینا گله‌رک اوْنو نۇپدیا (NuPedia.com) آدلاندیردیلار. نۇپدیا ویرچوال بیلیکلیگی (مجازی بیلیکلیگی)، ۲۰۰۰-جی ایلین مارس آییندان فعالیته باشلادی. اینگیلیس دیلینده یارادیلمیش بۇ بیلیکلیگی ویکی-سایت حساب ائتمک اوْلمازدی. اوْنون اساسینی عالیم و مۆتخصيصلر طرفیندن مقاله‌لرین دقیق یوْخلانماسی تشکیل ائدیردی. نۇپدیادا مقاله‌لرین داخیل اوْلونماسی پروْسه‌سی چوْخ لنگ گئتدیگیندن اوْنون باغلانماسی حاقیندا قرار قبول ائدیلدی و ۲۰۰۳-جۆ ایلین سنتیابر آییندا باغلاندی. نۇپدیا بیلیکلیگی باغلانارکن اوْرادا ۲۴ تاماملانمیش، ۷۴ یوْخلاما پروْسه‌سینده اوْلان مقاله وار ایدی.",
	"Tarihsel olarak Paris Komünü sayılmazsa, Marksist-Leninist ilkeler ilk olarak 1917 yılında gerçekleşen Ekim Devrimi'nden sonra Sovyetler Birliği'nde uygulandı ve ardından devletin resmi ideolojisi haline geldi. Ayrıca ülkede Marksizm-Leninizm Enstitüsü adında bir bilim akademisi bulunmakta ve birçok eser yayınlamaktaydı.",
	"蔡佩軒在2016年時開始在臉書直播自己自彈自唱[5]，定期將自己翻唱的歌曲上傳至網路，在2017年上半年爆紅[6]，當年6月Youtube就已經有14.6萬人[7][8]，當時她仍在學，即便有時差問題，一週仍直播三次，最多超過兩萬人同時觀看[9]；爆紅之後，她在2018年8月舉辦個人首場售票演唱會，演唱會名稱取自她大學畢業前發行的個人單曲《青春有你》，該首單曲吸引了理科太太翻唱[10]，演唱會的一千張門票全數售罄[11]；到了2019年時社群上超過100萬人關注，當年上半年共接獲近百個廠商的邀約[4]，此時她在YouTube影片總點閱次數已經上億，因此許多唱片公司都有合作意願，原本預定當年底發行首張專輯，但對於加拿大回臺灣氣候的調適不適應，導致身體出現不舒服狀況，使其專輯進度延宕[10]；在眾多唱片公司的意願下，她在2020年2月選擇加盟索尼音樂[12]，之後在當年五月釋放的單曲《記得捨不得》成為連續劇《浪漫輸給你》的片頭曲[13]，同月她也為連續劇《若是一個人》獻唱主題曲《愛到明仔載》；之後她在當年八月推出個人首張專輯《ARIEL》，並在當日晚上舉辦讚聲演唱會",
	"Το 1858, ο Κάρολος Δαρβίνος και ο Άλφρεντ Ράσελ Γουάλας δημοσίευσαν μια νέα εξελικτική θεωρία, η οποία εξηγούνταν λεπτομερώς στο έργο του Δαρβίνου, Καταγωγή των Ειδών (On the Origin of Species) (1859). Αντίθετα από τον Λαμάρκ, ο Δαρβίνος πρότεινε κοινή καταγωγή και διακλαδιζόμενο δέντρο της ζωής. Η θεωρία βασιζόταν στην ιδέα της φυσικής επιλογής, και συνέθετε ένα ευρύ φάσμα στοιχείων από την κτηνοτροφία, τη βιογεωγραφία, τη γεωλογία, τη μορφολογία και την εμβρυολογία.",
	"Таким образом, в 1209/10 году первым мужем Тамты стал Аль-Аухад Айюбид[en], сын Аль-Адиля и племянник Саладина. После скорой смерти Аль-Аухада Хлат перешёл под контроль его родного брата Аль-Ашрафа[en]. Тамта, как и Хлат, перешла к Аль-Ашрафу и стала одной из его жён. Тамте удалось добиться снижения налогов для монастырей. В 1230 году Джелал ад-Дин захватил Тамту в плен и сделал своей женой или наложницей.",
	"Брав активну участь у російській операції з анексії Курляндії-Семигалії (1794—1795): керував проросійською фракцією в ландтазі, ініціював прийняття резолюцій про розрив Курляндії з Польщею і приєднання до Російської імперії, очолював курляндську делегацію до Санкт-Петербургу, де оформив анексію і склав присягу на вірність Росії. Отримав від російської імператриці Катерини ІІ посаду таємного радника і маєтки, а від імператора Павла I —— сенаторство і Орден святої Анни 1-го ступеня.",
	"רוקד עם זאבים (אנגלית: Dances with Wolves) הוא סרט דרמה אמריקאי משנת 1990, בבימויו של קווין קוסטנר ובכיכובו. זהו עיבוד קולנועי של ספר בעל אותו השם, ואת התסריט של הסרט כתב מחבר הספר המקורי, מייקל בלייק. עלילת הסרט מתמקדת בחייל מצבא האיחוד שמואס בלחימה במלחמת האזרחים האמריקאית ובוחר לעבור לסְפָר, ובקשרים האמיצים שפיתח עם אינדיאנים משבט סו במערב הנידח. מלבד גרסת הסרט שיצאה בתחילה לקולנוע, התפרסמה כעבור כשנה גרסה ארוכה יותר באורך ארבע שעות.",
	"In den 1960er Jahren wandte er sich dem Kino zu. International bekannt wurde er durch seine Rollen in Filmen von François Truffaut – als Jeanne Moreaus drittes Mordopfer in Die Braut trug schwarz und als Schuhgeschäftsbesitzer Tabard in Geraubte Küsse. Für seine Darstellung des Inspektors Lebel in Der Schakal von Fred Zinnemann erhielt er 1973 eine Nominierung als Bester Nebendarsteller für den BAFTA Award.",
	"Efter krigen uddannede hun sig som komponist, og i begyndelsen af 1950'erne blev hun fanget af den konkrete musik under inspiration af Pierre Schaeffer. Hun komponerede Danmarks første værk inden for denne genre, En dag på Dyrehavsbakken, i 1955. Nogle år senere komponerede hun sit første elektroniske musikværk, Syv cirkler, inspireret af blandt andet Stockhausen, Ligetis og Boulez.",
	"ᐃᖃᓗᐃᑦ, ᓄᓇᕗᑦ (ᓯᑎᐱᕆ 21, 2020) – ᓘᑦᑖᖅ ᒪᐃᑯᓪ ᐸᑐᓴᓐ, ᓄᓇᕗᒻᒥ ᐋᓐᓂᐊᖃᕐᓇᖏᑦᑐᓕᕆᓂᕐᒧᑦ ᐊᖏᔪᖅᑳᖅ, ᐅᓪᓗᒥ ᓇᓗᓇᐃᖅᓯᖅᑲᐅᔪᖅ ᓇᓗᓇᐃᖅᑕᐅᓯᒪᔪᖃᕐᓂᖓᓂᒃ ᓄᕙᔾᔪᐊᕐᓇᖅ 19-ᓕᒻᒥᒃ ᓄᓘᔮᓂ ᐅᔭᕋᕐᓂᐊᕐᕕᖓᓂ ᐅᖓᓯᓐᓂᓕᒃ 176 ᑭᓚᒥᑕᑦ ᓂᒋᐊᑕ ᐱᖓᓐᓇᖓᓂ ᒥᑦᑎᒪᑕᓕᐅᑉ. ᑖᓐᓇ ᐊᐃᑦᑐᕐᓗᑦᑕᐅᓯᒪᙱᑦᑐᖅ ᓄᕙᔾᔪᐊᕐᓇᒥᒃ ᓄᓇᕗᒻᒥ ᐊᒻᒪᓗ ᓇᓗᓇᐃᖅᑕᐅᔪᖅ ᓈᓴᖅᑕᐅᓯᒪᖔᑐᐃᓐᓇᕐᓂᐊᖅᑐᖅ ᐊᖏᕐᕋᖓᑕ ᓄᓇᖓᓂ.",
}

var testZeroStrings []string = []string{
	"\x00", // 0
	"\u00bb",
	"ы\x00", // 11x + 0
	"Č\x00", //
	"蔡\x00", // 101x x x + 100x 0
	"𠀀",     // 101x x 0
	"𢠀",     // 101x 0 0
	"𢠐",     // 101x 0 x
	"𢳘𢤀",    // 101x x x + x 0
	"𢳘𢠀",    // 101x x x + 0 0
	"𢳘𢠐",    // 101x x x + 0 x
	"က",     // 100x 0
	"၀က",    // 100x x + 0
	"⁵",
	"🗯",
	"🗰",
}

func hexString(buf []byte) string {
	s := ""
	for _, v := range buf {
		if v < 16 {
			s += "0"
		}
		s += strconv.FormatInt(int64(v), 16)
		s += " "
	}
	return s
}

func TestUtfc(t *testing.T) {
	for _, test := range testStrings {
		name := test
		if len(name) > 20 {
			name = name[:20] + "…"
		}

		t.Run(name, func(t *testing.T) {
			utfc := UtfcEncode(test)
			t.Logf("String %v encoded as %v", strconv.Quote(test), hexString(utfc))

			ctrl := UtfcDecode(utfc)
			if ctrl != test {
				t.Errorf("String '%v' decoded back as '%v', bytes: %v", test, ctrl, hexString(utfc))
			}
		})
	}
}

func TestUtfcZeros(t *testing.T) {
	for _, test := range testZeroStrings {
		name := strconv.Quote(test)
		name = name[1 : len(name)-1]
		if len(name) > 20 {
			name = name[:20] + "…"
		}

		t.Run(name, func(t *testing.T) {
			utfc := UtfcEncode(test)
			t.Logf("String %v encoded as %v", strconv.Quote(test), hexString(utfc))

			ctrl := UtfcDecode(utfc)
			t.Logf("String %v decoded back as %v", strconv.Quote(test), strconv.Quote(ctrl))
			if ctrl != test {
				t.Fail()
			}
			for i := 0; i < len(utfc); i++ {
				if utfc[i] == 0 {
					t.Fatalf("Found 00 byte at position %d!", i)
					break
				}
			}
		})
	}

}
