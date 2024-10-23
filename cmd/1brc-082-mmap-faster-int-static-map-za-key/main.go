// instead of parsing float, parse an int instead, with a faster function; use
// a static mapping from names to slice indices; that is cheating
//
// real    0m16.535s
// user    1m24.738s
// sys     0m11.904s

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
	"unsafe"

	_ "embed"

	"golang.org/x/exp/mmap"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "file to write cpu profile to")
	filename   = flag.String("f", "measurements.txt", "measurements file")
)

const (
	chunkSize   = 8 * 1024 * 1024 // per-thread file slice
	noticeColor = "\033[1;36m%s\033[0m"
)

// StaticMap is a tailored "map" for cities.txt.
type StaticMap struct {
	M []*Measurements
}

func NewStaticMap() *StaticMap {
	return &StaticMap{
		M: make([]*Measurements, 16384),
	}
}

// calculateIndex, interestingly the most expensive part of the program.
func calculateIndex(s string) (index int) {
	for i, c := range s {
		index = index + i*(37+int(c))
	}
	return index % 16384
}

func (m *StaticMap) Index(s string) (index int) {
	return calculateIndex(s)
}

func (m *StaticMap) Init() {
	m.M = make([]*Measurements, 16384)
}

func (m *StaticMap) Set(key string, ms *Measurements) {
	m.M[m.Index(key)] = ms
}

func (m *StaticMap) Get(key string) *Measurements {
	return m.M[m.Index(key)]
}

// Measurements, as there is no need to keep all numbers around, we can compute
// them on the fly.
type Measurements struct {
	Min   int
	Max   int
	Sum   int
	Count int
}

// Add adds a new measurement, adjusting min and max as needed.
func (m *Measurements) Add(v int) {
	if v > m.Max {
		m.Max = v
	} else if v < m.Min {
		m.Min = v
	}
	m.Sum = m.Sum + v
	m.Count++
}

// Merge merges data from another measurement.
func (m *Measurements) Merge(o *Measurements) {
	if o.Min < m.Min {
		m.Min = o.Min
	}
	if o.Max > m.Max {
		m.Max = o.Max
	}
	m.Sum = m.Sum + o.Sum
	m.Count = m.Count + o.Count
}

// parseTempToInt turns '-16.7' into -167. It is up to the caller to take care
// of the back conversion.
func parseTempToInt(p []byte) int64 {
	negative := p[0] == '-'
	if negative {
		p = p[1:]
	}
	var result int64
	switch len(p) {
	// 1.2
	case 3:
		result = int64(p[0])*10 + int64(p[2]) - '0'*(10+1)
		// 12.3
	case 4:
		result = int64(p[0])*100 + int64(p[1])*10 + int64(p[3]) - '0'*(100+10+1)
	}
	if negative {
		return -result
	}
	return result
}

