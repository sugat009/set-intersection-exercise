package reader

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const dummyFile = `
key, col1, col2
a, 1, 2
b, 1, 2
`

func Test_ReadKeysFromCsvIntoChannel(t *testing.T) {
	outputChan := make(chan string)

	go func() {
		defer close(outputChan)
		err := ReadKeysFromCsvIntoChannel("key", strings.NewReader(dummyFile), outputChan)
		assert.NoError(t, err)
	}()

	output, more := <-outputChan
	assert.Equal(t, "a", output)
	assert.True(t, more)

	output, more = <-outputChan
	assert.Equal(t, "b", output)
	assert.True(t, more)

	_, more = <-outputChan
	assert.False(t, more)
}

func Test_ReadKeysFromCsvIntoChannel_Empty(t *testing.T) {
	outputChan := make(chan string)

	go func() {
		defer close(outputChan)
		err := ReadKeysFromCsvIntoChannel("key", strings.NewReader(""), outputChan)
		assert.NoError(t, err)
	}()
	_, more := <-outputChan
	assert.False(t, more)
}

func Test_ReadKeysFromCsvIntoChannel_Nil(t *testing.T) {
	outputChan := make(chan string)
	err := ReadKeysFromCsvIntoChannel("key", nil, outputChan)
	assert.Error(t, err)
}

func Test_ReadKeysFromCsvIntoChannel_NonExistentKey(t *testing.T) {
	err := ReadKeysFromCsvIntoChannel("non-existent", strings.NewReader(dummyFile), make(chan string))
	assert.Error(t, err)
}

func Test_ReadKeysFromCsvIntoChannel_EmptyKey(t *testing.T) {
	err := ReadKeysFromCsvIntoChannel("", strings.NewReader(dummyFile), make(chan string))
	assert.Error(t, err)
}

func Benchmark_ReadKeysFromCsvIntoChannel_1(b *testing.B) {
	benchmarkReadKeysFromCsvIntoChannel(1, b)
}

func Benchmark_ReadKeysFromCsvIntoChannel_2(b *testing.B) {
	benchmarkReadKeysFromCsvIntoChannel(2, b)
}

func Benchmark_ReadKeysFromCsvIntoChannel_4(b *testing.B) {
	benchmarkReadKeysFromCsvIntoChannel(4, b)
}

func Benchmark_ReadKeysFromCsvIntoChannel_8(b *testing.B) {
	benchmarkReadKeysFromCsvIntoChannel(8, b)
}

func Benchmark_ReadKeysFromCsvIntoChannel_16(b *testing.B) {
	benchmarkReadKeysFromCsvIntoChannel(16, b)
}

func Benchmark_ReadKeysFromCsvIntoChannel_32(b *testing.B) {
	benchmarkReadKeysFromCsvIntoChannel(32, b)
}

func Benchmark_ReadKeysFromCsvIntoChannel_64(b *testing.B) {
	benchmarkReadKeysFromCsvIntoChannel(64, b)
}

func Benchmark_ReadKeysFromCsvIntoChannel_128(b *testing.B) {
	benchmarkReadKeysFromCsvIntoChannel(128, b)
}

func Benchmark_ReadKeysFromCsvIntoChannel_256(b *testing.B) {
	benchmarkReadKeysFromCsvIntoChannel(256, b)
}

func Benchmark_ReadKeysFromCsvIntoChannel_512(b *testing.B) {
	benchmarkReadKeysFromCsvIntoChannel(512, b)
}

func benchmarkReadKeysFromCsvIntoChannel(buf int, b *testing.B) {
	for i := 0; i < b.N; i++ {

		outputChan := make(chan string, buf)
		go func() {
			for {
				_, more := <-outputChan
				if !more {
					break
				}

				//simulate delay in consumption
				time.Sleep(time.Duration(100 * int(time.Millisecond)))
			}
		}()

		_ = ReadKeysFromCsvIntoChannel("key", strings.NewReader(getinputFile1000()), outputChan)
		close(outputChan)
	}
}

