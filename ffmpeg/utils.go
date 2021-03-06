package ffmpeg

func GenReBoradcastParams(pullUrl string, pushUrl string, cookie string) []string {
	var params []string
	if cookie != "" {
		params = append(params, "-headers", `Cookie: `+cookie)
	}
	//params = append(params, "-progress", "udp://127.0.0.1:64750")
	params = append(params, "-re")
	params = append(params, "-i", pullUrl)
	params = append(params, "-c", "copy")
	params = append(params, "-f", "flv")
	params = append(params, pushUrl)

	return params
}
