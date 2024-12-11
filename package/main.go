package cloudfunction

import (
	"riccardotornesello.it/iracing-average-lap/sessions_downloader"
)

func SessionsDownloader() {
	sessions_downloader.Run()
}
