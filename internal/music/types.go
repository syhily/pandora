package music

// NeteaseMusic represents the music metadata structure
type NeteaseMusic struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Pic    string `json:"pic"`
	Lyric  string `json:"lyric"`
	URL    string `json:"url"`
}

// SongUrl represents the response from the song URL API
type SongUrl struct {
	Data []struct {
		ID            int         `json:"id"`
		URL           string      `json:"url"`
		Br            int         `json:"br"`
		Size          int         `json:"size"`
		Md5           string      `json:"md5"`
		Code          int         `json:"code"`
		Expi          int         `json:"expi"`
		Type          string      `json:"type"`
		Gain          float64     `json:"gain"`
		Peak          float64     `json:"peak"`
		ClosedGain    float64     `json:"closedGain"`
		ClosedPeak    float64     `json:"closedPeak"`
		Fee           int         `json:"fee"`
		Uf            interface{} `json:"uf"`
		Payed         int         `json:"payed"`
		Flag          int         `json:"flag"`
		CanExtend     bool        `json:"canExtend"`
		FreeTrialInfo struct {
			FragmentType int `json:"fragmentType"`
			Start        int `json:"start"`
			End          int `json:"end"`
			AlgData      struct {
				FragSource string `json:"fragSource"`
			} `json:"algData"`
		} `json:"freeTrialInfo"`
		Level              string      `json:"level"`
		EncodeType         string      `json:"encodeType"`
		ChannelLayout      interface{} `json:"channelLayout"`
		FreeTrialPrivilege struct {
			ResConsumable      bool        `json:"resConsumable"`
			UserConsumable     bool        `json:"userConsumable"`
			ListenType         interface{} `json:"listenType"`
			CannotListenReason interface{} `json:"cannotListenReason"`
			PlayReason         interface{} `json:"playReason"`
			FreeLimitTagType   interface{} `json:"freeLimitTagType"`
		} `json:"freeTrialPrivilege"`
		FreeTimeTrialPrivilege struct {
			ResConsumable  bool `json:"resConsumable"`
			UserConsumable bool `json:"userConsumable"`
			Type           int  `json:"type"`
			RemainTime     int  `json:"remainTime"`
		} `json:"freeTimeTrialPrivilege"`
		URLSource    int         `json:"urlSource"`
		RightSource  int         `json:"rightSource"`
		PodcastCtrp  interface{} `json:"podcastCtrp"`
		EffectTypes  interface{} `json:"effectTypes"`
		Time         int         `json:"time"`
		Message      interface{} `json:"message"`
		LevelConfuse interface{} `json:"levelConfuse"`
		MusicID      string      `json:"musicId"`
		Accompany    interface{} `json:"accompany"`
		Sr           int         `json:"sr"`
		AuEff        int         `json:"auEff"`
		ImmerseType  interface{} `json:"immerseType"`
	} `json:"data"`
	Code int `json:"code"`
}

