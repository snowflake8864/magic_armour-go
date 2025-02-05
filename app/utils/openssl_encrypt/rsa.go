package openssl_encrypt

import (
	goEncryptRsa "github.com/wumansgy/goEncrypt/rsa"
	"strings"
)

func CreateRsaFactory() *rsaCrypt {
	/**
	-----BEGIN PUBLIC KEY-----
	MIIBCgKCAQEAwlCBJWEwP97VDdmu1XVxAm+NY9xgLM7sSH6n2CCfgoWQtbNE01tj
	i/zndl/zVnd2Ci06aX8aqssTAgqe3MKV5RIhD/gx0XVpM+5/1igFh2BuxdrFPNO1
	+GfBu6qnnWfPeSCikKI0yY6TKvM9GOiQXca1Sym8oouwY8oV3eWC/Fg5voNmsh78
	iDdsHjB/2YbEEQX+Io++3j3sxzeekNe9K21k0eCH4NoumsZ8S2xsocaEVDNelPYX
	Zv4emKgIuJCO+Qi6LrjTQ5MDXFLc0kWTHAKWwO3iF00V0dee/HEj4DAYQcWFcNIl
	Nwx2FKxXL089ZZbWmrkPiEzXpjVIQwdClQIDAQAB
	-----END PUBLIC KEY-----
	*/
	base64PubKey2048 := "" +
		"MIIBCgKCAQEAwlCBJWEwP97VDdmu1XVxAm+NY9xgLM7sSH6n2CCfgoWQtbNE01tj" +
		"i/zndl/zVnd2Ci06aX8aqssTAgqe3MKV5RIhD/gx0XVpM+5/1igFh2BuxdrFPNO1" +
		"+GfBu6qnnWfPeSCikKI0yY6TKvM9GOiQXca1Sym8oouwY8oV3eWC/Fg5voNmsh78" +
		"iDdsHjB/2YbEEQX+Io++3j3sxzeekNe9K21k0eCH4NoumsZ8S2xsocaEVDNelPYX" +
		"Zv4emKgIuJCO+Qi6LrjTQ5MDXFLc0kWTHAKWwO3iF00V0dee/HEj4DAYQcWFcNIl" +
		"Nwx2FKxXL089ZZbWmrkPiEzXpjVIQwdClQIDAQAB"

	base64PriKey2048 := "" +
		"MIIEpQIBAAKCAQEAwlCBJWEwP97VDdmu1XVxAm+NY9xgLM7sSH6n2CCfgoWQtbNE" +
		"01tji/zndl/zVnd2Ci06aX8aqssTAgqe3MKV5RIhD/gx0XVpM+5/1igFh2BuxdrF" +
		"PNO1+GfBu6qnnWfPeSCikKI0yY6TKvM9GOiQXca1Sym8oouwY8oV3eWC/Fg5voNm" +
		"sh78iDdsHjB/2YbEEQX+Io++3j3sxzeekNe9K21k0eCH4NoumsZ8S2xsocaEVDNe" +
		"lPYXZv4emKgIuJCO+Qi6LrjTQ5MDXFLc0kWTHAKWwO3iF00V0dee/HEj4DAYQcWF" +
		"cNIlNwx2FKxXL089ZZbWmrkPiEzXpjVIQwdClQIDAQABAoIBAQCUyl5JyCTRE4ki" +
		"5OfqTdZuHgsvR0VhqkTlsg5GAZ0CrdHEQHRYP2d3fv2h+BD1yd3aJw3OUGSUGZIB" +
		"gy5oLpvm3z3v44YKNJSNDnhW/ShWpDHOURms3RQ/FPwU/rUAY8hLoQC6LWqsV/86" +
		"6SyEkThShmC2IExwTvSPROP/mMn2/SUbCHIcKOoYQ6WS/viMjQhkoAyckmKt19MH" +
		"9muNC0eqtsF6Vz3lZnyuZ90yl0NSu88o6VY1b5cQ0xU7fYn0qloFBayRlIpPXDKy" +
		"4AiF4KRt9xDlbZH9seklqP4uexLdcZaDozpWoqVYMRIg9npakzNfRwuu7cMo3YBy" +
		"2sIXDZthAoGBAPBiKBGV+yU6hcOM0SGWiBRi7IG5SzLsrlJz8DW9yYelo+J3juBV" +
		"mo1ADEQ/ddqHSxP4yAuuZ8DK2H/+fhQqbEJl63HlR7YWq4Nq7g4BZC5q2UT7jBmO" +
		"+u9hdA3DUBpluJFQ3IAN40jyPmMY978AAbvZVb4iCSOFJXlOKy9+dz9ZAoGBAM7w" +
		"K4mRo4f2K/qdRriNdLyESrpX4JmiZ/nyGfbWhSnQqYxi+jOAL9n5f5dksszMop4e" +
		"EXAThnbjmGy2oBg+O6ZGR5w1IgYFXxrw/K6/20E+sfwXeF1iAVx4KaMsHUAE54rV" +
		"DE/pl8kzaZTY+5KfwjwkUARl2ICm0krErm+0pJGdAoGBAIdg8Gn/CqML03uvGbSN" +
		"33ODxujEHxq9+usKynSpvw769VI8vSWQcGsub/obt1bsLRdBXABgJ6SHWOi0vfLD" +
		"dZbraW8Kjpl5KPzpQsp2dCJEILrd5Cgcwyv/uWlq3vkO4qRMYvJnX+KcMxCCZjyW" +
		"V2Og4c4jwvvTJ6iqol2sHqJ5AoGAMD3RCHsJMnLzAbOUfwWLBz+X5CUEWQyBTMVh" +
		"xCdyV7MrM20HguYJQblmQzbYb+Hsrwu3aPdES4b3S7CvVbd8VGDh5MqGKf0qlK9F" +
		"WL+lL2qcN3aJOsvUIWhp3H42O2xpRfYwrph4KBFyNZ731tCjiT214bCUbjl1oLSF" +
		"Lclg7AkCgYEAtfbgbkeOyKK1n6BFKwRW7GHkfD0FIRA/x1B6VoLeQ5/bLshnb49N" +
		"DVDmYCEeMDAdYejNRp6V/Ljua2uyWdld9dLA0VqBRspysgcHinJ7z2tjJp2YTVtR" +
		"Zsf2gGUj/YHasTI54jEAlOjt6uUor4hbtSPLy4TsjXmNyWQMDImta48="

	hexPubKey2048 := "" +
		"3082010a0282010100aeca846494013ffaa983e1b4ce118fb9939dcb1d25fd7fb" +
		"54c4db55942783ae8c13fe79b7ac4a108f20ec418ec0c39702ebaa6b8ef097f6d" +
		"3c7a3d9839637139739055f0e58be38b4cc32d156f09ca98226449879dbedf123" +
		"f6790af03f42e96fce6a0e6869f1574353178ab71994e083784f98f0f368d7514" +
		"f8d03998b0224f393bd547d9f21e5c3b0d1b4540d28a618864ceae4dbdff1d9e2" +
		"0a838b03e64acee5e82bd750403aab9ca263d6d0b55dc03f50ce5716497a26110" +
		"b417e1fd05d7fefc06c3cc79e61ae299b0e1e1bc1caa71243cb657572330580db" +
		"d983b1af3f1c8e56251eb1fc3430e61a7f2ffc4a734c744effb9f44a74260fe8e" +
		"4bc4038b6b0203010001"
	hexPriKey2048 := "" +
		"308204a40201000282010100aeca846494013ffaa983e1b4ce118fb9939dcb1d2" +
		"5fd7fb54c4db55942783ae8c13fe79b7ac4a108f20ec418ec0c39702ebaa6b8ef" +
		"097f6d3c7a3d9839637139739055f0e58be38b4cc32d156f09ca98226449879db" +
		"edf123f6790af03f42e96fce6a0e6869f1574353178ab71994e083784f98f0f36" +
		"8d7514f8d03998b0224f393bd547d9f21e5c3b0d1b4540d28a618864ceae4dbdf" +
		"f1d9e20a838b03e64acee5e82bd750403aab9ca263d6d0b55dc03f50ce5716497" +
		"a26110b417e1fd05d7fefc06c3cc79e61ae299b0e1e1bc1caa71243cb65757233" +
		"0580dbd983b1af3f1c8e56251eb1fc3430e61a7f2ffc4a734c744effb9f44a742" +
		"60fe8e4bc4038b6b02030100010282010100971b8b432c54bc4be65a9638e9e17" +
		"2c03b66b1debc26204ecaed6f370068d27cbffc902cf100a63b06904f654043fd" +
		"3007a5e3568388517a4d4b5c8e8a26770ef4e235f9f1687349cd0dfa9eca04cbf" +
		"9a3d9d15b85bb2a3c2b464c6591647b169b0282dd684a8590961ff3b33b9d0e41" +
		"ba45a23dd4281438c0edad8ded7b6fdc60bba9958c23aff0bda29c2bcb8fafa95" +
		"cc46aa46c747475bedd8c4f1bd4562a370aea56e3e4406a62b514101a6f67e925" +
		"b28038cec1e62dd370126e57b9011e00d11972506c9e9b748da0bcb70694fb87f" +
		"2fe52569ef7d13ff1ec7ed220c2592e50c110ca7bdcb45e86bd80b2f5d89ef340" +
		"08a2ffce16595496fcbe178d881102818100c3874a2a2a95c243f0263d90b93a3" +
		"c521fdcc54c8f880599b57eaf1e9f10156593cd5d647dc5f2e37cce19511bd616" +
		"94db6fd85f9cae44794f2a9e13aa0b47bfb50853450464c382962a15a4fe8671a" +
		"be5f666e2276588af7ccb1ac3c79d43f6cae11129eeebe339f0f48524376165b4" +
		"c6e58bdd44db7d875a35a5d882b0f51302818100e4d95f21fa8f4ba682ff139ba" +
		"ba8a5a3f55a7898274294d9395930890b18db33dd6a4dc8a20769f495a242f2f2" +
		"b07e9790b77599a3c648c0510fec312f64d61fb0ba389e98a97cfb2363e952d6c" +
		"e2454d63352ec9d36165338b6c55d84e00e4563b7bb8dadc03c02efee0bda299d" +
		"f97ed3c6a2079f9f9f701ffdda73a915d349028180014e882a5bfa803fdffd017" +
		"0542c95769f0be4adac4fe669ae92d82cc6375f4d3c7da4e0250f0b9e85990306" +
		"4dbf528b50242138125dffcb27c24b282a7cc4cffd2fbf0eb906cb5fb60d35f3d" +
		"257a7aed46c9f2767c7ffefdc4fa73023c7634ae4f4381f259fa1f7232b1f30ac" +
		"c613ddef6e5d4eef44e043e55072794df64ccd02818100b2c2dd9e143404faffc" +
		"eb2896250ad7ae2b816057c7d74fa2de2ee97ce544f3e3751c78ad68e8604a932" +
		"03abf5494ac9a820b6bb1a045465edac4b6aac43cf4317f08f33cbba71475780b" +
		"f0c9780953eb3c7899a0cbfe8367fd1bd5fece07c2555dc7aa4143a5a664cf71c" +
		"c900c346d5fc2568021e7d38b504386e80f46e7909028180523ec2b65b8489707" +
		"1360980345d27cd0dfa0fe8a92f5cf6bae10809e9fb01630c0096ade42d951f99" +
		"4b2e555809b1c8b04ed2618c1ddf2daa58075ade2302ace9698530073ae475153" +
		"8789db3de9723baebf3f09c0bfe60db28ea8a61079008b51b2853760a5edd8b5a" +
		"77dc4f7ba45c48628433bc76f9fa0c38fa30b6cb6b71"

	return &rsaCrypt{
		base64PubKey2048: base64PubKey2048,
		base64PriKey2048: base64PriKey2048,
		hexPubKey2048:    hexPubKey2048,
		hexPriKey2048:    hexPriKey2048,
	}
}

