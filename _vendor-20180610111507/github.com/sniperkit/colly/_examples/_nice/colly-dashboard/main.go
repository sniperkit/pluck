package main

import (
	"os"
	"path/filepath"

	// Logger
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	// collector - advanced sitemap parser
	sitemap "github.com/sniperkit/colly/plugins/data/format/sitemap"
)

// app params
var (
	version     string   = "0.0.1-alpha"
	workDir     string   = "."
	parentDir   string   = ".."
	configFiles []string = []string{
		// ./colly.yml
		"colly.yml",
		"colly.yaml",
		"colly.json",
		"colly.toml",
		// ../conf/app.yaml
		parentDir + "/conf/colly.yaml",
		parentDir + "/conf/app.yaml",
		parentDir + "/conf/cache.yaml",
		parentDir + "/conf/sitemap.yaml",
		parentDir + "/conf/collection.yaml",
		parentDir + "/conf/collector.yaml",
		parentDir + "/conf/debug.yaml",
		parentDir + "/conf/filters.yaml",
		parentDir + "/conf/outputs.yaml",
		parentDir + "/conf/proxy.yaml",
		parentDir + "/conf/transport.yaml",
		// ./conf/app.yaml
		workDir + "/conf/colly.yaml",
		workDir + "/conf/cache.yaml",
		workDir + "/conf/app.yaml",
		workDir + "/conf/sitemap.yaml",
		workDir + "/conf/collection.yaml",
		workDir + "/conf/collector.yaml",
		workDir + "/conf/debug.yaml",
		workDir + "/conf/filters.yaml",
		workDir + "/conf/outputs.yaml",
		workDir + "/conf/proxy.yaml",
		workDir + "/conf/transport.yaml",
	}
	log *logrus.Logger
)

// Initialize collector and other components
func init() {

	// Set the logger
	log = logrus.New()
	log.Formatter = new(prefixed.TextFormatter)
	log.Level = logrus.DebugLevel

}

func getWorkDir() string {
	// Get the current workdir
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func main() {

	masterCollector, err := newCollectorWithConfig(configFiles...)
	if err != nil {
		log.Println("could not instanciate the master collector.")
	}

	masterCollector = addCollectorEvents(masterCollector)

	// Initialize data collections for storing data/pattern extracted
	// or the sitemap urls by the collector into custom datasets
	initDataCollections()

	if appConfig.App.DashboardMode {
		appConfig.App.DebugMode = false
		appConfig.App.VerboseMode = false
		initDashboard()
		updateDashboard()
	}

	if appConfig.Debug.Tachymeter.Enabled {
		initTachymeter()
	}

	switch appConfig.Collector.CurrentMode {
	case "async":

		if appConfig.Collector.Sitemap.Enabled && appConfig.Collector.Sitemap.URL != "" {

			// Attach master collector to the sitemap collector
			sitemapCollector, err := sitemap.AttachCollector(appConfig.Collector.Sitemap.URL, masterCollector)
			if err != nil {
				log.Fatalln("could not instanciate the sitemap collector.")
			}
			sitemapCollector.VisitAll()
			sitemapCollector.Count()

		} else {

			masterCollector.Visit(appConfig.Collector.RootURL)

		}

		// Consume URLs
		masterCollector.Wait()

	case "queue":

		// Initialize collector queue
		collectorQueue, err := initCollectorQueue(appConfig.Collector.Modes.Queue.WorkersCount, appConfig.Collector.Modes.Queue.MaxSize, "InMemory")
		if err != nil {
			log.Fatalln("error: ", err)
		}

		if appConfig.Collector.Sitemap.Enabled && appConfig.Collector.Sitemap.URL != "" {

			// Attach queue and master collector to the sitemap collector
			sitemapCollector, err := sitemap.AttachQueue(appConfig.Collector.Sitemap.URL, collectorQueue)
			if err != nil {
				log.Fatalln("could not instanciate the sitemap collector.")
			}
			sitemapCollector.Count()
			// Enqueue all URLs found in the sitemap.txt
			sitemapCollector.EnqueueAll()

		} else {

			collectorQueue.AddURL(appConfig.Collector.RootURL)

		}

		// Consume URLs
		collectorQueue.Run(masterCollector)

	default:

		if appConfig.Collector.Sitemap.Enabled && appConfig.Collector.Sitemap.URL != "" {

			// Initalize new sitemap collector
			sitemapCollector, err := sitemap.New(appConfig.Collector.Sitemap.URL)
			if err != nil {
				log.Fatalln("could not instanciate the sitemap collector.")
			}
			sitemapCollector.Count()
			urls, _ := sitemapCollector.List()
			for _, url := range urls {
				masterCollector.Visit(url.String())
			}

		} else {

			masterCollector.Visit(appConfig.Collector.RootURL)

		}

	}

	uiWaitGroup.Wait()

	// if enableUI && !masterCollector.IsDebug() {
	if appConfig.App.DashboardMode {
		stopTheUI <- true
	}

	if appConfig.Debug.Tachymeter.Enabled {
		err := closeTachymeter("./shared/exports/stats/tachymeter")
		if err != nil {
			log.Fatalln("error: ", err)
		}
	}

}