// SongDetail represents the response from the song detail API
type SongDetail struct {
	Songs []struct {
		Name            string      `json:"name"`
		MainTitle       interface{} `json:"mainTitle"`
		AdditionalTitle interface{} `json:"additionalTitle"`
		ID              int         `json:"id"`
		Pst             int         `json:"pst"`
		T               int         `json:"t"`
		Ar              []struct {
			ID    int           `json:"id"`
			Name  string        `json:"name"`
			Tns   []interface{} `json:"tns"`
			Alias []interface{} `json:"alias"`
		} `json:"ar"`
		Alia []interface{} `json:"alia"`
		Pop  float64       `json:"pop"`
		St   int           `json:"st"`
		Rt   string        `json:"rt"`
		Fee  int           `json:"fee"`
		V    int           `json:"v"`
		Crbt interface{}   `json:"crbt"`
		Cf   string        `json:"cf"`
		Al   struct {
			ID     int           `json:"id"`
			Name   string        `json:"name"`
			PicURL string        `json:"picUrl"`
			Tns    []interface{} `json:"tns"`
			Pic    int64         `json:"pic"`
		} `json:"al"`
		Dt int `json:"dt"`
		H  struct {
			Br   int     `json:"br"`
			Fid  int     `json:"fid"`
			Size int     `json:"size"`
			Vd   float64 `json:"vd"`
			Sr   int     `json:"sr"`
		} `json:"h"`
		M struct {
			Br   int     `json:"br"`
			Fid  int     `json:"fid"`
			Size int     `json:"size"`
			Vd   float64 `json:"vd"`
			Sr   int     `json:"sr"`
		} `json:"m"`
		L struct {
			Br   int     `json:"br"`
			Fid  int     `json:"fid"`
			Size int     `json:"size"`
			Vd   float64 `json:"vd"`
			Sr   int     `json:"sr"`
		} `json:"l"`
		Sq struct {
			Br   int     `json:"br"`
			Fid  int     `json:"fid"`
			Size int     `json:"size"`
			Vd   float64 `json:"vd"`
			Sr   int     `json:"sr"`
		} `json:"sq"`
		Hr                   interface{}   `json:"hr"`
		A                    interface{}   `json:"a"`
		Cd                   string        `json:"cd"`
		No                   int           `json:"no"`
		RtURL                interface{}   `json:"rtUrl"`
		Ftype                int           `json:"ftype"`
		RtUrls               []interface{} `json:"rtUrls"`
		DjID                 int           `json:"djId"`
		Copyright            int           `json:"copyright"`
		SID                  int           `json:"s_id"`
		Mark                 int64         `json:"mark"`
		OriginCoverType      int           `json:"originCoverType"`
		OriginSongSimpleData interface{}   `json:"originSongSimpleData"`
		TagPicList           interface{}   `json:"tagPicList"`
		ResourceState        bool          `json:"resourceState"`
		Version              int           `json:"version"`
		SongJumpInfo         interface{}   `json:"songJumpInfo"`
		EntertainmentTags    interface{}   `json:"entertainmentTags"`
		AwardTags            interface{}   `json:"awardTags"`
		DisplayTags          interface{}   `json:"displayTags"`
		MarkTags             []interface{} `json:"markTags"`
		Single               int           `json:"single"`
		NoCopyrightRcmd      interface{}   `json:"noCopyrightRcmd"`
		Mv                   int           `json:"mv"`
		Rtype                int           `json:"rtype"`
		Rurl                 interface{}   `json:"rurl"`
		Mst                  int           `json:"mst"`
		Cp                   int           `json:"cp"`
		PublishTime          int64         `json:"publishTime"`
	} `json:"songs"`
	Privileges []struct {
		ID                 int         `json:"id"`
		Fee                int         `json:"fee"`
		Payed              int         `json:"payed"`
		St                 int         `json:"st"`
		Pl                 int         `json:"pl"`
		Dl                 int         `json:"dl"`
		Sp                 int         `json:"sp"`
		Cp                 int         `json:"cp"`
		Subp               int         `json:"subp"`
		Cs                 bool        `json:"cs"`
		Maxbr              int         `json:"maxbr"`
		Fl                 int         `json:"fl"`
		Toast              bool        `json:"toast"`
		Flag               int         `json:"flag"`
		PreSell            bool        `json:"preSell"`
		PlayMaxbr          int         `json:"playMaxbr"`
		DownloadMaxbr      int         `json:"downloadMaxbr"`
		MaxBrLevel         string      `json:"maxBrLevel"`
		PlayMaxBrLevel     string      `json:"playMaxBrLevel"`
		DownloadMaxBrLevel string      `json:"downloadMaxBrLevel"`
		PlLevel            string      `json:"plLevel"`
		DlLevel            string      `json:"dlLevel"`
		FlLevel            string      `json:"flLevel"`
		Rscl               interface{} `json:"rscl"`
		FreeTrialPrivilege struct {
			ResConsumable      bool        `json:"resConsumable"`
			UserConsumable     bool        `json:"userConsumable"`
			ListenType         interface{} `json:"listenType"`
			CannotListenReason interface{} `json:"cannotListenReason"`
			PlayReason         interface{} `json:"playReason"`
			FreeLimitTagType   interface{} `json:"freeLimitTagType"`
		} `json:"freeTrialPrivilege"`
		RightSource    int `json:"rightSource"`
		ChargeInfoList []struct {
			Rate          int         `json:"rate"`
			ChargeURL     interface{} `json:"chargeUrl"`
			ChargeMessage interface{} `json:"chargeMessage"`
			ChargeType    int         `json:"chargeType"`
		} `json:"chargeInfoList"`
		Code        int         `json:"code"`
		Message     interface{} `json:"message"`
		PlLevels    interface{} `json:"plLevels"`
		DlLevels    interface{} `json:"dlLevels"`
		IgnoreCache interface{} `json:"ignoreCache"`
		Bd          interface{} `json:"bd"`
	} `json:"privileges"`
	Code int `json:"code"`
}

// Lyric represents the response from the lyric API
type Lyric struct {
	Sgc bool `json:"sgc"`
	Sfy bool `json:"sfy"`
	Qfy bool `json:"qfy"`
	Lrc struct {
		Version int    `json:"version"`
		Lyric   string `json:"lyric"`
	} `json:"lrc"`
	Klyric struct {
		Version int    `json:"version"`
		Lyric   string `json:"lyric"`
	} `json:"klyric"`
	Tlyric struct {
		Version int    `json:"version"`
		Lyric   string `json:"lyric"`
	} `json:"tlyric"`
	Code int `json:"code"`
}