type rsaCrypt struct {
	base64PubKey2048 string
	base64PriKey2048 string
	hexPubKey2048    string
	hexPriKey2048    string
}

// EncryptToBase64 加密并返回Base64
func (a *rsaCrypt) EncryptToBase64(data string) (result string, err error) {
	result, err = goEncryptRsa.RsaEncryptToBase64([]byte(data), a.base64PubKey2048)
	if err != nil {
		return "", err
	}
	return result, nil
}

// DecryptByBase64 解密
func (a *rsaCrypt) DecryptByBase64(data string) (string, error) {
	data = strings.ReplaceAll(data, " ", "+")
	plainText, err := goEncryptRsa.RsaDecryptByBase64(data, a.base64PriKey2048)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

// EncryptToHex 加密并返回Hex
func (a *rsaCrypt) EncryptToHex(data string) (result string, err error) {
	result, err = goEncryptRsa.RsaEncryptToHex([]byte(data), a.hexPubKey2048)
	if err != nil {
		return "", err
	}
	return result, nil
}

// DecryptByHex 解密
func (a *rsaCrypt) DecryptByHex(data string) (string, error) {
	data = strings.ReplaceAll(data, " ", "+")
	plainText, err := goEncryptRsa.RsaDecryptByHex(data, a.hexPriKey2048)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}