func getinputFile1000() string {
	return `foo,bar,baz
ab0b5b4e,1b,27531
b9927308,c9,25846
fe9314b5,91,76981
9d2a533a,77,104653
85d6b74a,ed,149364
c1b81d87,ff,38840
56db39b2,be,76892
c221158e,e4,201386
b8183bf4,ee,81919
1008a421,05,181524
6e9b108e,83,172076
50dd628d,87,247282
e7631b77,5c,155392
ad439320,13,235523
b7c386ed,8f,468039
0cd80e63,3a,480350
425f7cd7,34,354127
8bec82fe,a3,423078
7d9b9ee3,46,88160
7c2fa2d4,83,397783
01d144cb,9b,11216
26eddff9,14,175096
3a42cef4,1e,103426
2a79c3ab,9b,632430
6413d0e4,60,128351
84fb0408,61,341591
e31699ff,6e,114830
1f9d4c66,6d,99832
1679896a,4d,949249
4ed9c989,22,214555
e118cabe,60,521040
90f5305c,b4,879872
76b7cafe,bd,662474
fda74c75,6d,329812
e88836c0,51,731195
d9b61858,82,618474
314ad261,5c,598427
e6cab166,69,1211283
a7c6afcb,60,373822
62fff94b,f8,1011033
f20fbf23,ce,707675
f1feb37a,3b,1059598
93e753b0,9f,563931
aecf9a18,dd,1285400
5d17b38b,eb,417764
dbf6aa05,83,531270
5b318121,c4,1243973
b0473344,24,1445503
dc454a16,82,112001
b4db1e14,a0,1555377
4500b1bb,ab,879026
d568f99b,c1,146633
e83effe3,01,333818
31b57dc5,9a,1173561
12b4ec75,5c,1604412
283eadd9,7f,640221
430c3677,c9,119857
5224073c,29,38054
c58f483e,b0,884880
75f36033,4e,124051
0036a90e,e6,476285
aba72023,8f,1971955
cfb9716b,2c,1862503
c431dac2,72,1784508
926939fa,a7,567976
91afb2a2,dc,1167333
5971d20c,a7,823750
d82609da,07,1694004
9b6d0bac,05,1158838
5d006b91,cd,1531598
73d985cb,59,1236797
8764673c,4a,92673
80e8e9e1,fd,1046857
37b4b419,f1,2259543
619556b7,3f,2287558
88637b5b,d7,1795436
af9b171c,e2,717310
6dfcb3e4,3b,1887622
1aa3e085,a0,1656695
b743cd26,46,928117
e4343cfa,1a,1825727
64c805b0,69,445968
84adbdc7,75,1196971
6a0edef9,5d,2422417
6d276b85,38,2309557
65a44ba7,89,930905
892060d8,d4,652746
c9eafd2d,48,2576114
76ec2afc,6a,1021773
c2958e94,07,2025071
560c4d21,47,2852081
78a52e2d,10,1774547
1026bbc0,1b,2003084
1d96167f,11,2644887
853fe18d,c9,1368332
4d91115b,1d,2906557
48ab0b3f,de,1266429
0f111b32,b1,2616431
d98b8ed5,c5,2219626
ab1c53e3,92,2985073
bd8591d1,41,1596835
9ed91408,fe,721359
c68e0f08,0f,3207200
fbaf009c,1a,3135679
02a2bc32,0d,508044
d05ed344,5f,3060288
c85e95d2,b6,2247741
1fce5ce7,2d,1528658
ab1c6a51,62,2213020
b1c2fc54,76,1013073
0bb2fe38,2b,2858884
e7d63dd8,a9,1128375
4f6760c9,71,1655268
9713ca15,6c,844633
99581ae1,cc,706684
88f18b85,e6,1049992
8fab1155,eb,2133325
91083a59,87,1610453
23f71860,6e,661364
8df93168,99,3565698
8f15954a,cb,409066
ad732671,c7,504010
61e4cafc,53,1996869
b0d63ade,57,3089987
5cf5464d,b9,4033542
5ea1cbb2,54,3860414
a3f6a719,82,2848342
0b2842e6,58,1607208
038b206b,1c,3169335
8fdcb108,66,1570447
874e8a15,bc,1262715
7f01b38a,fc,1004618
cc44e546,ec,2547285
e353c263,c4,1073194
57891d26,be,674123
9bfc6f3f,d8,3262781
fc6930a0,9a,563283
ab90baf7,d6,3588060
a2a196e3,12,747443
0fd2d780,46,3418029
900675b6,71,344349
0ab64cae,c0,4420887
49f3b2d5,a4,246142
6723ac47,0f,2461044
7087a289,64,837240
e269fe69,8b,1148489
b8889963,ce,3842910
46e0924a,cb,3553127
130f55c0,3d,3205627
5d428488,12,4754989
0e85967f,0f,3164019
8550db53,85,3784038
d19f5d58,11,468664
86c06c1c,4f,680754
539a2c41,38,2642167
eb5eecc7,a8,399907
21b23ed9,0a,359638
644ab84b,d1,1059569
f788e030,a9,2404052
8045912f,ec,4297469
e123e366,97,3024627
59f2fe0c,8a,4010937
a65964d7,54,277414
42cc64ed,2d,848047
76d3676f,da,5406685
7d752caf,c4,1111443
7ea05b83,d8,4870065
a3ee2530,40,690706
dd5300b9,c9,5525603
f1962502,79,301130
daea6a3a,71,4877920
f55d93c8,c6,407652
c3a2d1c0,f7,23591
665fbfc9,89,5263000
1b3eb128,ef,3405615
3f1883e8,08,1040237
811f4c6a,cd,946152
2eb6d3b5,1f,2284613
776b257d,27,5355332
d856f844,ce,4834758
a07f77a7,7c,2129797
9e9022f8,d0,3294896
b08e5978,c0,3474577
583f2969,53,2728720
fb45db3d,0c,4166997
753b55b1,72,607291
02b68e7c,cb,3252587
ad88ee1f,67,4665220
d73d603b,bd,1884546
87b535bc,70,6177535
19e25084,90,3611082
dcc24b17,9c,5521468
24a036ed,df,4729313
54ef03a6,11,3997975
a41f52c3,b6,226331
fe869744,ba,4802784
3c8e25cc,be,5378801
d66f9613,13,6003902
7c22c991,20,5694455
35b60335,f4,5446287
01cd5cb9,e9,6450913
99c9badf,07,4923387
546a361b,85,6009106
38e58949,fd,6575014
405e0e80,8c,4479728
0162534a,c3,3356598
e360ef3b,1d,1112191
1a2e57f6,0b,5657148
c681ef1d,30,6087976
5e15b25f,0d,529821
7ceca66f,ec,4492105
bfa46291,b2,1723116
8eb8576d,a0,4393506
78fdb1e3,d2,1606789
a6b75b58,d9,4935951
61c9b380,7a,2242750
7656130a,52,2337820
45a4cecb,be,1653187
546afe25,34,532193
72a1f7e2,03,4563792
353ae19e,87,1619658
b7d85c90,c3,4736657
3b56064b,f4,3731717
b34b390a,cd,7130589
a2ec36f7,99,2064693
f2987087,19,4044231
25c86a99,e5,5350161
865cc33c,49,846330
2837adb8,5d,3537951
d171102f,e4,4465760
94e35b6f,f8,7147928
dd2b4f06,6f,3427954
830231da,2d,2568022
1c010dbf,a4,6499800
c0bb89bb,8e,3345960
e0bb74ed,fe,24989
4e296a79,f3,2678832
d76fa208,b7,4667426
38c38e29,ea,6525588
ce1333fc,fd,1839399
a04d4cae,a7,5334293
21f9c82c,7a,3829726
2e2926e5,40,3837473
1c9f1434,e8,2438240
424be1a7,e9,5716750
9d00dde8,bc,1471569
c530cc13,12,5032844
160d1378,e6,332082
1d19a2f4,68,3377789
d3ef99ac,f1,5701499
bc40793f,64,5542971
81cf3023,38,5265336
80c22bb1,a2,2877696
4f857183,8c,1536626
f97ca916,e8,5089591
3462b6e6,ef,5260981
f26d6c45,ba,6153746
edb416ed,c6,2776124
55c54850,33,6284035
9521b83e,e9,1722788
20fd926a,9c,7876072
3a7304c2,bd,5878782
b631741a,ce,5628666
03474254,9f,2225539
2d921eaf,2e,4623461
9aca348a,4a,763927
f2615e63,eb,2279102
c8daa6eb,ed,7705031
739b1cd7,71,6047911
0f4b8225,bf,829353
954bf6e1,df,988141
c7b74bf9,89,3222911
dd0976da,5b,5158879
0e7a9a06,64,5326108
5d7374df,d7,6006474
b26333be,50,2611696
57bef898,34,7041434
e8d79151,e3,2626993
f03587d2,d9,3013677
db177e13,2c,1740972
0a05f606,c1,9063827
1a20e444,03,33068
8636854d,0c,7672682
db69dffe,7c,3084781
d75b7a73,b1,1757586
3c577841,21,4090696
8679540f,66,9015420
b8b0dc55,7e,8672115
c7613b6b,1d,7243295
eb17edd2,1f,6643117
d757978f,be,1155159
f67432e9,5b,6561774
5b0517c4,60,3685179
f358e704,54,7459189
444afa9f,45,9116061
91814f2a,de,8887231
cae269f4,a6,8388260
1cce632a,e8,1987620
355d242f,b6,7775957
abb6c270,81,5387472
29060037,c3,2932205
4bdde4c0,de,8955153
76ed63cf,bc,9031575
8cccde3d,d6,8706117
93cf1651,e8,4978569
3c1681f1,1d,5777557
3d1e6932,ef,1637296
05a0a047,8e,2764460
18f86553,a0,8754131
03d00082,ba,5001826
19c7b7f4,c6,4725110
f7eb2d80,f4,8679258
3dca21cf,8c,5086934
71fdce6b,d9,2994691
da35bebd,98,1862293
98349683,cd,7084454
73ef49f0,3c,7557396
3c687390,12,1449017
067a24df,bf,6304291
befd00b8,ec,5163415
633b1afd,f5,8815931
b0add1e1,ec,7641797
59e90297,ce,1886159
89ea1a35,17,2356570
1e02a211,48,5309096
07f37dbd,29,1295330
63cc7b5b,95,1481243
2976e573,88,3874002
830cee20,5a,3501637
d009c011,11,10076993
8147f0a6,fe,9853613
50a74d98,f3,6767761
ac9efa6b,68,9131250
85def5d7,14,8954007
c7b9a773,7e,5445547
264d449b,66,3688054
3b593b1a,81,4354482
86f6d5b2,95,7296939
1a56b087,cd,6763685
175360fb,d5,2884504
1ca3cc85,9f,1689826
1744357e,16,812986
f4366a09,13,1212152
32e97968,da,7295445
6c6700c2,ba,4110464
117f5069,a8,3268321
80c390d2,52,3768021
1fbf9254,90,1039391
b737dcca,dd,4886933
cfd20ca3,b1,10717532
7fa8aeef,c0,6711973
7659463a,02,3061920
0461f012,81,7619790
30810606,a5,8836525
5ccf064f,81,5670061
b241386c,92,1834643
19cd0454,39,10329917
c8dd04c9,5e,7339644
18f598f8,3e,6090246
73118e46,22,2451822
61d1e710,e0,6595522
9ba6f996,1f,5055583
7d512abb,6e,9871945
bfbcf038,a3,4704087
b90516e4,9b,2922228
586824dd,54,3909908
aa65e706,25,8771179
b2702fe8,49,7701290
879fc54f,df,11908142
ca0f6793,34,4100911
71e82e02,3c,10911575
e2cc3949,e5,1658708
83684c77,90,5020840
119d3575,c8,66283
43f5657e,77,9624969
95eceb47,11,9541033
f19490de,4c,3627962
d6c347e8,12,1420310
bb0643fa,fa,10749139
f5de1b40,07,8977848
7af40e60,9b,613774
caa122f5,93,5621617
0936023a,0c,12380304
189c5493,20,8907557
bde59e6f,49,2660431
8c0ae851,2f,5994023
cde9d3c8,be,10971549
0057302f,4e,1194091
6b0dc3f9,9a,1269057
304f234a,02,4893226
d9253b36,be,3866275
ba7a31de,cb,8440742
9afdf60d,e3,10419386
8eb070fb,8f,1700352
13af64ec,cc,666690
ac76e30a,0f,693214
813fd12e,d2,5954374
8c678b79,51,10183822
ce0e0a50,79,9048506
876a2a73,17,5800725
6435e652,e2,1565118
1f0d00e3,eb,7767128
fb927cef,e1,7641164
71bf5596,43,7014949
5e72e556,8c,7896869
83edef66,ae,4814872
c1dafbbd,23,4058126
119a02a5,0e,11881647
9f50ecf7,8a,6387240
9cc147da,b1,2281514
14cf5a9c,f4,8211617
e393b6f2,77,7097804
bb07dff2,56,8376031
3173e32a,ad,8087969
46264cd7,b3,3177425
e4b3582f,20,11311521
cf7a4d94,68,957730
26f9c539,56,1353778
cc91265d,3b,12682590
077399e5,89,2334811
f9c59a67,91,6645669
159b2e7b,8d,3118109
62eb0e4b,7d,11459738
a71ebbda,32,4040666
d80b47ba,61,4974413
1c0d85e0,81,12260025
47cf35cb,d2,4816737
91337af1,17,11428852
9f3b7d9b,2e,9265927
5fa41a8c,c8,511861
de907945,48,3636235
233136e8,ec,11016836
6a92e2e2,78,8881519
86cdeb51,14,11890492
61bf5d87,70,4392488
fefa844b,4b,3157527
87c06683,8d,2835578
f3165110,b0,8790002
0c13d2e7,a0,1578521
a028dc35,a7,9726404
a83112e4,79,11304228
e9e8cd91,66,10420293
33f9c4e4,34,2908359
e7948031,35,5836875
c444eb2a,67,4602849
d338d57a,3c,6342830
1217b188,bd,3383466
00cc040e,35,5662808
4d612be9,ed,7839673
b44e7aee,2c,2281494
8f645f68,fc,8207872
85a3f198,ed,215936
bf4f3cfe,9e,5643870
76330e9a,ad,5685375
b2054b4d,4a,4553462
71ff8e6f,e0,11018499
7f0cfd11,c3,3900155
33b8e1af,66,9749915
cb876baa,c1,8307121
2524ebd1,85,13861250
2dabc47d,f0,10362542
fd40de3d,30,12259190
c90755d0,94,10587462
3ddd577f,88,4743004
0847c120,e8,9841488
e8458084,eb,91692
ece34654,f5,8155701
18073ccd,cc,12941712
6ec35c57,13,9504408
21a214fd,cd,9898081
0c3c0ff8,07,8002617
9d87d1f6,bd,6197564
e27039a7,e3,5613125
2ca32dce,2b,11165468
6d9296f2,10,12481361
30cd5baf,fe,10572224
a585785c,d7,2389500
ae3d7016,ff,515201
c2df2550,6f,997638
d513f0bd,97,10785471
a0f76b82,3a,2957098
3766fcd9,b6,9775936
dae027ec,67,11083626
ecd6299b,5a,9005637
30345e9e,3a,17886
e08dbe11,2e,90919
7848dea5,4e,4871005
0c36e9d3,ee,4182714
c783f41c,49,10501324
1fc727b6,ee,13769180
15b6f24f,72,2914718
fd24fa28,97,5503267
a2f52cb5,ce,10780634
a0efc7d6,a3,14228870
9e64069d,3f,10596710
d962973d,bf,5090897
46940020,47,14413072
53f25a1d,59,3039555
77356282,45,2569419
69977993,c5,8248704
5c76220a,69,13608799
2c33b289,d5,11114195
a75ab8de,39,14902805
56c5a984,6b,3156222
b20746ab,8c,6528475
9e1efcd4,19,11707015
b4fced02,f1,14435761
51217f08,7a,9112050
27362bd4,e0,12325038
c16c83fa,79,15583384
dbbe84df,d2,3903447
e93ad2fe,96,15545362
6636ab11,d4,9269746
125fe089,21,15721361
175ffaf8,f6,8343373
61d89dc4,81,9342985
7a2da49b,21,15909839
4ddae271,47,13573045
faa817fa,4f,13845725
735317fd,92,10129853
c0299c4e,09,11219818
5de6039f,03,17024502
127d6796,de,16038312
180a68d8,25,5572516
5cd0f012,86,15040919
c94003ec,94,10154685
5082841e,81,11013498
8436842b,65,13137440
9ff91d3e,ce,13441366
e78ccb2c,cd,13805236
ffd8bd02,90,4572401
30954223,73,10535878
315f6abe,bd,8218565
35823588,0c,2921390
31b51478,ea,13945646
3d92903d,e5,15194057
d52da1da,33,15362046
06838099,ae,11713083
65434925,aa,7285247
009d8569,bb,10828827
53c7f037,53,10581519
dd1ee668,0a,11466047
54c2d95b,96,9582549
c203b65b,c8,2644312
43c39cd6,0e,10340490
d28187fb,b4,589758
26b217a2,7b,12563184
c0a9d977,f3,9304361
a6225e53,e2,14978331
aff3a6bd,51,9282420
1823d874,d9,2033877
7f8cce34,bf,8859651
7de30bd9,60,9247883
09d45eac,c8,880437
91ac5523,19,14809990
82a8fba4,61,7008094
5026e959,18,11638361
fb3edb39,da,8266851
d9bd6a86,f4,2637363
01b3521b,f8,7580011
76c2fe6d,a2,4541161
56ef9386,33,7490857
3a853233,67,322046
b5aed8a7,8d,13262014
900ed67c,a2,10621828
fb006b3f,df,15077474
2df7bf01,de,10825919
7822fb01,89,8314871
7c7a0f1d,8c,8900399
d8dcfca4,8e,18587450
76b7d441,69,1098742
f1efd277,26,1391885
30cee915,48,12031069
5ce155fa,3d,11234134
46c4b35a,13,4194184
3a997786,d3,4148356
44c862b7,21,11916093
92f65be0,ad,17491242
4598382e,c7,14000665
b8adb23a,5a,8796710
738f72a8,6a,8349430
6100261c,9d,16221513
9d435836,ed,18199415
7300ac56,f9,18158843
e6a6d891,a0,17234833
cf0098ee,ee,14728233
35b9e49f,77,6416142
fc62cec0,c4,10341731
9c14a26b,8d,4229403
2094b7ad,fd,9232555
341d5a0d,d9,18394114
ab0260ca,5c,9043044
60d6d64e,7d,17183554
e3af6656,1b,18827263
de6aa106,7a,3582865
af2bece6,c6,8945376
a5fb1b43,24,15263127
044ce340,55,15018742
0ce365d0,2b,17759080
70c03a91,1b,5064334
c88f1173,96,14998866
a8c5ca07,61,19006418
eeea8689,39,6556975
1eb4abe2,4d,7963829
51444b5a,fa,11117351
6d540709,06,11009903
23f4c485,21,12375003
cdb5f9bc,b6,3805833
49f1acdb,29,9538049
bafbaf4e,a7,7197953
9964f362,fc,13091491
5abfe8f3,9f,18379960
17aac4c7,3d,4232115
0708ec66,7d,12203356
8b51bbd0,c3,17440492
1804fd0d,42,2215872
fe8c171d,d9,7552498
968e7b91,d6,4029929
a267122a,5d,13113622
0b548d76,5d,12041223
9bedcd64,74,13767147
8f03878d,91,12154430
a82314f8,06,1201605
6a496a54,74,11468275
f52e5ffd,b7,11542885
9a666055,ab,4976927
5a0d6ee8,af,388531
db34126e,fc,7075672
07ac088a,29,187013
d2260aee,70,19068555
1f864153,f1,12435465
b29ff53f,07,15960423
802104d4,22,18402368
dbc29338,97,19388654
cb75a13b,43,3601073
979ff78c,ae,9336135
0b68d4d7,69,10180299
9c7e1f37,07,16625092
f8a46bf2,1b,13380055
e733a49e,ba,20251810
f862fb19,98,3262730
b76b6451,99,6161524
307c56a9,7c,18587888
b12ceee4,3a,7724138
7bd3c36b,9b,19009803
770ce9f5,19,15826091
fa8b18a5,cb,10087505
32930f35,3a,5796516
30aa381b,65,20131665
4e2ec830,32,2605441
6071b950,81,18466637
532d8be3,b0,13314385
f592e128,fb,15377665
c0bcd817,2c,19813133
11c0acf3,0e,3950620
4f68e970,48,6067917
44ed78e3,aa,3598831
9ad00b3d,d4,4376456
401164fc,54,13520616
4f0483df,4a,3811495
6f2fb35a,0b,2743611
680933dd,75,4936150
bf376906,60,20572112
03996c06,bb,301701
90b5416f,35,3504479
ce38e7d9,90,2618585
8ca49cb8,da,10094263
058627db,7e,14196019
3bbc88f4,5c,20063275
2780b7d1,63,2214270
e008e9ae,d2,13505204
21ad384f,f5,1553671
37a1a8dc,24,8683228
bafbfb93,8d,10963991
7cfce090,18,9656559
1b60c5fa,80,6494221
2ca19fe3,63,5414404
087010af,eb,20270310
fcdb92a6,18,12596861
bc1cc7e1,ad,4249428
4cc64abe,84,774695
22e614e8,c6,9650742
b6c9e60c,3d,18215606
cf27890d,30,16885801
1c9dcffd,6f,7999401
69ca97c5,90,22433213
e594a632,7b,802909
4ba4efe8,97,11803075
ffd0bcfc,48,4536808
23fc91be,84,14963054
34fab898,71,15845062
db067c80,21,7428844
a5b7f910,94,20195410
4924d751,98,14706868
b812dd05,9c,7776837
e2ce10f6,e5,1144155
ed5314ee,50,17510914
28739ce6,04,18373818
76ab4246,38,16003904
414dffc8,03,15641647
2ccd3756,d8,20769511
798f983b,d0,7198594
347b8739,c0,17334027
2e8d8fde,cf,6872923
4e00e44d,de,18697628
52f48d8b,51,4373836
06e32861,f0,13694245
81caa2f1,f2,1239780
1f104831,dd,2357034
adfa91f7,79,3659049
ae18b3d0,c8,5688186
8cc28b6f,8d,3176969
e1df7472,c7,13763922
702cef54,9d,1358211
35ddd73d,5d,20841448
55ba6f43,03,22183253
98eb7a60,79,1316221
2b4df980,5f,21768034
a5d1bf8d,4e,11050908
f79994ba,46,6062660
fd97becf,05,13868567
959da6ec,d4,3994433
99c5c823,71,13849561
ba3dfaf8,d0,11301053
43f6c4ac,11,19374834
aa8d8f11,3c,22030746
435c7b1d,11,12542997
3331335a,eb,13889637
42afd4b6,b9,17423758
a53149e5,b0,5387778
65fdab9c,e2,6327816
7db333fa,c4,15197277
a3e84be5,19,12926401
7a94b593,97,400474
32bed985,28,22435128
a06b5cb3,46,8381851
1ead56bb,99,4968001
2afdf8a8,42,12638616
2dd34445,64,9710081
3215a3e4,6e,7448316
d7b5de05,e0,16508802
27a73016,b0,15685694
ea0d578c,ed,10793804
149c6156,5a,6558601
20bde9e9,f0,17166295
dbc47bbe,f1,8137522
a8a975c8,a7,5251053
2e9d2f3b,3c,18608550
668f43ed,69,6334894
9966fbdf,64,16799514
60517ae7,af,398107
4402d4f0,0a,20825125
5dbabc5d,d9,21032652
154301c0,ce,14843085
d0fbd323,2d,7953303
26ce3726,20,16547115
7fefd1b2,16,13068220
e3a54132,be,21062748
ed0afb44,a2,6222958
20fa3a6f,b4,6381017
eb88f476,b1,1826237
a47ffd39,f6,12843719
5aac0af5,08,22247083
ff758228,66,15306521
9abe1128,e1,13311674
a5f32f3f,18,20615639
5b124184,e3,24101517
3e235506,8f,18540844
b84bfd9d,67,8667293
0d1c8d2d,eb,9081832
7eccb390,50,1109481
02e28832,53,603664
bcba3f83,7d,128578
e573d772,e7,12357950
7a83cf91,34,7437887
6a072ef3,4e,18009674
90e99f2d,ee,20881526
0955b08c,2e,12935647
5fccbc5e,69,11932788
5a586605,13,2000196
48547f39,41,4887614
21e55e34,1b,12393283
85bc73a3,52,23691497
36a28e4f,2f,1114815
1124fea7,17,2171293
4a6b39d4,02,6306484
5b8de5f0,e1,18344715
daef0a1b,70,15782971
7b474592,e8,2400756
1de3c158,f4,24891768
69c5bc69,e9,22484354
b99c32c6,b1,4310515
a80820e7,3e,12367503
02c9f0a3,da,19702800
788c5af7,80,20254440
e45632b6,54,182064
b01e5765,57,15111101
0236adca,84,19257641
d4ef797f,02,19471988
d6fda188,d0,24184671
51d1bdd7,08,2530335
c5381fe6,c6,20699191
28df0bc6,b0,24903269
b5bda587,8d,2673547
22c52422,ce,7251286
3ebd3a81,c2,6320702
4f831611,96,21412656
dcb51582,f1,2516403
0ed2dc9e,13,19795307
d1162673,db,7357826
c261be19,e9,4605437
34d79961,4d,24950330
b1723d97,79,20269655
166c47fe,27,2578797
b0c21cd6,e2,26201407
150f7d74,33,22601490
8f2cd2c1,e6,9143514
9e5768ad,f8,18560945
bda96392,3b,12251467
2b77a663,d3,11672172
fe40be2f,e5,17595023
762227e0,e6,8726584
3cc4948d,f5,16194110
1e6befc4,c5,3509459
469e39e8,35,2196886
e485a340,d5,10231009
aac955b8,c8,3711392
d9d85026,3e,17903735
6b335177,13,3109242
bc8fd09e,2d,23952327
4dd4c4a1,c0,15860224
f08b6166,ed,5748757
e6dbfc67,f7,18242443
7cb3085d,a4,14453635
a99937b0,f3,8546117
cdf588a1,3e,25838671
bbf81960,4b,21069434
17e7c5d8,c1,3352437
4bcb43b3,44,1051990
d624185d,ec,14173673
b66ce3c3,5b,11003861
8e53ccbf,7c,5837229
32f04e76,36,12503738
b573839a,fa,4429468
d1bbdae8,70,8535034
1bb381d0,7f,12024532
c642c4ca,d0,150842
2b8a07d8,5a,18055790
8344fd7e,21,3511496
687a7c52,d9,12866574
cd01d68e,3b,2347536
e703d104,8f,21783296
f4642f0a,4b,21967675
1225d9b5,cc,19159907
4b74d184,02,25501488
77ff2f6e,ce,24324372
4c95e4ce,8e,1760020
35ca7593,38,1326536
30cf94e6,1c,14835954
843e5288,93,4991692
9027a150,7b,26178107
1f631ba1,19,3093633
ddc1cbb1,53,10972058
a1e6ddba,a5,16879101
b536f38a,6d,18095281
7b85f93c,ee,19873598
0915132e,c9,15324657
8f378f16,59,11565199
74ecd327,f5,23418898
37ed3133,53,16468544
4b17a53d,66,26301920
795209fa,23,6335524
e36707a4,37,22577472
a2404f8e,5c,10716824
b3d12104,43,10949482
84793b23,72,2798444
23b864de,20,23218758
05f5781e,14,11143387
fcf74a66,92,21523876
5f666d66,87,26937897
5d4fb4c8,4d,24517289
64659f71,e4,24030389
6d0105a2,bc,20671468
4a18853b,99,18399896
ba969b9a,b2,14967853
1446655f,9d,18135355
ccec964e,4e,14602401
a15f9cf3,48,16842277
12a6e56b,41,19571226
e4f0f1d9,4c,859630
c25f9514,cc,22070802
1f5fe8ec,5f,17528856
427b1c78,cd,4067417
6183d56e,44,4216711
1adfec44,6b,5745167
54a347f6,80,22824864
17089972,9e,24815870
d887bc4f,fb,21637023
da94f840,e0,5421096
a8a3046a,d8,19661992
7b5f90ae,b5,9240693
afd48c3d,50,3120490
75fced94,b0,26288287
9a7aaf3a,f7,3028381
43e45e50,4e,14227810
43c5af78,bf,8025223
ab263220,a5,5935907
4ac8b830,f2,8561120
fbdaf473,16,19588629
76768881,e2,28238508
045cfd78,35,6624604
daa1ba89,9a,15148848
f71d74b4,ea,23289412
cc845bf0,4a,28043613
b0b403c4,82,4261101
2940ca72,85,8842601
bf8664d9,c1,16857292
ccaf46cc,79,19387816
2529c118,47,26301802
9a435eb9,0e,7019569
f5dee3cb,bd,20340748
489a3d78,2a,19014264
2ccc171d,5b,25180502
7212b174,b0,24615620
0a221ec6,7c,23439060
115bd9d5,eb,863340
6c326482,df,17934168
b8d60165,ef,18779873
7de01cfd,6d,23250095
610ab46d,f3,23600238
f94a4b6c,5b,8672815
34d13bca,64,2344030
60e8703c,f6,26904104
f95cc010,00,5286214
f429cd9a,97,5483807
39871c7e,47,11034937
1b4667bc,9f,13603010
13df6380,e7,11640190
babfecbb,b2,19923710
24efe3b3,f9,3101312
668f36ed,66,10043210
031f0a7d,b7,26837394
27ef0a0d,dd,18774986
2a28bffd,e5,3222419
18b57f26,10,24951868
667e83b6,83,23240182
b4c46e8d,76,12374520
d57e0fcd,7b,11394184
33eae5dc,61,12259576
c2ae527d,f4,8473144
b753befa,75,18686341
5658f2ac,0b,2128594
d3a7a256,99,28153473
605d3b35,6c,13510834
8323c143,c9,27576770
15c9992d,ad,21140477
09d9e60a,33,14445564
7d294afa,6f,14805564
77e5c3ea,d2,9189524
c89386cd,14,7059500
fd5fe8d4,ab,7748794
46044125,44,18179870
1b516216,b0,9503559
192988db,01,3982656
2f98029d,1f,23698232
9502bc05,25,15198788
da93b921,15,15395461
fcdcfcea,c2,6105538
97eba1c5,3e,27270877
f5c8ab6c,01,4239732
6ffcb36c,69,9328050
49843aae,84,5878882
9abe93f1,54,90159
31949f3b,08,28749061
e4b7265f,b7,9225267
b19e2e1d,22,25860832
688fdff0,f7,20825028
89304c05,94,22033252
ccf4c32c,bb,5627229
f31a1609,8c,1436381
b8870616,eb,30850475
190b1d06,13,24939938
6c77b94d,7d,3638686
8a12c5a1,66,27770380
e20c946b,5d,6689515
e5803541,18,32121724
a46ece1f,d2,17339026
429bc37a,b5,21623387
406a675d,65,15096564
cad534a9,00,26878307
31058083,ed,28976509
05abbcd7,7a,23141113
a23cb07d,05,13186841
0244e4fa,40,6302325
a106ba1b,29,27321826
ba540170,6b,5049275
fee2d945,7a,22007880
2054c31b,1e,10590746
aa9d6332,80,11365551
`
}