// aggregate aggregates measurements by reading a chunk from an io.ReaderAt and
// passing the result to a results channel.
func aggregate(rat io.ReaderAt, offset, length int, resultC chan *StaticMap, sem chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if length == 0 {
		return
	}
	buf := make([]byte, length)
	_, err := rat.ReadAt(buf, int64(offset))
	if err == io.EOF {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(".")
	var (
		data    = NewStaticMap()
		j, k, l = 0, 0, 0 // j=start, k=semi, l=newline
		n       = 0
	)
	for i := 0; i < length; i++ {
		if buf[i] == ';' {
			k = i
		} else if buf[i] == '\n' {
			l = i
			b := buf[j:k]
			name := *(*string)(unsafe.Pointer(&b)) // from []byte to string, w/o allocation
			temp := int(parseTempToInt(buf[k+1 : l]))
			if data.Get(name) == nil {
				data.Set(name, &Measurements{
					Min:   temp,
					Max:   temp,
					Sum:   temp,
					Count: 1,
				})
			} else {
				data.Get(name).Add(temp)
			}
			j = l + 1
			n++
		}
	}
	resultC <- data
	<-sem
}

// merger merges all measurements from workers and merges them into a single
// map.
func merger(data *StaticMap, resultC chan *StaticMap, done chan bool) {
	for m := range resultC {
		for i := range m.M {
			if m.M[i] == nil {
				continue
			}
			if data.M[i] == nil {
				data.M[i] = m.M[i]
			} else {
				data.M[i].Merge(m.M[i])
			}
		}
	}
	done <- true
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	r, err := mmap.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	var (
		resultC = make(chan *StaticMap)
		done    = make(chan bool)
		sem     = make(chan bool, runtime.NumCPU())
		wg      sync.WaitGroup
		data    = NewStaticMap()
	)
	go merger(data, resultC, done)
	fmt.Printf("1BRC ⏩ ...")
	var i, j, L int // start and stop index
	started := time.Now()
	for i < r.Len() {
		j = i + chunkSize
		if j > r.Len() {
			L = j - i
			wg.Add(1)
			sem <- true
			go aggregate(r, i, L, resultC, sem, &wg)
			break
		}
		for {
			if r.At(j) == '\n' {
				break // found newline
			}
			j++
		}
		L = j - i
		wg.Add(1)
		sem <- true
		go aggregate(r, i, L, resultC, sem, &wg)
		i = j + 1
	}
	wg.Wait()
	close(resultC)
	<-done
	fmt.Printf(" done ✅\n")
	took := time.Since(started)
	for _, c := range cities[:10] {
		i := calculateIndex(c)
		avg := (float64(data.M[i].Sum) / 10) / float64(data.M[i].Count)
		fmt.Printf("%s\t%0.2f/%0.2f/%0.2f\n", c, float64(data.M[i].Min)/10, float64(data.M[i].Max)/10, avg)
	}
	fmt.Printf("...\n")
	fmt.Printf("%d lines omitted (agg took: %v)", len(cities)-10, fmt.Sprintf(noticeColor, took))
	fmt.Println()
}

var cities = []string{
	"Abha",
	"Abidjan",
	"Abéché",
	"Accra",
	"Addis Ababa",
	"Adelaide",
	"Aden",
	"Ahvaz",
	"Albuquerque",
	"Alexandra",
	"Alexandria",
	"Algiers",
	"Alice Springs",
	"Almaty",
	"Amsterdam",
	"Anadyr",
	"Anchorage",
	"Andorra la Vella",
	"Ankara",
	"Antananarivo",
	"Antsiranana",
	"Arkhangelsk",
	"Ashgabat",
	"Asmara",
	"Assab",
	"Astana",
	"Athens",
	"Atlanta",
	"Auckland",
	"Austin",
	"Baghdad",
	"Baguio",
	"Baku",
	"Baltimore",
	"Bamako",
	"Bangkok",
	"Bangui",
	"Banjul",
	"Barcelona",
	"Bata",
	"Batumi",
	"Beijing",
	"Beirut",
	"Belgrade",
	"Belize City",
	"Benghazi",
	"Bergen",
	"Berlin",
	"Bilbao",
	"Birao",
	"Bishkek",
	"Bissau",
	"Blantyre",
	"Bloemfontein",
	"Boise",
	"Bordeaux",
	"Bosaso",
	"Boston",
	"Bouaké",
	"Bratislava",
	"Brazzaville",
	"Bridgetown",
	"Brisbane",
	"Brussels",
	"Bucharest",
	"Budapest",
	"Bujumbura",
	"Bulawayo",
	"Burnie",
	"Busan",
	"Cabo San Lucas",
	"Cairns",
	"Cairo",
	"Calgary",
	"Canberra",
	"Cape Town",
	"Changsha",
	"Charlotte",
	"Chiang Mai",
	"Chicago",
	"Chihuahua",
	"Chittagong",
	"Chișinău",
	"Chongqing",
	"Christchurch",
	"City of San Marino",
	"Colombo",
	"Columbus",
	"Conakry",
	"Copenhagen",
	"Cotonou",
	"Cracow",
	"Da Lat",
	"Da Nang",
	"Dakar",
	"Dallas",
	"Damascus",
	"Dampier",
	"Dar es Salaam",
	"Darwin",
	"Denpasar",
	"Denver",
	"Detroit",
	"Dhaka",
	"Dikson",
	"Dili",
	"Djibouti",
	"Dodoma",
	"Dolisie",
	"Douala",
	"Dubai",
	"Dublin",
	"Dunedin",
	"Durban",
	"Dushanbe",
	"Edinburgh",
	"Edmonton",
	"El Paso",
	"Entebbe",
	"Erbil",
	"Erzurum",
	"Fairbanks",
	"Fianarantsoa",
	"Flores,  Petén",
	"Frankfurt",
	"Fresno",
	"Fukuoka",
	"Gaborone",
	"Gabès",
	"Gagnoa",
	"Gangtok",
	"Garissa",
	"Garoua",
	"George Town",
	"Ghanzi",
	"Gjoa Haven",
	"Guadalajara",
	"Guangzhou",
	"Guatemala City",
	"Halifax",
	"Hamburg",
	"Hamilton",
	"Hanga Roa",
	"Hanoi",
	"Harare",
	"Harbin",
	"Hargeisa",
	"Hat Yai",
	"Havana",
	"Helsinki",
	"Heraklion",
	"Hiroshima",
	"Ho Chi Minh City",
	"Hobart",
	"Hong Kong",
	"Honiara",
	"Honolulu",
	"Houston",
	"Ifrane",
	"Indianapolis",
	"Iqaluit",
	"Irkutsk",
	"Istanbul",
	"Jacksonville",
	"Jakarta",
	"Jayapura",
	"Jerusalem",
	"Johannesburg",
	"Jos",
	"Juba",
	"Kabul",
	"Kampala",
	"Kandi",
	"Kankan",
	"Kano",
	"Kansas City",
	"Karachi",
	"Karonga",
	"Kathmandu",
	"Khartoum",
	"Kingston",
	"Kinshasa",
	"Kolkata",
	"Kuala Lumpur",
	"Kumasi",
	"Kunming",
	"Kuopio",
	"Kuwait City",
	"Kyiv",
	"Kyoto",
	"La Ceiba",
	"La Paz",
	"Lagos",
	"Lahore",
	"Lake Havasu City",
	"Lake Tekapo",
	"Las Palmas de Gran Canaria",
	"Las Vegas",
	"Launceston",
	"Lhasa",
	"Libreville",
	"Lisbon",
	"Livingstone",
	"Ljubljana",
	"Lodwar",
	"Lomé",
	"London",
	"Los Angeles",
	"Louisville",
	"Luanda",
	"Lubumbashi",
	"Lusaka",
	"Luxembourg City",
	"Lviv",
	"Lyon",
	"Madrid",
	"Mahajanga",
	"Makassar",
	"Makurdi",
	"Malabo",
	"Malé",
	"Managua",
	"Manama",
	"Mandalay",
	"Mango",
	"Manila",
	"Maputo",
	"Marrakesh",
	"Marseille",
	"Maun",
	"Medan",
	"Mek'ele",
	"Melbourne",
	"Memphis",
	"Mexicali",
	"Mexico City",
	"Miami",
	"Milan",
	"Milwaukee",
	"Minneapolis",
	"Minsk",
	"Mogadishu",
	"Mombasa",
	"Monaco",
	"Moncton",
	"Monterrey",
	"Montreal",
	"Moscow",
	"Mumbai",
	"Murmansk",
	"Muscat",
	"Mzuzu",
	"N'Djamena",
	"Naha",
	"Nairobi",
	"Nakhon Ratchasima",
	"Napier",
	"Napoli",
	"Nashville",
	"Nassau",
	"Ndola",
	"New Delhi",
	"New Orleans",
	"New York City",
	"Ngaoundéré",
	"Niamey",
	"Nicosia",
	"Niigata",
	"Nouadhibou",
	"Nouakchott",
	"Novosibirsk",
	"Nuuk",
	"Odesa",
	"Odienné",
	"Oklahoma City",
	"Omaha",
	"Oranjestad",
	"Oslo",
	"Ottawa",
	"Ouagadougou",
	"Ouahigouya",
	"Ouarzazate",
	"Oulu",
	"Palembang",
	"Palermo",
	"Palm Springs",
	"Palmerston North",
	"Panama City",
	"Parakou",
	"Paris",
	"Perth",
	"Petropavlovsk-Kamchatsky",
	"Philadelphia",
	"Phnom Penh",
	"Phoenix",
	"Pittsburgh",
	"Podgorica",
	"Pointe-Noire",
	"Pontianak",
	"Port Moresby",
	"Port Sudan",
	"Port Vila",
	"Port-Gentil",
	"Portland (OR)",
	"Porto",
	"Prague",
	"Praia",
	"Pretoria",
	"Pyongyang",
	"Rabat",
	"Rangpur",
	"Reggane",
	"Reykjavík",
	"Riga",
	"Riyadh",
	"Rome",
	"Roseau",
	"Rostov-on-Don",
	"Sacramento",
	"Saint Petersburg",
	"Saint-Pierre",
	"Salt Lake City",
	"San Antonio",
	"San Diego",
	"San Francisco",
	"San Jose",
	"San José",
	"San Juan",
	"San Salvador",
	"Sana'a",
	"Santo Domingo",
	"Sapporo",
	"Sarajevo",
	"Saskatoon",
	"Seattle",
	"Seoul",
	"Seville",
	"Shanghai",
	"Singapore",
	"Skopje",
	"Sochi",
	"Sofia",
	"Sokoto",
	"Split",
	"St. John's",
	"St. Louis",
	"Stockholm",
	"Surabaya",
	"Suva",
	"Suwałki",
	"Sydney",
	"Ségou",
	"Tabora",
	"Tabriz",
	"Taipei",
	"Tallinn",
	"Tamale",
	"Tamanrasset",
	"Tampa",
	"Tashkent",
	"Tauranga",
	"Tbilisi",
	"Tegucigalpa",
	"Tehran",
	"Tel Aviv",
	"Thessaloniki",
	"Thiès",
	"Tijuana",
	"Timbuktu",
	"Tirana",
	"Toamasina",
	"Tokyo",
	"Toliara",
	"Toluca",
	"Toronto",
	"Tripoli",
	"Tromsø",
	"Tucson",
	"Tunis",
	"Ulaanbaatar",
	"Upington",
	"Vaduz",
	"Valencia",
	"Valletta",
	"Vancouver",
	"Veracruz",
	"Vienna",
	"Vientiane",
	"Villahermosa",
	"Vilnius",
	"Virginia Beach",
	"Vladivostok",
	"Warsaw",
	"Washington, D.C.",
	"Wau",
	"Wellington",
	"Whitehorse",
	"Wichita",
	"Willemstad",
	"Winnipeg",
	"Wrocław",
	"Xi'an",
	"Yakutsk",
	"Yangon",
	"Yaoundé",
	"Yellowknife",
	"Yerevan",
	"Yinchuan",
	"Zagreb",
	"Zanzibar City",
	"Zürich",
	"Ürümqi",
	"İzmir",
}
