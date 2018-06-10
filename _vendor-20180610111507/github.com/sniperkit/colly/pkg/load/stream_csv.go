package load

import (
	"net/http"
	"os"
)

func linksFromCSV(filePath string) ([]string, error) {

	isRemote := isRemoteURL(filePath)
	var reader *csv.Reader
	if !isRemote {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return nil, err
		}
		if !enable_ui {
			log.Infoln("reading file:", filePath)
		}
		file, err := os.Open(filePath)
		if err != nil {
			if !enable_ui {
				log.Fatalln("failed to open file, error: ", err)
			}
			return nil, err
		}
		defer file.Close()
		reader = csv.NewReader(file)

	} else {
		if !enable_ui {
			log.Infoln("loading remote:", filePath)
		}
		resp, err := http.Get(filePath)
		if err != nil {
			log.Fatalln("failed to fetch content, error: ", err)
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			log.Fatalln("failed with status", resp.Status)
			return nil, err
		}
		reader = csv.NewReader(resp.Body)

	}

	lines := streamCsv(reader, csvStreamBuffer)
	links := make([]string, len(lines))
	for line := range lines {
		links = append(links, line.GetByKey(0))
		if !enable_ui {
			log.Infoln("[LIST-ROW] col[0]=", line.GetByKey(0))
		}
	}
	if !enable_ui {
		log.Infoln("links pre-loaded:", len(links))
	}
	return links, nil

}
