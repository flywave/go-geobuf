package geobuf

import (
	"io/ioutil"
	"testing"

	"github.com/flywave/go-geobuf/io"
	"github.com/flywave/go-geom/general"
)

var feature_ss = `{"type":"Feature","geometry":{"type":"LineString","coordinates":[[-80.214562,39.722209],[-80.214657,39.722396],[-80.214843,39.723198],[-80.214837,39.724068],[-80.214739,39.724274],[-80.214631,39.725024],[-80.21468,39.725342],[-80.21488,39.725888],[-80.215006,39.726014],[-80.215262,39.726269],[-80.215999,39.726711],[-80.216555,39.727157],[-80.21686,39.727487],[-80.217047,39.727899],[-80.217151,39.72824],[-80.217104,39.72881],[-80.216884,39.729378],[-80.21635,39.730418],[-80.216273,39.730568],[-80.216228,39.730975],[-80.216338,39.731524],[-80.216544,39.732006],[-80.21691,39.732364],[-80.217309,39.732628],[-80.219144,39.733516],[-80.219498,39.733732],[-80.220198,39.734404],[-80.220879,39.735298],[-80.221058,39.735533],[-80.221408,39.736025],[-80.221463,39.736107],[-80.221734,39.736573],[-80.221738,39.736591],[-80.221864,39.73707],[-80.221867,39.737267],[-80.221892,39.738659],[-80.221936,39.739471],[-80.222007,39.739985],[-80.222066,39.740175],[-80.223055,39.74256],[-80.223267,39.743234],[-80.223395,39.743875],[-80.223437,39.744605],[-80.223422,39.744884],[-80.223416,39.744998],[-80.223346,39.746255],[-80.223172,39.746815],[-80.223023,39.747097],[-80.222815,39.747356],[-80.222444,39.747666],[-80.221809,39.748196],[-80.221703,39.748276],[-80.221598,39.748357],[-80.220978,39.748866],[-80.220364,39.74937],[-80.219497,39.750167],[-80.219348,39.750323],[-80.219098,39.750644],[-80.218902,39.750979],[-80.218747,39.751573],[-80.218751,39.752086],[-80.218796,39.752328],[-80.218818,39.752444],[-80.21901,39.75287],[-80.219317,39.753294],[-80.219808,39.753728],[-80.219926,39.753821],[-80.220321,39.75411],[-80.221006,39.754664],[-80.221156,39.754786],[-80.221495,39.755132],[-80.222028,39.755874],[-80.222291,39.756361],[-80.222725,39.757647],[-80.223267,39.759193],[-80.223352,39.759354],[-80.223394,39.759489],[-80.223511,39.759859],[-80.22355,39.76015],[-80.22359,39.76044],[-80.22352,39.761136],[-80.223342,39.761705],[-80.222967,39.762512],[-80.222778,39.762953],[-80.222481,39.764355],[-80.222485,39.764733],[-80.22259,39.765106],[-80.2227,39.765312],[-80.223593,39.766273],[-80.224527,39.767457],[-80.22491,39.768156],[-80.225213,39.768901],[-80.225204,39.769409],[-80.225137,39.769654],[-80.224937,39.769984],[-80.223607,39.771721],[-80.223191,39.772308],[-80.222733,39.772971],[-80.221743,39.774451],[-80.220086,39.776907],[-80.219642,39.777601],[-80.21911,39.778224],[-80.218956,39.778391],[-80.218015,39.779173],[-80.217186,39.779865],[-80.216548,39.780397],[-80.216464,39.780467],[-80.216136,39.780767],[-80.215832,39.780987],[-80.21567,39.781123],[-80.214983,39.781791],[-80.214563,39.782303],[-80.214141,39.783066],[-80.21387,39.78364],[-80.213363,39.784822],[-80.212761,39.786174],[-80.212219,39.787266],[-80.211468,39.788353],[-80.211335,39.788539],[-80.210167,39.790171],[-80.209946,39.790537],[-80.20952,39.791243],[-80.20895,39.792516],[-80.208048,39.7945],[-80.207743,39.795277],[-80.207555,39.796097],[-80.207568,39.796672],[-80.207593,39.796852],[-80.207605,39.796942],[-80.208259,39.800099],[-80.208261,39.800735],[-80.208184,39.801196],[-80.20723,39.804022],[-80.207086,39.804714],[-80.207066,39.805311],[-80.207247,39.805925],[-80.20739,39.806172],[-80.20806,39.807122],[-80.20918,39.808585],[-80.209402,39.808925],[-80.209624,39.809264],[-80.21012,39.809882],[-80.210301,39.810107],[-80.21041,39.810243],[-80.210582,39.810484],[-80.210659,39.810781],[-80.210772,39.811075],[-80.21099,39.811305],[-80.211285,39.811414],[-80.211836,39.811511],[-80.212159,39.811656],[-80.212371,39.811907],[-80.212452,39.812275],[-80.212571,39.813316],[-80.212653,39.813726],[-80.212736,39.813959],[-80.21323,39.814736],[-80.213304,39.814896],[-80.213388,39.815234],[-80.213302,39.815783],[-80.213376,39.816146],[-80.213664,39.817035],[-80.214152,39.817777],[-80.21437,39.818028],[-80.214425,39.818158],[-80.214526,39.818396],[-80.214489,39.81863],[-80.214467,39.818772],[-80.214176,39.819936],[-80.214149,39.820328],[-80.214373,39.821069],[-80.214355,39.821335],[-80.214045,39.822109],[-80.213723,39.82295],[-80.213392,39.824185],[-80.213124,39.8247],[-80.212926,39.824967],[-80.212438,39.825468],[-80.212266,39.825644],[-80.212221,39.825687],[-80.212153,39.825751],[-80.211353,39.826509],[-80.210875,39.826809],[-80.210464,39.82702],[-80.208427,39.827905],[-80.207891,39.828147],[-80.20699,39.828553],[-80.206563,39.828962],[-80.206328,39.829531],[-80.206108,39.830249],[-80.20606,39.830872],[-80.206132,39.831358],[-80.206145,39.831446],[-80.206195,39.832956],[-80.206223,39.833507],[-80.206239,39.833807],[-80.206377,39.834253],[-80.206412,39.834364],[-80.206458,39.836312],[-80.206469,39.837273],[-80.206386,39.837619],[-80.206259,39.837873],[-80.205936,39.838365],[-80.205392,39.839167],[-80.203674,39.84132],[-80.202788,39.842544],[-80.20262,39.842872],[-80.202026,39.844199],[-80.201819,39.844817],[-80.201693,39.845137],[-80.201687,39.845625],[-80.201742,39.846001],[-80.201857,39.846484],[-80.201789,39.847015],[-80.201525,39.848364],[-80.200512,39.85078],[-80.200453,39.851096],[-80.200384,39.852495],[-80.200142,39.85414],[-80.200134,39.854197],[-80.199986,39.855744],[-80.199926,39.856372],[-80.19995,39.857389],[-80.199847,39.857831],[-80.199707,39.858073],[-80.19942,39.858914],[-80.199467,39.859455],[-80.199614,39.859791],[-80.200183,39.860703],[-80.200377,39.861283],[-80.200469,39.861713],[-80.200875,39.863604],[-80.200747,39.864206],[-80.200357,39.865205],[-80.200308,39.865331],[-80.19991,39.866727],[-80.199745,39.867251],[-80.199559,39.867559],[-80.198822,39.869298],[-80.198071,39.870487],[-80.198,39.870738],[-80.197886,39.871508],[-80.197701,39.873858],[-80.197669,39.874267],[-80.197665,39.874316],[-80.19719,39.875096],[-80.196755,39.875549],[-80.196555,39.875861],[-80.196484,39.876123],[-80.196331,39.876518],[-80.196311,39.876569],[-80.196108,39.876929],[-80.195605,39.877616],[-80.195315,39.878558],[-80.195232,39.87881],[-80.194979,39.87958],[-80.194681,39.880488],[-80.194504,39.880781],[-80.193921,39.881422],[-80.192096,39.88346],[-80.191504,39.883879],[-80.191452,39.883916],[-80.190092,39.884847],[-80.189823,39.885116],[-80.189596,39.8854],[-80.189083,39.886042],[-80.18762,39.887268],[-80.187146,39.887711],[-80.18698,39.888075],[-80.186937,39.88817],[-80.186694,39.889092],[-80.186444,39.889902],[-80.186344,39.890446],[-80.186299,39.890692]]},"properties":{"shit":199}}`
var feature_buf = []byte{0x12, 0xb, 0xa, 0x4, 0x73, 0x68, 0x69, 0x74, 0x12, 0x3, 0x30, 0x8e, 0x3, 0x18, 0x2, 0x22, 0x96, 0x9, 0xa7, 0x95, 0xfe, 0xfc, 0x5, 0x94, 0x83, 0xe9, 0xfa, 0x2, 0xeb, 0xe, 0x9c, 0x1d, 0x87, 0x1d, 0xa6, 0x7d, 0x78, 0xfa, 0x87, 0x1, 0xa8, 0xf, 0x98, 0x20, 0xf0, 0x10, 0x98, 0x75, 0xd3, 0x7, 0xd8, 0x31, 0x9d, 0x1f, 0xa8, 0x55, 0xd9, 0x13, 0xd8, 0x13, 0xff, 0x27, 0xec, 0x27, 0x93, 0x73, 0x88, 0x45, 0xef, 0x56, 0xd8, 0x45, 0xd3, 0x2f, 0xc8, 0x33, 0x99, 0x1d, 0xb0, 0x40, 0xa1, 0x10, 0xa4, 0x35, 0xac, 0x7, 0x88, 0x59, 0xb2, 0x22, 0xde, 0x58, 0xb6, 0x53, 0xc2, 0xa2, 0x1, 0x84, 0xc, 0xb8, 0x17, 0x84, 0x7, 0xcc, 0x3f, 0x95, 0x11, 0xe4, 0x55, 0x99, 0x20, 0xa8, 0x4b, 0x97, 0x39, 0xf6, 0x37, 0xab, 0x3e, 0xa2, 0x29, 0xdb, 0x9e, 0x2, 0xe0, 0x8a, 0x1, 0xa7, 0x37, 0xe0, 0x21, 0xaf, 0x6d, 0x80, 0x69, 0xb3, 0x6a, 0xd8, 0x8b, 0x1, 0xfb, 0x1b, 0xda, 0x24, 0xd7, 0x36, 0xf2, 0x4c, 0xcb, 0x8, 0xe6, 0xc, 0xab, 0x2a, 0xea, 0x48, 0x4f, 0xe8, 0x2, 0xd7, 0x13, 0xec, 0x4a, 0x3b, 0xe4, 0x1e, 0xf3, 0x3, 0xc0, 0xd9, 0x1, 0xef, 0x6, 0xf0, 0x7e, 0x8b, 0xb, 0xa8, 0x50, 0x9b, 0x9, 0xd8, 0x1d, 0xc3, 0x9a, 0x1, 0xd4, 0xf4, 0x2, 0x8f, 0x21, 0xa8, 0x69, 0xff, 0x13, 0x94, 0x64, 0xc7, 0x6, 0x88, 0x72, 0xac, 0x2, 0xcc, 0x2b, 0x78, 0xe8, 0x11, 0xf8, 0xa, 0xb4, 0xc4, 0x1, 0x98, 0x1b, 0xc0, 0x57, 0xa4, 0x17, 0x86, 0x2c, 0xc0, 0x20, 0xbe, 0x28, 0xfc, 0x39, 0xb8, 0x30, 0x9e, 0x63, 0xe8, 0x52, 0xc6, 0x10, 0xc0, 0xc, 0xb4, 0x10, 0xd4, 0xc, 0xf0, 0x60, 0xc4, 0x4f, 0xf8, 0x5f, 0xe0, 0x4e, 0xbc, 0x87, 0x1, 0xc4, 0x7c, 0xa4, 0x17, 0xb0, 0x18, 0x88, 0x27, 0x94, 0x32, 0xd0, 0x1e, 0xac, 0x34, 0x9e, 0x18, 0xe8, 0x5c, 0x51, 0x94, 0x50, 0x83, 0x7, 0xe8, 0x25, 0xb7, 0x3, 0x8e, 0x12, 0xff, 0x1d, 0xca, 0x42, 0xfb, 0x2f, 0x9e, 0x42, 0xdb, 0x4c, 0xea, 0x43, 0xb7, 0x12, 0xc4, 0xe, 0xdb, 0x3d, 0x94, 0x2d, 0x83, 0x6b, 0xc8, 0x56, 0xb5, 0x17, 0x88, 0x13, 0xfd, 0x34, 0x88, 0x36, 0xa3, 0x53, 0xf8, 0x73, 0x8b, 0x29, 0x8c, 0x4c, 0xe7, 0x43, 0xf8, 0xc8, 0x1, 0xd7, 0x54, 0xc8, 0xf1, 0x1, 0xa3, 0xd, 0x94, 0x19, 0xc7, 0x6, 0x8c, 0x15, 0xa3, 0x12, 0xe8, 0x39, 0x8b, 0x6, 0xbc, 0x2d, 0x9f, 0x6, 0xa8, 0x2d, 0xfa, 0xa, 0xe0, 0x6c, 0xe6, 0x1b, 0xf4, 0x58, 0xcc, 0x3a, 0x8c, 0x7e, 0xc4, 0x1d, 0xf4, 0x44, 0xb4, 0x2e, 0x88, 0xdb, 0x1, 0x4f, 0x88, 0x3b, 0xb3, 0x10, 0xa4, 0x3a, 0x97, 0x11, 0x98, 0x20, 0xc1, 0x8b, 0x1, 0x94, 0x96, 0x1, 0xf9, 0x91, 0x1, 0x80, 0xb9, 0x1, 0xeb, 0x3b, 0x9c, 0x6d, 0xab, 0x2f, 0xb4, 0x74, 0xb4, 0x1, 0xb0, 0x4f, 0xbc, 0xa, 0xa4, 0x26, 0xa0, 0x1f, 0xc8, 0x33, 0xe8, 0xcf, 0x1, 0xb4, 0x8f, 0x2, 0x80, 0x41, 0xdc, 0x5b, 0xc8, 0x47, 0xcc, 0x67, 0xd8, 0x9a, 0x1, 0xa0, 0xe7, 0x1, 0xf4, 0x82, 0x2, 0xe0, 0xff, 0x2, 0xb2, 0x45, 0xb8, 0x6c, 0x8e, 0x53, 0xac, 0x61, 0x88, 0x18, 0x8c, 0x1a, 0x84, 0x93, 0x1, 0x98, 0x7a, 0xc4, 0x81, 0x1, 0x90, 0x6c, 0xd8, 0x63, 0x90, 0x53, 0x90, 0xd, 0xf8, 0xa, 0xa0, 0x33, 0xf0, 0x2e, 0xc0, 0x2f, 0xb0, 0x22, 0xa8, 0x19, 0xa0, 0x15, 0xac, 0x6b, 0xb0, 0x68, 0xd0, 0x41, 0x80, 0x50, 0xf8, 0x41, 0x9c, 0x77, 0xac, 0x2a, 0xd8, 0x59, 0x9c, 0x4f, 0xd8, 0xb8, 0x1, 0x88, 0x5e, 0xa0, 0xd3, 0x1, 0xd8, 0x54, 0xd0, 0xaa, 0x1, 0xac, 0x75, 0xec, 0xa9, 0x1, 0xe4, 0x14, 0x88, 0x1d, 0xc0, 0xb6, 0x1, 0x80, 0xff, 0x1, 0xc4, 0x22, 0x98, 0x39, 0xc8, 0x42, 0xa8, 0x6e, 0x88, 0x59, 0xf4, 0xc6, 0x1, 0xf8, 0x8c, 0x1, 0x80, 0xb6, 0x2, 0xd6, 0x2f, 0xb4, 0x79, 0xae, 0x1d, 0x90, 0x80, 0x1, 0x83, 0x2, 0xec, 0x59, 0xf3, 0x3, 0x90, 0x1c, 0xef, 0x1, 0x88, 0xe, 0x97, 0x66, 0xa4, 0xed, 0x3, 0x25, 0xb0, 0x63, 0x82, 0xc, 0x84, 0x48, 0x88, 0x95, 0x1, 0xc8, 0xb9, 0x3, 0xc0, 0x16, 0x90, 0x6c, 0x90, 0x3, 0xa4, 0x5d, 0xa3, 0x1c, 0xf8, 0x5f, 0xab, 0x16, 0xca, 0x26, 0xd7, 0x68, 0xba, 0x94, 0x1, 0xff, 0xae, 0x1, 0xcc, 0xe4, 0x1, 0xd7, 0x22, 0x90, 0x35, 0xd7, 0x22, 0xfc, 0x34, 0xbf, 0x4d, 0xc8, 0x60, 0xa3, 0x1c, 0x94, 0x23, 0x83, 0x11, 0xa0, 0x15, 0xef, 0x1a, 0xd4, 0x25, 0x83, 0xc, 0xb4, 0x2e, 0xd3, 0x11, 0xf8, 0x2d, 0x87, 0x22, 0xf8, 0x23, 0x8b, 0x2e, 0x84, 0x11, 0x8b, 0x56, 0x94, 0xf, 0xbb, 0x32, 0xd4, 0x16, 0x8f, 0x21, 0x9c, 0x27, 0xd3, 0xc, 0xc0, 0x39, 0xcb, 0x12, 0xd4, 0xa2, 0x1, 0xe7, 0xc, 0x88, 0x40, 0xfb, 0xc, 0xb2, 0x24, 0x97, 0x4d, 0xb6, 0x79, 0xc5, 0xb, 0x80, 0x19, 0x91, 0xd, 0xe6, 0x34, 0xb8, 0xd, 0xe6, 0x55, 0xc7, 0xb, 0xdc, 0x38, 0xff, 0x2c, 0xf4, 0x8a, 0x1, 0x9f, 0x4c, 0xf8, 0x73, 0x87, 0x22, 0x9c, 0x27, 0xcb, 0x8, 0xa6, 0x14, 0xe3, 0xf, 0x9a, 0x25, 0xe4, 0x5, 0xc8, 0x24, 0xb8, 0x3, 0x98, 0x16, 0xbc, 0x2d, 0xf0, 0xb5, 0x1, 0x9c, 0x4, 0xa0, 0x3d, 0xff, 0x22, 0xe4, 0x73, 0xe8, 0x2, 0xc8, 0x29, 0xb8, 0x30, 0xf8, 0x78, 0xa8, 0x32, 0xb4, 0x83, 0x1, 0xdc, 0x33, 0xfc, 0xc0, 0x1, 0xf2, 0x29, 0xbc, 0x50, 0xf6, 0x1e, 0xdc, 0x29, 0xa0, 0x4c, 0xa4, 0x4e, 0xf0, 0x1a, 0xbe, 0x1b, 0x84, 0x7, 0xde, 0x6, 0xd0, 0xa, 0xfe, 0x9, 0x80, 0x7d, 0xba, 0x76, 0xd8, 0x4a, 0xf0, 0x2e, 0x9c, 0x40, 0xfc, 0x20, 0xa4, 0xbe, 0x2, 0xa4, 0x8a, 0x1, 0xe0, 0x53, 0xe8, 0x25, 0xe4, 0x8c, 0x1, 0xb8, 0x3f, 0xdc, 0x42, 0xf4, 0x3f, 0xdc, 0x24, 0xf4, 0x58, 0xb0, 0x22, 0x98, 0x70, 0xc2, 0x7, 0xac, 0x61, 0xa1, 0xb, 0xf8, 0x4b, 0x83, 0x2, 0xe0, 0xd, 0xe5, 0x7, 0xf8, 0xeb, 0x1, 0xb1, 0x4, 0x8c, 0x56, 0xbf, 0x2, 0xf0, 0x2e, 0xc7, 0x15, 0xd6, 0x45, 0xbb, 0x5, 0xae, 0x11, 0x97, 0x7, 0xb0, 0xb0, 0x2, 0xdb, 0x1, 0x94, 0x96, 0x1, 0xfc, 0xc, 0x86, 0x36, 0xec, 0x13, 0xda, 0x27, 0xbc, 0x32, 0xf0, 0x4c, 0x80, 0x55, 0xa8, 0x7d, 0xb8, 0x8c, 0x2, 0xb4, 0xd0, 0x2, 0xb8, 0x8a, 0x1, 0x9e, 0xbf, 0x1, 0xa0, 0x1a, 0xa2, 0x33, 0xe8, 0x5c, 0xac, 0xcf, 0x1, 0xac, 0x20, 0xc8, 0x60, 0xd8, 0x13, 0x80, 0x32, 0x78, 0xa0, 0x4c, 0xcb, 0x8, 0xe0, 0x3a, 0xfb, 0x11, 0xba, 0x4b, 0xd0, 0xa, 0xfe, 0x52, 0xa0, 0x29, 0xe2, 0xd2, 0x1, 0xa4, 0x9e, 0x1, 0xc2, 0xf9, 0x2, 0x9c, 0x9, 0xb0, 0x31, 0xe4, 0xa, 0xcc, 0xda, 0x1, 0xe8, 0x25, 0x84, 0x81, 0x2, 0xa0, 0x1, 0xf4, 0x8, 0x90, 0x17, 0xdc, 0xf1, 0x1, 0xb0, 0x9, 0x90, 0x62, 0xdf, 0x3, 0xf4, 0x9e, 0x1, 0x8c, 0x10, 0x88, 0x45, 0xf0, 0x15, 0xe8, 0x25, 0xec, 0x2c, 0xb4, 0x83, 0x1, 0xab, 0x7, 0xc2, 0x54, 0xfb, 0x16, 0xc2, 0x34, 0xf3, 0x58, 0xc0, 0x8e, 0x1, 0xa7, 0x1e, 0xd0, 0x5a, 0xaf, 0xe, 0x98, 0x43, 0xb7, 0x3f, 0xbc, 0xa7, 0x2, 0x80, 0x14, 0x88, 0x5e, 0xf8, 0x3c, 0x8c, 0x9c, 0x1, 0xd4, 0x7, 0xd8, 0x13, 0x98, 0x3e, 0x90, 0xda, 0x1, 0xe6, 0x19, 0xf0, 0x51, 0x88, 0x1d, 0x90, 0x30, 0x92, 0x73, 0xdc, 0x8f, 0x2, 0xac, 0x75, 0xe4, 0xb9, 0x1, 0x8e, 0xb, 0x9c, 0x27, 0xe6, 0x11, 0xa8, 0x78, 0xf4, 0x1c, 0x98, 0xef, 0x2, 0x80, 0x5, 0xf4, 0x3f, 0x50, 0xd4, 0x7, 0x9c, 0x4a, 0xf0, 0x79, 0xfc, 0x43, 0xe4, 0x46, 0xa0, 0x1f, 0xe0, 0x30, 0x8c, 0xb, 0xf8, 0x28, 0xf4, 0x17, 0xdc, 0x3d, 0x90, 0x3, 0xfc, 0x7, 0xdc, 0x1f, 0x9e, 0x38, 0xcc, 0x4e, 0xae, 0x6b, 0xaa, 0x2d, 0x98, 0x93, 0x1, 0xfa, 0xc, 0xb0, 0x27, 0xc4, 0x27, 0xa8, 0x78, 0xc8, 0x2e, 0xf0, 0x8d, 0x1, 0xd4, 0x1b, 0xe4, 0x2d, 0x8c, 0x5b, 0x94, 0x64, 0x94, 0x9d, 0x2, 0xb8, 0xbe, 0x2, 0xc0, 0x5c, 0xbc, 0x41, 0x90, 0x8, 0xe4, 0x5, 0xc0, 0xd4, 0x1, 0xbc, 0x91, 0x1, 0x84, 0x2a, 0x82, 0x2a, 0xbc, 0x23, 0xb2, 0x2c, 0x94, 0x50, 0xa8, 0x64, 0xcc, 0xe4, 0x1, 0xc8, 0xbf, 0x1, 0x88, 0x4a, 0x9c, 0x45, 0xf8, 0x19, 0xf0, 0x38, 0xdc, 0x6, 0xec, 0xe, 0xfc, 0x25, 0x88, 0x90, 0x1, 0x88, 0x27, 0xc8, 0x7e, 0xd0, 0xf, 0x80, 0x55, 0x84, 0x7, 0xb8, 0x26, 0x2a, 0x1e, 0x9e, 0xb5, 0xba, 0x81, 0xfd, 0xff, 0xff, 0xff, 0xff, 0x1, 0xca, 0xc1, 0xb4, 0xbd, 0x1, 0xb2, 0x95, 0xd2, 0x81, 0xfd, 0xff, 0xff, 0xff, 0xff, 0x1, 0xa8, 0xac, 0x9b, 0xbe, 0x1}

func BenchmarkReadFeatureCollection(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		bytevals, _ := ioutil.ReadFile("test_data/county.geojson")
		general.UnmarshalFeatureCollection(bytevals)
	}
}

func BenchmarkReadFeatureCollectionGeobuf(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		geobuf := ReaderFile("test_data/county.geobuf")
		for geobuf.Next() {

			geobuf.Feature()
		}
	}
}

func Benchmark_Read_Feature_Benchmark_Old(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		general.UnmarshalFeature([]byte(feature_ss))
	}
}

func Benchmark_Read_Feature_Benchmark_New(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		io.ReadFeature(feature_buf)
	}
}
